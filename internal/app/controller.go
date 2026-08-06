// Package app owns interface-independent chord navigation state.
package app

import (
	"errors"

	"github.com/jake-mckenzie/alt-tab/internal/chords"
)

// Kind identifies the selected member of a natural chord family.
type Kind int

const (
	Accidental Kind = -1
	Base       Kind = 0
	Minor      Kind = 1
)

// Family groups one natural chord with its available variants.
type Family struct {
	Base       string
	Accidental string
	Minor      string
}

// Controller owns chord selection shared by terminal and graphical frontends.
type Controller struct {
	catalog       chords.Catalog
	families      []Family
	selected      int
	kind          Kind
	voicingNumber int
	voicing       chords.Voicing
	err           error
}

// NewController initializes the first available base chord and voicing.
func NewController(catalog chords.Catalog) *Controller {
	controller := &Controller{catalog: catalog, voicingNumber: 1}
	if catalog == nil {
		controller.err = errors.New("chord catalog is unavailable")
		return controller
	}

	controller.families = BuildFamilies(catalog.Names())
	controller.load()
	return controller
}

// BuildFamilies preserves natural-chord order and attaches known variants.
func BuildFamilies(names []string) []Family {
	families := make([]Family, 0, len(names))
	indexes := make(map[string]int)
	for _, name := range names {
		if len(name) != 1 || name[0] < 'A' || name[0] > 'G' {
			continue
		}
		indexes[name] = len(families)
		families = append(families, Family{Base: name})
	}

	for _, name := range names {
		if len(name) != 2 {
			continue
		}
		index, exists := indexes[name[:1]]
		if !exists {
			continue
		}
		switch name[1] {
		case 'b', '#':
			families[index].Accidental = name
		case 'm':
			families[index].Minor = name
		}
	}
	return families
}

// Families returns a caller-owned view of the chord dial.
func (controller *Controller) Families() []Family {
	return append([]Family(nil), controller.families...)
}

// SelectedIndex returns the active natural-chord family index.
func (controller *Controller) SelectedIndex() int {
	return controller.selected
}

// Kind returns the selected base, accidental, or minor position.
func (controller *Controller) Kind() Kind {
	return controller.kind
}

// Voicing returns the currently loaded caller-owned chord value.
func (controller *Controller) Voicing() chords.Voicing {
	return controller.voicing
}

// VoicingCount returns the number of positions for the active chord.
func (controller *Controller) VoicingCount() int {
	if controller.catalog == nil {
		return 0
	}
	return controller.catalog.VoicingCount(controller.Name())
}

// Err reports catalog initialization or lookup failures.
func (controller *Controller) Err() error {
	return controller.err
}

// Name resolves the highlighted family position to a catalog name.
func (controller *Controller) Name() string {
	if len(controller.families) == 0 {
		return ""
	}
	family := controller.families[controller.selected]
	switch controller.kind {
	case Accidental:
		return family.Accidental
	case Minor:
		return family.Minor
	default:
		return family.Base
	}
}

// MoveChord rotates the natural-chord dial and resets its voicing.
func (controller *Controller) MoveChord(delta int) {
	if len(controller.families) == 0 {
		return
	}
	controller.selected = Wrap(controller.selected+delta, len(controller.families))
	controller.kind = Base
	controller.voicingNumber = 1
	controller.load()
}

// MoveKind selects an available accidental or minor chord.
func (controller *Controller) MoveKind(delta int) {
	if len(controller.families) == 0 {
		return
	}
	family := controller.families[controller.selected]
	next := controller.kind
	if delta < 0 {
		if controller.kind == Minor {
			next = Base
		} else if controller.kind == Base && family.Accidental != "" {
			next = Accidental
		}
	} else if delta > 0 {
		if controller.kind == Accidental {
			next = Base
		} else if controller.kind == Base && family.Minor != "" {
			next = Minor
		}
	}
	if next == controller.kind {
		return
	}
	controller.kind = next
	controller.voicingNumber = 1
	controller.load()
}

// CycleVoicing advances through the active chord's stored positions.
func (controller *Controller) CycleVoicing() {
	count := controller.VoicingCount()
	if count == 0 {
		return
	}
	controller.voicingNumber = Wrap(controller.voicingNumber, count) + 1
	controller.load()
}

// load refreshes the active voicing after a navigation change.
func (controller *Controller) load() {
	if len(controller.families) == 0 {
		controller.err = errors.New("no chords are available")
		return
	}
	controller.voicing, controller.err = controller.catalog.Load(
		controller.Name(),
		controller.voicingNumber,
	)
}

// Wrap confines a possibly negative index to a non-empty collection.
func Wrap(value, size int) int {
	if size <= 0 {
		return 0
	}
	value %= size
	if value < 0 {
		value += size
	}
	return value
}
