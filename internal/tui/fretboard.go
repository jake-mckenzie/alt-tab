package tui

import (
	"fmt"
	"strings"

	"github.com/jake-mckenzie/alt-tab/internal/chords"
)

const (
	fullNeckLastFret             = 27
	fullNeckCellWidth            = 3
	fullNeckMinimumTerminalWidth = 98
	compactCellWidth             = 5
)

// stringNames follows standard tablature order from high e to low E.
var stringNames = [chords.StringCount]string{"e", "B", "G", "D", "A", "E"}

// renderFretboard draws standard high-e-to-low-E horizontal tablature.
func renderFretboard(
	voicing chords.Voicing,
	fullNeck bool,
) string {
	firstFret, lastFret, cellWidth := fretRange(voicing, fullNeck)
	var output strings.Builder

	writeFretLabels(&output, firstFret, lastFret, cellWidth)

	for index, placement := range voicing.Strings {
		fmt.Fprintf(&output, "%s  %c", stringNames[index], stringBoundary(placement))
		for fret := firstFret; fret <= lastFret; fret++ {
			marker := ""
			if placement.Fret == fret {
				marker = fmt.Sprintf("%d", placement.Finger)
			}
			output.WriteString(centerText(marker, cellWidth, '-'))
		}
		output.WriteString("|\n")
	}

	output.WriteString("\nFingers: 1 index  2 middle")
	output.WriteString("\n         3 ring   4 little")
	output.WriteString("\n\nSymbols: O open  X muted")
	return output.String()
}

// fretRange keeps common voicings compact while preserving exact fret labels.
func fretRange(
	voicing chords.Voicing,
	fullNeck bool,
) (int, int, int) {
	if fullNeck {
		return 1, fullNeckLastFret, fullNeckCellWidth
	}

	lowest := 0
	highest := 0
	for _, placement := range voicing.Strings {
		if placement.Fret > 0 && (lowest == 0 || placement.Fret < lowest) {
			lowest = placement.Fret
		}
		if placement.Fret > highest {
			highest = placement.Fret
		}
	}

	if highest <= 4 {
		return 1, 4, compactCellWidth
	}
	if lowest == 0 {
		lowest = 1
	}

	return lowest, max(lowest+3, highest), compactCellWidth
}

// writeFretLabels places every complete fret number above its string cell.
func writeFretLabels(output *strings.Builder, first, last, cellWidth int) {
	output.WriteString("    ")
	for fret := first; fret <= last; fret++ {
		output.WriteString(centerText(fmt.Sprintf("%d", fret), cellWidth, ' '))
	}
	output.WriteByte('\n')
}

// stringBoundary replaces the nut marker for open and muted strings.
func stringBoundary(placement chords.StringPlacement) byte {
	switch placement.Fret {
	case -1:
		return 'X'
	case 0:
		return 'O'
	default:
		return '|'
	}
}

// centerText pads a short ASCII marker to a fixed-width fret cell.
func centerText(value string, width int, fill byte) string {
	if len(value) >= width {
		return value
	}

	// Give the left side the extra byte so markers align visually with labels.
	left := (width - len(value) + 1) / 2
	right := width - len(value) - left
	return strings.Repeat(string(fill), left) +
		value +
		strings.Repeat(string(fill), right)
}
