package output

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestColorForPercent(t *testing.T) {
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

func TestRenderProgressBar(t *testing.T) {
	bar := RenderProgressBar(50, 20)

	// Должен содержать [ и ]
	if !strings.HasPrefix(bar, "[") || !strings.HasSuffix(bar, "]") {
		t.Errorf("progress bar should be wrapped in [], got: %s", bar)
	}

	// При 0% не должно быть заполненных блоков (█)
	bar0 := RenderProgressBar(0, 10)
	if strings.Contains(bar0, "█") {
		t.Errorf("0%% bar should have no filled blocks")
	}

	// При 100% не должно быть пустых блоков (░)
	bar100 := RenderProgressBar(100, 10)
	if strings.Contains(bar100, "░") {
		t.Errorf("100%% bar should have no empty blocks")
	}
}

func TestRenderProgressBarClamp(t *testing.T) {
	// Не должен паниковать при отрицательных или >100 значениях
	bar := RenderProgressBar(-10, 20)
	if bar == "" {
		t.Error("negative percent should still render")
	}

	bar = RenderProgressBar(150, 20)
	if bar == "" {
		t.Error(">100 percent should still render")
	}
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
