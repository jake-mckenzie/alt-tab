package app

import (
	"testing"

	"github.com/jake-mckenzie/alt-tab/internal/chords"
)

// TestControllerNavigatesCatalog checks family, kind, and voicing movement.
func TestControllerNavigatesCatalog(t *testing.T) {
	controller := NewController(chords.NewCatalog())
	if controller.Name() != "A" || controller.Voicing().Number != 1 {
		t.Fatalf("initial selection = %s:%d", controller.Name(), controller.Voicing().Number)
	}
	controller.MoveKind(1)
	if controller.Name() != "Am" {
		t.Fatalf("minor selection = %s, want Am", controller.Name())
	}
	controller.MoveChord(1)
	controller.MoveKind(-1)
	if controller.Name() != "Bb" {
		t.Fatalf("accidental selection = %s, want Bb", controller.Name())
	}
	controller.CycleVoicing()
	if controller.Voicing().Number != 2 {
		t.Fatalf("cycled voicing = %d, want 2", controller.Voicing().Number)
	}
}

// TestBuildFamiliesDoesNotInventVariants checks catalog-backed dial choices.
func TestBuildFamiliesDoesNotInventVariants(t *testing.T) {
	families := BuildFamilies([]string{"A", "Am", "B", "Bb", "C"})
	if len(families) != 3 || families[0].Minor != "Am" ||
		families[1].Accidental != "Bb" || families[2].Minor != "" {
		t.Fatalf("families = %+v", families)
	}
}
