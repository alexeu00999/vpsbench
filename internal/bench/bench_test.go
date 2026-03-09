package bench

import "testing"

func TestResult(t *testing.T) {
	r := Result{
		Name:    "Single-core",
		Value:   1240,
		Unit:    "ops/s",
		Percent: 45,
	}

	if r.Name != "Single-core" {
		t.Errorf("expected Name=Single-core, got %s", r.Name)
	}
	if r.Value != 1240 {
		t.Errorf("expected Value=1240, got %f", r.Value)
	}
	if r.Unit != "ops/s" {
		t.Errorf("expected Unit=ops/s, got %s", r.Unit)
	}
	if r.Percent != 45 {
		t.Errorf("expected Percent=45, got %d", r.Percent)
	}
}

func TestModuleResult(t *testing.T) {
	mr := ModuleResult{
		Module: "CPU",
		Info:   "Test CPU",
		Results: []Result{
			{Name: "test", Value: 100, Unit: "ops/s"},
		},
		Err: nil,
	}

	if mr.Module != "CPU" {
		t.Errorf("expected Module=CPU, got %s", mr.Module)
	}
	if len(mr.Results) != 1 {
		t.Errorf("expected 1 result, got %d", len(mr.Results))
	}
	if mr.Err != nil {
		t.Errorf("expected nil error, got %v", mr.Err)
	}
}
