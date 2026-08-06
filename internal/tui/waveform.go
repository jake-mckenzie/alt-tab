package tui

import (
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
)

const (
	waveformHeight   = 5
	waveformInterval = 120 * time.Millisecond
)

// waveformPattern describes one repeating wave from crest to trough.
var waveformPattern = [...]int{2, 1, 0, 1, 2, 3, 4, 3}

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

// renderWaveform draws a moving five-row waveform using portable ASCII.
func renderWaveform(width, frame int) string {
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

	for column := 0; column < width; column++ {
		patternIndex := wrap(column+frame, len(waveformPattern))
		rows[waveformPattern[patternIndex]][column] = '*'
	}

	lines := make([]string, len(rows))
	for row := range rows {
		lines[row] = string(rows[row])
	}
	return strings.Join(lines, "\n")
}
