package pctl

import "testing"

var dummy float64

func BenchmarkPIDLoop(b *testing.B) {
	ctl := PID{}
	for n := 0; n < b.N; n++ {
		dummy = ctl.Update(0)
	}
}
