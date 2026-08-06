package chords

import "strings"

// mutedOffset marks a string that a shape does not sound.
const mutedOffset int8 = -1

// shapeID identifies one stored open or movable fingering pattern.
type shapeID uint8

const (
	openA shapeID = iota
	openAMinor
	openC
	openD
	openDMinor
	openE
	openEMinor
	openG
	eMajorBarre
	aMajorBarre
	eMinorBarre
	aMinorBarre
	dMajorMovable
	dMinorMovable
)

// chordShape stores compact fret offsets and finger assignments in tab order.
type chordShape struct {
	frets   [StringCount]int8
	fingers [StringCount]uint8
	movable bool
}

// voicingSpec reuses a shape at one fret instead of repeating six placements.
type voicingSpec struct {
	shape    shapeID
	baseFret int8
}

// chordDefinition groups a name with its ordered voicings.
type chordDefinition struct {
	name     string
	voicings []voicingSpec
}

// chordShapes contains open forms plus reusable major/minor E, A, and D forms.
var chordShapes = [...]chordShape{
	openA: {
		frets:   [StringCount]int8{0, 2, 2, 2, 0, mutedOffset},
		fingers: [StringCount]uint8{0, 3, 2, 1, 0, 0},
	},
	openAMinor: {
		frets:   [StringCount]int8{0, 1, 2, 2, 0, mutedOffset},
		fingers: [StringCount]uint8{0, 1, 3, 2, 0, 0},
	},
	openC: {
		frets:   [StringCount]int8{0, 1, 0, 2, 3, mutedOffset},
		fingers: [StringCount]uint8{0, 1, 0, 2, 3, 0},
	},
	openD: {
		frets:   [StringCount]int8{2, 3, 2, 0, mutedOffset, mutedOffset},
		fingers: [StringCount]uint8{2, 3, 1, 0, 0, 0},
	},
	openDMinor: {
		frets:   [StringCount]int8{1, 3, 2, 0, mutedOffset, mutedOffset},
		fingers: [StringCount]uint8{1, 3, 2, 0, 0, 0},
	},
	openE: {
		frets:   [StringCount]int8{0, 0, 1, 2, 2, 0},
		fingers: [StringCount]uint8{0, 0, 1, 3, 2, 0},
	},
	openEMinor: {
		frets:   [StringCount]int8{0, 0, 0, 2, 2, 0},
		fingers: [StringCount]uint8{0, 0, 0, 3, 2, 0},
	},
	openG: {
		frets:   [StringCount]int8{3, 0, 0, 0, 2, 3},
		fingers: [StringCount]uint8{3, 0, 0, 0, 1, 2},
	},
	eMajorBarre: {
		frets:   [StringCount]int8{0, 0, 1, 2, 2, 0},
		fingers: [StringCount]uint8{1, 1, 2, 4, 3, 1},
		movable: true,
	},
	aMajorBarre: {
		frets:   [StringCount]int8{0, 2, 2, 2, 0, mutedOffset},
		fingers: [StringCount]uint8{1, 4, 3, 2, 1, 0},
		movable: true,
	},
	eMinorBarre: {
		frets:   [StringCount]int8{0, 0, 0, 2, 2, 0},
		fingers: [StringCount]uint8{1, 1, 1, 4, 3, 1},
		movable: true,
	},
	aMinorBarre: {
		frets:   [StringCount]int8{0, 1, 2, 2, 0, mutedOffset},
		fingers: [StringCount]uint8{1, 2, 4, 3, 1, 0},
		movable: true,
	},
	dMajorMovable: {
		frets:   [StringCount]int8{2, 3, 2, 0, mutedOffset, mutedOffset},
		fingers: [StringCount]uint8{3, 4, 2, 1, 0, 0},
		movable: true,
	},
	dMinorMovable: {
		frets:   [StringCount]int8{1, 3, 2, 0, mutedOffset, mutedOffset},
		fingers: [StringCount]uint8{2, 4, 3, 1, 0, 0},
		movable: true,
	},
}

// chordDefinitions keeps names unique and derives voicing numbers from order.
var chordDefinitions = [...]chordDefinition{
	{"A", []voicingSpec{{openA, 0}, {eMajorBarre, 5}, {dMajorMovable, 7}}},
	{"Am", []voicingSpec{{openAMinor, 0}, {eMinorBarre, 5}, {dMinorMovable, 7}}},
	{"B", []voicingSpec{{aMajorBarre, 2}, {eMajorBarre, 7}, {dMajorMovable, 9}}},
	{"Bb", []voicingSpec{{aMajorBarre, 1}, {eMajorBarre, 6}, {dMajorMovable, 8}}},
	{"C", []voicingSpec{{openC, 0}, {aMajorBarre, 3}, {dMajorMovable, 10}}},
	{"Cm", []voicingSpec{{aMinorBarre, 3}, {eMinorBarre, 8}, {dMinorMovable, 10}}},
	{"D", []voicingSpec{{openD, 0}, {aMajorBarre, 5}, {dMajorMovable, 12}}},
	{"Dm", []voicingSpec{{openDMinor, 0}, {aMinorBarre, 5}, {dMinorMovable, 12}}},
	{"E", []voicingSpec{{openE, 0}, {aMajorBarre, 7}, {dMajorMovable, 2}}},
	{"Em", []voicingSpec{{openEMinor, 0}, {aMinorBarre, 7}, {dMinorMovable, 2}}},
	{"F", []voicingSpec{{eMajorBarre, 1}, {aMajorBarre, 8}, {dMajorMovable, 3}}},
	{"F#", []voicingSpec{{eMajorBarre, 2}, {aMajorBarre, 9}, {dMajorMovable, 4}}},
	{"G", []voicingSpec{{openG, 0}, {eMajorBarre, 3}, {dMajorMovable, 5}}},
	{"Gm", []voicingSpec{{eMinorBarre, 3}, {aMinorBarre, 10}, {dMinorMovable, 5}}},
}

// expand materializes one compact specification for rendering and audio math.
func (spec voicingSpec) expand(name string, number int) Voicing {
	shape := chordShapes[spec.shape]
	voicing := Voicing{Name: name, Number: number}
	for index, offset := range shape.frets {
		fret := int(offset)
		if shape.movable && offset != mutedOffset {
			fret += int(spec.baseFret)
		}
		voicing.Strings[index] = StringPlacement{
			Fret:   fret,
			Finger: int(shape.fingers[index]),
		}
	}
	return voicing
}

// findDefinition performs a small case-insensitive lookup without a second map.
func findDefinition(name string) *chordDefinition {
	for index := range chordDefinitions {
		if strings.EqualFold(chordDefinitions[index].name, name) {
			return &chordDefinitions[index]
		}
	}
	return nil
}
