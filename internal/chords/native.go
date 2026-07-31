package chords

/*
#cgo CFLAGS: -std=c11 -Wall -Wextra -pedantic -I${SRCDIR}/../../include
#include <stdlib.h>
#include "backend/chord_api.h"
*/
import "C"

import (
	"strings"
	"unsafe"
)

// NativeCatalog reads the built-in C chord library through its public API.
type NativeCatalog struct{}

// NewNativeCatalog returns the production chord-data adapter.
func NewNativeCatalog() NativeCatalog {
	return NativeCatalog{}
}

// Names returns a Go-owned copy of every supported chord name.
func (NativeCatalog) Names() []string {
	count := int(C.altTabChordCount())
	names := make([]string, 0, count)

	for index := 0; index < count; index++ {
		name := C.altTabChordNameAt(C.size_t(index))
		if name != nil {
			names = append(names, C.GoString(name))
		}
	}

	return names
}

// VariationCount returns the number of stored voicings for a chord.
func (NativeCatalog) VariationCount(name string) int {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	return int(C.altTabChordVariationCount(cName))
}

// Load copies one native voicing into Go-owned memory.
func (catalog NativeCatalog) Load(name string, variation int) (Voicing, error) {
	if variation < 1 || variation > catalog.VariationCount(name) {
		return Voicing{}, ErrChordNotFound
	}

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	var native C.AltTabChordVoicing
	if C.altTabChordLoad(cName, C.int(variation), &native) == 0 {
		return Voicing{}, ErrChordNotFound
	}

	voicing := Voicing{
		Name:      catalog.canonicalName(name),
		Variation: int(native.variation),
	}
	for index := range voicing.Strings {
		voicing.Strings[index] = StringPlacement{
			Fret:   int(native.strings[index].fret),
			Finger: int(native.strings[index].finger),
		}
	}

	return voicing, nil
}

// canonicalName preserves the library's display casing after flexible lookup.
func (catalog NativeCatalog) canonicalName(name string) string {
	for _, candidate := range catalog.Names() {
		if strings.EqualFold(candidate, name) {
			return candidate
		}
	}

	return name
}
