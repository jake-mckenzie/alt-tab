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
	tabLineWidth                 = 21
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

	writeFretboardLegend(
		&output,
		4+(lastFret-firstFret+1)*cellWidth+1,
		fullNeck,
	)
	return output.String()
}

// renderTabDiagram places each string's fret number in one aligned tab column.
func renderTabDiagram(voicing chords.Voicing) string {
	var output strings.Builder
	for index, placement := range voicing.Strings {
		marker := "X"
		if placement.Fret >= 0 {
			marker = fmt.Sprintf("%d", placement.Fret)
		}
		fmt.Fprintf(
			&output,
			"%s  %s",
			stringNames[index],
			centerText(marker, tabLineWidth, '-'),
		)
		if index < len(voicing.Strings)-1 {
			output.WriteByte('\n')
		}
	}
	return output.String()
}

// writeFretboardLegend explains markers and centers the full-neck legend.
func writeFretboardLegend(output *strings.Builder, width int, centered bool) {
	fingerLines := []string{
		"Fingers: 1 index  2 middle",
		"         3 ring   4 little",
	}
	symbolLine := "Symbols: O open  X muted"
	if !centered {
		fmt.Fprintf(output, "\n%s\n%s\n\n%s", fingerLines[0], fingerLines[1], symbolLine)
		return
	}

	// Keep the two finger lines aligned while centering them as one block.
	fingerWidth := max(len(fingerLines[0]), len(fingerLines[1]))
	fingerPadding := strings.Repeat(" ", max(0, width-fingerWidth)/2)
	symbolPadding := strings.Repeat(" ", max(0, width-len(symbolLine))/2)
	fmt.Fprintf(
		output,
		"\n%s%s\n%s%s\n\n%s%s",
		fingerPadding,
		fingerLines[0],
		fingerPadding,
		fingerLines[1],
		symbolPadding,
		symbolLine,
	)
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
