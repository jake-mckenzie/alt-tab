package tui

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestViewContainsApplicationName(t *testing.T) {
	view := New().View()

	if !strings.Contains(view.Content, "ALT-TAB") {
		t.Fatal("view does not contain application name")
	}
	if !view.AltScreen {
		t.Fatal("view does not request the alternate screen")
	}
}

func TestQuitKeyReturnsCommand(t *testing.T) {
	model := New()
	_, command := model.Update(tea.KeyPressMsg{Code: 'q', Text: "q"})

	if command == nil {
		t.Fatal("q did not return a quit command")
	}
}
