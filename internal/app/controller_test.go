package app

import (
	"fmt"
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

// TestChordChangesRetainAvailableVoicing covers arrows, direct keys, and variants.
func TestChordChangesRetainAvailableVoicing(t *testing.T) {
	for _, change := range []struct {
		name string
		move func(*Controller)
		want string
	}{
		{"next chord", func(c *Controller) { c.MoveChord(1) }, "B"},
		{"previous chord wraps", func(c *Controller) { c.MoveChord(-1) }, "B"},
		{"direct selection", func(c *Controller) { c.SetChord(1) }, "B"},
		{"same chord", func(c *Controller) { c.SetChord(0) }, "A"},
		{"minor", func(c *Controller) { c.MoveKind(1) }, "Am"},
		{"accidental", func(c *Controller) { c.MoveKind(-1) }, "Ab"},
	} {
		for _, count := range []int{1, 2, 3, 4} {
			t.Run(fmt.Sprintf("%s/%d voicings", change.name, count), func(t *testing.T) {
				catalog := navigationCatalog{targetCount: count}
				controller := NewController(catalog)
				controller.CycleVoicing()
				controller.CycleVoicing()
				change.move(controller)
				wantNumber := 3
				if change.want != "A" && count < wantNumber {
					wantNumber = 1
				}
				if controller.Err() != nil || controller.Name() != change.want || controller.Voicing().Number != wantNumber {
					t.Fatalf("selection = %s:%d, error = %v; want %s:%d",
						controller.Name(), controller.Voicing().Number, controller.Err(), change.want, wantNumber)
				}
				// Returning to a chord with more positions keeps the current selection,
				// including a fallback to 1, rather than restoring a hidden preference.
				controller.SetChord(0)
				if controller.Err() != nil || controller.Voicing().Number != wantNumber {
					t.Fatalf("return to A = %d, error = %v; want %d", controller.Voicing().Number, controller.Err(), wantNumber)
				}
				controller.CycleVoicing()
				if got, want := controller.Voicing().Number, wantNumber%3+1; got != want {
					t.Fatalf("next voicing = %d, want %d", got, want)
				}
			})
		}
	}
}

// navigationCatalog supplies unequal counts absent from the built-in catalog.
type navigationCatalog struct{ targetCount int }

func (navigationCatalog) Names() []string { return []string{"A", "Am", "Ab", "B"} }

func (catalog navigationCatalog) VoicingCount(name string) int {
	if name == "A" {
		return 3
	}
	return catalog.targetCount
}

func (catalog navigationCatalog) Load(name string, number int) (chords.Voicing, error) {
	if number < 1 || number > catalog.VoicingCount(name) {
		return chords.Voicing{}, chords.ErrChordNotFound
	}
	return chords.Voicing{Name: name, Number: number}, nil
}

func TestSetVoicing(t *testing.T) {
	controller := NewController(navigationCatalog{targetCount: 2})
	for _, number := range []int{3, 1, 2, 2} {
		controller.SetVoicing(number)
		if controller.Err() != nil || controller.Voicing().Number != number {
			t.Fatalf("select %d: got %+v, error %v", number, controller.Voicing(), controller.Err())
		}
	}
	for _, number := range []int{-1, 0, 4, 9} {
		controller.SetVoicing(number)
		if controller.Err() != nil || controller.Voicing().Number != 2 {
			t.Fatalf("invalid selection %d changed voicing: %+v, error %v", number, controller.Voicing(), controller.Err())
		}
	}
	controller.MoveChord(1)
	if controller.Voicing().Number != 2 {
		t.Fatal("directly selected voicing was not retained on chord change")
	}
	controller.SetVoicing(3)
	if controller.Err() != nil || controller.Voicing().Number != 2 {
		t.Fatal("selection did not respect the new chord's voicing count")
	}
	controller.CycleVoicing()
	if controller.Voicing().Number != 1 {
		t.Fatal("cycling after direct selection did not wrap")
	}
}

func TestSetVoicingWithoutCatalog(t *testing.T) {
	controller := NewController(nil)
	controller.SetVoicing(1)
	if controller.Err() == nil || controller.Voicing() != (chords.Voicing{}) {
		t.Fatal("selecting without a catalog changed initialization state")
	}
}
