package ram

import (
	"context"
	"log/slog"
	"time"
)

const (
	// bufferSize — размер буфера для тестов (128 MB).
	bufferSize = 128 * 1024 * 1024 // байты

	// bufferElements — количество uint64 элементов в буфере.
	bufferElements = bufferSize / 8
)

// workloadWrite измеряет скорость последовательной записи в память.
// Заполняет буфер []uint64 целиком в цикле, возвращает MB/s.
func workloadWrite(ctx context.Context, duration time.Duration) float64 {
	slog.Debug("[ram] workloadWrite started",
		"buffer_size_mb", bufferSize/1024/1024,
		"elements", bufferElements,
		"duration", duration,
	)

	buf := make([]uint64, bufferElements)
	start := time.Now()
	deadline := start.Add(duration)
	iterations := 0

	for time.Now().Before(deadline) {
		// Последовательная запись — заполняем весь буфер
		for i := 0; i < bufferElements; i++ {
			buf[i] = uint64(i) ^ 0xDEADBEEFCAFEBABE
		}
		iterations++

		if ctx.Err() != nil {
			slog.Debug("[ram] workloadWrite cancelled by context")
			break
		}
	}

	elapsed := time.Since(start).Seconds()
	totalBytes := float64(iterations) * float64(bufferSize)
	mbPerSec := totalBytes / elapsed / (1024 * 1024)

	// Предотвращаем оптимизацию — используем последний элемент
	sinkValue = buf[bufferElements-1]

	slog.Debug("[ram] workloadWrite finished",
		"iterations", iterations,
		"total_bytes_mb", int(totalBytes/1024/1024),
		"mb_per_sec", mbPerSec,
		"elapsed", elapsed,
	)

	return mbPerSec
}

// workloadRead измеряет скорость последовательного чтения из памяти.
// Читает весь буфер []uint64, суммируя значения (предотвращение оптимизации).
// Возвращает MB/s.
func workloadRead(ctx context.Context, duration time.Duration) float64 {
	slog.Debug("[ram] workloadRead started",
		"buffer_size_mb", bufferSize/1024/1024,
		"elements", bufferElements,
		"duration", duration,
	)

	// Предварительно заполняем буфер (чтобы читать реальные данные)
	buf := make([]uint64, bufferElements)
	for i := 0; i < bufferElements; i++ {
		buf[i] = uint64(i) ^ 0xCAFEBABEDEADBEEF
	}

	start := time.Now()
	deadline := start.Add(duration)
	iterations := 0
	var sum uint64

	for time.Now().Before(deadline) {
		// Последовательное чтение — читаем весь буфер
		var localSum uint64
		for i := 0; i < bufferElements; i++ {
			localSum += buf[i]
		}
		sum += localSum
		iterations++

		if ctx.Err() != nil {
			slog.Debug("[ram] workloadRead cancelled by context")
			break
		}
	}

	elapsed := time.Since(start).Seconds()
	totalBytes := float64(iterations) * float64(bufferSize)
	mbPerSec := totalBytes / elapsed / (1024 * 1024)

	// Предотвращаем оптимизацию
	sinkValue = sum

	slog.Debug("[ram] workloadRead finished",
		"iterations", iterations,
		"total_bytes_mb", int(totalBytes/1024/1024),
		"mb_per_sec", mbPerSec,
		"elapsed", elapsed,
	)

	return mbPerSec
}

// sinkValue предотвращает оптимизацию компилятором (dead code elimination).
var sinkValue uint64
