/*Package pctl facilitates fluent and instrumented control of measured processes
 */
package pctl

import (
	"sync"
	"time"
)

// PIDCtl is a Proportional, Integral, Derivative controller
//
// If the loop is running and the Locking variable is changed, the thread
// may deadlock.  Acquire the lock when turning off locking while the loop is
// running to prevent this.
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
	Locking bool
}

// LastObs returns the read-only last observation time.
func (pid *PIDCtl) LastObs() time.Time {
	return pid.lastObs
}

/*Loop runs the PID loop.  It takes a channel of measurements to read from
and a channel of outputs to write to.  To stop the loop, simply close m.  The
remaining values in the channel will be exhausted, then the loop exited.

If you desire close (stop) to be immediate, use an unbuffered m channel.

The first observation in m is used to seed the loop with a fresh measurement.
The existing output value is sent on o to maintain synchronization.

use go Loop(...) to avoid blocking the calling thread.

To do multiple things with the output value, fan out the channel
or provide that logic in the consuming function.

The struct fields or lastObs() may be accessed at any time while the loop
is running.  To guarantee that they are in sync, ensure pid.Locking == true
and acquire the lock during your read.
*/
func (pid *PIDCtl) Loop(m <-chan float64, o chan<- float64) {
	var (
		// prevErr is the previous error, integral is the integral error
		prevErr  float64 = 0
		integral float64 = 0
		started          = false
	)
	for pVal := range m {
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
		if pid.Locking {
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
		if pid.Locking {
			pid.Unlock()
		}
	}
}
