// Package chords provides UI-independent access to chord voicings.
package chords

import "errors"

const StringCount = 6

var ErrChordNotFound = errors.New("chord variation not found")

// StringPlacement describes one string in high-e-to-low-E tab order.
type StringPlacement struct {
	Fret   int
	Finger int
}

// Voicing contains one numbered fingering for a named chord.
type Voicing struct {
	Name      string
	Variation int
	Strings   [StringCount]StringPlacement
}

// Catalog supplies chord data without exposing its storage implementation.
type Catalog interface {
	Names() []string
	VariationCount(name string) int
	Load(name string, variation int) (Voicing, error)
}
