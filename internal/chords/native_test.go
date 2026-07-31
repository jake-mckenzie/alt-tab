package chords

import (
	"errors"
	"reflect"
	"testing"
)

func TestNativeCatalogNames(t *testing.T) {
	catalog := NewNativeCatalog()
	expected := []string{
		"A", "Am", "B", "Bb", "C", "Cm", "D",
		"Dm", "E", "Em", "F", "F#", "G", "Gm",
	}

	if names := catalog.Names(); !reflect.DeepEqual(names, expected) {
		t.Fatalf("Names() = %v, want %v", names, expected)
	}
}

func TestNativeCatalogLoadsOwnedVoicing(t *testing.T) {
	catalog := NewNativeCatalog()
	voicing, err := catalog.Load("c", 2)

	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if voicing.Name != "C" || voicing.Variation != 2 {
		t.Fatalf("Load() identity = %q:%d, want C:2", voicing.Name, voicing.Variation)
	}
	if voicing.Strings[0] != (StringPlacement{Fret: 3, Finger: 1}) {
		t.Fatalf("high e placement = %+v, want fret 3 finger 1", voicing.Strings[0])
	}
	if voicing.Strings[5].Fret != -1 {
		t.Fatalf("low E fret = %d, want muted", voicing.Strings[5].Fret)
	}
}

func TestNativeCatalogRejectsMissingVariation(t *testing.T) {
	catalog := NewNativeCatalog()

	if count := catalog.VariationCount("C"); count != 2 {
		t.Fatalf("VariationCount() = %d, want 2", count)
	}
	if _, err := catalog.Load("C", 3); !errors.Is(err, ErrChordNotFound) {
		t.Fatalf("Load() error = %v, want ErrChordNotFound", err)
	}
}
