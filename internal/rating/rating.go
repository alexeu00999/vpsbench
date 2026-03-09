package rating

import (
	"log/slog"

	"github.com/user/vpsbench/internal/bench"
)

// Baseline содержит эталонные значения для расчёта процентов.
type Baseline struct {
	// Ключ: "CPU:Single-core", "DISK:Seq. Read" и т.д.
	// Значение: эталонное значение для 100%.
	Values map[string]float64
}

// DefaultBaseline возвращает эталонную конфигурацию:
// 8-core Ryzen 9, NVMe Gen4, 10 Gbps, 16 GB DDR5.
func DefaultBaseline() Baseline {
	return Baseline{
		Values: map[string]float64{
			"CPU:Single-core":  5000,    // ops/s для одного ядра Ryzen 9
			"CPU:Multi-core":   40000,   // ops/s для 8 ядер
			"RAM:Speed (R/W)":  50000,   // MB/s DDR5
			"DISK:Seq. Read":   7000,    // MB/s NVMe Gen4
			"DISK:Seq. Write":  5000,    // MB/s NVMe Gen4
			"DISK:Rand 4K IOPS": 1000000, // IOPS NVMe Gen4
			"NETWORK:Ping (DNS)": 1.0,   // ms (инвертированная шкала — меньше = лучше)
		},
	}
}

// Calculate заполняет поле Percent в каждом Result на основе baseline.
func Calculate(results []bench.ModuleResult, baseline Baseline) []bench.ModuleResult {
	slog.Debug("[rating] calculating percentages", "modules", len(results))

	for i, mr := range results {
		if mr.Err != nil {
			slog.Debug("[rating] skipping module with error", "module", mr.Module)
			continue
		}
		for j, r := range mr.Results {
			key := mr.Module + ":" + r.Name
			baseVal, ok := baseline.Values[key]
			if !ok {
				slog.Debug("[rating] no baseline for metric", "key", key)
				continue
			}

			var pct int
			if r.Unit == "ms" {
				// Для латентности: меньше = лучше
				if r.Value <= 0 {
					pct = 100
				} else {
					pct = int(baseVal / r.Value * 100)
				}
			} else {
				// Для пропускной способности: больше = лучше
				if baseVal <= 0 {
					pct = 0
				} else {
					pct = int(r.Value / baseVal * 100)
				}
			}

			// Clamp 0-100
			if pct < 0 {
				pct = 0
			}
			if pct > 100 {
				pct = 100
			}

			results[i].Results[j].Percent = pct
			slog.Debug("[rating] calculated", "key", key, "value", r.Value, "baseline", baseVal, "percent", pct)
		}
	}

	return results
}

// OverallRating вычисляет средний рейтинг системы (0-100).
func OverallRating(results []bench.ModuleResult) int {
	total := 0
	count := 0

	for _, mr := range results {
		if mr.Err != nil {
			continue
		}
		for _, r := range mr.Results {
			if r.Percent > 0 {
				total += r.Percent
				count++
			}
		}
	}

	if count == 0 {
		return 0
	}

	overall := total / count
	slog.Debug("[rating] overall rating", "total", total, "count", count, "overall", overall)
	return overall
}
