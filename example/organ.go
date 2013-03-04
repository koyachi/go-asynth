package main

import (
	"../../go-asynth"
	"math"
)

func main() {
	s := asynth.New(func(note asynth.Note, t float64) float64 {
		freq := 440.0 * math.Pow(2, float64(note.Key-49)/12)
		x := math.Sin(2.0 * math.Pi * t * freq)
		y := math.Sin(2.0 * math.Pi * t * freq * 2)
		z := math.Sin(2.0 * math.Pi * t * freq / 2)
		return x*0.6 + y*0.2 + z*0.2
	})
	s.Play(nil)
}
