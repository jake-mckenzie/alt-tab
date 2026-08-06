//go:build raylib

package rayui

import (
	"reflect"
	"testing"

	"github.com/jake-mckenzie/alt-tab/internal/chords"
	"github.com/jake-mckenzie/alt-tab/internal/signal"
)

// TestPalettesMatchApprovedSet keeps the runtime cycle aligned with documentation.
func TestPalettesMatchApprovedSet(t *testing.T) {
	want := []string{
		"Super Famicom", "Atomic Grape", "Paper Terminal", "Glacier Circuit",
		"Haunted Cartridge", "Oxide Industrial", "Cassette Future",
		"Royal Terminal", "Sakura Console", "CRT Amber",
	}
	got := make([]string, len(palettes))
	for index, theme := range palettes {
		got[index] = theme.name
		if theme.background.A != 255 || theme.panel.A != 255 || theme.signal.A != 255 {
			t.Fatalf("%s contains a partially transparent core color", theme.name)
		}
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("palette names = %v, want %v", got, want)
	}
}

// TestComputeLayout verifies that compact and full-neck panels remain visible.
func TestComputeLayout(t *testing.T) {
	for _, fullNeck := range []bool{false, true} {
		layout := computeLayout(minimumWindowWidth, minimumWindowHeight, fullNeck)
		for name, bounds := range map[string]float32{
			"diagram":  layout.diagram.Height,
			"waveform": layout.waveform.Height,
			"spectrum": layout.spectrum.Height,
		} {
			if bounds <= 0 {
				t.Fatalf("%s height is not positive in full-neck=%v", name, fullNeck)
			}
		}
	}

	full := computeLayout(minimumWindowWidth, minimumWindowHeight, true)
	if full.diagram.Height-96 < 75 ||
		full.waveform.Height-graphTopPadding-waveBottomPadding < 20 ||
		full.spectrum.Height-graphTopPadding-spectrumBottomPad < 20 {
		t.Fatalf("minimum full-neck layout leaves unusable content: %+v", full)
	}
}

// TestFretRange keeps open and upper-position diagrams correctly numbered.
func TestFretRange(t *testing.T) {
	open := chords.Voicing{Strings: [chords.StringCount]chords.StringPlacement{{Fret: 3}}}
	if first, last := fretRange(open, false); first != 1 || last != 4 {
		t.Fatalf("open range = %d-%d, want 1-4", first, last)
	}
	high := chords.Voicing{Strings: [chords.StringCount]chords.StringPlacement{{Fret: 12}, {Fret: 14}}}
	if first, last := fretRange(high, false); first != 12 || last != 15 {
		t.Fatalf("high range = %d-%d, want 12-15", first, last)
	}
	if first, last := fretRange(high, true); first != 1 || last != fullNeckLastFret {
		t.Fatalf("full range = %d-%d, want 1-%d", first, last, fullNeckLastFret)
	}
}

// TestAggregateNotes verifies that unison strings become one stronger peak.
func TestAggregateNotes(t *testing.T) {
	notes := []signal.Note{{Name: "A3", MIDI: 57}, {Name: "A3", MIDI: 57}, {Name: "E4", MIDI: 64}}
	peaks := aggregateNotes(notes)
	if len(peaks) != 2 || peaks[0].count != 2 || peaks[1].count != 1 {
		t.Fatalf("aggregateNotes() = %#v", peaks)
	}
}

// TestRuntimeResourceSettings guards the low-resource rendering and audio targets.
func TestRuntimeResourceSettings(t *testing.T) {
	if audioSampleRate != 44100 {
		t.Fatalf("audio sample rate = %d, want 44100", audioSampleRate)
	}
	if idleFrameRate != 30 || activeFrameRate != 60 {
		t.Fatalf("frame rates = %d/%d, want 30/60", idleFrameRate, activeFrameRate)
	}
}
