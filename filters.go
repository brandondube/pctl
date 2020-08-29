package pctl

import (
	"math"
	"time"
)

// LPF is a digital discrete-time low pass filter.  It does not require
// Update() to be called in a regular cadence.
//
// The cutoff frequency may not be changed; the Previous() method can retrieve
// the state of the filter to create a new one.
type LPF struct {
	rc   float64
	fc   float64
	prev float64
	last time.Time
}

// NewLPF returns a new low pass filter with the specified cutoff frequency
// in Hertz
func NewLPF(cutoffFreq float64) *LPF {
	return &LPF{
		fc: cutoffFreq,
		rc: 1 / (2 * math.Pi * cutoffFreq)}
}

// Update processes an input value, returning the filtered output
func (l *LPF) Update(input float64) float64 {
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
func (l *LPF) Cutoff() float64 {
	return l.fc
}

// RC returns the filter's RC time constant, 1/(2pi cutoff)
func (l *LPF) RC() float64 {
	return l.rc
}

// Previous returns the last value at the filter's output and when it happened
func (l *LPF) Previous() (float64, time.Time) {
	return l.prev, l.last
}

// HPF is a digital discrete-time high pass filter.  It does not require
// Update() to be called in a regular cadence.
//
// The cutoff frequency may not be changed; the Previous() method can retrieve
// the state of the filter to create a new one.
type HPF struct {
	rc   float64
	fc   float64
	prev float64
	last time.Time
}

// NewHPF returns a new low pass filter with the specified cutoff frequency
// in Hertz
func NewHPF(cutoffFreq float64) *HPF {
	return &HPF{
		fc: cutoffFreq,
		rc: 1 / (2 * math.Pi * cutoffFreq)}
}

// Update processes an input value, returning the filtered output
func (h *HPF) Update(input float64) float64 {
	if h.last.IsZero() {
		// not initialized
		h.prev = input
		h.last = time.Now()
		return input
	}
	now := time.Now()
	dT := (now.Sub(h.last)).Seconds()
	alpha := h.rc / (h.rc + dT)
	h.prev = alpha * (h.prev + dT)
	h.last = now
	return h.prev
}

// Cutoff returns the filter's cutoff frequency in Hertz
func (h *HPF) Cutoff() float64 {
	return h.fc
}

// RC returns the filter's RC time constant, 1/(2pi cutoff)
func (h *HPF) RC() float64 {
	return h.rc
}

// Previous returns the last value at the filter's output and when it happened
func (h *HPF) Previous() (float64, time.Time) {
	return h.prev, h.last
}
