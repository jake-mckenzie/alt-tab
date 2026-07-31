package tui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

const wideLayoutMinimum = 78

// render composes responsive wide and narrow layouts.
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
	useWideLayout := width >= wideLayoutMinimum &&
		model.height >= 24 &&
		(!model.fullNeck || width >= 88)
	if model.showHelp {
		body = model.styles.panel.
			Width(contentWidth - 2).
			Render(model.renderHelp())
	} else if useWideLayout {
		body = model.renderWide(contentWidth)
	} else {
		body = model.renderNarrow(contentWidth)
	}

	footer := model.styles.muted.Render(
		"↑/↓ chord  ←/→ variation  f full neck  ? help  q quit",
	)
	return lipgloss.NewStyle().
		Padding(1, 2).
		Render(header + "\n\n" + body + "\n\n" + footer)
}

func (model Model) renderWide(width int) string {
	const sidebarWidth = 15
	gap := 2
	detailWidth := width - sidebarWidth - gap - 6

	sidebar := model.styles.panel.
		Width(sidebarWidth).
		Render(model.renderChordList())
	detail := model.styles.panel.
		Width(detailWidth).
		Render(model.renderDetail(detailWidth))

	return lipgloss.JoinHorizontal(lipgloss.Top, sidebar, strings.Repeat(" ", gap), detail)
}

func (model Model) renderNarrow(width int) string {
	selector := model.styles.panel.
		Width(width - 2).
		Render(model.renderChordSelector())
	detail := model.styles.panel.
		Width(width - 2).
		Render(model.renderDetail(width - 2))

	return lipgloss.JoinVertical(lipgloss.Left, selector, detail)
}

func (model Model) renderChordList() string {
	var output strings.Builder
	output.WriteString(model.styles.heading.Render("CHORDS"))
	output.WriteString("\n\n")

	for index, name := range model.names {
		if index == model.selected {
			output.WriteString(model.styles.selected.Render("› " + name))
		} else {
			output.WriteString(model.styles.normal.Render("  " + name))
		}
		output.WriteByte('\n')
	}

	return strings.TrimSuffix(output.String(), "\n")
}

func (model Model) renderChordSelector() string {
	if len(model.names) == 0 {
		return model.styles.err.Render("No chords available")
	}

	previous := model.names[wrap(model.selected-1, len(model.names))]
	current := model.names[model.selected]
	next := model.names[wrap(model.selected+1, len(model.names))]
	return model.styles.muted.Render(previous) + "   " +
		model.styles.selected.Render("‹  "+current+"  ›") + "   " +
		model.styles.muted.Render(next)
}

func (model Model) renderDetail(availableWidth int) string {
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
	return model.styles.heading.Render(heading) + "\n" +
		model.styles.accent.Render(mode) + "\n\n" +
		model.styles.normal.Render(
			renderFretboard(model.voicing, model.fullNeck, availableWidth),
		)
}

func (model Model) renderHelp() string {
	return model.styles.heading.Render("KEYBOARD HELP") + "\n\n" +
		"↑ / k     Previous chord\n" +
		"↓ / j     Next chord\n" +
		"← / h     Previous variation\n" +
		"→ / l     Next variation\n" +
		"f         Toggle compact/full neck\n" +
		"?         Open or close help\n" +
		"esc       Close help\n" +
		"q         Quit\n\n" +
		model.styles.muted.Render("Press ? or esc to return")
}
