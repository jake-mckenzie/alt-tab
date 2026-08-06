package tui

import (
	"errors"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/jake-mckenzie/alt-tab/internal/chords"
)

// fakeCatalog isolates model behavior from the native chord table.
type fakeCatalog struct{}

// Names supplies the small ordered chord set used by model tests.
func (fakeCatalog) Names() []string {
	return []string{"C", "G"}
}

// VariationCount gives each fake chord two navigable voicings.
func (fakeCatalog) VariationCount(string) int {
	return 2
}

// Load constructs a predictable voicing without invoking the C catalog.
func (fakeCatalog) Load(name string, variation int) (chords.Voicing, error) {
	if variation < 1 || variation > 2 {
		return chords.Voicing{}, errors.New("missing variation")
	}

	return chords.Voicing{
		Name:      name,
		Variation: variation,
		Strings: [chords.StringCount]chords.StringPlacement{
			{Fret: 0},
			{Fret: 1, Finger: 1},
			{Fret: 0},
			{Fret: 2, Finger: 2},
			{Fret: 3, Finger: 3},
			{Fret: -1},
		},
	}, nil
}

// TestViewContainsChordAndFingering checks the essential rendered content.
func TestViewContainsChordAndFingering(t *testing.T) {
	view := New(fakeCatalog{}).View()

	for _, expected := range []string{"ALT-TAB", "C", "variation 1 of 2", "--1--", "X"} {
		if !strings.Contains(view.Content, expected) {
			t.Fatalf("view does not contain %q", expected)
		}
	}
	if view.AltScreen {
		t.Fatal("view uses the alternate screen, which converts scrolling to keys")
	}
}

// TestChordListIsHorizontalAndAboveBothFretboardModes checks shared layout.
func TestChordListIsHorizontalAndAboveBothFretboardModes(t *testing.T) {
	for _, fullNeck := range []bool{false, true} {
		model := New(fakeCatalog{})
		model.fullNeck = fullNeck
		plain := ansiSequence.ReplaceAllString(model.View().Content, "")

		if !strings.Contains(plain, "‹ C ›   G") {
			t.Fatalf("fullNeck=%t does not show the horizontal chord list", fullNeck)
		}
		if strings.Index(plain, "CHORD SELECTOR") >
			strings.Index(plain, "variation 1 of 2") {
			t.Fatalf("fullNeck=%t places the chord list below the fretboard", fullNeck)
		}
	}
}

// TestSectionsAreClearlyLabeledAndOrdered checks the complete screen hierarchy.
func TestSectionsAreClearlyLabeledAndOrdered(t *testing.T) {
	plain := ansiSequence.ReplaceAllString(New(fakeCatalog{}).View().Content, "")
	titles := []string{"CHORD SELECTOR", "CONTROLS", "CHORD DIAGRAM", "NOTE WAVEFORM"}
	previous := -1
	for _, title := range titles {
		position := strings.Index(plain, title)
		if position <= previous {
			t.Fatalf("section %q is missing or out of order", title)
		}
		previous = position
	}
}

// TestCompactPanelsShrinkToContent checks compact and full-neck width policy.
func TestCompactPanelsShrinkToContent(t *testing.T) {
	model := New(fakeCatalog{})
	const maximumWidth = 76

	selector := model.renderPanel(model.renderChordSelector(), maximumWidth, true)
	diagram := model.renderPanel(model.renderChordDiagram(), maximumWidth, true)
	if lipgloss.Width(selector) >= maximumWidth || lipgloss.Width(diagram) >= maximumWidth {
		t.Fatal("compact selector or diagram retains unnecessary trailing width")
	}

	model.fullNeck = true
	diagram = model.renderPanel(model.renderChordDiagram(), maximumWidth, false)
	if lipgloss.Width(diagram) != maximumWidth {
		t.Fatalf("full-neck diagram width = %d, want %d", lipgloss.Width(diagram), maximumWidth)
	}
}

// lineWith returns the rendered line containing a section title.
func lineWith(output, title string) string {
	for _, line := range strings.Split(output, "\n") {
		if strings.Contains(line, title) {
			return line
		}
	}
	return ""
}

// TestCompactLayoutUsesColumnsAndFullNeckStacks checks responsive grouping.
func TestCompactLayoutUsesColumnsAndFullNeckStacks(t *testing.T) {
	model := New(fakeCatalog{})
	plain := ansiSequence.ReplaceAllString(model.View().Content, "")
	if !strings.Contains(lineWith(plain, "CHORD SELECTOR"), "CONTROLS") {
		t.Fatal("selector and controls are not side-by-side")
	}
	if !strings.Contains(lineWith(plain, "CHORD DIAGRAM"), "NOTE WAVEFORM") {
		t.Fatal("compact diagram and waveform are not side-by-side")
	}

	model.fullNeck = true
	plain = ansiSequence.ReplaceAllString(model.View().Content, "")
	if strings.Contains(lineWith(plain, "CHORD DIAGRAM"), "NOTE WAVEFORM") {
		t.Fatal("full-neck diagram and waveform remain side-by-side")
	}
}

// TestTopPanelsShareHeightAndControlsStayOnOneLine checks the compact header row.
func TestTopPanelsShareHeightAndControlsStayOnOneLine(t *testing.T) {
	model := New(chords.NewCatalog())
	selector, controls := model.renderMatchedPanels(
		model.renderChordSelector(),
		37,
		model.renderControls(),
		37,
	)
	if lipgloss.Height(selector) != lipgloss.Height(controls) {
		t.Fatalf(
			"panel heights differ: selector=%d controls=%d",
			lipgloss.Height(selector),
			lipgloss.Height(controls),
		)
	}
	if strings.Count(model.renderControls(), "\n") != 2 {
		t.Fatal("controls are not rendered on one line")
	}
}

// TestResizeRedrawsWithinTerminalWidth prevents wrapping outside Bubble Tea's frame.
func TestResizeRedrawsWithinTerminalWidth(t *testing.T) {
	for _, fullNeck := range []bool{false, true} {
		model := New(chords.NewCatalog())
		model.fullNeck = fullNeck
		for _, width := range []int{120, 80, 40, 20, 98, 30, 64, 39, 97, 79} {
			updated, command := model.Update(tea.WindowSizeMsg{Width: width, Height: 40})
			model = updated.(Model)
			if command == nil {
				t.Fatalf("width %d did not request a clean redraw", width)
			}

			plain := ansiSequence.ReplaceAllString(model.View().Content, "")
			for lineNumber, line := range strings.Split(plain, "\n") {
				if lineWidth := lipgloss.Width(line); lineWidth > width {
					t.Fatalf(
						"fullNeck=%t width %d line %d renders %d columns: %q\n%s",
						fullNeck,
						width,
						lineNumber+1,
						lineWidth,
						line,
						plain,
					)
				}
			}
		}
	}
}

// TestNavigationUpdatesSelection checks the remapped directional controls.
func TestNavigationUpdatesSelection(t *testing.T) {
	model := New(fakeCatalog{})
	updated, _ := model.Update(tea.KeyPressMsg{Code: tea.KeyRight})
	model = updated.(Model)

	if model.voicing.Name != "G" || model.voicing.Variation != 1 {
		t.Fatalf("right selected %s:%d, want G:1", model.voicing.Name, model.voicing.Variation)
	}

	updated, _ = model.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	model = updated.(Model)
	if model.voicing.Variation != 2 {
		t.Fatalf("down selected variation %d, want 2", model.voicing.Variation)
	}
}

// TestFullNeckAndHelpToggles checks both single-key mode switches.
func TestFullNeckAndHelpToggles(t *testing.T) {
	model := New(fakeCatalog{})
	updated, _ := model.Update(tea.KeyPressMsg{Code: 'f', Text: "f"})
	model = updated.(Model)
	if !model.fullNeck ||
		!strings.Contains(model.View().Content, " 1  2  3  4  5  6  7  8  9") ||
		!strings.Contains(model.View().Content, " 25 26 27") {
		t.Fatal("f did not enable the full-neck view")
	}

	updated, _ = model.Update(tea.KeyPressMsg{Code: '?', Text: "?"})
	model = updated.(Model)
	if !model.showHelp || !strings.Contains(model.View().Content, "KEYBOARD HELP") {
		t.Fatal("? did not open help")
	}
}

// TestThemeToggleCyclesPalettes checks runtime switching and wraparound.
func TestThemeToggleCyclesPalettes(t *testing.T) {
	model := New(fakeCatalog{})
	if model.theme != 0 || !strings.Contains(model.View().Content, "Synthwave") {
		t.Fatal("new model does not use the Synthwave theme")
	}

	for index := 1; index <= len(palettes); index++ {
		updated, _ := model.Update(tea.KeyPressMsg{Code: 't', Text: "t"})
		model = updated.(Model)
		want := index % len(palettes)
		if model.theme != want ||
			!strings.Contains(model.View().Content, paletteAt(want).name) {
			t.Fatalf("theme index = %d, want %d", model.theme, want)
		}
	}
}

// TestFullNeckRequiresNinetyEightColumns checks both sides of the boundary.
func TestFullNeckRequiresNinetyEightColumns(t *testing.T) {
	model := New(fakeCatalog{})
	model.width = fullNeckMinimumTerminalWidth - 1
	model.fullNeck = true

	if !strings.Contains(model.View().Content, "Widen the terminal") {
		t.Fatal("narrow full-neck view does not request a wider terminal")
	}

	model.width = fullNeckMinimumTerminalWidth
	if strings.Contains(model.View().Content, "Widen the terminal") {
		t.Fatal("full-neck view rejects a 98-column terminal")
	}
	if !strings.Contains(model.View().Content, " 25 26 27") {
		t.Fatal("full-neck view does not render at 98 columns")
	}
}

// TestQuitKeyReturnsCommand checks that quitting is delegated to Bubble Tea.
func TestQuitKeyReturnsCommand(t *testing.T) {
	model := New(fakeCatalog{})
	_, command := model.Update(tea.KeyPressMsg{Code: 'q', Text: "q"})

	if command == nil {
		t.Fatal("q did not return a quit command")
	}
}
