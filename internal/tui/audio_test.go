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
