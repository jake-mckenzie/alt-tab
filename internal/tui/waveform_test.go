package tui

import (
	"math"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/jake-mckenzie/alt-tab/internal/chords"
)

// TestRenderWaveformUsesRequestedDimensions checks stable ASCII geometry.
func TestRenderWaveformUsesRequestedDimensions(t *testing.T) {
	voicing, _ := fakeCatalog{}.Load("C", 1)
	output := renderWaveform(16, 0, voicing)
	lines := strings.Split(output, "\n")

	if len(lines) != waveformHeight {
		t.Fatalf("waveform height = %d, want %d", len(lines), waveformHeight)
	}
	for _, line := range lines {
		if len(line) != 16 {
			t.Fatalf("waveform width = %d, want 16", len(line))
		}
	}
	if strings.ContainsAny(output, "0123456789") || !strings.Contains(output, "*") {
		t.Fatalf("waveform is not portable ASCII:\n%s", output)
	}
}

// TestWaveformUsesVoicingFrequencies checks note conversion and chord identity.
func TestWaveformUsesVoicingFrequencies(t *testing.T) {
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
	frequencies := voicingFrequencies(voicing)
	if len(frequencies) != 1 || math.Abs(frequencies[0]-329.63) > 0.001 {
		t.Fatalf("high-e frequencies = %v, want [329.63]", frequencies)
	}

	cMajor, _ := fakeCatalog{}.Load("C", 1)
	gMajor := cMajor
	gMajor.Strings[0].Fret = 3
	if renderWaveform(40, 0, cMajor) == renderWaveform(40, 0, gMajor) {
		t.Fatal("different chord notes produced identical waveforms")
	}
}

// TestWaveformToggleControlsAnimation checks timer invalidation and restart.
func TestWaveformToggleControlsAnimation(t *testing.T) {
	model := New(fakeCatalog{})
	updated, command := model.Update(tea.KeyPressMsg{Code: 'w', Text: "w"})
	model = updated.(Model)
	if model.waveform || command != nil || strings.Contains(model.View().Content, "*") {
		t.Fatal("w did not disable the waveform and its timer")
	}

	staleGeneration := model.waveTimer - 1
	updated, command = model.Update(waveformTickMsg{generation: staleGeneration})
	model = updated.(Model)
	if command != nil || model.waveFrame != 0 {
		t.Fatal("stale waveform timer advanced the animation")
	}

	updated, command = model.Update(tea.KeyPressMsg{Code: 'w', Text: "w"})
	model = updated.(Model)
	if !model.waveform || command == nil {
		t.Fatal("w did not restart the waveform timer")
	}

	updated, command = model.Update(waveformTickMsg{generation: model.waveTimer})
	model = updated.(Model)
	if model.waveFrame != 1 || command == nil {
		t.Fatal("active waveform timer did not advance and reschedule")
	}
}
