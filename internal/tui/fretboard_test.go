package tui

import (
	"strings"
	"testing"

	"alt-tab/internal/chords"
)

func TestCompactFretboardUsesExactFretsAndFingers(t *testing.T) {
	voicing := chords.Voicing{
		Name:      "C",
		Variation: 1,
		Strings: [chords.StringCount]chords.StringPlacement{
			{Fret: 0},
			{Fret: 1, Finger: 1},
			{Fret: 0},
			{Fret: 2, Finger: 2},
			{Fret: 3, Finger: 3},
			{Fret: -1},
		},
	}
	output := renderFretboard(voicing, false, 80)

	for _, expected := range []string{
		"  1  ", "  4  ",
		"e  O |--------------------|",
		"B    |--1-----------------|",
		"D    |-------2------------|",
		"A    |------------3-------|",
		"E  X |--------------------|",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("fretboard does not contain %q:\n%s", expected, output)
		}
	}
}

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
	output := renderFretboard(voicing, false, 80)

	if !strings.Contains(output, "  8  ") || !strings.Contains(output, "  11 ") {
		t.Fatalf("higher-position labels are missing:\n%s", output)
	}
	if strings.Contains(output, "  1  ") {
		t.Fatalf("higher-position diagram unexpectedly starts at fret 1:\n%s", output)
	}
}

func TestFullFretboardLabelsAllTwentySevenFrets(t *testing.T) {
	output := renderFretboard(chords.Voicing{}, true, 80)

	if !strings.Contains(output, "Fret   1 2 3 4 5 6 7 8 9 0 1") ||
		!strings.Contains(output, " 5 6 7") {
		t.Fatalf("full-neck labels are incomplete:\n%s", output)
	}
	if !strings.Contains(
		output,
		"e  O |"+strings.Repeat("-", fullNeckLastFret*2)+"|",
	) {
		t.Fatalf("full-neck string does not span 27 frets:\n%s", output)
	}
}

func TestFullFretboardFallsBackAtNarrowWidths(t *testing.T) {
	output := renderFretboard(chords.Voicing{}, true, 50)

	if !strings.Contains(output, "111111111122222222") ||
		!strings.Contains(output, "Fret  12345678901234567") {
		t.Fatalf("narrow full-neck labels are incomplete:\n%s", output)
	}
	if !strings.Contains(output, "e  O |"+strings.Repeat("-", 27)+"|") {
		t.Fatalf("narrow full-neck string does not span 27 frets:\n%s", output)
	}
}
