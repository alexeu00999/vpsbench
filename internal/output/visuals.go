package output

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Braille characters constants.
const (
	brailleBase = 0x2800
)

// brailleOffsetMap maps y index (0-3) and x index (0-1) to dot position bit.
// 0 3
// 1 4
// 2 5
// 6 7
var brailleOffsetMap = [4][2]int{
	{0, 3},
	{1, 4},
	{2, 5},
	{6, 7},
}

// RenderSparkline рисует мини-график с использованием шрифта Брайля.
// data: значения от 0.0 до 1.0. width: количество символов в ширину.
func RenderSparkline(data []float64, width int) string {
	if len(data) == 0 {
		return ""
	}

	// Каждый символ Брайля имеет ширину 2 точки.
	// Итоговое количество точек по горизонтали = width * 2.
	dotsX := width * 2
	height := 4 // 4 точки по вертикали в одном символе

	// Масштабируем данные под dotsX.
	scaledData := make([]float64, dotsX)
	for i := 0; i < dotsX; i++ {
		idx := int(float64(i) * float64(len(data)) / float64(dotsX))
		if idx >= len(data) {
			idx = len(data) - 1
		}
		scaledData[i] = data[idx]
	}

	var result string
	for i := 0; i < width; i++ {
		char := rune(brailleBase)
		// Для каждой колонки в символе (левая и правая)
		for x := 0; x < 2; x++ {
			val := scaledData[i*2+x]
			// Количество закрашенных точек по вертикали (0-4)
			filledDots := int(math.Round(val * float64(height)))
			if filledDots < 0 {
				filledDots = 0
			}
			if filledDots > height {
				filledDots = height
			}
			for y := 0; y < filledDots; y++ {
				// y=0 это нижняя точка в Брайле для графиков обычно, 
				// но в шрифте Брайля биты расположены сверху вниз.
				// Инвертируем y, чтобы график рос снизу вверх.
				bit := brailleOffsetMap[3-y][x]
				char |= rune(1 << bit)
			}
		}
		result += string(char)
	}

	return result
}

// InterpolateColor вычисляет промежуточный цвет между двумя цветами.
func InterpolateColor(start, end string, ratio float64) lipgloss.Color {
	var r1, g1, b1, r2, g2, b2 int
	fmt.Sscanf(start, "#%02x%02x%02x", &r1, &g1, &b1)
	fmt.Sscanf(end, "#%02x%02x%02x", &r2, &g2, &b2)

	r := int(float64(r1) + ratio*float64(r2-r1))
	g := int(float64(g1) + ratio*float64(g2-g1))
	b := int(float64(b1) + ratio*float64(b2-b1))

	return lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", r, g, b))
}

// RenderGradientBar рисует прогресс-бар с градиентом от зеленого к красному.
func RenderGradientBar(percent int, width int) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}

	filled := width * percent / 100
	empty := width - filled

	var sb string
	// Цвета btop-style: Green -> Yellow -> Red
	colorGreen := "#00FF00"
	colorYellow := "#FFFF00"
	colorRed := "#FF4444"

	for i := 0; i < filled; i++ {
		ratio := float64(i) / float64(width)
		var c lipgloss.Color
		if ratio < 0.5 {
			c = InterpolateColor(colorGreen, colorYellow, ratio*2)
		} else {
			c = InterpolateColor(colorYellow, colorRed, (ratio-0.5)*2)
		}
		sb += lipgloss.NewStyle().Foreground(c).Render("█")
	}

	if empty > 0 {
		dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#333333"))
		sb += dimStyle.Render(strings.Repeat("░", empty))
	}

	return " " + sb + " "
}
