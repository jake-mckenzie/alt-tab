//go:build !linux

package audio

import "errors"

// Player reserves the same backend contract on platforms awaiting audio support.
type Player struct{ synth }

// Open reports that no native audio backend is available on this platform.
func Open() (*Player, error) {
	return nil, errors.New("audio currently requires Linux with PipeWire/PulseAudio")
}

// Close satisfies the output contract when no native resources exist.
func (*Player) Close() {}
