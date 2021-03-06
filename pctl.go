/*Package pctl provides the building blocks for high performance control systems
 */
package pctl

// Updater is the essential building block of a DSP or control system
type Updater interface {
	Update(float64) float64
}

// Cascade applies a chain of updaters in the sequence given
func Cascade(input float64, chain ...Updater) float64 {
	for _, elem := range chain {
		// re-assigning input avoids "seeding" the loop
		input = elem.Update(input)
	}
	return input
}

// Setpoint implements Updater and returns the process error.
type Setpoint float64

// Update computes the process error, meas - setpt
func (s *Setpoint) Update(meas float64) float64 {
	return meas - float64(*s)
}
