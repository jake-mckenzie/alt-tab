package signal

import (
	"math"
	"testing"

	"github.com/jake-mckenzie/alt-tab/internal/chords"
)

// TestNotesUsesStandardTuning checks pitch names and equal temperament.
func TestNotesUsesStandardTuning(t *testing.T) {
	voicing := chords.Voicing{Strings: [chords.StringCount]chords.StringPlacement{
		{Fret: 5}, {Fret: -1}, {Fret: -1}, {Fret: -1}, {Fret: -1}, {Fret: -1},
	}}
	notes := Notes(voicing)
	if len(notes) != 1 || notes[0].Name != "A4" || math.Abs(notes[0].Frequency-440) > 1e-12 {
		t.Fatalf("notes = %+v, want A4 at 440 Hz", notes)
	}
}

// TestCompositeSampleStartsAtZero checks deterministic oscillator phase.
func TestCompositeSampleStartsAtZero(t *testing.T) {
	if sample := CompositeSample(0, []Note{{Frequency: 440}}); sample != 0 {
		t.Fatalf("initial sample = %f, want 0", sample)
	}
}
