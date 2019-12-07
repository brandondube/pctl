# pctl

pctl, "process control" is a package for process control in Go.  This applies to industrial processes, not the computer science term.

It contains a streaming PID controller that communicates using channels.  Usage
is simple and idiomatic Go:

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

go funcThatPollsSensor(measC)
go funcThatProcessesOutputs(outC)
go controller.Loop(measC, outC)
```
