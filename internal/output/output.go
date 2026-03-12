package output

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"

	"github.com/user/vpsbench/internal/bench"
)

// HeaderData содержит данные для шапки отчёта.
// Не зависит от sysinfo напрямую — данные передаются из main.go.
type HeaderData struct {
	OSVersion string // "Ubuntu 22.04", "macOS 15.3"
	Kernel    string // "6.1.0-18-amd64"
	Arch      string // "amd64", "arm64"
	CPUModel  string // "AMD Ryzen 9 5900X"
	CPUCores  int
	RAM       string // "16 GB"
	Disks     string // "NVMe 500 GB, SSD 1 TB"
	Location  string // "Frankfurt, Germany"
	PublicIP  string // "203.0.113.42"
}

var (
	colorRed    = lipgloss.Color("#FF4444")
	colorOrange = lipgloss.Color("#FF8800")
	colorWhite  = lipgloss.Color("#CCCCCC")
	colorYellow = lipgloss.Color("#FFFF00")
	colorGreen  = lipgloss.Color("#00FF00")
	colorDim    = lipgloss.Color("#555555")
	colorCyan   = lipgloss.Color("#00CCCC")
)

var noColor bool

// SetNoColor включает или выключает использование ANSI-цветов.
func SetNoColor(v bool) {
	noColor = v
	if v {
		lipgloss.SetColorProfile(termenv.Ascii)
	}
	slog.Debug("[output] no-color mode set", "value", v)
}

// ColorForPercent возвращает цвет по шкале 0-100%.
func ColorForPercent(percent int) lipgloss.Color {
	if noColor {
		return lipgloss.Color("")
	}
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

	if noColor {
		return "[" + strings.Repeat("█", filled) + strings.Repeat("░", empty) + "]"
	}

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
func RenderHeader(data HeaderData, overallRating int) string {
	width := 70
	titleText := "VPSBENCH v1.0"

	var sb strings.Builder

	if !noColor {
		titleStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(colorCyan).
			Width(width).
			Align(lipgloss.Center)

		headerBox := lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder(), true, false).
			BorderForeground(colorCyan).
			Padding(0, 1)

		sb.WriteString(headerBox.Render(titleStyle.Render(titleText)) + "\n")
	} else {
		line := strings.Repeat("=", width)
		sb.WriteString(line + "\n")
		sb.WriteString(fmt.Sprintf(" %s\n", titleText))
		sb.WriteString(line + "\n")
	}

	// Системная информация
	sysInfoStyle := lipgloss.NewStyle().PaddingLeft(1)
	sb.WriteString(sysInfoStyle.Render(fmt.Sprintf("OS:       %-20s Kernel: %-15s Arch: %s", data.OSVersion, data.Kernel, data.Arch)) + "\n")
	sb.WriteString(sysInfoStyle.Render(fmt.Sprintf("CPU:      %s (%d Cores)", data.CPUModel, data.CPUCores)) + "\n")
	sb.WriteString(sysInfoStyle.Render(fmt.Sprintf("RAM:      %s", data.RAM)) + "\n")
	sb.WriteString(sysInfoStyle.Render(fmt.Sprintf("Disk:     %s", data.Disks)) + "\n")

	locationText := data.Location
	if data.PublicIP != "" {
		if !noColor {
			dimStyle := lipgloss.NewStyle().Foreground(colorDim)
			locationText = fmt.Sprintf("%s %s", data.Location, dimStyle.Render("("+data.PublicIP+")"))
		} else {
			locationText = fmt.Sprintf("%s (%s)", data.Location, data.PublicIP)
		}
	}
	sb.WriteString(sysInfoStyle.Render(fmt.Sprintf("Location: %s", locationText)) + "\n")

	separator := strings.Repeat("-", width)
	sb.WriteString(separator + "\n")

	// Общий рейтинг
	ratingText := fmt.Sprintf("SYSTEM RATING: %d%%", overallRating)
	if !noColor {
		ratingColor := ColorForPercent(overallRating)
		ratingStyle := lipgloss.NewStyle().Bold(true).Foreground(ratingColor).PaddingLeft(1)
		sb.WriteString(ratingStyle.Render(ratingText) + "\n")
	} else {
		sb.WriteString(fmt.Sprintf(" %s\n", ratingText))
	}
	sb.WriteString(sysInfoStyle.Render("Baseline 100%: 8-core Ryzen 9, NVMe Gen4, 10Gbps, 16GB DDR5") + "\n")
	sb.WriteString(separator)

	return sb.String()
}

// RenderModuleResult рисует результаты одного модуля.
func RenderModuleResult(mr bench.ModuleResult) string {
	if mr.Err != nil {
		errText := fmt.Sprintf("  ERROR: %v", mr.Err)
		if !noColor {
			errStyle := lipgloss.NewStyle().Foreground(colorRed)
			errText = errStyle.Render(errText)
		}
		return fmt.Sprintf("\n [ %s ] %s\n%s",
			mr.Module,
			mr.Info,
			errText,
		)
	}

	var sb strings.Builder

	moduleText := fmt.Sprintf("[ %s ]", mr.Module)
	if !noColor {
		moduleStyle := lipgloss.NewStyle().Bold(true)
		moduleText = moduleStyle.Render(moduleText)
	}
	sb.WriteString(fmt.Sprintf("\n %s %s\n", moduleText, mr.Info))

	for _, r := range mr.Results {
		bar := RenderProgressBar(r.Percent, 20)

		pctText := fmt.Sprintf("%d%%", r.Percent)
		if !noColor {
			pctColor := ColorForPercent(r.Percent)
			pctStyle := lipgloss.NewStyle().Foreground(pctColor)
			pctText = pctStyle.Render(pctText)
		}

		label := fmt.Sprintf("%-12s", r.Name)
		value := fmt.Sprintf(": %s %s", formatValue(r.Value), r.Unit)

		sb.WriteString(fmt.Sprintf("  %-28s %s %s\n",
			label+value,
			bar,
			pctText,
		))
	}

	return sb.String()
}

// RenderFooter рисует нижнюю рамку.
func RenderFooter() string {
	width := 70
	if !noColor {
		footerBox := lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder(), false, false, true, false).
			BorderForeground(colorCyan).
			Width(width + 2) // +2 за счет Padding(0, 1) в хедере (визуальное соответствие)
		return footerBox.Render("")
	}
	return strings.Repeat("=", width)
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
