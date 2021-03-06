# pctl

pctl, "process control" is a package for process control in Go.  This applies to industrial processes, not the computer science term.

It contains a PID controller, as well as several building blocks for filters:
- single-pole low and highpass
- Biquad
- State-space

Biquads and the single pole filters can be user-cascaded to make second-order section filters.

The package declares the top-level `Cascade` function, which takes a sequence of interfaces that are met by all types in the package to facilitate SOS and other "fluent" designs.

Its types are not concurrent safe, and use double precision, which is low cost on most software platforms.  Tinygo may perform relatively worse, although it should not matter much.  Most methods in this package takes approximately 20 clocks to execute, which corresponds to MHz rep rates on an ordinary CPU, even a Raspberry Pi.


Filter design is outside the scope of this package, which exists to assemble control systems in Go.


## Usage

### Biquad filter on measurement with PID controller

```go
// Biquad, 1k sample rate, 50Hz corner freq, -6dB/octave gain
a0 := 0.0201
a1 := 0.0402
a2 := a0
b1 := -1.5610
b2 := 0.6414
inputFilter := pctl.NewBiquad(a0, a1, a2, b1, b2)
controller := pctl.PID{P: 1, I: 0.5, Setpt: 50, DT: 1e-3}
for {
    input := getInput()
    process := pctl.Cascade(input, inputFilter, controller)
    applyControl(process)
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
    process := pctl.Cascade(input, setpt, FBFilter)
    applyControl(process)
}
```

### Single-pole filters and PID discrete time parameters

The implementation of the PID controller and single-pole filters, unlike the Biquad and state-space filters, are not relative to the sample rate, but relative to a concept of wall time.  In other words, a LPF is designed with a corner frequency and a sampling period, the two of which are distinct.  Biquads and state-space filters are truly relative to their sample rate, not an idea of wall time.

The concept of sample rate is introduced through the DT field.  It is a usage error to forget it and leave it as zero.  The types do not check this on `Update` as it reduces performance.

```go
// wrong
ctl := &pctl.PID{P: 1, Setpt: 2}

// right
ctl := &pctl.PID{P: 1, Setpt: 2, DT: 1e-3}
```

The `NewLPF` and `NewHPF` methods make it difficult to make this error using single-pole filters.


### Shaped controller response, control setpoint change stability


The previous examples lack prefilters on the setpoint, so the system can be destabilized by large setpoint changes.  A prefilter can be added that operates on `*setpt` to remedy this.

Opening or closing the control loop independent of measurement is also not possible.  The latter can be achieved by simply adding one line:

```go
for {
    // ...
    if controlLoopClosed {
        applyControl(process)
    }
}
```

Manipulating of this variable is outside the scope of pctl.  It could be e.g. a struct member, or simply a pointer to a bool that is dereferenced at the if.  The "size" of the solution can scale with the "size" of the processor and problem.


## Performance

See `pctl_test.go` for a benchmark suite.  Most basic operations require approximately 20 clocks.  A two-state statespace filter is approximately 5x as expensive, and computational complexity increases with the number of states.

## Design

Several designs have been iterated in this repository.  An early design used channels to communicate, which took about 500ns per update.  This was less composable than methods/functions.

An intermediate design maintained clocks inside each control element.  This was less performant, but more importantly could not be used in a simulation capacity running at any speed other than real time.  Explicitly including dT (fielded as DT) in the structs allows these controllers to be used in simulation studies as well.  The nearly 10x increase in performance and better friendliness to tinygo platforms are also nice benefits.

The current design has been released as v1 (guaranteed stable) and is unlikely to change for marginal improvements in favor of API stability.

## Expansion

This library is dependency-free and easily portable to tiny platforms, even if a float32 type-change would be required (this is as simply as ctrl+F).  Future additions shall not disturb that property.  LQR/LQG, Kalman filtering, etc, may be implemented here if the the implementations do not require a dependency on e.g. Gonum.
