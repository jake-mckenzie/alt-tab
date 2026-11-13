package audio

import "github.com/jfreymuth/pulse"

// Player sends Go-generated samples to PipeWire's PulseAudio service or PulseAudio.
// Future platform backends implement Open, Set and Close with the same synth.
type Player struct {
	synth
	client *pulse.Client
	stream *pulse.PlaybackStream
}

// Open connects to PulseAudio or PipeWire and starts a silent playback stream.
func Open() (*Player, error) {
	client, err := pulse.NewClient(pulse.ClientApplicationName("Alt-Tab"))
	if err != nil {
		return nil, err
	}
	p := &Player{client: client}
	p.stream, err = client.NewPlayback(pulse.Float32Reader(p.read),
		pulse.PlaybackSampleRate(sampleRate), pulse.PlaybackLatency(0.025))
	if err != nil {
		client.Close()
		return nil, err
	}
	p.stream.Start()
	return p, nil
}

// Close releases the output connection after the terminal program stops.
func (p *Player) Close() {
	p.stream.Close()
	p.client.Close()
}
