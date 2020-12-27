# pctl

pctl, "process control" is a package for process control in Go.  This applies to industrial processes, not the computer science term.

It contains a PID controller as well as simple tools for filtering measurement data with first order low and high-pass filters, or Biquads.

Its types are not concurrent safe, and use double precision, which is low cost on most software platforms.  Tinygo may perform relatively worse, although it should not matter much.  A system with a biquad and PID will require ~40 clocks per update, sufficient for 25MHz at 1GHz CPU clock, ignoring other loop elements.

The user should do biquad filter design with another tool (Scipy, Matlab, online calculators, ...).  This package exists
to assemble control systems in Go, not design them.

## Usage

This package is designed with a homogenious API; all types implement `pctl.Updater` and can be combined elegantly with the `pctl.Cascade` function.  For example, suppose you want to low-pass filter your measurements before sending them to your controller.  You may write the following snippet:

```go
// Biquad, 1k sample rate, 50Hz corner freq, Q=sqrt(2)/2, -6dB gain
a0 := 0.0201
a1 := 0.0402
a2 := a0
b1 := -1.5610
b2 := 0.6414
inputFilter := pctl.NewBiquad(a0, a1, a2, b1, b2)
controller := pctl.PID{P: 1, I: 0.5, Setpt: 50}
for {
    input := getInput()
    process := pctl.Cascade(input, inputFilter, controller)
    applyControl(process)
}
```

You can also imagine cascading filter for bandpass control, etc.


## Performance

See `pctl_test.go` for a benchmark suite.  Each of these types require 2-4 ns per `Update()` on an intel i7-9700k, Windows 10 x64 platform.

## Design

Several designs have been iterated in this repository.  An early design used channels to communicate, which took about 500ns per update.  This was less composable than methods/functions.

An intermediate design maintained clocks inside each control element.  This was less performant, but more importantly could not be used in a simulation capacity running at any speed other than real time.  Explicitly including dT (fielded as DT) in the structs allows these controllers to be used in simulation studies as well.  The nearly 10x increase in performance and better friendliness to tinygo platforms are also nice benefits.
