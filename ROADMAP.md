# Roadmap

This document outlines plans and conditions for changes to pctl

## Invariants

These are design elements or decisions for the library which I consider inviolable.  Without exceptionally good justification, changes which break an invariant will be rejected:

### pctl must be reasonably tinygo friendly.  The compiled code size should stay reasonably small

  It will necessarily grow over time, but "large microcontroller" program memory sizes do, too.  As it stands in mid 2022 with Go 1.18.2, an empty `func main() { return }` program is 8,128 bytes when built for GOOS=linux GOARCH=amd64 with tinygo, after running strip.  Changing that program to:

  ```go
  package main

import "github.com/brandondube/pctl"

func main() {
	s := new(pctl.Setpoint)
	*s = 42
	f := pctl.NewBiquadLowShelf(1000, 100, 0.707, 6)
	out := pctl.Cascade(12, s, f)
	_ = out
}
```

  raises the binary size to 9,040 bytes (+1K).  So, it is possible to write programs with pctl even for extremely low-end hardware, such as the ATmega32U.  Ideally pctl remain small enough to be suitable for this class of hardware.  It must absolutely stay small enough to support boards like the Arduino Mega with 128K of program memory.

### pctl must have a reasonably deterministic runtime performance

The underlying requirement of real-time software is that doing the same thing many times must take about the same amount of time on each run.  That is, if a program completes one loop iteration in 100 microseconds, it should always take about 100 microseconds, and not have a random iteration require 10 milliseconds.  The predominant reason for large latency spikes _that are within the control of the programmer_ is a substantial operation by the garbage collector.  Provided that the "real time" functions such as the `Updater` interface in pctl do not allocate, then pctl itself will not cause pressure on the GC, and programs will remain suitable for real-time use.

### pctl must be reasonably performant

Modern desktop or even laptop class processors are fast enough to run 100kHz class controls systems, limited by (for example) the ~10 microsecond PCI-express bus rep rate.  Unless using extremely large filters, pctl is three orders of magnitude faster than this as of version 1.3.  However, when a 4 Gigahert processor that can perform multiple floating point operations per clock is replaced with a microcontroller that has to simulate floating point in software, a tax of 10x on top of the clock speed might be felt.  E.g., a 32MHz ARM processor without floating point hardware might take 100 clocks, 3.125 microseconds, to perform tasks like the Biquad update in pctl.  To the extent possible, we want to keep the algorithms highly optimized so that pctl can be used at even kHz on commodity single board computers and microcontrollers.

### Real-time only

pctl is not meant to be an equivalent to, say, the filter design toolbox in matlab or python-control.  This is a package for building real-time systems.  Reasonably lightweight filter synthesis is in scope, such as the various Biquad constructors.  Design of complex multi-band FIR filters that would require 100s of taps is not in scope.  Actuator transfer function inversion routines or latency estimators and compensators are not in scope.

## Future Work

## Alternate Number Formats

There are not 32-bit floating point implementations of the types, for example.  Nor are there fixed point implementations (which would favor a different Biquad calculation method, as well).  These would be welcome additions.  Presupposing that tinygo adopts generics, it is OK to polymorphize over the signed integer types using them.


## Adaptive/Predictive Control

A subset of predictive controllers would also be welcome, provided that they satisfy the invariants outlined above.  The State-Space filter can be used as an example for how to implement (relatively) simple operations without allocation.  The [brandondube/linalg](https://github.com/brandondube/linalg) package is substantially smaller than Gonum/mat and implements more functionality than those copied/vendored into pctl.  There are not, however, any matrix inversion or similar routines.  Often, predictive controllers can be designed in such a way that there are no inversions or other difficult operations required in real-time, under some given assumptions e.g. that the system is in steady-state.  Pure feedback can often be used to initially stabilize a system and then a secondary predictive controller enabled only thereafter.  If predictor types were added to pctl which have these sorts of "boundary conditions" for their use, they should be extremely clearly stated.


## Control utilities

Today pctl only provides controllers and filters, but does not provide other routines which may be required to construct a controls system, for example a rotation matrix, or simplified inverse kinematic models such as those which compute commands for three actuators arranged in a tripod from a two dimensional control offset.  Quaternions, etc, are more examples.

These are all welcome.  Three should be a debate or performance test between the use of arrays and slices in those cases.  While the State-Space filter uses slices due to the unknowable number of states to be used, these transformations will generally be 2x2, 2x3, 3x2, and so on.  Arrays may be faster in those cases, due to the compiler doing the bounds checking and not the runtime.  See [brandondube/goray/linalg.go](https://github.com/brandondube/goray/blob/main/linalg.go) for basic operations of this nature.  It should be benchmarked first, since it's likely that the ergonomics of initialization will be worse with arrays.

----

More may be added to this list over time.
