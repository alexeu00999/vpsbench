package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestModel_Update(t *testing.T) {
	m := Model{
		Items: []Item{
			{ID: "CPU", Label: "CPU", Selected: true},
			{ID: "RAM", Label: "RAM", Selected: false},
		},
		Cursor: 0,
		Timer:  3,
	}

	// 1. Тест перемещения курсора вниз
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")}
	newModel, _ := m.Update(msg)
	m = newModel.(Model)
	if m.Cursor != 1 {
		t.Errorf("expected cursor at 1, got %d", m.Cursor)
	}
	if m.Timer != -1 {
		t.Errorf("timer should be disabled after keypress, got %d", m.Timer)
	}

	// 2. Тест переключения выбора (Space)
	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(" ")}
	newModel, _ = m.Update(msg)
	m = newModel.(Model)
	if m.Items[1].Selected != true {
		t.Errorf("expected RAM to be selected after space")
	}

	// 3. Тест выхода (Enter)
	msg = tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ = m.Update(msg)
	m = newModel.(Model)
	if m.Finished != true {
		t.Errorf("expected Finished to be true after enter")
	}
}

func TestModel_TimerTick(t *testing.T) {
	m := Model{
		Items: []Item{{ID: "test"}},
		Timer: 1,
	}

	// Тик таймера 1 -> 0
	msg := tickMsg{}
	newModel, cmd := m.Update(msg)
	m = newModel.(Model)
	if m.Timer != 0 {
		t.Errorf("expected timer 0, got %d", m.Timer)
	}
	if cmd == nil {
		t.Errorf("expected next tick cmd")
	}

	// Тик таймера 0 -> Finish
	newModel, cmd = m.Update(msg)
	m = newModel.(Model)
	if m.Finished != true {
		t.Errorf("expected Finished=true after timer reached 0")
	}
}
