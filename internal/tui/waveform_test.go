package tui

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

// TestRenderWaveformUsesRequestedDimensions checks stable ASCII geometry.
func TestRenderWaveformUsesRequestedDimensions(t *testing.T) {
	output := renderWaveform(16, 0)
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
