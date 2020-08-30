package pctl

import (
	"math"
	"testing"
	"time"
)

func TestLowPassFilterAsymptotic(t *testing.T) {
	lpf := NewLPF(1e6)
	// 1Mhz low-pass filter has corner frequency of a MHz, corresponds to
	// microseconds to reach nearly steady state
	process := lpf.Update(0.) // seed the filter
	target := 1.
	for i := 0; i < 5; i++ {
		process = lpf.Update(target)
		time.Sleep(time.Millisecond) // test over 5ms ensures that it shall converge
	}
	err := 1 - process
	if math.Abs(err) > 1e-5 {
		t.Errorf("process of %f has error of %f, expected to converge to 1 along step", process, err)
	}
}

func TestHighPassRejectsOscillation(t *testing.T) {
	hpf := NewHPF(1e6)
	// cutoff of 1Mhz means that inputs at ms should be rejected
	process := hpf.Update(0.)
	target := 1.
	time.Sleep(time.Millisecond) // early sleep to ensure we don't update immediately
	for i := 0; i < 5; i++ {
		process = hpf.Update(target)
		time.Sleep(time.Millisecond)
	}
	err := 0 - process
	if math.Abs(err) > 1e-5 {
		t.Errorf("process of %f has error of %f, expected to resist motion in ms domain", process, err)
	}
}
