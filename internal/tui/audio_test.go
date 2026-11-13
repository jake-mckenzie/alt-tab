package tui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/jake-mckenzie/alt-tab/internal/chords"
	"testing"
)

type testPlayer struct {
	voicing chords.Voicing
	active  bool
}

func (p *testPlayer) Set(v chords.Voicing, active bool) { p.voicing, p.active = v, active }

func TestAudioControls(t *testing.T) {
	p := &testPlayer{}
	m := New(chords.NewCatalog()).WithAudio(p)
	press := func(key string) {
		next, _ := m.Update(tea.KeyPressMsg{Code: rune(key[0]), Text: key})
		m = next.(Model)
	}
	press("p")
	if !p.active || p.voicing.Name != "A" {
		t.Fatal("p did not play selection")
	}
	press("v")
	if p.voicing.Number != 2 {
		t.Fatal("playing voicing did not update")
	}
	press("c")
	if !p.active || p.voicing.Name != "C" {
		t.Fatal("c did not select and play C")
	}
	press("c")
	if p.active {
		t.Fatal("legacy terminal chord key did not toggle off")
	}
	next, _ := m.Update(tea.KeyboardEnhancementsMsg{Flags: ansi.KittyReportEventTypes})
	m = next.(Model)
	press("g")
	press("g")
	if !p.active {
		t.Fatal("repeated held key stopped playback")
	}
	next, _ = m.Update(tea.KeyReleaseMsg{Code: 'c', Text: "c"})
	m = next.(Model)
	if !p.active {
		t.Fatal("unrelated key release stopped playback")
	}
	next, _ = m.Update(tea.KeyReleaseMsg{Code: 'g', Text: "g"})
	m = next.(Model)
	if p.active {
		t.Fatal("release did not stop held chord")
	}
	press("p")
	next, _ = m.Update(tea.BlurMsg{})
	m = next.(Model)
	if p.active {
		t.Fatal("focus loss did not stop playback")
	}
	press("p")
	press("?")
	if p.active {
		t.Fatal("help did not stop playback")
	}
	press("esc")
	press("p")
	press("q")
	if p.active {
		t.Fatal("quit did not stop playback")
	}
}

// Selection changes must keep the diagram and audio on the same voicing.
func TestSelectionKeepsDiagramAndAudioTogether(t *testing.T) {
	p := &testPlayer{}
	m := New(chords.NewCatalog()).WithAudio(p)
	for _, step := range []struct {
		key, chord string
		number     int
	}{
		{"p", "A", 1},
		{"3", "A", 3},
		{"9", "A", 3}, // unavailable positions leave the selection intact
		{"c", "C", 3}, // letter shortcuts retain the selected voicing
		{"l", "D", 1}, // dial movement resets to the first voicing
		{"2", "D", 2},
		{"j", "Dm", 1}, // changing chord kind also resets the voicing
		{"v", "Dm", 2},
		{"v", "Dm", 3},
		{"v", "Dm", 1},
	} {
		next, _ := m.Update(tea.KeyPressMsg{Code: rune(step.key[0]), Text: step.key})
		m = next.(Model)
		if m.err != nil || m.voicing.Name != step.chord || m.voicing.Number != step.number {
			t.Fatalf("%s selected %s:%d (%v), want %s:%d", step.key, m.voicing.Name, m.voicing.Number, m.err, step.chord, step.number)
		}
		if !p.active || p.voicing != m.voicing {
			t.Fatalf("%s left audio out of sync", step.key)
		}
	}
}
