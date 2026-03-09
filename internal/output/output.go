package output

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/user/vpsbench/internal/bench"
	"github.com/user/vpsbench/internal/sysinfo"
)

var (
	colorRed    = lipgloss.Color("#FF4444")
	colorOrange = lipgloss.Color("#FF8800")
	colorWhite  = lipgloss.Color("#CCCCCC")
	colorYellow = lipgloss.Color("#FFFF00")
	colorGreen  = lipgloss.Color("#00FF00")
	colorDim    = lipgloss.Color("#555555")
	colorCyan   = lipgloss.Color("#00CCCC")
)

// ColorForPercent возвращает цвет по шкале 0-100%.
func ColorForPercent(percent int) lipgloss.Color {
	switch {
	case percent <= 20:
		return colorRed
	case percent <= 40:
		return colorOrange
	case percent <= 60:
		return colorWhite
	case percent <= 80:
		return colorYellow
	default:
		return colorGreen
	}
}

// RenderProgressBar рисует прогресс-бар [█████░░░░░] с цветом.
func RenderProgressBar(percent int, width int) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}

	filled := width * percent / 100
	empty := width - filled

	color := ColorForPercent(percent)
	filledStyle := lipgloss.NewStyle().Foreground(color)
	emptyStyle := lipgloss.NewStyle().Foreground(colorDim)

	bar := "[" +
		filledStyle.Render(strings.Repeat("█", filled)) +
		emptyStyle.Render(strings.Repeat("░", empty)) +
		"]"

	slog.Debug("[output] rendered progress bar", "percent", percent, "width", width, "filled", filled)
	return bar
}

// RenderHeader рисует шапку отчёта.
func RenderHeader(info sysinfo.SystemInfo, overallRating int) string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(colorCyan)
	line := strings.Repeat("=", 70)

	header := fmt.Sprintf("%s\n %s | Location: %s | OS: %s %s\n%s\n",
		line,
		titleStyle.Render("SUPER-BENCH v1.0"),
		info.Location,
		info.OS,
		info.Arch,
		line,
	)

	ratingColor := ColorForPercent(overallRating)
	ratingStyle := lipgloss.NewStyle().Bold(true).Foreground(ratingColor)
	header += fmt.Sprintf(" SYSTEM RATING: %s\n", ratingStyle.Render(fmt.Sprintf("%d%%", overallRating)))
	header += " Baseline 100%: 8-core Ryzen 9, NVMe Gen4, 10Gbps, 16GB DDR5\n"
	header += strings.Repeat("-", 70)

	return header
}

// RenderModuleResult рисует результаты одного модуля.
func RenderModuleResult(mr bench.ModuleResult) string {
	if mr.Err != nil {
		errStyle := lipgloss.NewStyle().Foreground(colorRed)
		return fmt.Sprintf("\n[ %s ] %s\n%s",
			mr.Module,
			mr.Info,
			errStyle.Render(fmt.Sprintf("  ERROR: %v", mr.Err)),
		)
	}

	moduleStyle := lipgloss.NewStyle().Bold(true)
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("\n%s %s\n", moduleStyle.Render(fmt.Sprintf("[ %s ]", mr.Module)), mr.Info))

	for _, r := range mr.Results {
		bar := RenderProgressBar(r.Percent, 20)
		pctColor := ColorForPercent(r.Percent)
		pctStyle := lipgloss.NewStyle().Foreground(pctColor)

		label := fmt.Sprintf("%-12s", r.Name)
		value := fmt.Sprintf(": %s %s", formatValue(r.Value), r.Unit)

		sb.WriteString(fmt.Sprintf("%-28s %s %s\n",
			label+value,
			bar,
			pctStyle.Render(fmt.Sprintf("%d%%", r.Percent)),
		))
	}

	return sb.String()
}

// RenderFooter рисует нижнюю рамку.
func RenderFooter() string {
	return strings.Repeat("=", 70)
}

// formatValue форматирует число с разделителями тысяч.
func formatValue(v float64) string {
	if v >= 1000000 {
		return fmt.Sprintf("%.0f", v)
	}
	if v >= 1000 {
		whole := int(v)
		thousands := whole / 1000
		remainder := whole % 1000
		return fmt.Sprintf("%d %03d", thousands, remainder)
	}
	if v == float64(int(v)) {
		return fmt.Sprintf("%.0f", v)
	}
	return fmt.Sprintf("%.1f", v)
}
