package tui

import (
	"os"
	"regexp"
	"testing"
)

// ansiSequence strips terminal styling before comparing text snapshots.
var ansiSequence = regexp.MustCompile(`\x1b\[[0-?]*[ -/]*[@-~]`)

// TestReadmeContainsRaylibPreview keeps the documented desktop preview available.
func TestReadmeContainsRaylibPreview(t *testing.T) {
	readme, err := os.ReadFile("../../README.md")
	if err != nil {
		t.Fatalf("read README: %v", err)
	}

	const preview = "docs/images/raylib-desktop-preview.png"
	if !regexp.MustCompile(`!\[Alt-Tab Raylib desktop interface\]\(` + regexp.QuoteMeta(preview) + `\)`).Match(readme) {
		t.Fatalf("README Raylib preview reference %q is missing", preview)
	}

	if _, err := os.Stat("../../" + preview); err != nil {
		t.Fatalf("read README Raylib preview: %v", err)
	}
}
