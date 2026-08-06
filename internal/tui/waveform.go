package tui

import (
	"math"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/jake-mckenzie/alt-tab/internal/chords"
)

const (
	waveformHeight        = 7
	waveformInterval      = 40 * time.Millisecond
	waveformFrameCount    = 4096
	waveformSpatialCycles = 1.25
	waveformPhaseStep     = 0.08
)

// standardTuningHz lists open-string frequencies from high e to low E.
var standardTuningHz = [chords.StringCount]float64{
	329.63,
	246.94,
	196.00,
	146.83,
	110.00,
	82.41,
}

// waveformTickMsg advances one animation that still owns its timer generation.
type waveformTickMsg struct {
	generation uint64
}

// waveformTick schedules exactly one future animation frame.
func waveformTick(generation uint64) tea.Cmd {
	return tea.Tick(waveformInterval, func(time.Time) tea.Msg {
		return waveformTickMsg{generation: generation}
	})
}

// voicingFrequencies converts each sounding string and fret to hertz.
func voicingFrequencies(voicing chords.Voicing) []float64 {
	frequencies := make([]float64, 0, chords.StringCount)
	for index, placement := range voicing.Strings {
		if placement.Fret < 0 {
			continue
		}
		semitones := float64(placement.Fret) / 12
		frequencies = append(
			frequencies,
			standardTuningHz[index]*math.Pow(2, semitones),
		)
	}
	return frequencies
}

// compositeWaveSample combines the chord's note frequencies at one position.
func compositeWaveSample(
	position float64,
	phase float64,
	frequencies []float64,
) float64 {
	if len(frequencies) == 0 {
		return 0
	}

	lowest := frequencies[0]
	for _, frequency := range frequencies[1:] {
		lowest = min(lowest, frequency)
	}

	var sample float64
	for _, frequency := range frequencies {
		ratio := frequency / lowest
		angle := 2*math.Pi*waveformSpatialCycles*ratio*position - phase*ratio
		sample += math.Sin(angle)
	}

	// Square-root scaling keeps chords visible without clipping every peak.
	return max(-1, min(1, sample/math.Sqrt(float64(len(frequencies)))))
}

// renderWaveform draws the active chord's animated composite wave in ASCII.
func renderWaveform(width, frame int, voicing chords.Voicing) string {
	if width <= 0 {
		return ""
	}

	rows := make([][]byte, waveformHeight)
	for row := range rows {
		fill := byte(' ')
		if row == waveformHeight/2 {
			fill = '-'
		}
		rows[row] = []byte(strings.Repeat(string(fill), width))
	}

	frequencies := voicingFrequencies(voicing)
	phase := float64(frame) * waveformPhaseStep
	for column := 0; column < width; column++ {
		position := float64(column) / float64(max(1, width-1))
		sample := compositeWaveSample(position, phase, frequencies)
		center := float64(waveformHeight-1) / 2
		row := int(math.Round(center - sample*center))
		rows[row][column] = '*'
	}

	lines := make([]string, len(rows))
	for row := range rows {
		lines[row] = string(rows[row])
	}
	return strings.Join(lines, "\n")
}
