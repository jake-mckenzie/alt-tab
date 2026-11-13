package audio

import (
	"math"
	"sync"
	"testing"

	"github.com/jake-mckenzie/alt-tab/internal/chords"
	"github.com/jake-mckenzie/alt-tab/internal/signal"
)

func TestSynthMatchesRaylibMixAndFadesToSilence(t *testing.T) {
	v, _ := chords.NewCatalog().Load("A", 1)
	var s synth
	s.Set(v, true)
	buffer := make([]float32, 4096)
	s.read(buffer)
	notes := signal.Notes(v)
	gain := 0.0
	for i, got := range buffer {
		gain += (0.28*0.65 - gain) * 0.0025
		want := signal.CompositeSample(float64(i)/sampleRate, notes) * gain
		if math.Abs(float64(got)-want) > 1e-7 {
			t.Fatalf("sample %d = %g, want %g", i, got, want)
		}
	}
	s.Set(v, false)
	for range 3 {
		s.read(buffer)
	}
	s.read(buffer)
	for _, got := range buffer {
		if got != 0 {
			t.Fatal("stopped synth did not settle to silence")
		}
	}
}

func TestSynthPreservesPhaseAcrossBuffers(t *testing.T) {
	v, _ := chords.NewCatalog().Load("C", 2)
	var whole, split synth
	whole.Set(v, true)
	split.Set(v, true)
	want, got := make([]float32, 2048), make([]float32, 2048)
	whole.read(want)
	split.read(got[:777])
	split.Set(v, true)
	split.read(got[777:])
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("phase discontinuity at %d", i)
		}
	}
}

func TestSynthConcurrentSelection(t *testing.T) {
	var s synth
	v, _ := chords.NewCatalog().Load("G", 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			s.Set(v, i%2 == 0)
		}
	}()
	buffer := make([]float32, 128)
	for range 100 {
		s.read(buffer)
	}
	wg.Wait()
}
