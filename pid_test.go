package pctl

import (
	"fmt"
	"testing"
	"time"
)

var dummy float64

func BenchmarkPIDLoop(b *testing.B) {
	ctl := PID{}
	for n := 0; n < b.N; n++ {
		dummy = ctl.Update(0)
	}
}

// ExampleControlLoop shows an example of a PI controller performing a step change
// in output through a lowpass filter
func ExamplePID() {
	pid := PID{
		P:     1,
		I:     0,
		Setpt: 100,
	}
	lpf := NewLPF(200.) // 200Hz filter
	systemOutput := 0.  // start at "rest"
	var cmdhistory []float64
	var outhistory []float64
	tick := time.NewTicker(time.Millisecond) // run loop at 1kHz
	const updates = 50
	ticks := 0
	for {
		<-tick.C // do nothing with the time
		command := pid.Update(systemOutput)
		systemOutput = lpf.Update(command)
		cmdhistory = append(cmdhistory, command)
		outhistory = append(outhistory, systemOutput)
		ticks++
		if ticks == updates {
			tick.Stop()
			break
		}
	}
	fmt.Println("PID can at 1kHz for 50 updates, commands were:")
	fmt.Println(cmdhistory)
	fmt.Println("outputs were:")
	fmt.Println(outhistory)
}
