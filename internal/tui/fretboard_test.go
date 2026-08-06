package tui

import (
	"strings"
	"testing"

	"github.com/jake-mckenzie/alt-tab/internal/chords"
)

// TestCompactFretboardUsesExactFretsAndFingers checks its complete tab output.
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
	output := renderFretboard(voicing, false)

	for _, expected := range []string{
		"  1  ", "  4  ",
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
