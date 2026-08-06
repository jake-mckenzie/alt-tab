package tui

import (
	"math"
	"strings"
	"testing"
	"unicode/utf8"

	tea "charm.land/bubbletea/v2"

	"github.com/jake-mckenzie/alt-tab/internal/chords"
)

// isBrailleRune identifies one Unicode Braille pattern character.
func isBrailleRune(character rune) bool {
	return character >= '\u2800' && character <= '\u28ff'
}

// TestRenderWaveformUsesRequestedDimensions checks detailed ASCII geometry.
func TestRenderWaveformUsesRequestedDimensions(t *testing.T) {
	voicing, _ := fakeCatalog{}.Load("C", 1)
	output := renderWaveform(40, voicing)
	lines := strings.Split(output, "\n")

	if len(lines) != waveformHeight+3 {
		t.Fatalf("waveform height = %d, want %d", len(lines), waveformHeight+3)
	}
	for _, line := range lines[:waveformHeight+3] {
		if utf8.RuneCountInString(line) != 40 {
			t.Fatalf(
				"waveform width = %d, want 40",
				utf8.RuneCountInString(line),
			)
		}
	}
	for _, expected := range []string{
		waveformTitle, "+1.0 |", " 0.0 +", "-1.0 |", "25 ms",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("waveform is missing %q:\n%s", expected, output)
		}
	}
	if !strings.ContainsFunc(output, isBrailleRune) {
		t.Fatal("waveform does not contain Braille subcells")
	}
	waveformTitleLine := lineIndex(lines, waveformTitle)
	if waveformTitleLine < 0 || waveformTitleLine+1 >= len(lines) ||
		strings.TrimSpace(lines[waveformTitleLine+1]) !=
			strings.Repeat("─", len(waveformTitle)) {
		t.Fatalf("%q does not have an underline below it", waveformTitle)
	}
	for row, line := range lines[2 : waveformHeight+2] {
		wantBoundary := "|"
		if row == waveformHeight/2 {
			wantBoundary = "+"
		}
		if !strings.HasSuffix(line, wantBoundary) {
			t.Fatalf("waveform row %d has no matching right boundary: %q", row, line)
		}
	}
}

// lineIndex finds an exact visible line while ignoring centering spaces.
func lineIndex(lines []string, value string) int {
	for index, line := range lines {
		if strings.TrimSpace(line) == value {
			return index
		}
	}
	return -1
}

// TestSpectrumMarksEverySoundingNote checks its peaks and annotations.
func TestSpectrumMarksEverySoundingNote(t *testing.T) {
	voicing, _ := fakeCatalog{}.Load("C", 1)
	spectrum := renderSpectrum(40, voicing)
	for _, expected := range []string{
		"1.0 |", "0.5 |", "0.0 +", "C3", "E3", "G3", "C4", "E4", "█", "Notes:",
	} {
		if !strings.Contains(spectrum, expected) {
			t.Fatalf("frequency spectrum is missing %q:\n%s", expected, spectrum)
		}
	}
	lines := strings.Split(spectrum, "\n")
	legend := lines[len(lines)-1]
	legendText := strings.TrimSpace(legend)
	left := strings.Index(legend, legendText)
	right := utf8.RuneCountInString(legend) - left - len(legendText)
	if absoluteInt(left-right) > 1 {
		t.Fatalf("note legend is not centered: left=%d right=%d", left, right)
	}
}

// TestWaveformUsesVoicingNotes checks exact note conversion and chord identity.
func TestWaveformUsesVoicingNotes(t *testing.T) {
	voicing := chords.Voicing{
		Strings: [chords.StringCount]chords.StringPlacement{
			{Fret: 0},
			{Fret: -1},
			{Fret: -1},
			{Fret: -1},
			{Fret: -1},
			{Fret: -1},
		},
	}
	notes := voicingNotes(voicing)
	if len(notes) != 1 || notes[0].name != "E4" ||
		math.Abs(notes[0].frequency-329.63) > 0.001 {
		t.Fatalf("high-e notes = %v, want E4 at 329.63 Hz", notes)
	}

	cMajor, _ := fakeCatalog{}.Load("C", 1)
	gMajor := cMajor
	gMajor.Strings[0].Fret = 3
	if renderWaveform(60, cMajor) == renderWaveform(60, gMajor) {
		t.Fatal("different chord notes produced identical waveforms")
	}
}

// TestWaveformIsStationary checks deterministic output and its runtime toggle.
func TestWaveformIsStationary(t *testing.T) {
	model := New(fakeCatalog{})
	before := model.View().Content
	if before != model.View().Content {
		t.Fatal("stationary waveform changed without a chord change")
	}

	updated, command := model.Update(tea.KeyPressMsg{Code: 'w', Text: "w"})
	model = updated.(Model)
	if model.waveform || command != nil || strings.Contains(model.View().Content, "+1.0 |") {
		t.Fatal("w did not hide the stationary waveform")
	}

	updated, command = model.Update(tea.KeyPressMsg{Code: 'w', Text: "w"})
	model = updated.(Model)
	if !model.waveform || command != nil || !strings.Contains(model.View().Content, "+1.0 |") {
		t.Fatal("w did not restore the stationary waveform")
	}
}
