package bench

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestReportJSON(t *testing.T) {
	report := Report{
		Timestamp: "2026-03-10T12:00:00Z",
		SystemInfo: SystemInfo{
			OSVersion: "Ubuntu 22.04",
			CPUModel:  "AMD Ryzen 9",
			CPUCores:  12,
		},
		ModuleResults: []ModuleResult{
			{
				Module: "CPU",
				Results: []Result{
					{Name: "Single-core", Value: 1000, Unit: "ops/s", Percent: 80},
				},
			},
		},
		OverallRating: 75,
	}

	data, err := json.Marshal(report)
	if err != nil {
		t.Fatalf("failed to marshal report: %v", err)
	}

	jsonStr := string(data)

	// Проверяем наличие ключевых полей
	expectedFields := []string{
		"\"timestamp\":\"2026-03-10T12:00:00Z\"",
		"\"os_version\":\"Ubuntu 22.04\"",
		"\"cpu_model\":\"AMD Ryzen 9\"",
		"\"cpu_cores\":12",
		"\"module\":\"CPU\"",
		"\"percent\":80",
		"\"overall_rating\":75",
	}

	for _, field := range expectedFields {
		if !strings.Contains(jsonStr, field) {
			t.Errorf("expected JSON to contain %s, but it didn't: %s", field, jsonStr)
		}
	}
}

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
