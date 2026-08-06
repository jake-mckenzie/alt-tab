package chords

import (
	"errors"
	"reflect"
	"testing"
)

// TestNativeCatalogNames verifies the ordered, de-duplicated name list.
func TestNativeCatalogNames(t *testing.T) {
	catalog := NewCatalog()
	expected := []string{
		"A", "Am", "B", "Bb", "C", "Cm", "D",
		"Dm", "E", "Em", "F", "F#", "G", "Gm",
	}

	if names := catalog.Names(); !reflect.DeepEqual(names, expected) {
		t.Fatalf("Names() = %v, want %v", names, expected)
	}
}

// TestNativeCatalogLoadsOwnedVoicing verifies lookup, casing, and C-to-Go copying.
func TestNativeCatalogLoadsOwnedVoicing(t *testing.T) {
	catalog := NewCatalog()
	voicing, err := catalog.Load("c", 2)

	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if voicing.Name != "C" || voicing.Number != 2 {
		t.Fatalf("Load() identity = %q:%d, want C:2", voicing.Name, voicing.Number)
	}
	if voicing.Strings[0] != (StringPlacement{Fret: 3, Finger: 1}) {
		t.Fatalf("high e placement = %+v, want fret 3 finger 1", voicing.Strings[0])
	}
	if voicing.Strings[5].Fret != -1 {
		t.Fatalf("low E fret = %d, want muted", voicing.Strings[5].Fret)
	}
}

// TestNativeCatalogRejectsMissingVariation verifies invalid voicings are rejected.
func TestNativeCatalogRejectsMissingVariation(t *testing.T) {
	catalog := NewCatalog()

	if count := catalog.VoicingCount("C"); count != 2 {
		t.Fatalf("VoicingCount() = %d, want 2", count)
	}
	if _, err := catalog.Load("C", 3); !errors.Is(err, ErrChordNotFound) {
		t.Fatalf("Load() error = %v, want ErrChordNotFound", err)
	}
}
