//go:build raylib

package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/jake-mckenzie/alt-tab/internal/chords"
	"github.com/jake-mckenzie/alt-tab/internal/rayui"
)

// init keeps Raylib's window and graphics context on the process main thread.
func init() {
	runtime.LockOSThread()
}

// main launches the alternate graphical chord viewer.
func main() {
	if err := rayui.Run(chords.NewCatalog()); err != nil {
		fmt.Fprintf(os.Stderr, "alt-tab-raylib: %v\n", err)
		os.Exit(1)
	}
}
