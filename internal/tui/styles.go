package tui

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

// colorPair provides one color for light terminals and one for dark terminals.
type colorPair struct {
	light string
	dark  string
}

// palette defines every color needed to render one complete theme.
type palette struct {
	name       string
	accent     colorPair
	text       colorPair
	muted      colorPair
	border     colorPair
	selection  colorPair
	errorColor colorPair
}

// palettes are ordered exactly as the runtime theme cycle presents them.
var palettes = []palette{
	{
		name:       "Synthwave",
		accent:     colorPair{light: "#A21CAF", dark: "#FF70D9"},
		text:       colorPair{light: "#241331", dark: "#FFF7FF"},
		muted:      colorPair{light: "#6B5A78", dark: "#BCAAC8"},
		border:     colorPair{light: "#7C3AED", dark: "#A78BFA"},
		selection:  colorPair{light: "#F3E8FF", dark: "#3B155E"},
		errorColor: colorPair{light: "#B91C1C", dark: "#FF6B81"},
	},
	{
		name:       "Tidal",
		accent:     colorPair{light: "#0369A1", dark: "#38BDF8"},
		text:       colorPair{light: "#0C2533", dark: "#E6F8FF"},
		muted:      colorPair{light: "#526A78", dark: "#8EB6C7"},
		border:     colorPair{light: "#0891B2", dark: "#22D3EE"},
		selection:  colorPair{light: "#CFFAFE", dark: "#123B4A"},
		errorColor: colorPair{light: "#B91C1C", dark: "#FF7A90"},
	},
	{
		name:       "Ember",
		accent:     colorPair{light: "#C2410C", dark: "#FB923C"},
		text:       colorPair{light: "#321A0F", dark: "#FFF7ED"},
		muted:      colorPair{light: "#786052", dark: "#C6A58E"},
		border:     colorPair{light: "#D97706", dark: "#FBBF24"},
		selection:  colorPair{light: "#FFEDD5", dark: "#4A2615"},
		errorColor: colorPair{light: "#B91C1C", dark: "#FF6B6B"},
	},
	{
		name:       "Evergreen",
		accent:     colorPair{light: "#047857", dark: "#34D399"},
		text:       colorPair{light: "#102A22", dark: "#ECFDF5"},
		muted:      colorPair{light: "#536B62", dark: "#91B8A9"},
		border:     colorPair{light: "#059669", dark: "#6EE7B7"},
		selection:  colorPair{light: "#D1FAE5", dark: "#173D31"},
		errorColor: colorPair{light: "#B91C1C", dark: "#FF7A90"},
	},
}

// color selects the variant that contrasts with the reported terminal mode.
func (pair colorPair) color(dark bool) color.Color {
	if dark {
		return lipgloss.Color(pair.dark)
	}
	return lipgloss.Color(pair.light)
}

// paletteAt returns a valid palette for any positive or negative index.
func paletteAt(index int) palette {
	return palettes[wrap(index, len(palettes))]
}

// styles groups the palette-dependent presentation rules used by each view.
type styles struct {
	title    lipgloss.Style
	subtitle lipgloss.Style
	panel    lipgloss.Style
	heading  lipgloss.Style
	selected lipgloss.Style
	normal   lipgloss.Style
	muted    lipgloss.Style
	accent   lipgloss.Style
	err      lipgloss.Style
}

// newStyles creates presentation rules for one theme and terminal mode.
func newStyles(dark bool, themeIndex int) styles {
	palette := paletteAt(themeIndex)
	accent := palette.accent.color(dark)
	text := palette.text.color(dark)
	muted := palette.muted.color(dark)
	border := palette.border.color(dark)
	selection := palette.selection.color(dark)

	return styles{
		title: lipgloss.NewStyle().
			Bold(true).
			Foreground(accent),
		subtitle: lipgloss.NewStyle().
			Italic(true).
			Foreground(muted),
		panel: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(border).
			Padding(0, 1),
		heading: lipgloss.NewStyle().
			Bold(true).
			Foreground(text),
		selected: lipgloss.NewStyle().
			Bold(true).
			Foreground(accent).
			Background(selection),
		normal: lipgloss.NewStyle().
			Foreground(text),
		muted: lipgloss.NewStyle().
			Foreground(muted),
		accent: lipgloss.NewStyle().
			Foreground(accent),
		err: lipgloss.NewStyle().
			Bold(true).
			Foreground(palette.errorColor.color(dark)),
	}
}
