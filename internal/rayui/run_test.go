//go:build raylib

package rayui

import (
	"testing"

	"github.com/jake-mckenzie/alt-tab/internal/chords"
	"github.com/jake-mckenzie/alt-tab/internal/signal"
)

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
