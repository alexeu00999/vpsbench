package output

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/user/vpsbench/internal/bench"
)

func TestColorForPercent(t *testing.T) {
	// Сбрасываем noColor для чистоты теста
	SetNoColor(false)

	tests := []struct {
		percent int
		want    lipgloss.Color
	}{
		{0, colorRed},
		{10, colorRed},
		{20, colorRed},
		{21, colorOrange},
		{40, colorOrange},
		{41, colorWhite},
		{60, colorWhite},
		{61, colorYellow},
		{80, colorYellow},
		{81, colorGreen},
		{100, colorGreen},
	}

	for _, tt := range tests {
		got := ColorForPercent(tt.percent)
		if got != tt.want {
			t.Errorf("ColorForPercent(%d) = %v, want %v", tt.percent, got, tt.want)
		}
	}
}

func TestNoColorMode(t *testing.T) {
	SetNoColor(true)
	defer SetNoColor(false)

	// 1. Проверка ColorForPercent
	got := ColorForPercent(50)
	if string(got) != "" {
		t.Errorf("ColorForPercent in no-color mode should return empty string, got %q", got)
	}

	// 2. Проверка RenderProgressBar (не должно быть ANSI кодов \x1b)
	bar := RenderProgressBar(50, 10)
	if strings.Contains(bar, "\x1b") {
		t.Errorf("RenderProgressBar contains ANSI codes in no-color mode: %q", bar)
	}

	// 3. Проверка RenderHeader
	header := RenderHeader(HeaderData{OSVersion: "Linux"}, 50)
	if strings.Contains(header, "\x1b") {
		t.Errorf("RenderHeader contains ANSI codes in no-color mode")
	}

	// 4. Проверка RenderModuleResult
	mr := bench.ModuleResult{
		Module: "TEST",
		Results: []bench.Result{
			{Name: "Metric", Value: 100, Unit: "ops", Percent: 50},
		},
	}
	res := RenderModuleResult(mr)
	if strings.Contains(res, "\x1b") {
		t.Errorf("RenderModuleResult contains ANSI codes in no-color mode")
	}
}

func TestRenderProgressBar(t *testing.T) {
	SetNoColor(false)
	bar := RenderProgressBar(50, 20)

	// Должен содержать [ и ]
	if !strings.HasPrefix(bar, "[") || !strings.HasSuffix(bar, "]") {
		t.Errorf("progress bar should be wrapped in [], got: %s", bar)
	}

	// При 0% не должно быть заполненных блоков (█) в версии без цвета (проверим без цвета для надежности логики)
	SetNoColor(true)
	bar0 := RenderProgressBar(0, 10)
	if strings.Contains(bar0, "█") {
		t.Errorf("0%% bar should have no filled blocks")
	}

	// При 100% не должно быть пустых блоков (░)
	bar100 := RenderProgressBar(100, 10)
	if strings.Contains(bar100, "░") {
		t.Errorf("100%% bar should have no empty blocks")
	}
	SetNoColor(false)
}

func TestFormatValue(t *testing.T) {
	tests := []struct {
		value float64
		want  string
	}{
		{100, "100"},
		{1500, "1 500"},
		{45000, "45 000"},
		{1.2, "1.2"},
		{0, "0"},
		{1000000, "1000000"},
	}

	for _, tt := range tests {
		got := formatValue(tt.value)
		if got != tt.want {
			t.Errorf("formatValue(%v) = %q, want %q", tt.value, got, tt.want)
		}
	}
}
