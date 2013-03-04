package asynth

import (
	"fmt"
	"github.com/koyachi/go-baudio"
	"github.com/youpy/go-coremidi"
	"math"
)

type Note struct {
	Key     byte
	Start   float64
	Elapsed float64
	Up      float64
	Down    float64
}

type Asynth struct {
	notes []*Note
	now   float64
	b     *baudio.B
}

func New(fn func(note Note, t float64) float64) *Asynth {
	a := &Asynth{
		now: 0.0,
	}
	a.initCoreMidi()
	a.b = a.initBaudio(fn)
	return a
}

func (a *Asynth) initCoreMidi() {
	client, err := coremidi.NewClient("go-asynth client")
	if err != nil {
		panic(err)
	}

	sources, err := coremidi.AllSources()
	if err != nil {
		panic(err)
	}

	port, err := coremidi.NewInputPort(client, "test", func(source coremidi.Source, value []byte) {
		fmt.Printf("source: %v, manufacturer: %v, value: %v\n", source.Name(), source.Manufacturer(), value)
		note := value[1]
		if value[2] == 0 {
			i := 0
			for i = 0; i < len(a.notes) && a.notes[i].Key != note; i++ {
			}
			if len(a.notes) >= i {
				a.notes[i].Up = a.now
			}
		} else {
			a.notes = append(a.notes, &Note{Key: note, Down: a.now})
		}
		return
	})
	if err != nil {
		panic(err)
	}

	for _, source := range sources {
		func(source coremidi.Source) {
			port.Connect(source)
		}(source)
	}
}

func (a *Asynth) initBaudio(fn func(Note, float64) float64) *baudio.B {
	opt := baudio.NewBOptions()
	opt.Size = 16
	opt.Rate = 44000
	b := baudio.New(opt, nil)
	b.Push(func(t float64, i int) float64 {
		a.now = t
		sum := 0.0
		for i := 0; i < len(a.notes); i++ {
			note := a.notes[i]
			if note.Start == 0.0 {
				note.Start = t
			}
			elapsed := t - note.Start
			note.Elapsed = elapsed
			if note.Up != 0.0 && elapsed >= note.Up-note.Down {
				a.notes = append(a.notes[:i], a.notes[i+1:]...)
				i--
				continue
			}
			sum += fn(*note, t)
		}
		if len(a.notes) > 0 {
			return sum / math.Sqrt(float64(len(a.notes)))
		}
		return 0
	})
	return b
}

func (a *Asynth) Play() {
	// TODO opts for b.play
	a.b.Play(nil)
}
