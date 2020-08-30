package pctl

import (
	"math"
	"testing"
)

var dummy float64

func BenchmarkPIDLoop(b *testing.B) {
	ctl := PID{}
	for n := 0; n < b.N; n++ {
		dummy = ctl.Update(0)
	}
}

func TestPIDConverges(t *testing.T) {
	ctl := PID{P: 1}
	// loss free process
	state := ctl.Update(0)
	ctl.Setpt = 1.
	var cmd float64
	for i := 0; i < 5; i++ {
		cmd = ctl.Update(state)
		state += cmd
	}
	cmdErr := math.Abs(cmd - 0)
	stateErr := math.Abs(state - ctl.Setpt)
	if cmdErr > 1e-5 {
		t.Errorf("expected controller to go to zero output, got: %f", cmdErr)
	}
	if stateErr > 1e-5 {
		t.Errorf("expected state to converge to setpoint, got %f with error of %f", state, stateErr)
	}
}
