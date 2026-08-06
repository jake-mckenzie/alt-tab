package tui

import (
	"strings"
	"testing"

	"github.com/jake-mckenzie/alt-tab/internal/chords"
)

// TestCompactFretboardUsesExactFretsAndFingers checks its complete tab output.
func TestCompactFretboardUsesExactFretsAndFingers(t *testing.T) {
	voicing := chords.Voicing{
		Name:   "C",
		Number: 1,
		Strings: [chords.StringCount]chords.StringPlacement{
			{Fret: 0},
			{Fret: 1, Finger: 1},
			{Fret: 0},
			{Fret: 2, Finger: 2},
			{Fret: 3, Finger: 3},
			{Fret: -1},
		},
	}
	output := renderFretboard(voicing, false)

	for _, expected := range []string{
		"  1  ", "  4  ",
		"    " + strings.Repeat("─", 4*compactCellWidth) + "\n",
		"e  O--------------------|",
		"B  |--1-----------------|",
		"D  |-------2------------|",
		"A  |------------3-------|",
		"E  X--------------------|",
		"4 little\n\nSymbols:",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("fretboard does not contain %q:\n%s", expected, output)
		}
	}
}

// TestCompactFretboardStartsAtHigherPosition checks movable chord windows.
func TestCompactFretboardStartsAtHigherPosition(t *testing.T) {
	voicing := chords.Voicing{
		Strings: [chords.StringCount]chords.StringPlacement{
			{Fret: 8, Finger: 1},
			{Fret: 8, Finger: 1},
			{Fret: 9, Finger: 2},
			{Fret: 10, Finger: 4},
			{Fret: 10, Finger: 3},
			{Fret: 8, Finger: 1},
		},
	}
	output := renderFretboard(voicing, false)

	if !strings.Contains(output, "  8  ") || !strings.Contains(output, "  11 ") {
		t.Fatalf("higher-position labels are missing:\n%s", output)
	}
	if strings.Contains(output, "  1  ") {
		t.Fatalf("higher-position diagram unexpectedly starts at fret 1:\n%s", output)
	}
}

// TestTabDiagramUsesFretNumbers checks numbering within compact neck geometry.
func TestTabDiagramUsesFretNumbers(t *testing.T) {
	voicing := chords.Voicing{
		Strings: [chords.StringCount]chords.StringPlacement{
			{Fret: 0},
			{Fret: 1, Finger: 1},
			{Fret: 2, Finger: 3},
			{Fret: 2, Finger: 2},
			{Fret: 0},
			{Fret: -1},
		},
	}
	output := renderTabDiagram(voicing)
	for _, expected := range []string{
		"e  O--------------------|",
		"B  |--1-----------------|",
		"G  |-------2------------|",
		"D  |-------2------------|",
		"A  O--------------------|",
		"E  X--------------------|",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("tab diagram does not contain %q:\n%s", expected, output)
		}
	}
	if strings.Contains(output, "Fingers:") || strings.Contains(output, "Symbols:") {
		t.Fatal("tab diagram includes fretboard-only legends")
	}
	if strings.Contains(output, "  1  ") || strings.Contains(output, "────") {
		t.Fatal("tab diagram includes the separate fret scale")
	}
}

// TestFullFretboardLabelsAllTwentySevenFrets checks the full neck dimensions.
func TestFullFretboardLabelsAllTwentySevenFrets(t *testing.T) {
	output := renderFretboard(chords.Voicing{}, true)

	if !strings.Contains(output, " 1  2  3  4  5  6  7  8  9  10 11 12") ||
		!strings.Contains(output, " 25 26 27") {
		t.Fatalf("full-neck labels are incomplete:\n%s", output)
	}
	if !strings.Contains(
		output,
		"e  O"+strings.Repeat("-", fullNeckLastFret*fullNeckCellWidth)+"|",
	) {
		t.Fatalf("full-neck string does not span 27 frets:\n%s", output)
	}
	if !strings.Contains(
		output,
		"    "+strings.Repeat("─", fullNeckLastFret*fullNeckCellWidth)+"\n",
	) {
		t.Fatalf("full-neck fret separator is incomplete:\n%s", output)
	}
}

// TestFullFretboardAlignsFingerAtFinalFret guards the rightmost cell alignment.
func TestFullFretboardAlignsFingerAtFinalFret(t *testing.T) {
	voicing := chords.Voicing{
		Strings: [chords.StringCount]chords.StringPlacement{
			{Fret: 27, Finger: 4},
		},
	}
	output := renderFretboard(voicing, true)

	if !strings.Contains(
		output,
		"e  |"+strings.Repeat("-", fullNeckLastFret*fullNeckCellWidth-2)+"4-|",
	) {
		t.Fatalf("fret 27 marker is misaligned:\n%s", output)
	}
}

// TestFullFretboardCentersLegend checks its compact side-by-side footer.
func TestFullFretboardCentersLegend(t *testing.T) {
	output := renderFretboard(chords.Voicing{}, true)
	width := 4 + fullNeckLastFret*fullNeckCellWidth + 1
	legend := "Fingers: 1 index  2 middle    Symbols: O open  X muted"
	line := lineWith(output, "Fingers:")
	if !strings.Contains(line, legend) {
		t.Fatalf("full-neck legends are not side-by-side:\n%s", output)
	}
	left := strings.Index(line, legend)
	right := width - left - len(legend)
	if absoluteInt(left-right) > 1 {
		t.Fatalf("full-neck legend is not centered: left=%d right=%d", left, right)
	}
	if strings.Count(output, "\n") != 10 {
		t.Fatalf("full-neck output uses unexpected vertical space:\n%s", output)
	}
}
