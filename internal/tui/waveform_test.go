package tui

import (
	"fmt"
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

	if len(lines) != waveformHeight+1 {
		t.Fatalf("waveform height = %d, want %d", len(lines), waveformHeight+1)
	}
	for _, line := range lines {
		if utf8.RuneCountInString(line) != 40 {
			t.Fatalf(
				"waveform width = %d, want 40",
				utf8.RuneCountInString(line),
			)
		}
	}
	for _, expected := range []string{
		"+1.0 |", " 0.0 +", "-1.0 |", "25 ms",
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
	for row, line := range lines[:spectrumHeight] {
		if !strings.HasSuffix(line, "|") {
			t.Fatalf("spectrum row %d has no right boundary: %q", row, line)
		}
	}
	if !strings.HasSuffix(lines[spectrumHeight], "+") {
		t.Fatalf("spectrum axis has no right endpoint: %q", lines[spectrumHeight])
	}
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
		math.Abs(notes[0].frequency-midiFrequency(64)) > 1e-9 {
		t.Fatalf("high-e notes = %v, want E4 at %.6f Hz", notes, midiFrequency(64))
	}

	cMajor, _ := fakeCatalog{}.Load("C", 1)
	gMajor := cMajor
	gMajor.Strings[0].Fret = 3
	if renderWaveform(60, cMajor) == renderWaveform(60, gMajor) {
		t.Fatal("different chord notes produced identical waveforms")
	}
}

// TestCatalogVoicingsProduceExactPlotNotes validates every audio-data input.
func TestCatalogVoicingsProduceExactPlotNotes(t *testing.T) {
	catalog := chords.NewCatalog()
	for _, name := range catalog.Names() {
		for number := 1; number <= catalog.VoicingCount(name); number++ {
			voicing, err := catalog.Load(name, number)
			if err != nil {
				t.Fatalf("Load(%s, %d): %v", name, number, err)
			}

			notes := voicingNotes(voicing)
			noteIndex := 0
			for stringIndex, placement := range voicing.Strings {
				if placement.Fret < 0 {
					continue
				}
				midi := standardTuningMIDI[stringIndex] + placement.Fret
				expectedName := fmt.Sprintf(
					"%s%d",
					noteNames[wrap(midi, len(noteNames))],
					midi/12-1,
				)
				if notes[noteIndex].name != expectedName || notes[noteIndex].midi != midi {
					t.Fatalf("%s:%d string %d note = %+v, want %s MIDI %d",
						name, number, stringIndex, notes[noteIndex], expectedName, midi)
				}
				if math.Abs(notes[noteIndex].frequency-midiFrequency(midi)) > 1e-9 {
					t.Fatalf("%s:%d %s frequency = %.9f, want %.9f",
						name, number, expectedName, notes[noteIndex].frequency, midiFrequency(midi))
				}
				noteIndex++
			}
			if noteIndex != len(notes) {
				t.Fatalf("%s:%d produced %d notes, checked %d", name, number, len(notes), noteIndex)
			}

			spectrum := renderSpectrum(80, voicing)
			if !strings.Contains(spectrum, renderNoteLegend(notes)) {
				t.Fatalf("%s:%d spectrum omits its exact note legend", name, number)
			}
			if !strings.ContainsFunc(renderWaveform(80, voicing), isBrailleRune) {
				t.Fatalf("%s:%d waveform has no rendered signal", name, number)
			}
		}
	}
	if math.Abs(midiFrequency(69)-440) > 1e-12 {
		t.Fatalf("A4 frequency = %.12f, want 440", midiFrequency(69))
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
