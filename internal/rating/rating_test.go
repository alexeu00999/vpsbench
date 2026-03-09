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

func TestCalculateClamp(t *testing.T) {
	baseline := Baseline{
		Values: map[string]float64{
			"CPU:test": 100,
		},
	}

	results := []bench.ModuleResult{
		{
			Module: "CPU",
			Results: []bench.Result{
				{Name: "test", Value: 200, Unit: "ops/s"}, // 200% → clamped to 100
			},
		},
	}

	calculated := Calculate(results, baseline)

	if calculated[0].Results[0].Percent != 100 {
		t.Errorf("expected 100%% (clamped), got %d%%", calculated[0].Results[0].Percent)
	}
}

func TestCalculateLatency(t *testing.T) {
	baseline := Baseline{
		Values: map[string]float64{
			"NETWORK:Ping (DNS)": 1.0, // 1ms = 100%
		},
	}

	results := []bench.ModuleResult{
		{
			Module: "NETWORK",
			Results: []bench.Result{
				{Name: "Ping (DNS)", Value: 2.0, Unit: "ms"}, // 2ms → 50%
			},
		},
	}

	calculated := Calculate(results, baseline)

	if calculated[0].Results[0].Percent != 50 {
		t.Errorf("expected 50%% for latency, got %d%%", calculated[0].Results[0].Percent)
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

func TestOverallRatingSkipsErrors(t *testing.T) {
	results := []bench.ModuleResult{
		{
			Module: "OK",
			Results: []bench.Result{
				{Name: "test", Percent: 50},
			},
		},
		{
			Module: "FAIL",
			Err:    errTest,
			Results: []bench.Result{
				{Name: "test", Percent: 90}, // должен быть проигнорирован
			},
		},
	}

	overall := OverallRating(results)
	if overall != 50 {
		t.Errorf("expected 50 (skip errored module), got %d", overall)
	}
}

var errTest = benchError("test error")

type benchError string

func (e benchError) Error() string { return string(e) }
