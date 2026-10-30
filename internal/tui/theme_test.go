package tui

import (
	"testing"

	uv "github.com/charmbracelet/ultraviolet"

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

func TestNeonDarkPanelBackgroundIsUniform(t *testing.T) {
	model := New(fakeCatalog{})
	model.theme = 1
	model.styles = newStyles(true, model.theme)
	content := model.styles.heading.Render("TITLE") + " trailing\n\n" +
		model.styles.selected.Render("A") + " after selection"
	for name, rendered := range map[string]string{
		"centered": model.renderCenteredPanelAtHeight(content, 40, 8),
		"plain":    model.renderPanel(content, 40, false),
	} {
		width, height := lipgloss.Width(rendered), lipgloss.Height(rendered)
		screen := uv.NewScreenBuffer(width, height)
		uv.NewStyledString(rendered).Draw(screen, uv.Rect(0, 0, width, height))
		selected := false
		for y := 1; y < height-1; y++ {
			for x := 1; x < width-1; x++ {
				cell := screen.CellAt(x, y)
				want := lipgloss.Color("#0E1012")
				if cell.Content == "A" {
					want = lipgloss.Color("#151A18")
					selected = true
				}
				if cell.Style.Bg == nil || cell.Style.Bg != want {
					t.Fatalf("%s panel cell (%d,%d) %q: background %v, want %v", name, x, y, cell.Content, cell.Style.Bg, want)
				}
			}
		}
		if !selected {
			t.Fatal("selected chord missing from panel")
		}
	}
}
