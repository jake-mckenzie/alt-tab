package tui

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func TestNeonDarkBackgroundAndThemeSwitch(t *testing.T) {
	for _, dark := range []bool{false, true} {
		model := New(fakeCatalog{})
		model.dark = dark
		updated, _ := model.Update(tea.KeyPressMsg{Code: 't', Text: "t"})
		model = updated.(Model)
		if paletteAt(model.theme).name != "Neon Dark" {
			t.Fatal("first theme after default is not Neon Dark")
		}
		if model.View().BackgroundColor != lipgloss.Color("#070809") ||
			model.styles.panel.GetBackground() != lipgloss.Color("#0E1012") ||
			model.styles.accent.GetForeground() != lipgloss.Color("#39FF88") ||
			model.styles.signal.GetForeground() != lipgloss.Color("#83FFE8") {
			t.Fatalf("Neon Dark colors changed with terminal dark=%v", dark)
		}
		updated, _ = model.Update(tea.KeyPressMsg{Code: 't', Text: "t"})
		model = updated.(Model)
		if model.View().BackgroundColor != nil {
			t.Fatal("switching away did not restore the terminal background")
		}
		if model.styles.signal.GetForeground() != model.styles.accent.GetForeground() {
			t.Fatal("switching away did not restore the original graph color")
		}
	}
}
