// Package chords provides UI-independent access to chord voicings.
package chords

import "errors"

// StringCount is the number of strings in a standard guitar voicing.
const StringCount = 6

// ErrChordNotFound reports a missing name or numbered variation.
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
	// Names returns supported chord names in display order.
	Names() []string
	// VariationCount returns the number of voicings for one chord.
	VariationCount(name string) int
	// Load returns one chord voicing by its one-based variation number.
	Load(name string, variation int) (Voicing, error)
}
