package pctl

import "testing"

func BenchmarkPIDLoop(b *testing.B) {
	ctl := PID{}
	for n := 0; n < b.N; n++ {
		ctl.Update(0)
	}
}

func BenchmarkLPF(b *testing.B) {
	lpf := LPF{}
	for n := 0; n < b.N; n++ {
		lpf.Update(0)
	}
}

func BenchmarkHPF(b *testing.B) {
	hpf := HPF{}
	for n := 0; n < b.N; n++ {
		hpf.Update(0)
	}
}

func BenchmarkBiquad(b *testing.B) {
	bq := Biquad{}
	for n := 0; n < b.N; n++ {
		bq.Update(0)
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
		filt.Update(0)
	}
}

func BenchmarkSetpoint(b *testing.B) {
	s := Setpoint(0)
	for n := 0; n < b.N; n++ {
		s.Update(0)
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
