package pctl

import "testing"

var dummy float64

func BenchmarkPIDLoopNonLocking(b *testing.B) {
	ctl := PIDCtl{}
	in := make(chan float64)
	out := make(chan float64)
	go ctl.Start(in, out)
	for n := 0; n < b.N; n++ {
		in <- 0
		dummy = <-out
	}
}

func BenchmarkPIDLoopLocking(b *testing.B) {
	ctl := PIDCtl{locking: true}
	in := make(chan float64)
	out := make(chan float64)
	go ctl.Start(in, out)
	for n := 0; n < b.N; n++ {
		in <- 0
		dummy = <-out
	}
}
