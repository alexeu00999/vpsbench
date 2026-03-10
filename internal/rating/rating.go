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
			// CPU: Ryzen 9 5900X (или аналогичный 8-ядерный современный CPU)
			"CPU:Single-core": 200000000,  // 200M ops/s
			"CPU:Multi-core":  1500000000, // 1.5B ops/s

			// RAM: DDR5-5600 dual-channel
			"RAM:Write": 40000, // 40 GB/s
			"RAM:Read":  45000, // 45 GB/s

			// DISK: NVMe Gen4
			"DISK:Seq. Read":   7000,    // 7 GB/s
			"DISK:Seq. Write":  5000,    // 5 GB/s
			"DISK:Rand 4K IOPS": 1000000, // 1M IOPS

			// NETWORK: 10 Gbps Port (Калибровано под крупные облака)
			"NETWORK:Ping (EU)":     1.0,   // ms (инвертированная шкала)
			"NETWORK:Ping (US)":     80.0,  // ms (типичный трансатлантик из EU)
			"NETWORK:Ping (ASIA)":   180.0, // ms (типичный пинг в Азию из EU)
			"NETWORK:Download (EU)": 8000,  // Mbps (~8 Gbps на 10G порту)
			"NETWORK:Upload (EU)":   6000,  // Mbps
			"NETWORK:Download (US)": 1500,  // Mbps (через океан)
			"NETWORK:Upload (US)":   1000,  // Mbps
			"NETWORK:Download (ASIA)": 500,  // Mbps
			"NETWORK:Upload (ASIA)":   300,  // Mbps
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
				// Для латентности: меньше = лучше (инвертированная шкала)
				if r.Value <= 0 {
					pct = 100
				} else {
					// Если задержка в 2 раза меньше базовой — это 200%
					pct = int(baseVal / r.Value * 100)
				}
			} else {
				// Для пропускной способности и операций: больше = лучше
				if baseVal <= 0 {
					pct = 0
				} else {
					pct = int(r.Value / baseVal * 100)
				}
			}

			// Clamp только снизу (отрицательный процент не имеет смысла)
			if pct < 0 {
				pct = 0
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

// GetRatingLabel возвращает текстовое описание для процента.
func GetRatingLabel(percent int) string {
	switch {
	case percent <= 20:
		return "Poor"
	case percent <= 40:
		return "Below Average"
	case percent <= 60:
		return "Average"
	case percent <= 80:
		return "Good"
	default:
		return "Excellent"
	}
}
