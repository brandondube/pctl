# pctl

pctl, "process control" is a package for process control in Go.  This applies to industrial processes, not the computer science term.

It contains a streaming PID controller that communicates using channels.

## Usage

```go
import pctl

controller := pctl.PIDCtl{
    P: 10,
    I: 1,
    D: 0, // PI, not PID.
}
controller.Setpt = 50
measC := make(chan float64) // unbuffered or buffered makes no difference
outC := make(chan float64)

go pollSensor(measC) // this functions writes f64s to measC
go processesOutputs(outC) // this function reads f64s from outC
go controller.Loop(measC, outC)
```
While the controller is running, you can adjust the setpoint or gains, query the measurement and output values, or the clock:
```go
v1 := controller.Meas
v2 := controller.Setpt
controller.Setpt = 55
```
In simple terms, you can access the struct fields or the `LastObs()` method.  By default, the controller does not use a lock during the loop.  Accessing the fields may result
in desyncronization in this case, though it will not make a meaningful difference in the numerical values if the system has stabilized.  If synchronicity is more important than performance, enable locking: `controller.EnableLocking()`.

To monitor (instrument) the process and controller at a regular interval, you can use a `time.Timer`:

```go
import time

pollFreq := time.Second()
timer := time.NewTimer(pollFreq)
go func() {
    ts := <-timer.C
    pV := controller.Meas
    pO := controller.Output
    // do something with ts, pV, pO
    // write to CSV file or database, etc.
}()

// wait some time
close(measC) // this stops the controller
timer.Stop()
```
