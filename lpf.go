package pctl

import (
	"math"
	"time"
)

// LowPass is a digital discrete-time low pass filter.  It does not require
// Update() to be called in a regular cadence.
//
// The cutoff frequency may not be changed; the Previous() method can retrieve
// the state of the filter to create a new one.
type LowPass struct {
	rc   float64
	fc   float64
	prev float64
	last time.Time
}

// NewLowPass returns a new low pass filter with the specified cutoff frequency
// in Hertz
func NewLowPass(cutoffFreq float64) *LowPass {
	return &LowPass{
		fc: cutoffFreq,
		rc: 1 / (2 * math.Pi * cutoffFreq)}
}

// Update processes an input value, returning the filtered output
func (l *LowPass) Update(input float64) float64 {
	if l.last.IsZero() {
		// not initialized
		l.prev = input
		l.last = time.Now()
		return input
	}
	now := time.Now()
	dT := (now.Sub(l.last)).Seconds()
	alpha := dT / (l.rc + dT)
	l.prev = l.prev + alpha*(input-l.prev)
	l.last = now
	return l.prev
}

// Cutoff returns the filter's cutoff frequency in Hertz
func (l *LowPass) Cutoff() float64 {
	return l.fc
}

// RC returns the filter's RC time constant, 1/(2pi cutoff)
func (l *LowPass) RC() float64 {
	return l.rc
}

// Previous returns the last value at the filter's output and when it happened
func (l *LowPass) Previous() (float64, time.Time) {
	return l.prev, l.last
}
