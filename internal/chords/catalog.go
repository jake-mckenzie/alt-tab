// Package chords provides UI-independent access to chord voicings.
package chords

import "errors"

// StringCount is the number of strings in a standard guitar voicing.
const StringCount = 6

// ErrChordNotFound reports a missing name or numbered voicing.
var ErrChordNotFound = errors.New("chord voicing not found")

// StringPlacement describes one string in high-e-to-low-E tab order.
type StringPlacement struct {
	Fret   int
	Finger int
}

// Voicing contains one numbered fingering for a named chord.
type Voicing struct {
	Name    string
	Number  int
	Strings [StringCount]StringPlacement
}

// Catalog supplies chord data without exposing its storage implementation.
type Catalog interface {
	// Names returns supported chord names in display order.
	Names() []string
	// VoicingCount returns the number of voicings for one chord.
	VoicingCount(name string) int
	// Load returns one chord voicing by its one-based number.
	Load(name string, number int) (Voicing, error)
}
