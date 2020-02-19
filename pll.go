package pctl

import "time"

// pll is a part of pctl that contains a phase locked loop

// PhaseLock is a struct which uses variable sleeps to control jitter
// in an update process
type PhaseLock struct {
	// Interval is the desired update
	Interval time.Duration

	// started is the time the current pair of Start() Stop() calls was started
	started time.Time
}

// Start should be called at the beginning of a loop to be stabilized
func (pl *PhaseLock) Start() {
	pl.started = time.Now()
}

// Stop computes the time since Start was called and sleeps if necessary
// to ensure the full interval elapses
func (pl *PhaseLock) Stop() {
	dT := time.Now().Sub(pl.started)
	if dT > pl.Interval {
		return
	}
	time.Sleep(pl.Interval - dT)
}
