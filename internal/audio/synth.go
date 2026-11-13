// Package audio synthesizes chords independently of the platform output backend.
package audio

import (
	"math"
	"sync"

	"github.com/jake-mckenzie/alt-tab/internal/chords"
	"github.com/jake-mckenzie/alt-tab/internal/signal"
)

const sampleRate = 44100

// synth is shared by every output backend. Its callback and UI updates may run concurrently.
type synth struct {
	mu           sync.Mutex
	voicing      chords.Voicing
	notes        []signal.Note
	phases       []float64
	gain, target float64
}

// Set updates the chord and playback gain under the audio callback lock.
func (s *synth) Set(v chords.Voicing, active bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.notes == nil || s.voicing != v {
		s.voicing = v
		s.notes = signal.Notes(v)
		s.phases = make([]float64, len(s.notes))
	}
	s.target = 0
	if active {
		s.target = 0.28 * 0.65
	}
}

// read fills a 44.1 kHz buffer with phase-continuous sine waves and smoothed gain.
func (s *synth) read(buffer []float32) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.notes) == 0 || (s.target == 0 && s.gain < 0.000001) {
		s.gain = 0
		clear(buffer)
		return len(buffer), nil
	}
	for frame := range buffer {
		s.gain += (s.target - s.gain) * 0.0025
		var sample float64
		for i, note := range s.notes {
			sample += math.Sin(s.phases[i])
			s.phases[i] = math.Mod(s.phases[i]+2*math.Pi*note.Frequency/sampleRate, 2*math.Pi)
		}
		buffer[frame] = float32(sample / float64(len(s.notes)) * s.gain)
	}
	return len(buffer), nil
}
