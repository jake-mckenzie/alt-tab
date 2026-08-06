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

// catalog exposes the built-in definitions without leaking their compact form.
type catalog struct{}

// NewCatalog returns the immutable built-in chord catalog.
func NewCatalog() Catalog {
	return catalog{}
}

// Names returns a caller-owned list in display order.
func (catalog) Names() []string {
	names := make([]string, len(chordDefinitions))
	for index, definition := range chordDefinitions {
		names[index] = definition.name
	}
	return names
}

// VoicingCount reports the number of generated and explicit chord positions.
func (catalog) VoicingCount(name string) int {
	definition := findDefinition(name)
	if definition == nil {
		return 0
	}
	return len(definition.voicings)
}

// Load expands a one-based voicing into caller-owned display data.
func (catalog) Load(name string, number int) (Voicing, error) {
	definition := findDefinition(name)
	if definition == nil || number < 1 || number > len(definition.voicings) {
		return Voicing{}, ErrChordNotFound
	}

	return definition.voicings[number-1].expand(definition.name, number), nil
}
