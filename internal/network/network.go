package network

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/user/vpsbench/internal/bench"
)

// NetworkBench реализует bench.Benchmark для тестирования сети.
type NetworkBench struct{}

// New создаёт новый экземпляр NetworkBench.
func New() *NetworkBench {
	return &NetworkBench{}
}

// Name возвращает имя модуля.
func (n *NetworkBench) Name() string {
	return "NETWORK"
}

// Run выполняет сетевой бенчмарк (заглушка с пинг-тестом).
func (n *NetworkBench) Run(ctx context.Context) bench.ModuleResult {
	slog.Info("[network] starting benchmark")
	result := bench.ModuleResult{
		Module:  "NETWORK",
		Info:    "Network detection pending", // TODO: реальный детект интерфейсов
	}

	// Ping test (HTTP latency to well-known DNS)
	slog.Debug("[network] running ping test")
	pingMs, err := measurePing(ctx, "https://1.1.1.1")
	if err != nil {
		slog.Error("[network] ping test failed", "error", err)
		result.Err = fmt.Errorf("ping test: %w", err)
		return result
	}
	result.Results = append(result.Results, bench.Result{
		Name:  "Ping (DNS)",
		Value: pingMs,
		Unit:  "ms",
	})
	slog.Debug("[network] ping result", "latency_ms", pingMs)

	slog.Info("[network] benchmark completed")
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
