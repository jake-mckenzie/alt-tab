package tui

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/jake-mckenzie/alt-tab/internal/chords"
)

const (
	waveformHeight        = 13
	waveformLabelWidth    = 6
	waveformWindowSeconds = 0.025
	spectrumHeight        = 5
	spectrumLabelWidth    = 6
	braillePixelWidth     = 2
	braillePixelHeight    = 4
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

// standardTuningMIDI lists matching note numbers for waveform labels.
var standardTuningMIDI = [chords.StringCount]int{64, 59, 55, 50, 45, 40}

// noteNames maps a MIDI pitch class to its sharp note spelling.
var noteNames = [...]string{
	"C", "C#", "D", "D#", "E", "F",
	"F#", "G", "G#", "A", "A#", "B",
}

// noteSample pairs one sounding string's display name and frequency.
type noteSample struct {
	name      string
	frequency float64
	midi      int
}

// voicingNotes converts each sounding string and fret to a named frequency.
func voicingNotes(voicing chords.Voicing) []noteSample {
	notes := make([]noteSample, 0, chords.StringCount)
	for index, placement := range voicing.Strings {
		if placement.Fret < 0 {
			continue
		}

		midi := standardTuningMIDI[index] + placement.Fret
		semitones := float64(placement.Fret) / 12
		notes = append(notes, noteSample{
			name: fmt.Sprintf(
				"%s%d",
				noteNames[wrap(midi, len(noteNames))],
				midi/12-1,
			),
			frequency: standardTuningHz[index] * math.Pow(2, semitones),
			midi:      midi,
		})
	}
	return notes
}

// compositeWaveSample averages ideal equal-amplitude sine waves at one time.
func compositeWaveSample(seconds float64, notes []noteSample) float64 {
	if len(notes) == 0 {
		return 0
	}

	var sample float64
	for _, note := range notes {
		sample += math.Sin(2 * math.Pi * note.frequency * seconds)
	}
	return sample / float64(len(notes))
}

// waveformPixelRow maps a normalized amplitude to a Braille pixel row.
func waveformPixelRow(sample float64, height int) int {
	center := float64(height-1) / 2
	return int(math.Round(center - sample*center))
}

// absoluteInt returns the magnitude needed by the integer line rasterizer.
func absoluteInt(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

// drawBrailleLine connects two sampled points with Bresenham's line algorithm.
func drawBrailleLine(canvas [][]bool, x0, y0, x1, y1 int) {
	deltaX := absoluteInt(x1 - x0)
	stepX := -1
	if x0 < x1 {
		stepX = 1
	}
	deltaY := -absoluteInt(y1 - y0)
	stepY := -1
	if y0 < y1 {
		stepY = 1
	}
	errorTerm := deltaX + deltaY

	for {
		canvas[y0][x0] = true
		if x0 == x1 && y0 == y1 {
			return
		}
		twiceError := 2 * errorTerm
		if twiceError >= deltaY {
			errorTerm += deltaY
			x0 += stepX
		}
		if twiceError <= deltaX {
			errorTerm += deltaX
			y0 += stepY
		}
	}
}

// brailleCell encodes one two-by-four pixel block as a terminal character.
func brailleCell(canvas [][]bool, cellX, cellY int) rune {
	dotBits := [braillePixelHeight][braillePixelWidth]rune{
		{0x01, 0x08},
		{0x02, 0x10},
		{0x04, 0x20},
		{0x40, 0x80},
	}

	var bits rune
	for pixelY := 0; pixelY < braillePixelHeight; pixelY++ {
		for pixelX := 0; pixelX < braillePixelWidth; pixelX++ {
			if canvas[cellY*braillePixelHeight+pixelY][cellX*braillePixelWidth+pixelX] {
				bits |= dotBits[pixelY][pixelX]
			}
		}
	}
	if bits == 0 {
		return ' '
	}
	return '\u2800' + bits
}

// renderBrailleCanvas compresses subcell pixels into fixed-width text rows.
func renderBrailleCanvas(canvas [][]bool, width int) []string {
	lines := make([]string, waveformHeight)
	for cellY := 0; cellY < waveformHeight; cellY++ {
		var line strings.Builder
		for cellX := 0; cellX < width; cellX++ {
			line.WriteRune(brailleCell(canvas, cellX, cellY))
		}
		lines[cellY] = line.String()
	}
	return lines
}

// renderWaveform draws an ideal 25 ms composite signal for the active voicing.
func renderWaveform(width int, voicing chords.Voicing) string {
	// Reserve the last cell for a right boundary matching the y-axis.
	plotWidth := width - waveformLabelWidth - 1
	if plotWidth <= 0 {
		return ""
	}

	pixelWidth := plotWidth * braillePixelWidth
	pixelHeight := waveformHeight * braillePixelHeight
	canvas := make([][]bool, pixelHeight)
	for row := range canvas {
		canvas[row] = make([]bool, pixelWidth)
	}

	// A dotted subcell line preserves the zero-amplitude reference axis.
	center := pixelHeight / 2
	for column := 0; column < pixelWidth; column++ {
		canvas[center][column] = true
	}

	notes := voicingNotes(voicing)
	previousRow := center
	for column := 0; column < pixelWidth; column++ {
		position := float64(column) / float64(max(1, pixelWidth-1))
		row := waveformPixelRow(
			compositeWaveSample(position*waveformWindowSeconds, notes),
			pixelHeight,
		)
		drawBrailleLine(canvas, max(0, column-1), previousRow, column, row)
		previousRow = row
	}

	plotLines := renderBrailleCanvas(canvas, plotWidth)
	lines := make([]string, 0, waveformHeight+1)
	for row, plotLine := range plotLines {
		label := "     |"
		boundary := "|"
		switch row {
		case 0:
			label = "+1.0 |"
		case waveformHeight / 2:
			label = " 0.0 +"
			boundary = "+"
		case waveformHeight - 1:
			label = "-1.0 |"
		}
		lines = append(lines, label+plotLine+boundary)
	}
	lines = append(lines, renderTimeAxis(width))
	return strings.Join(lines, "\n")
}

// renderTimeAxis labels the exact signal window represented by the plot.
func renderTimeAxis(width int) string {
	left := strings.Repeat(" ", waveformLabelWidth) + "0 ms"
	right := fmt.Sprintf("%.0f ms", waveformWindowSeconds*1000)
	gap := max(1, width-len(left)-len(right))
	return left + strings.Repeat(" ", gap) + right
}

// renderSpectrum plots normalized note peaks across a frequency scale.
func renderSpectrum(width int, voicing chords.Voicing) string {
	// Reserve one cell for a closing axis so the graph has a complete frame.
	plotWidth := width - spectrumLabelWidth - 1
	notes := voicingNotes(voicing)
	if len(notes) == 0 || plotWidth < 10 {
		return ""
	}

	ordered := append([]noteSample(nil), notes...)
	sort.Slice(ordered, func(left, right int) bool {
		return ordered[left].midi < ordered[right].midi
	})
	lowest := ordered[0]
	highest := ordered[len(ordered)-1]
	span := max(1, highest.midi-lowest.midi)
	positions := make([]int, len(ordered))
	labels := []rune(strings.Repeat(" ", plotWidth))

	for index, note := range ordered {
		positions[index] = int(math.Round(
			float64(note.midi-lowest.midi) / float64(span) * float64(plotWidth-1),
		))
		placeScaleLabel(labels, note.name, positions[index])
	}

	lines := make([]string, 0, spectrumHeight+5)
	for row := 0; row < spectrumHeight; row++ {
		plot := []rune(strings.Repeat(" ", plotWidth))
		for _, position := range positions {
			plot[position] = '█'
		}
		label := "    |"
		switch row {
		case 0:
			label = "1.0 |"
		case spectrumHeight / 2:
			label = "0.5 |"
		}
		lines = append(lines, label+string(plot)+"|")
	}

	axis := []rune(strings.Repeat("─", plotWidth))
	for _, position := range positions {
		axis[position] = '┼'
	}
	lines = append(lines, "0.0 +"+string(axis)+"+")
	prefix := strings.Repeat(" ", spectrumLabelWidth)
	lines = append(lines, prefix+string(labels))
	leftFrequency := fmt.Sprintf("%.0f Hz", lowest.frequency)
	rightFrequency := fmt.Sprintf("%.0f Hz", highest.frequency)
	frequencyGap := max(1, plotWidth-len(leftFrequency)-len(rightFrequency))
	lines = append(
		lines,
		prefix+leftFrequency+strings.Repeat(" ", frequencyGap)+rightFrequency,
		"",
		centerDisplayText(renderNoteLegend(notes), width),
	)
	return strings.Join(lines, "\n")
}

// placeScaleLabel centers an ASCII note name while avoiding scale boundaries.
func placeScaleLabel(line []rune, label string, position int) {
	characters := []rune(label)
	start := min(max(0, position-len(characters)/2), len(line)-len(characters))
	for index, character := range characters {
		line[start+index] = character
	}
}

// renderNoteLegend lists the pitches combined into the displayed signal.
func renderNoteLegend(notes []noteSample) string {
	names := make([]string, len(notes))
	for index, note := range notes {
		names[index] = note.name
	}
	return "Notes: " + strings.Join(names, "  ")
}
