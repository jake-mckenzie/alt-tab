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
	board, legend := renderFretboardParts(voicing, fullNeck, false)
	return board + "\n\n" + legend
}

// renderFretboardParts separates the neck from its position legend for layout.
func renderFretboardParts(
	voicing chords.Voicing,
	fullNeck bool,
	tabNumbers bool,
) (string, string) {
	firstFret, lastFret, cellWidth := fretRange(voicing, fullNeck)
	var output strings.Builder

	if !tabNumbers {
		writeFretLabels(&output, firstFret, lastFret, cellWidth)
	}

	for index, placement := range voicing.Strings {
		fmt.Fprintf(&output, "%s  %c", stringNames[index], stringBoundary(placement))
		for fret := firstFret; fret <= lastFret; fret++ {
			marker := ""
			if placement.Fret == fret {
				marker = fmt.Sprintf("%d", placement.Finger)
				if tabNumbers {
					marker = fmt.Sprintf("%d", placement.Fret)
				}
			}
			output.WriteString(centerText(marker, cellWidth, '-'))
		}
		output.WriteByte('|')
		if index < len(voicing.Strings)-1 {
			output.WriteByte('\n')
		}
	}

	legend := renderFretboardLegend(
		4+(lastFret-firstFret+1)*cellWidth+1,
		fullNeck,
		!tabNumbers,
	)
	return output.String(), legend
}

// renderTabDiagram preserves the neck while replacing fingers with fret numbers.
func renderTabDiagram(voicing chords.Voicing) string {
	board, _ := renderFretboardParts(voicing, false, true)
	return board
}

// renderFretboardLegend explains markers and centers the full-neck legend.
func renderFretboardLegend(width int, centered, includeFingers bool) string {
	fingerLines := []string{
		"Fingers: 1 index  2 middle",
		"         3 ring   4 little",
	}
	symbolLine := "Symbols: O open  X muted"
	if !includeFingers {
		return symbolLine
	}
	if !centered {
		fingerWidth := max(len(fingerLines[0]), len(fingerLines[1]))
		fingerLines[1] += strings.Repeat(" ", fingerWidth-len(fingerLines[1]))
		return fingerLines[0] + "\n" + fingerLines[1] + "\n\n" + symbolLine
	}

	// Keep the two finger lines aligned while centering them as one block.
	fingerWidth := max(len(fingerLines[0]), len(fingerLines[1]))
	fingerPadding := strings.Repeat(" ", max(0, width-fingerWidth)/2)
	symbolPadding := strings.Repeat(" ", max(0, width-len(symbolLine))/2)
	return fingerPadding + fingerLines[0] + "\n" +
		fingerPadding + fingerLines[1] + "\n\n" +
		symbolPadding + symbolLine
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
	output.WriteString("    ")
	output.WriteString(strings.Repeat("─", (last-first+1)*cellWidth))
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
