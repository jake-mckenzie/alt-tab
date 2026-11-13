package tui

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/jake-mckenzie/alt-tab/internal/chords"
)

// ansiSequence strips terminal styling before comparing text snapshots.
var ansiSequence = regexp.MustCompile(`\x1b\[[0-?]*[ -/]*[@-~]`)

// TestTerminalPreviewMatchesSnapshot guards the terminal layout against regressions.
func TestTerminalPreviewMatchesSnapshot(t *testing.T) {
	snapshot, err := os.ReadFile("testdata/preview.txt")
	if err != nil {
		t.Fatalf("read terminal preview: %v", err)
	}

	model := New(chords.NewCatalog())
	model.width = 80
	actual := ansiSequence.ReplaceAllString(model.View().Content, "")

	if normalizeSnapshot(string(snapshot)) != normalizeSnapshot(actual) {
		t.Fatalf(
			"terminal preview differs from snapshot\n--- documented ---\n%s\n--- actual ---\n%s",
			snapshot,
			actual,
		)
	}
}

// normalizeSnapshot removes irrelevant outer and line-ending whitespace.
func normalizeSnapshot(value string) string {
	lines := strings.Split(strings.TrimSpace(value), "\n")
	for index := range lines {
		lines[index] = strings.TrimRight(lines[index], " ")
	}
	return strings.Join(lines, "\n")
}
