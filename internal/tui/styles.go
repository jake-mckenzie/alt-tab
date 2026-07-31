package tui

import "charm.land/lipgloss/v2"

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

// newStyles creates a readable palette for the reported terminal background.
func newStyles(dark bool) styles {
	choose := lipgloss.LightDark(dark)
	accent := choose(lipgloss.Color("#5B21B6"), lipgloss.Color("#C4B5FD"))
	text := choose(lipgloss.Color("#1F2937"), lipgloss.Color("#F4F4F5"))
	muted := choose(lipgloss.Color("#6B7280"), lipgloss.Color("#A1A1AA"))
	border := choose(lipgloss.Color("#C4B5FD"), lipgloss.Color("#5B4B8A"))

	return styles{
		title: lipgloss.NewStyle().
			Bold(true).
			Foreground(accent),
		subtitle: lipgloss.NewStyle().
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
			Foreground(accent),
		normal: lipgloss.NewStyle().
			Foreground(text),
		muted: lipgloss.NewStyle().
			Foreground(muted),
		accent: lipgloss.NewStyle().
			Foreground(accent),
		err: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#EF4444")),
	}
}
