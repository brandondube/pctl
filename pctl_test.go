package pctl

import (
	"math/rand"
	"testing"
)

func BenchmarkPIDLoop(b *testing.B) {
	ctl := PID{}
	for n := 0; n < b.N; n++ {
		ctl.Update(3.14)
	}
}

func BenchmarkLPF(b *testing.B) {
	lpf := LPF{}
	for n := 0; n < b.N; n++ {
		lpf.Update(3.14)
	}
}

func BenchmarkHPF(b *testing.B) {
	hpf := HPF{}
	for n := 0; n < b.N; n++ {
		hpf.Update(3.14)
	}
}

func BenchmarkBiquad(b *testing.B) {
	bq := Biquad{}
	for n := 0; n < b.N; n++ {
		bq.Update(3.14)
	}
}

func BenchmarkStateSpace(b *testing.B) {
	A := [][]float64{
		{2, -1},
		{1, 0},
	}
	B := []float64{5e-5, 0}
	C := []float64{4, 0.02}
	D := 5e-5
	filt := NewStateSpaceFilter(A, B, C, D, nil)
	for n := 0; n < b.N; n++ {
		filt.Update(3.14)
	}
}

func BenchmarkSetpoint(b *testing.B) {
	s := Setpoint(0)
	for n := 0; n < b.N; n++ {
		s.Update(3.14)
	}
}

func BenchmarkFIRFilter(b *testing.B) {
	const filterSize = 32
	coefs := make([]float64, filterSize)
	for i := 0; i < filterSize; i++ {
		coefs[i] = rand.Float64()
	}
	f := NewFIRFilter(coefs)
	for n := 0; n < b.N; n++ {
		f.Update(3.14)
	}
}

func TestSetpointCorrect(t *testing.T) {
	s := Setpoint(0)
	meas := 2. // 2 is exactly representable in fp
	err := s.Update(meas)
	if err != 2 {
		t.Error("Setpoint improperly computed error")
	}
}
