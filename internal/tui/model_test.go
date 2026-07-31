package tui

import (
	"errors"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/jake-mckenzie/alt-tab/internal/chords"
)

type fakeCatalog struct{}

func (fakeCatalog) Names() []string {
	return []string{"C", "G"}
}

func (fakeCatalog) VariationCount(string) int {
	return 2
}

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

func TestViewContainsChordAndFingering(t *testing.T) {
	view := New(fakeCatalog{}).View()

	for _, expected := range []string{"ALT-TAB", "C", "variation 1 of 2", "--1--", "X"} {
		if !strings.Contains(view.Content, expected) {
			t.Fatalf("view does not contain %q", expected)
		}
	}
	if !view.AltScreen {
		t.Fatal("view does not request the alternate screen")
	}
}

func TestChordListIsVerticalAndAboveBothFretboardModes(t *testing.T) {
	for _, fullNeck := range []bool{false, true} {
		model := New(fakeCatalog{})
		model.fullNeck = fullNeck
		plain := ansiSequence.ReplaceAllString(model.View().Content, "")

		if !strings.Contains(plain, "│ › C") || !strings.Contains(plain, "│   G") {
			t.Fatalf("fullNeck=%t does not show the vertical chord list", fullNeck)
		}
		if strings.Index(plain, "CHORDS") > strings.Index(plain, "variation 1 of 2") {
			t.Fatalf("fullNeck=%t places the chord list below the fretboard", fullNeck)
		}
	}
}

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

func TestQuitKeyReturnsCommand(t *testing.T) {
	model := New(fakeCatalog{})
	_, command := model.Update(tea.KeyPressMsg{Code: 'q', Text: "q"})

	if command == nil {
		t.Fatal("q did not return a quit command")
	}
}
