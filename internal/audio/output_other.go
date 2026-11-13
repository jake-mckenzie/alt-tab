//go:build !linux

package audio

import "errors"

// Player reserves the same backend contract on platforms awaiting audio support.
type Player struct{ synth }

func Open() (*Player, error) {
	return nil, errors.New("audio currently requires Linux with PipeWire/PulseAudio")
}
func (*Player) Close() {}
