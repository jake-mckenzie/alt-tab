package tui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

// render composes the stacked chord-list and fretboard layout.
func (model Model) render() string {
	width := model.width
	if width < 40 {
		width = 40
	}
	contentWidth := width - 4

	header := lipgloss.JoinHorizontal(
		lipgloss.Bottom,
		model.styles.title.Render("ALT-TAB"),
		"  ",
		model.styles.subtitle.Render("Guitar Chord Viewer"),
	)

	var body string
	if model.showHelp {
		body = model.styles.panel.
			Width(contentWidth - 2).
			Render(model.renderHelp())
	} else {
		body = model.renderStacked(contentWidth)
	}

	footer := model.styles.muted.Render(
		"←/→ chord  ↑/↓ variation  f full neck  ? help  q quit",
	)
	return lipgloss.NewStyle().
		Padding(1, 2).
		Render(header + "\n\n" + body + "\n\n" + footer)
}

func (model Model) renderStacked(width int) string {
	chordList := model.styles.panel.
		Width(width - 2).
		Render(model.renderChordList())
	detail := model.styles.panel.
		Width(width - 2).
		Render(model.renderDetail())

	return lipgloss.JoinVertical(lipgloss.Left, chordList, "", detail)
}

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

func (model Model) renderDetail() string {
	if model.err != nil {
		return model.styles.err.Render(model.err.Error())
	}

	variationCount := model.catalog.VariationCount(model.voicing.Name)
	mode := "compact"
	if model.fullNeck {
		mode = "full neck · frets 1–27"
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

func (model Model) renderHelp() string {
	return model.styles.heading.Render("KEYBOARD HELP") + "\n\n" +
		"← / h     Previous chord\n" +
		"→ / l     Next chord\n" +
		"↑ / k     Previous variation\n" +
		"↓ / j     Next variation\n" +
		"f         Toggle compact/full neck\n" +
		"?         Open or close help\n" +
		"esc       Close help\n" +
		"q         Quit\n\n" +
		model.styles.muted.Render("Press ? or esc to return")
}
