// Package signal derives exact pitch and waveform data from guitar voicings.
package signal

import (
	"fmt"
	"math"

	"github.com/jake-mckenzie/alt-tab/internal/chords"
)

// StandardTuningMIDI lists open strings from high e to low E.
var StandardTuningMIDI = [chords.StringCount]int{64, 59, 55, 50, 45, 40}

// Note contains one sounding string's display and synthesis data.
type Note struct {
	Name      string
	Frequency float64
	MIDI      int
}

// noteNames maps a MIDI pitch class to its sharp note spelling.
var noteNames = [...]string{
	"C", "C#", "D", "D#", "E", "F",
	"F#", "G", "G#", "A", "A#", "B",
}

// Notes converts every sounding string and fret to a named frequency.
func Notes(voicing chords.Voicing) []Note {
	notes := make([]Note, 0, chords.StringCount)
	for index, placement := range voicing.Strings {
		if placement.Fret < 0 {
			continue
		}
		midi := StandardTuningMIDI[index] + placement.Fret
		notes = append(notes, Note{
			Name:      fmt.Sprintf("%s%d", noteNames[wrap(midi, len(noteNames))], midi/12-1),
			Frequency: MIDIFrequency(midi),
			MIDI:      midi,
		})
	}
	return notes
}

// MIDIFrequency converts a MIDI note to equal temperament at A4=440 Hz.
func MIDIFrequency(midi int) float64 {
	return 440 * math.Pow(2, float64(midi-69)/12)
}

// CompositeSample averages ideal equal-amplitude sine waves at one time.
func CompositeSample(seconds float64, notes []Note) float64 {
	if len(notes) == 0 {
		return 0
	}
	var sample float64
	for _, note := range notes {
		sample += math.Sin(2 * math.Pi * note.Frequency * seconds)
	}
	return sample / float64(len(notes))
}

// wrap confines a pitch-class index to its non-empty name table.
func wrap(value, size int) int {
	value %= size
	if value < 0 {
		value += size
	}
	return value
}
