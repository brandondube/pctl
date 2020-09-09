package pctl

import (
	"math"
)

// LPF is a digital discrete-time low pass filter.  It does not require
// Update() to be called in a regular cadence.
//
// The cutoff frequency may not be changed; the Previous() method can retrieve
// the state of the filter to create a new one.
type LPF struct {
	// DT is the inter-update time in seconds
	DT   float64
	rc   float64
	fc   float64
	prev float64
}

// NewLPF returns a new low pass filter with the specified corner frequency
// in Hertz
func NewLPF(cutoffFreq, dT float64) *LPF {
	return &LPF{
		fc: cutoffFreq,
		rc: 1 / (2 * math.Pi * cutoffFreq),
		DT: dT}
}

// Update processes an input value, returning the filtered output
func (l *LPF) Update(input float64) float64 {
	alpha := l.DT / (l.rc + l.DT)
	l.prev = l.prev + alpha*(input-l.prev)
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

// HPF is a digital discrete-time high pass filter.  It does not require
// Update() to be called in a regular cadence.
//
// The cutoff frequency may not be changed; the Previous() method can retrieve
// the state of the filter to create a new one.
type HPF struct {
	// DT is the inter-update time in seconds
	DT   float64
	rc   float64
	fc   float64
	prev float64
}

// NewHPF returns a new low pass filter with the specified corner frequency
// in Hertz
func NewHPF(cutoffFreq, dT float64) *HPF {
	return &HPF{
		fc: cutoffFreq,
		rc: 1 / (2 * math.Pi * cutoffFreq),
		DT: dT}
}

// Update processes an input value, returning the filtered output
func (h *HPF) Update(input float64) float64 {
	alpha := h.rc / (h.rc + h.DT)
	h.prev = alpha * (h.prev + h.DT)
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
