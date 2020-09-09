# pctl

pctl, "process control" is a package for process control in Go.  This applies to industrial processes, not the computer science term.

It contains a PID controller as well as low and high pass filters.  These types are optimized for usability, not for
maximum performance.  They are not concurrent safe, no guarantees are provided about state or safety during `Update` invocations.  The PID controller is capable of running at a few MHz on a laptop.  This should translate into 1s to 10s of kHz on embedded hardware with tinygo.

## Usage

Here we show using this package to tune a PID loop for a simplistic system ("plant").  The plant loses a fixed 10 units per interval of time.  It has no input loss, and gains exactly the units we put into it.  We desire the plant to hold 100 units, then later decrease to 25.  We will use [asciigraph](https://github.com/guptarohit/asciigraph) to show the process over time.  We will use clamping and a 100Hz low pass filter behind the controller to regulate the output.  Pay special attention to the fact that the intergal error is reset when we move the setpoint.  See the members of `pctl.PID` for an anti-windup variable.  We reset when changing the setpoint to drain the integrator.

```go
package main

import (
	"fmt"

	"github.com/brandondube/pctl"

	"github.com/guptarohit/asciigraph"
)

// ExamplePID shows an example of a PI controller performing a step change
// in output through a lowpass filter
func ExamplePID(P, I, D float64) {
	pid := pctl.PID{
		P:     P,
		I:     I,
		D:     D,
		DT:    1/1000,
		Setpt: 100,
	}
	lpf := pctl.NewLPF(100., 1/1000) // 200Hz filter
	const (
		updates = 100
		cmdMax  = 200
		cmdMin  = 0
	)
	var (
		cmdhistory   []float64
		outhistory   []float64
		systemOutput = 0.
	)
	for ticks := 0; ticks < updates; ticks++ {
		command := pid.Update(systemOutput)
		if command > cmdMax {
			command = cmdMax
		} else if command < cmdMin {
			command = cmdMin
		}
		systemOutput = (lpf.Update(systemOutput+command) - 10)
		cmdhistory = append(cmdhistory, command)
		outhistory = append(outhistory, systemOutput)
		if ticks == updates/2 {
			pid.Setpt = 25
			pid.IntegralReset()
		}
	}
	fmt.Println("commands:")
	graph := asciigraph.Plot(cmdhistory, asciigraph.Height(10))
	fmt.Println(graph)
	fmt.Println("output:")
	graph = asciigraph.Plot(outhistory, asciigraph.Height(10))
	fmt.Println(graph)
}

func main() {
	ExamplePID(1, 0, 0)
}

```

This shows the system with a simple P controller with unity gain.  The output is:
```
go run main.go
commands:
 110 ┼╭╮
  99 ┤││
  88 ┤││
  77 ┤││
  66 ┤│╰╮
  55 ┤│ │
  44 ┤│ │
  33 ┤│ ╰╮
  22 ┤│  ╰─╮
  11 ┤│    ╰───────────────────────────────────────────╮              ╭──────────────────────────────────
   0 ┼╯                                                ╰──────────────╯
output:
  90.00 ┤    ╭────────────────────────────────────────────╮
  80.00 ┼  ╭─╯                                            ╰──╮
  70.00 ┤  │                                                 ╰─╮
  60.00 ┤ ╭╯                                                   ╰─╮
  50.00 ┤ │                                                      ╰─╮
  40.00 ┤ │                                                        ╰─╮
  30.00 ┤╭╯                                                          ╰─╮
  20.00 ┤│                                                             ╰────────────────────────────────────
  10.00 ┤│
  -0.00 ┤│
 -10.00 ┼╯

```

We can see the system has no overshoot, but does not reach the setpoint of 100.  If we double the gain, the droop is halved, but the commands become pulse train like:
```
# ExamplePID(2, 0, 0)
go run main.go
commands:
 200 ┼╭╮
 180 ┤││
 160 ┤││
 140 ┤││
 120 ┤││
 100 ┤││
  80 ┤││
  60 ┤││
  40 ┤││
  20 ┤││╭╮╭─╮ ╭╮╭╮╭╮ ╭─────────────────────────────────╮              ╭╮╭╮ ╭─╮╭╮╭╮╭──╮╭──────────────────
   0 ┼╯╰╯╰╯ ╰─╯╰╯╰╯╰─╯                                 ╰──────────────╯╰╯╰─╯ ╰╯╰╯╰╯  ╰╯
output:
  95.75 ┤╭─────────────────────────────────────────────────╮
  85.18 ┼│                                                 ╰─╮
  74.60 ┤│                                                   ╰─╮
  64.03 ┤│                                                     ╰─╮
  53.45 ┤│                                                       ╰─╮
  42.88 ┤│                                                         ╰─╮
  32.30 ┤│                                                           ╰─╮
  21.73 ┤│                                                             ╰────────────────────────────────────
  11.15 ┤│
   0.58 ┤│
 -10.00 ┼╯
```

By adding integral gain, we can remove the droop, at the expense of a small undershoot when moving the setpoint:
```
# ExamplePID(4, 20, 0)
go run main.go
commands:
 200 ┼╭╮
 180 ┤││
 160 ┤││
 140 ┤││
 120 ┤││
 100 ┤││
  80 ┤││
  60 ┤││
  40 ┤││
  20 ┤│╰╮╭╮╭╮╭╮╭╮╭╮╭╮╭╮╭╮╭╮╭╮╭╮╭╮╭╮╭╮╭╮╭╮╭╮╭╮╭╮╭╮╭╮╭╮╭─╮                 ╭╮╭╮╭╮╭╮╭╮╭╮╭╮╭╮╭╮╭─╮╭─────────╮
   0 ┼╯ ╰╯╰╯╰╯╰╯╰╯╰╯╰╯╰╯╰╯╰╯╰╯╰╯╰╯╰╯╰╯╰╯╰╯╰╯╰╯╰╯╰╯╰╯╰╯ ╰─────────────────╯╰╯╰╯╰╯╰╯╰╯╰╯╰╯╰╯╰╯ ╰╯         ╰
output:
 103 ┤ ╭───────────────────────────────────────────────╮
  92 ┼╭╯                                               ╰──╮
  80 ┤│                                                   ╰─╮
  69 ┤│                                                     ╰──╮
  58 ┤│                                                        ╰─╮
  47 ┤│                                                          ╰─╮
  35 ┤│                                                            ╰──╮
  24 ┤│                                                               ╰─╮╭───────────────────────────────
  13 ┤│                                                                 ╰╯
   1 ┤│
 -10 ┼╯

```

Note that because these are discrete time control elements with their own clocks, your results may differ based on
the relatively noisy/imprecise behavior of `time.Ticker`.  We can see the process control works well, but the controller output is flickering on and off.  Further tuning is left to the user.


## Performance

To demonstrate that this controller is capable of running at MHz, we show a benchmark performed on a windows 10 computer with an i7-9700k processor:
```
Running tool: C:\Go\bin\go.exe test -benchmem -run=^$ github.com/brandondube/pctl -bench ^(BenchmarkPIDLoop)$

goos: windows
goarch: amd64
pkg: github.com/brandondube/pctl
BenchmarkPIDLoop-8   	564444776	         2.12 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/brandondube/pctl	1.538s
```

The reciprocal of 14.9 nanoseconds is ~470MHz.  At 4.5GHz, this is is ~100 clocks.  Assuming your device has double precision FPUs, you can assume the controller can run at approx 1/100th the clock speed.

## Design

Several designs have been iterated in this repository.  An early design used channels to communicate, which took about 500ns per update.  This was less composable than methods/functions.

An intermediate design maintained clocks inside each control element.  This was less performant, but more importantly could not be used in a simulation capacity running at any speed other than real time.  Explicitly including dT (fielded as DT) in the structs allows these controllers to be used in simulation studies as well.  The nearly 10x increase in performance and better friendliness to tinygo platforms are also nice benefits.
