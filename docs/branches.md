# Maintained branches

| Branch | Purpose | Run |
| --- | --- | --- |
| `master` | Original pure-Go Bubble Tea terminal application | `make run` |
| `raylib` | Go application with Raylib desktop UI and synthesized audio | `make run-raylib` |

Switch with `git switch master` or `git switch raylib`. The Raylib branch also
retains the TUI and shares its catalog, navigation, and signal packages; its
existing `make run` target still launches the terminal application.

## Raylib history

Remote tips were verified on 2026-10-30. The two implementations diverged from
`e2b9e38` (2026-08-06); neither branch contains the other's later commits.

| Original branch | Last commit before organization | Differences |
| --- | --- | --- |
| `feature/raylib-gui` | `44e67ca`, 2026-08-08 | Go with Raylib bindings; newer and more complete. Adds font sizing changes and A–G chord selection/playback after the fork. |
| `feat/c-raylib-ui` | `9835b9d`, 2026-08-06 | Early native C port (`d342cb3`) plus native-build CI. Keeps the older Go source but changes the default build to C. |

The Go desktop implementation has synthesized audio, chord-derived waveform and
spectrum plots, ten themes, compact/full-neck views, tab labels, and help.
The C prototype implements catalog generation and navigation, but its waveform
and spectrum are hard-coded examples. It has no audio playback; its fretboard
only draws fretted notes at frets 1–4, omitting higher positions. Several controls
advertised by its README and window are not implemented. Merging the C branch
would replace the Go build/CI without delivering feature parity.

The latest Go shortcuts also bind F to both F-chord playback and full-neck
view. This existing conflict remains a follow-up issue.

## Preserving the C experiment

The local C branch is archived as the annotated tag `archive/c-raylib-ui`
pointing to `9835b9d`. Inspect it with:

```sh
git show archive/c-raylib-ui:src/main.c
```

Resume it without changing either maintained branch:

```sh
git switch -c experiment/c-raylib archive/c-raylib-ui
```

The local `feature/raylib-gui` branch was renamed to `raylib`. Remote branches
are unchanged; this organization does not publish or delete remote refs.

## Code boundaries

On `master`, `cmd/alt-tab` launches `internal/tui`, backed by `internal/chords`.
On `raylib`, `cmd/alt-tab-raylib` launches `internal/rayui`; `internal/app`,
`internal/chords`, and `internal/signal` supply shared behavior. Keep graphical
input, rendering, and audio in `internal/rayui`. Port shared fixes deliberately
between branches: the TUI branch predates the shared app/signal extraction.

## Verification

Run `make check` on each maintained branch. The Raylib branch additionally
requires cgo and the platform graphics development libraries. After Go 1.27.3 was installed, both branches passed formatting, unit tests,
race checks, vet, and builds in temporary copies on 2026-10-30:

- `master`: `make check build`
- `raylib`: `make check build build-raylib`

The Raylib Go module supplies its native source; the desktop build succeeded
without a standalone `raylib.pc`. Graphical interaction and audio playback
were not manually exercised.
