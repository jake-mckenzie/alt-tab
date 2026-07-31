package tui

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"alt-tab/internal/chords"
)

var ansiSequence = regexp.MustCompile(`\x1b\[[0-?]*[ -/]*[@-~]`)

// Keeps the documented wide-layout snapshot synchronized with the real view.
func TestReadmeContainsCurrentTUIOutput(t *testing.T) {
	readme, err := os.ReadFile("../../README.md")
	if err != nil {
		t.Fatalf("read README: %v", err)
	}

	snapshot := textBetween(
		string(readme),
		"<!-- BEGIN VERIFIED TUI OUTPUT -->\n```text\n",
		"\n```\n<!-- END VERIFIED TUI OUTPUT -->",
	)
	if snapshot == "" {
		t.Fatal("README verified TUI output block is missing")
	}

	model := New(chords.NewNativeCatalog())
	model.width = 80
	model.height = 30
	actual := ansiSequence.ReplaceAllString(model.View().Content, "")

	if normalizeSnapshot(snapshot) != normalizeSnapshot(actual) {
		t.Fatalf(
			"README TUI output is stale\n--- documented ---\n%s\n--- actual ---\n%s",
			snapshot,
			actual,
		)
	}
}

func normalizeSnapshot(value string) string {
	lines := strings.Split(strings.TrimSpace(value), "\n")
	for index := range lines {
		lines[index] = strings.TrimRight(lines[index], " ")
	}
	return strings.Join(lines, "\n")
}

func textBetween(value, start, end string) string {
	startIndex := strings.Index(value, start)
	if startIndex < 0 {
		return ""
	}
	after := value[startIndex+len(start):]
	endIndex := strings.Index(after, end)
	if endIndex < 0 {
		return ""
	}
	return after[:endIndex]
}
