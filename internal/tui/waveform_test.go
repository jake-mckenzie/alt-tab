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

	if len(lines) != waveformHeight+10 {
		t.Fatalf("waveform height = %d, want %d", len(lines), waveformHeight+10)
	}
	for _, line := range lines[:waveformHeight+1] {
		if utf8.RuneCountInString(line) != 40 {
			t.Fatalf(
				"waveform width = %d, want 40",
				utf8.RuneCountInString(line),
			)
		}
	}
	for _, expected := range []string{
		waveformCaption, pitchCaption, "+1.0 |", " 0.0 +", "-1.0 |", "25 ms", "Notes:",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("waveform is missing %q:\n%s", expected, output)
		}
	}
	if !strings.ContainsFunc(output, isBrailleRune) {
		t.Fatal("waveform does not contain Braille subcells")
	}
	for row, line := range lines[:waveformHeight] {
		wantBoundary := "|"
		if row == waveformHeight/2 {
			wantBoundary = "+"
		}
		if !strings.HasSuffix(line, wantBoundary) {
			t.Fatalf("waveform row %d has no matching right boundary: %q", row, line)
		}
	}
}

// TestPitchScaleMarksEverySoundingNote checks semitone-positioned annotations.
func TestPitchScaleMarksEverySoundingNote(t *testing.T) {
	voicing, _ := fakeCatalog{}.Load("C", 1)
	scale := renderPitchScale(40, voicingNotes(voicing))
	for _, expected := range []string{"C3", "E3", "G3", "C4", "E4", "│", pitchCaption} {
		if !strings.Contains(scale, expected) {
			t.Fatalf("pitch scale is missing %q:\n%s", expected, scale)
		}
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
