package chords

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

// standardStringPitchClasses follows high-e-to-low-E standard tuning.
var standardStringPitchClasses = [StringCount]int{4, 11, 7, 2, 9, 4}

// TestCatalogNames verifies the ordered, caller-owned name list.
func TestCatalogNames(t *testing.T) {
	catalog := NewCatalog()
	expected := []string{
		"A", "Am", "B", "Bb", "C", "Cm", "D",
		"Dm", "E", "Em", "F", "F#", "G", "Gm",
	}

	names := catalog.Names()
	if !reflect.DeepEqual(names, expected) {
		t.Fatalf("Names() = %v, want %v", names, expected)
	}
	names[0] = "changed"
	if catalog.Names()[0] != "A" {
		t.Fatal("Names() exposed mutable catalog storage")
	}
}

// TestCatalogLoadsGeneratedVoicing verifies lookup and shape expansion.
func TestCatalogLoadsGeneratedVoicing(t *testing.T) {
	catalog := NewCatalog()
	voicing, err := catalog.Load("c", 3)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if voicing.Name != "C" || voicing.Number != 3 {
		t.Fatalf("Load() identity = %q:%d, want C:3", voicing.Name, voicing.Number)
	}

	expected := [StringCount]StringPlacement{
		{Fret: 12, Finger: 3},
		{Fret: 13, Finger: 4},
		{Fret: 12, Finger: 2},
		{Fret: 10, Finger: 1},
		{Fret: -1, Finger: 0},
		{Fret: -1, Finger: 0},
	}
	if voicing.Strings != expected {
		t.Fatalf("C:3 strings = %+v, want %+v", voicing.Strings, expected)
	}
}

// TestCatalogRejectsMissingVoicing verifies invalid names and numbers.
func TestCatalogRejectsMissingVoicing(t *testing.T) {
	catalog := NewCatalog()
	if count := catalog.VoicingCount("C"); count != 3 {
		t.Fatalf("VoicingCount() = %d, want 3", count)
	}
	for _, number := range []int{0, 4} {
		if _, err := catalog.Load("C", number); !errors.Is(err, ErrChordNotFound) {
			t.Fatalf("Load(C, %d) error = %v, want ErrChordNotFound", number, err)
		}
	}
	if _, err := catalog.Load("missing", 1); !errors.Is(err, ErrChordNotFound) {
		t.Fatalf("Load(missing, 1) error = %v, want ErrChordNotFound", err)
	}
}

// TestEveryVoicingMatchesItsChord validates pitch content and playable fingers.
func TestEveryVoicingMatchesItsChord(t *testing.T) {
	catalog := NewCatalog()
	for _, name := range catalog.Names() {
		root := chordRootPitchClass(name)
		third := 4
		if strings.HasSuffix(name, "m") {
			third = 3
		}
		wanted := map[int]bool{0: true, third: true, 7: true}

		for number := 1; number <= catalog.VoicingCount(name); number++ {
			voicing, err := catalog.Load(name, number)
			if err != nil {
				t.Fatalf("Load(%s, %d): %v", name, number, err)
			}
			assertValidVoicing(t, voicing, root, wanted)
		}
	}
}

// chordRootPitchClass converts the supported chord spelling to C-based pitch.
func chordRootPitchClass(name string) int {
	roots := map[byte]int{'C': 0, 'D': 2, 'E': 4, 'F': 5, 'G': 7, 'A': 9, 'B': 11}
	root := roots[name[0]]
	if len(name) > 1 && name[1] == 'b' {
		root--
	}
	if len(name) > 1 && name[1] == '#' {
		root++
	}
	return (root + 12) % 12
}

// assertValidVoicing checks every string and requires all three chord tones.
func assertValidVoicing(
	t *testing.T,
	voicing Voicing,
	root int,
	wanted map[int]bool,
) {
	t.Helper()
	seen := make(map[int]bool)
	fingerFrets := make(map[int]int)
	for index, placement := range voicing.Strings {
		if placement.Fret < 0 || placement.Fret == 0 {
			if placement.Finger != 0 {
				t.Fatalf("%s:%d string %d uses finger %d at fret %d",
					voicing.Name, voicing.Number, index, placement.Finger, placement.Fret)
			}
			if placement.Fret < 0 {
				continue
			}
		} else {
			if placement.Finger < 1 || placement.Finger > 4 {
				t.Fatalf("%s:%d string %d has invalid finger %d",
					voicing.Name, voicing.Number, index, placement.Finger)
			}
			if fret, exists := fingerFrets[placement.Finger]; exists && fret != placement.Fret {
				t.Fatalf("%s:%d finger %d spans frets %d and %d",
					voicing.Name, voicing.Number, placement.Finger, fret, placement.Fret)
			}
			fingerFrets[placement.Finger] = placement.Fret
		}
		if placement.Fret > 27 {
			t.Fatalf("%s:%d exceeds the displayed neck at fret %d",
				voicing.Name, voicing.Number, placement.Fret)
		}

		interval := (standardStringPitchClasses[index] + placement.Fret - root + 24) % 12
		if !wanted[interval] {
			t.Fatalf("%s:%d string %d produces non-chord interval %d",
				voicing.Name, voicing.Number, index, interval)
		}
		seen[interval] = true
	}
	for interval := range wanted {
		if !seen[interval] {
			t.Fatalf("%s:%d omits chord interval %d", voicing.Name, voicing.Number, interval)
		}
	}
}
