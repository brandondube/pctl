# pctl

pctl, "process control" is a package for industrial control in Go.

It contains an implementation of the classic PID controller with integral
anti-windup, as well as many filter types that can be used for loop shaping:
- single pole low pass
- single pole high pass
- Biquads
- State-Space filters with an arbitrary number of states
- FIR filters with an arbitrary number of taps

The package declares the top-level `Cascade` function, which takes a sequence of
interfaces that are met by all types in the package to facilitate SOS and other
"fluent" designs.  Use of Cascade will be somewhat slower or less efficient than
manually writing a chain of function calls due to the virtualization implied by
interfaces.

Its types are not concurrent safe, and use double precision, which is low cost
on most software platforms.  Tinygo may perform relatively worse, although it
should not matter much.  The implementations of each type in this repository are
relatively optimized, easily able to function at up to MHz on even a raspberry
pi.

For Biquads, design methods are included to synthesize common filter types from
corner frequencies, etc, in applications where detailed analysis of the transfer
functions or plant response are not required.

## Usage

### Biquad filter on measurement with PID controller

```go
// Biquad, 1k sample rate, 50Hz corner freq, maximally flat in band
// 6 = gain; unused for LPF; see NewBiquad interface
// or bring your own a0, a1, a2, b1, b2 coefs
inputFilter := pctl.NewBiquadLowPass(1000, 50, math.Sqrt(2), 6)
controller := pctl.PID{P: 1, I: 0.5, Setpt: 50, DT: 1e-3}
for {
    input := getInput()
    controlCommand := pctl.Cascade(input, inputFilter, controller)
    applyControl(controlCommand)
}
```

### State-Space filtering the error signal for control shaping

```go
// State-space second order lowpass filter,
// 900Hz sample rate, 2Hz corner freq, -6dB/octave
A := [][]float64{
    {2, -1},
    {1, 0},
}
B := []float64{5e-5, 0}
C := []float64{4, 0.02}
D := 5e-5
setpt := pctl.Setpoint(50)
// FB = feedback
FBFilter := pctl.NewStateSpaceFilter(A, B, C, D, nil)
for {
    input := getInput()
    controlCommand := pctl.Cascade(input, setpt, FBFilter)
    applyControl(controlCommand)
}
```

### Shaped controller response, control setpoint change stability


The previous examples lack prefilters on the setpoint, so the system can be
destabilized by large setpoint changes.  A prefilter can be added that operates
on `*setpt` to remedy this.

Opening or closing the control loop independent of measurement is also not
possible.  The latter can be achieved by simply adding one line:

```go
for {
    // ...
    if controlLoopClosed {
        applyControl(process)
    }
}
```

Manipulating of this variable is outside the scope of pctl.  It could be e.g. a
struct member, or simply a pointer to a bool that is dereferenced at the if.
The "size" of the solution can scale with the "size" of the processor and
problem.


## Performance

See `pctl_test.go` for a benchmark suite.  The FIR filter in the benchmark has
32 taps.

### Mac M1 Pro

M1 Pro Boost frequency = 3.2GHz; 1 clock ~=0.3125 ns.

```sh
name           time/op
PIDLoop-10     3.50ns ± 1%
LPF-10         4.52ns ± 2%
HPF-10         4.49ns ± 2%
Biquad-10      4.89ns ± 1%
StateSpace-10  12.5ns ± 3%
Setpoint-10    0.32ns ± 1%
FIRFilter-10   11.8ns ± 1%
```
A reasonable average is the Biquad filter, 15.6 clocks.

### Intel i7-9700k

This CPU boosts to 4.6GHz during the benchmark; 1 clock ~=0.217 ns.
```sh
name          time/op
PIDLoop-8     1.99ns ± 2%
LPF-8         3.74ns ± 1%
HPF-8         2.80ns ± 1%
Biquad-8      3.65ns ± 1%
StateSpace-8  9.72ns ± 1%
Setpoint-8    0.21ns ± 3%
FIRFilter-8   8.66ns ± 2%
```

The Biquad filter takes 16.8 clocks.  Broadly comparable to the ARM64 M1.

### AMD 7950X (Windows)

This CPU boosts to 5.3GHz during the benchmark; 1 clock ~= 0.189 ns.  cTDP 105w
eco mode is enabled.
```sh
name            time/op
PIDLoop-8       3.444n ± 0%
LPF-8           4.317n ± 0%
HPF-8           4.312n ± 3%
Biquad-8        4.694n ± 3%
StateSpace-8    10.90n ± 1%
Setpoint-8       0.29n ± 1%
FIRFilter-8     10.88n ± 0%
```

Despite having a considerably higher clockspeed, this CPU takes more time to
perform the functions within pctl.

### AMD 7950X (WSL)


```sh
name             time/op
PIDLoop-32       2.168n ± 1%
LPF-32           3.190n ± 0%
HPF-32           2.654n ± 0%
Biquad-32        3.191n ± 0%
StateSpace-32    7.301n ± 1%
Setpoint-32     0.1801n ± 1%
FIRFilter-32     6.943n ± 2%
```

Performance is ~50% higher in Windows subsystem for Linux / Ubuntu.

## Design

Several designs have been iterated in this repository.  An early design used
channels to communicate, which took about 500ns per update.  This was less
composable than methods/functions.

An intermediate design maintained clocks inside each control element.  This was
less performant, but more importantly could not be used in a simulation capacity
running at any speed other than real time.  Explicitly including dT (fielded as
DT) in the structs allows these controllers to be used in simulation studies as
well.  The nearly 10x increase in performance and better friendliness to tinygo
platforms are also nice benefits.

The current design has been released as v1 (guaranteed stable) and is unlikely
to change for marginal improvements in favor of API stability.

## Expansion

This library is dependency-free outside stdlib/math and easily portable to tiny
platforms, even if a float32 type-change would be required (this is as simply as
ctrl+F).  Future additions shall not disturb that property.  LQR/LQG, Kalman
filtering, etc, may be implemented here if the the implementations do not
require a dependency on e.g. Gonum.
