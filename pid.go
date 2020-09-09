package pctl

// PID is a Proportional, Integral, Derivative controller.
// Use IErrMax for anti windup
type PID struct {
	// P is the proportional gain, unitless
	P float64

	// I is the integral gain, units of reciprocal seconds
	I float64

	// D is the derivative gain, units of seconds
	D float64

	// DT is the inter-update time in seconds.  If DT == 0 and I != 0 || D != 0,
	// output behavior is undefined.
	DT float64

	// IErrMax is the cap to the integral error term
	// if zero, it is ignored
	IErrMax float64

	// Setpt is the setpoint, in process units
	Setpt float64

	// input is the measured value, in process units
	input float64

	// Output is the computed output value, in process units
	output float64

	// prevErr holds the error on the previous iteration
	prevErr float64

	// integralErr is the accumulated error
	integralErr float64
}

// Update runs the loop once and returns the new output value.
// If the value is not used, or is desired again before the
// next update, it can be retrieved with pid.Output().
// if the input is desired, it can be retrieved with pid.Input().
func (pid *PID) Update(input float64) float64 {
	// update the clock and measurement
	pid.input = input

	err := pid.Setpt - input
	pid.integralErr += err * pid.DT
	if pid.IErrMax != 0 && pid.integralErr > pid.IErrMax {
		pid.integralErr = pid.IErrMax
	}
	derivative := (err - pid.prevErr) / pid.DT
	pid.output = pid.P*err + pid.I*pid.integralErr + pid.D*derivative

	pid.prevErr = err
	return pid.output
}

// Input returns the last input value
func (pid *PID) Input() float64 {
	return pid.input
}

// Output returns the last output value
func (pid *PID) Output() float64 {
	return pid.output
}

// IErr is the integral error.  You will only need to query this
// if you need to debug or tune the loop
func (pid *PID) IErr() float64 {
	return pid.integralErr
}

// IntegralReset zeros the integral error
func (pid *PID) IntegralReset() {
	pid.integralErr = 0
}
