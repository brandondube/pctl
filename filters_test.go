package pctl

import (
	"math"
	"testing"
)

func TestLowPassFilterAsymptotic(t *testing.T) {
	lpf := NewLPF(1e6, 1e-3)
	// 1Mhz low-pass filter has corner frequency of a MHz, corresponds to
	// microseconds to reach nearly steady state
	process := lpf.Update(0) // seed the filter
	target := 1.
	for i := 0; i < 5; i++ {
		process = lpf.Update(target)
	}
	err := 1 - process
	if math.Abs(err) > 1e-5 {
		t.Errorf("process of %f has error of %f, expected to converge to 1 along step", process, err)
	}
}

func TestHighPassRejectsOscillation(t *testing.T) {
	hpf := NewHPF(1e6, 1e-3)
	// cutoff of 1Mhz means that inputs at ms should be rejected
	process := hpf.Update(0)
	target := 1.
	for i := 0; i < 5; i++ {
		process = hpf.Update(target)
	}
	err := 0 - process
	if math.Abs(err) > 1e-5 {
		t.Errorf("process of %f has error of %f, expected to resist motion in ms domain", process, err)
	}
}

func TestBiquadFilterAsymptotic(t *testing.T) {
	// Bq at sample rate 1kHz, 250Hz corner, Q=sqrt(2)/2, -6dB gain
	a0 := 0.2928920553392428
	a1 := 0.5857841106784856
	a2 := a0
	b1 := -1.3007020142696517e-16
	b2 := 0.17156822135697122
	bq := NewBiquad(a0, a1, a2, b1, b2)
	process := bq.Update(0)
	target := 1.
	for i := 0; i < 100; i++ {
		process = bq.Update(target)
		t.Logf("%f %f\n", process, target)
	}
	err := target - process
	if math.Abs(err) > 1e-5 {
		t.Errorf("process of %f has error of %f, expected to converge to target=1", process, err)
	}
}
