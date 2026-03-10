package network

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/user/vpsbench/internal/bench"
)

// Server представляет конфигурацию сервера для тестирования.
type Server struct {
	ID          string // "hetzner-de", "aws-us-east-1"
	Region      string // "EU", "US", "ASIA"
	Provider    string // "Hetzner", "AWS", "DigitalOcean"
	URL         string // URL для загрузки/отгрузки
	LatencyHost string // Хост для пинга (если отличается от URL)
}

// RegionResult содержит результаты тестов для одного региона.
type RegionResult struct {
	Region   string
	Latency  float64
	Download float64
	Upload   float64
}

// defaultServers — список серверов для бенчмарка по умолчанию.
var defaultServers = []Server{
	{
		ID:       "hetzner-de",
		Region:   "EU",
		Provider: "Hetzner",
		URL:      "https://nbg1-speed.hetzner.com/100MB.bin",
	},
	{
		ID:       "aws-us-east-1",
		Region:   "US",
		Provider: "AWS",
		URL:      "https://s3.amazonaws.com/speedtest-east-1/100mb.bin",
	},
	{
		ID:       "do-sgp",
		Region:   "ASIA",
		Provider: "DigitalOcean",
		URL:      "http://speedtest-sgp1.digitalocean.com/100mb.test",
	},
}

// NetworkBench реализует bench.Benchmark для тестирования сети.
type NetworkBench struct {
	servers []Server
}

// New создаёт новый экземпляр NetworkBench.
func New() *NetworkBench {
	slog.Debug("[network] initializing with default servers", "count", len(defaultServers))
	return &NetworkBench{
		servers: defaultServers,
	}
}

// Name возвращает имя модуля.
func (n *NetworkBench) Name() string {
	return "NETWORK"
}

// Run выполняет сетевой бенчмарк для всех настроенных регионов.
func (n *NetworkBench) Run(ctx context.Context) bench.ModuleResult {
	slog.Info("[network] starting benchmark", "servers", len(n.servers))
	result := bench.ModuleResult{
		Module: "NETWORK",
		Info:   "Multi-region speedtest (EU, US, Asia)",
	}

	for _, srv := range n.servers {
		select {
		case <-ctx.Done():
			result.Err = ctx.Err()
			return result
		default:
		}

		slog.Info("[network] testing region", "region", srv.Region, "provider", srv.Provider)

		// 1. Пинг
		latencyMs, err := measurePing(ctx, srv.URL)
		if err != nil {
			slog.Warn("[network] ping failed", "region", srv.Region, "error", err)
			continue
		}
		result.Results = append(result.Results, bench.Result{
			Name:  fmt.Sprintf("Ping (%s)", srv.Region),
			Value: latencyMs,
			Unit:  "ms",
		})

		// 2. Загрузка
		dlMbps, err := measureDownload(ctx, srv.URL)
		if err != nil {
			slog.Warn("[network] download failed", "region", srv.Region, "error", err)
		} else {
			result.Results = append(result.Results, bench.Result{
				Name:  fmt.Sprintf("Download (%s)", srv.Region),
				Value: dlMbps,
				Unit:  "Mbps",
			})
		}

		// 3. Отдача
		ulMbps, err := measureUpload(ctx, srv.URL)
		if err != nil {
			slog.Warn("[network] upload failed", "region", srv.Region, "error", err)
		} else {
			result.Results = append(result.Results, bench.Result{
				Name:  fmt.Sprintf("Upload (%s)", srv.Region),
				Value: ulMbps,
				Unit:  "Mbps",
			})
		}
	}

	slog.Info("[network] benchmark completed", "results", len(result.Results))
	return result
}

// measurePing измеряет HTTP latency до указанного URL.
func measurePing(ctx context.Context, url string) (float64, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	if err != nil {
		return 0, err
	}

	// Среднее из 3 попыток
	var totalMs float64
	attempts := 3
	success := 0

	for i := 0; i < attempts; i++ {
		start := time.Now()
		resp, err := client.Do(req)
		elapsed := time.Since(start)

		if err != nil {
			slog.Debug("[network] ping attempt failed", "attempt", i+1, "error", err)
			continue
		}
		resp.Body.Close()

		ms := float64(elapsed.Microseconds()) / 1000.0
		totalMs += ms
		success++
		slog.Debug("[network] ping attempt", "attempt", i+1, "latency_ms", ms)
	}

	if success == 0 {
		return 0, fmt.Errorf("all %d ping attempts failed", attempts)
	}

	return totalMs / float64(success), nil
}

// measureDownload измеряет скорость загрузки в Mbps.
func measureDownload(ctx context.Context, url string) (float64, error) {
	slog.Debug("[network] starting download test", "url", url)
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, fmt.Errorf("create request: %w", err)
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("get url: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("bad status: %s", resp.Status)
	}

	// Читаем первые 10MB или пока не кончится
	const limit = 10 * 1024 * 1024 // 10MB
	n, err := io.Copy(io.Discard, io.LimitReader(resp.Body, limit))
	if err != nil && err != io.EOF {
		return 0, fmt.Errorf("read body: %w", err)
	}

	elapsed := time.Since(start)
	slog.Debug("[network] download completed", "bytes", n, "duration", elapsed)

	// Скорость в Mbps
	// n * 8 (bits) / elapsed (seconds) / 1,000,000
	mbps := (float64(n) * 8) / elapsed.Seconds() / 1_000_000

	return mbps, nil
}

// measureUpload измеряет скорость отдачи в Mbps.
func measureUpload(ctx context.Context, url string) (float64, error) {
	slog.Debug("[network] starting upload test", "url", url)

	// Подготавливаем 2MB данных
	const size = 2 * 1024 * 1024
	data := make([]byte, size)
	body := bytes.NewReader(data)

	start := time.Now()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return 0, fmt.Errorf("create request: %w", err)
	}
	req.ContentLength = int64(size)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		// Некоторые серверы могут запрещать POST, тогда пробуем другой подход или пропускаем
		return 0, fmt.Errorf("post data: %w", err)
	}
	defer resp.Body.Close()

	// Нам не важно, какой статус, если данные ушли (для бенчмарка)
	// Но обычно 200, 201, 204 или даже 405 (Method Not Allowed) после того как тело прочитано.

	elapsed := time.Since(start)
	slog.Debug("[network] upload completed", "bytes", size, "duration", elapsed)

	// Скорость в Mbps
	mbps := (float64(size) * 8) / elapsed.Seconds() / 1_000_000

	return mbps, nil
}
