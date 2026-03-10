package ui

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/user/vpsbench/internal/sysinfo"
)

// Item представляет элемент выбора (модуль или конкретный диск).
type Item struct {
	ID       string
	Label    string
	Selected bool
	IsDisk   bool // Если true, это конкретный диск
}

// Model реализует интерфейс tea.Model для выбора компонентов.
type Model struct {
	Items    []Item
	Cursor   int
	Finished bool
	Timer    int
	Quitting bool
}

// SelectionResult содержит итоги выбора пользователя.
type SelectionResult struct {
	Modules []string
	Disks   []string
}

// tickMsg отправляется каждую секунду для обновления таймера.
type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// SelectComponents запускает интерактивный выбор компонентов и дисков.
func SelectComponents(availableModules []string, disks []sysinfo.DiskInfo) (SelectionResult, error) {
	slog.Info("[ui] starting interactive selector", "modules", availableModules, "disks", len(disks))

	var items []Item
	for _, name := range availableModules {
		if name == "DISK" && len(disks) > 0 {
			// Если есть диски, добавляем их как подпункты
			for i, d := range disks {
				selected := false
				if i == 0 {
					selected = true // Первый диск выбран по умолчанию
				}
				label := fmt.Sprintf("Диск: %s (%s)", d.Device, formatBytes(d.Size))
				items = append(items, Item{ID: d.Device, Label: label, Selected: selected, IsDisk: true})
			}
			continue
		}
		items = append(items, Item{ID: name, Label: name, Selected: true})
	}

	m := Model{
		Items: items,
		Timer: 3,
	}

	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return SelectionResult{}, err
	}

	m = finalModel.(Model)
	if m.Quitting {
		return SelectionResult{}, fmt.Errorf("interrupted by user")
	}

	res := SelectionResult{}
	hasAnyDisk := false
	for _, item := range m.Items {
		if item.Selected {
			if item.IsDisk {
				res.Disks = append(res.Disks, item.ID)
				hasAnyDisk = true
			} else {
				res.Modules = append(res.Modules, item.ID)
			}
		}
	}

	if hasAnyDisk {
		res.Modules = append(res.Modules, "DISK")
	}

	slog.Debug("[ui] selection complete", "modules", res.Modules, "disks", res.Disks)
	return res, nil
}

func formatBytes(bytes uint64) string {
	const gb = 1024 * 1024 * 1024
	if bytes == 0 {
		return "0 GB"
	}
	return fmt.Sprintf("%.0f GB", float64(bytes)/float64(gb))
}

func (m Model) Init() tea.Cmd {
	return tick()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Сбрасываем таймер при любой активности пользователя
		m.Timer = -1

		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.Quitting = true
			return m, tea.Quit

		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}

		case "down", "j":
			if m.Cursor < len(m.Items)-1 {
				m.Cursor++
			}

		case " ", "enter":
			if msg.String() == "enter" {
				m.Finished = true
				return m, tea.Quit
			}
			m.Items[m.Cursor].Selected = !m.Items[m.Cursor].Selected
		}

	case tickMsg:
		if m.Timer > 0 {
			m.Timer--
			return m, tick()
		} else if m.Timer == 0 {
			m.Finished = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.Quitting {
		return ""
	}

	var sb strings.Builder

	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00CCCC"))
	sb.WriteString(headerStyle.Render("\n 🛠  Выберите компоненты для тестирования\n"))

	if m.Timer >= 0 {
		timerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF8800"))
		sb.WriteString(fmt.Sprintf(" Автостарт через: %s...\n", timerStyle.Render(fmt.Sprintf("%d сек", m.Timer))))
	} else {
		sb.WriteString(" (стрелки + Пробел для выбора, Enter для старта)\n")
	}
	sb.WriteString("\n")

	for i, item := range m.Items {
		cursor := " "
		if m.Cursor == i {
			cursor = ">"
		}

		checked := " "
		if item.Selected {
			checked = "x"
		}

		style := lipgloss.NewStyle()
		if m.Cursor == i {
			style = style.Foreground(lipgloss.Color("#FFFF00")).Bold(true)
		}

		sb.WriteString(fmt.Sprintf(" %s [%s] %s\n", cursor, checked, style.Render(item.Label)))
	}

	sb.WriteString("\n")
	return sb.String()
}
