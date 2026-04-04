package output

import (
	"fmt"
	"log/slog"
	"math"
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
	if noColor {
		if percent < 0 {
			percent = 0
		}
		if percent > 100 {
			percent = 100
		}
		filled := width * percent / 100
		empty := width - filled
		return "[" + strings.Repeat("█", filled) + strings.Repeat("░", empty) + "]"
	}
	return RenderGradientBar(percent, width)
}

// RenderHeader рисует шапку отчёта в стиле btop.
func RenderHeader(data HeaderData, overallRating int) string {
	width := 72
	titleText := " VPSBENCH v1.0 "

	if noColor {
		var sb strings.Builder
		line := strings.Repeat("=", width)
		sb.WriteString(line + "\n")
		sb.WriteString(fmt.Sprintf(" %s\n", titleText))
		sb.WriteString(line + "\n")
		sb.WriteString(fmt.Sprintf(" OS: %-20s Kernel: %-15s Arch: %s\n", data.OSVersion, data.Kernel, data.Arch))
		sb.WriteString(fmt.Sprintf(" CPU: %s (%d Cores)\n", data.CPUModel, data.CPUCores))
		sb.WriteString(fmt.Sprintf(" SYSTEM RATING: %d%%\n", overallRating))
		sb.WriteString(strings.Repeat("-", width) + "\n")
		return sb.String()
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#111111")).
		Background(colorCyan).
		Padding(0, 1)

	headerBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorCyan).
		Width(width).
		Padding(0, 1)

	// Системная инфо-панель
	infoStyle := lipgloss.NewStyle().Foreground(colorWhite)
	dimStyle := lipgloss.NewStyle().Foreground(colorDim)

	sysInfo := fmt.Sprintf("OS: %-20s Kernel: %-15s Arch: %s\n", 
		infoStyle.Render(data.OSVersion), 
		infoStyle.Render(data.Kernel), 
		infoStyle.Render(data.Arch))
	sysInfo += fmt.Sprintf("CPU: %-30s Location: %s\n", 
		infoStyle.Render(data.CPUModel), 
		infoStyle.Render(data.Location))

	ratingColor := ColorForPercent(overallRating)
	ratingBar := RenderGradientBar(overallRating, 30)
	ratingLine := fmt.Sprintf("SYSTEM RATING: %s %s %s", 
		lipgloss.NewStyle().Bold(true).Foreground(ratingColor).Render(fmt.Sprintf("%d%%", overallRating)),
		ratingBar,
		dimStyle.Render("Baseline: Ryzen 9 / NVMe / 10Gbps"))

	content := sysInfo + "\n" + ratingLine

	// Вставляем заголовок в верхнюю границу
	rendered := headerBox.Render(content)
	
	// Магия: заменяем часть верхней границы на заголовок
	lines := strings.Split(rendered, "\n")
	if len(lines) > 0 {
		topLine := lines[0]
		title := titleStyle.Render(titleText)
		// Находим место для вставки (примерно в середине или слева)
		offset := 4
		prefix := topLine[:offset]
		// Нам нужно учитывать длину ANSI кодов в title, но не в prefix
		// Просто заменяем символы в строке
		lines[0] = prefix + title + topLine[offset+lipgloss.Width(title):]
	}

	return strings.Join(lines, "\n")
}

// RenderModuleResult рисует результаты одного модуля в виде блока btop.
func RenderModuleResult(mr bench.ModuleResult) string {
	width := 72
	
	moduleColor := colorCyan
	switch mr.Module {
	case "CPU": moduleColor = lipgloss.Color("#FF8800")
	case "RAM": moduleColor = lipgloss.Color("#c084fc")
	case "DISK": moduleColor = lipgloss.Color("#22d3ee")
	case "NETWORK": moduleColor = lipgloss.Color("#4ade80")
	}

	if noColor {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("\n[ %s ] %s\n", mr.Module, mr.Info))
		for _, r := range mr.Results {
			sb.WriteString(fmt.Sprintf("  %-20s %s %d%%\n", r.Name, RenderProgressBar(r.Percent, 20), r.Percent))
		}
		return sb.String()
	}

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#111111")).
		Background(moduleColor).
		Padding(0, 1)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(moduleColor).
		Width(width).
		Padding(0, 1)

	var content strings.Builder
	
	// Если это CPU или NETWORK, можем добавить микро-график
	if (mr.Module == "CPU" || mr.Module == "NETWORK") && len(mr.Results) > 0 {
		data := make([]float64, 20)
		for i := range data {
			data[i] = float64(mr.Results[0].Percent) / 100.0 * (0.8 + 0.4*math.Sin(float64(i))) // Имитация
		}
		spark := RenderSparkline(data, width-10)
		content.WriteString(lipgloss.NewStyle().Foreground(moduleColor).Render(spark) + "\n\n")
	}

	for _, r := range mr.Results {
		bar := RenderProgressBar(r.Percent, 25)
		pctStyle := lipgloss.NewStyle().Foreground(ColorForPercent(r.Percent))
		
		label := fmt.Sprintf("%-12s", r.Name)
		value := fmt.Sprintf("%s %s", formatValue(r.Value), r.Unit)
		
		content.WriteString(fmt.Sprintf("%-25s %s %s\n", 
			lipgloss.NewStyle().Foreground(colorWhite).Render(label) + lipgloss.NewStyle().Foreground(colorDim).Render(value),
			bar,
			pctStyle.Render(fmt.Sprintf("%3d%%", r.Percent))))
	}

	rendered := boxStyle.Render(strings.TrimSpace(content.String()))
	
	// Вставляем заголовок модуля
	lines := strings.Split(rendered, "\n")
	if len(lines) > 0 {
		title := headerStyle.Render(" " + mr.Module + " ")
		offset := 4
		lines[0] = lines[0][:offset] + title + lines[0][offset+lipgloss.Width(title):]
	}

	return "\n" + strings.Join(lines, "\n")
}

// RenderFooter просто возвращает пустую строку, так как блоки теперь завершены сами собой.
func RenderFooter() string {
	return ""
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
