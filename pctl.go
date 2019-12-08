/*Package pctl facilitates fluent and instrumented control of measured processes
 */
package pctl

import (
	"sync"
	"time"
)

// PIDCtl is a Proportional, Integral, Derivative controller
type PIDCtl struct {
	sync.Mutex

	// P is the proportional gain, unitless
	P float64

	// I is the integral gain, units of reciprocal seconds
	I float64

	// D is the derivative gain, units of seconds
	D float64

	// Setpt is the setpoint, in process units
	Setpt float64

	// Meas is the measured value, in process units
	Meas float64

	// Output is the computed output value, in process units
	Output float64

	// lastObs is the last observation time
	lastObs time.Time

	// if Locking is true, the embedded mutex is acquired during the loop
	locking bool

	// Running indicates if the loop is currently running
	Running bool

	// signal is the channel used to send commands to the PID loops
	signal chan string
}

// LastObs returns the read-only last observation time.
func (pid *PIDCtl) LastObs() time.Time {
	return pid.lastObs
}

// EnableLocking enables the lock and prevents a deadlock inside a concurrently running Loop()
func (pid *PIDCtl) EnableLocking() {
	pid.Lock()
	pid.locking = true
	pid.Unlock()
}

// DisableLocking disables the lock and prevents a deadlock inside a concurrently running Loop()
func (pid *PIDCtl) DisableLocking() {
	pid.Lock()
	pid.locking = false
	pid.Unlock()
}

// Locking indicates if the controller uses locks
func (pid *PIDCtl) Locking() bool {
	return pid.locking
}

// Pause temporarely suspends the PID loop.  The behavior due to a chain of
// calls Pause() Stop() is undefined.  Always Resume() or exit the program
// after calling Pause().
func (pid *PIDCtl) Pause() {
	pid.signal <- "pause"
	pid.Running = false
	// the value of pause and resume isn't really needed for the machine,
	// but it is more clear for the ape behind the keyboard to understand
	// what is going on
}

// Resume restarts the PID loop from suspension without zeroing the state
func (pid *PIDCtl) Resume() {
	pid.signal <- "resume"
	pid.Running = true
}

// Stop the PID loop.
//
// The behavior due to a chain of calls Pause() Stop() is undefined.  Always
// Resume() or exit the program after calling Pause().
func (pid *PIDCtl) Stop() {
	pid.signal <- "stop"
	pid.Running = false
}

/*Start runs the PID loop.  It takes a channel of measurements to read from
and a channel of outputs to write to.

The first observation in m is used to seed the loop with a fresh measurement.
The existing output value is sent on o to maintain synchronization.

use go Start(...) to avoid blocking the calling thread.

To do multiple things with the output value, fan out the channel
or provide that logic in the consuming function.

The struct fields or LastObs() may be accessed at any time while the loop
is running.  To guarantee that they are in sync, ensure pid.Locking() == true
and acquire the lock during your read.
*/
func (pid *PIDCtl) Start(m <-chan float64, o chan<- float64) {
	pid.loop(m, o)
}

func (pid *PIDCtl) loop(m <-chan float64, o chan<- float64) {
	// guard against the uninitialized case
	if pid.signal == nil {
		pid.signal = make(chan string)
	}
	pid.Running = true
	defer func() { pid.Running = false }()
	var (
		// prevErr is the previous error, integral is the integral error
		prevErr  float64 = 0
		integral float64 = 0
		started          = false
	)
	for {
		select {
		case cmd := <-pid.signal:
			if cmd == "pause" {
				<-pid.signal // wait for the resume command
			} else if cmd == "stop" {
				return
			}

		case pVal := <-m:
			// update the clock and measurement
			updateT := time.Now()
			pid.Meas = pVal
			// zero time, just update the measurement and output
			// to seed the loop and don't try to control
			if !started {
				pid.lastObs = updateT
				pid.Meas = pVal
				started = true
				o <- pid.Output
				continue
			}
			if pid.locking {
				pid.Lock()
			}
			dt := updateT.Sub(pid.lastObs).Seconds()
			err := pid.Setpt - pVal
			integral += err * dt
			derivative := (err - prevErr) / dt
			pid.Output = pid.P*err + pid.I*integral + pid.D*derivative
			pid.Meas = pVal
			pid.lastObs = updateT
			prevErr = err
			o <- pid.Output
			if pid.locking {
				pid.Unlock()
			}
		}
	}
}
