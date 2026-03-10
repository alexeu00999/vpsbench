package rating

import (
	"testing"

	"github.com/user/vpsbench/internal/bench"
)

func TestCalculate(t *testing.T) {
	baseline := Baseline{
		Values: map[string]float64{
			"CPU:Single-core": 1000,
			"CPU:Multi-core":  4000,
		},
	}

	results := []bench.ModuleResult{
		{
			Module: "CPU",
			Results: []bench.Result{
				{Name: "Single-core", Value: 500, Unit: "ops/s"},  // 50%
				{Name: "Multi-core", Value: 3000, Unit: "ops/s"},  // 75%
			},
		},
	}

	calculated := Calculate(results, baseline)

	if calculated[0].Results[0].Percent != 50 {
		t.Errorf("expected 50%%, got %d%%", calculated[0].Results[0].Percent)
	}
	if calculated[0].Results[1].Percent != 75 {
		t.Errorf("expected 75%%, got %d%%", calculated[0].Results[1].Percent)
	}
}

func TestCalculateNoClamp(t *testing.T) {
	baseline := Baseline{
		Values: map[string]float64{
			"CPU:test": 100,
		},
	}

	results := []bench.ModuleResult{
		{
			Module: "CPU",
			Results: []bench.Result{
				{Name: "test", Value: 250, Unit: "ops/s"}, // 250%
			},
		},
	}

	calculated := Calculate(results, baseline)

	if calculated[0].Results[0].Percent != 250 {
		t.Errorf("expected 250%% (no clamp), got %d%%", calculated[0].Results[0].Percent)
	}
}

func TestCalculateLatency(t *testing.T) {
	baseline := Baseline{
		Values: map[string]float64{
			"NETWORK:Ping (EU)": 1.0, // 1ms = 100%
		},
	}

	results := []bench.ModuleResult{
		{
			Module: "NETWORK",
			Results: []bench.Result{
				{Name: "Ping (EU)", Value: 0.5, Unit: "ms"}, // 0.5ms → 200%
				{Name: "Ping (EU)", Value: 2.0, Unit: "ms"}, // 2ms → 50%
			},
		},
	}

	calculated := Calculate(results, baseline)

	if calculated[0].Results[0].Percent != 200 {
		t.Errorf("expected 200%% for low latency, got %d%%", calculated[0].Results[0].Percent)
	}
	if calculated[0].Results[1].Percent != 50 {
		t.Errorf("expected 50%% for high latency, got %d%%", calculated[0].Results[1].Percent)
	}
}

func TestOverallRating(t *testing.T) {
	results := []bench.ModuleResult{
		{
			Module: "CPU",
			Results: []bench.Result{
				{Name: "a", Percent: 40},
				{Name: "b", Percent: 60},
			},
		},
		{
			Module: "RAM",
			Results: []bench.Result{
				{Name: "c", Percent: 80},
			},
		},
	}

	overall := OverallRating(results)
	// (40 + 60 + 80) / 3 = 60
	if overall != 60 {
		t.Errorf("expected overall 60, got %d", overall)
	}
}

func TestGetRatingLabel(t *testing.T) {
	tests := []struct {
		percent  int
		expected string
	}{
		{10, "Poor"},
		{30, "Below Average"},
		{50, "Average"},
		{70, "Good"},
		{90, "Excellent"},
		{150, "Excellent"},
	}

	for _, tt := range tests {
		got := GetRatingLabel(tt.percent)
		if got != tt.expected {
			t.Errorf("GetRatingLabel(%d) = %s; want %s", tt.percent, got, tt.expected)
		}
	}
}

func TestDefaultBaselineContents(t *testing.T) {
	bl := DefaultBaseline()
	required := []string{
		"CPU:Single-core",
		"CPU:Multi-core",
		"NETWORK:Download (EU)",
		"DISK:Seq. Read",
	}

	for _, req := range required {
		if _, ok := bl.Values[req]; !ok {
			t.Errorf("DefaultBaseline missing required key: %s", req)
		}
	}
}
