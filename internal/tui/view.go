package tui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

const minimumTerminalWidth = 40

// render composes the stacked chord-list and fretboard layout.
func (model Model) render() string {
	width := model.width
	if width < minimumTerminalWidth {
		width = minimumTerminalWidth
	}
	// The outer style contributes two columns of padding on both sides.
	contentWidth := width - 4

	header := lipgloss.JoinHorizontal(
		lipgloss.Bottom,
		model.styles.title.Render("ALT-TAB"),
		"  ",
		model.styles.subtitle.Render(
			fmt.Sprintf("Guitar Chord Viewer · %s", paletteAt(model.theme).name),
		),
	)
	if model.waveform {
		header += "\n" + model.styles.accent.Render(
			renderWaveform(contentWidth, model.waveFrame),
		)
	}

	var body string
	if model.showHelp {
		body = model.styles.panel.
			Width(contentWidth - 2).
			Render(model.renderHelp())
	} else {
		body = model.renderStacked(contentWidth)
	}

	footer := model.styles.muted.Render(
		"←/→ chord  ↑/↓ variation  f neck  t theme  w wave  ? help  q quit",
	)
	return lipgloss.NewStyle().
		Padding(1, 2).
		Render(header + "\n\n" + body + "\n\n" + footer)
}

// renderStacked places the shared chord selector above the active voicing.
func (model Model) renderStacked(width int) string {
	// Lip Gloss adds one border column to each side of the requested width.
	panel := model.styles.panel.Width(width - 2)
	chordList := panel.Render(model.renderChordList())
	detail := panel.Render(model.renderDetail())

	return lipgloss.JoinVertical(lipgloss.Left, chordList, "", detail)
}

// renderChordList draws every chord in one horizontal navigation row.
func (model Model) renderChordList() string {
	var output strings.Builder
	output.WriteString(model.styles.heading.Render("CHORDS"))
	output.WriteString("\n\n")

	for index, name := range model.names {
		if index > 0 {
			output.WriteString("   ")
		}
		if index == model.selected {
			output.WriteString(model.styles.selected.Render("‹ " + name + " ›"))
		} else {
			output.WriteString(model.styles.normal.Render(name))
		}
	}

	return output.String()
}

// renderDetail draws the current chord heading, mode, and fretboard.
func (model Model) renderDetail() string {
	if model.err != nil {
		return model.styles.err.Render(model.err.Error())
	}

	variationCount := model.catalog.VariationCount(model.voicing.Name)
	mode := "compact"
	if model.fullNeck {
		mode = fmt.Sprintf("full neck · frets 1–%d", fullNeckLastFret)
	}

	heading := fmt.Sprintf(
		"%s  ·  variation %d of %d",
		model.voicing.Name,
		model.voicing.Variation,
		variationCount,
	)
	header := model.styles.heading.Render(heading) + "\n" +
		model.styles.accent.Render(mode)
	if model.fullNeck && model.width < fullNeckMinimumTerminalWidth {
		return header + "\n\n" + model.styles.err.Render(
			fmt.Sprintf(
				"Widen the terminal to at least %d columns (currently %d).",
				fullNeckMinimumTerminalWidth,
				model.width,
			),
		)
	}

	return header + "\n\n" +
		model.styles.normal.Render(renderFretboard(model.voicing, model.fullNeck))
}

// renderHelp lists every keyboard command in a dedicated panel.
func (model Model) renderHelp() string {
	return model.styles.heading.Render("KEYBOARD HELP") + "\n\n" +
		"← / h     Previous chord\n" +
		"→ / l     Next chord\n" +
		"↑ / k     Previous variation\n" +
		"↓ / j     Next variation\n" +
		"f         Toggle compact/full neck\n" +
		"t         Cycle color theme\n" +
		"w         Toggle animated waveform\n" +
		"?         Open or close help\n" +
		"esc       Close help\n" +
		"q         Quit\n\n" +
		model.styles.muted.Render("Press ? or esc to return")
}
