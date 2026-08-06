package chords

/*
#cgo CFLAGS: -std=c11 -Wall -Wextra -pedantic
#include <stdlib.h>
#include "chords.h"
*/
import "C"

import "unsafe"

// nativeCatalog reads the built-in C chord table through its small C API.
type nativeCatalog struct{}

// NewCatalog returns the built-in chord catalog.
func NewCatalog() Catalog {
	return nativeCatalog{}
}

// Names returns a Go-owned copy of every supported chord name.
func (nativeCatalog) Names() []string {
	count := int(C.altTabChordNameCount())
	names := make([]string, 0, count)

	for index := 0; index < count; index++ {
		name := C.altTabChordNameAt(C.size_t(index))
		names = append(names, C.GoString(name))
	}

	return names
}

// VoicingCount returns the number of stored voicings for a chord.
func (nativeCatalog) VoicingCount(name string) int {
	cName := C.CString(name)
	// C.CString allocates outside Go's heap and must be released explicitly.
	defer C.free(unsafe.Pointer(cName))

	return int(C.altTabChordNamedVoicingCount(cName))
}

// Load copies one native voicing into Go-owned memory.
func (nativeCatalog) Load(
	name string,
	number int,
) (Voicing, error) {
	cName := C.CString(name)
	// Keep the temporary lookup string alive until the C call has returned.
	defer C.free(unsafe.Pointer(cName))

	native := C.altTabChordFind(cName, C.int(number))
	if native == nil {
		return Voicing{}, ErrChordNotFound
	}

	// Copy the immutable C record so no C pointer escapes this method.
	voicing := Voicing{
		Name:   C.GoString(native.name),
		Number: int(native.number),
	}
	for index := range voicing.Strings {
		voicing.Strings[index] = StringPlacement{
			Fret:   int(native.strings[index].fret),
			Finger: int(native.strings[index].finger),
		}
	}

	return voicing, nil
}
