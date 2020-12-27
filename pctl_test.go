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
