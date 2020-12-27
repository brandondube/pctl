package pctl

import (
	"math"
)

// LPF is a digital discrete-time single pole / first order low pass filter.
//
// It requires 1 division, 1 multiply, ~ 11 clocks per update
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

// HPF is a digital discrete-time single pole / first order high pass filter.
//
// It requires 1 division, 1 multiply, and two additions, ~ 10 clocks per update
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

// Biquad is a digital discrete-time Biquad filter.  It is implemented using the
// "type 2 transposed"method which accumulates the least floating point error.
//
// The variable naming convention follows Digital Audio Signal Processing, Zölzer
// with a in the numerator and b in the denominator. Coefficients should be
// normalized by b0.
//
// For more information see e.g.
//
// http://www.earlevel.com/main/2013/10/13/biquad-calculator-v2/
//
// http://www.earlevel.com/main/2003/02/28/biquads/
//
// https://www.earlevel.com/main/2012/11/26/biquad-c-source-code/
//
// Biquads require 5 multiplies and 4 additions/subtractions ~ 24 clocks per update
type Biquad struct {
	a0      float64
	a1      float64
	a2      float64
	b1      float64
	b2      float64
	z1      float64
	z2      float64
	prevIn  float64
	prevOut float64
}

// NewBiquad returns a new biquad filter
func NewBiquad(a0, a1, a2, b1, b2 float64) *Biquad {
	return &Biquad{
		a0: a0,
		a1: a1,
		a2: a2,
		b1: b1,
		b2: b2,
	} // z1..prevOut init to 0
}

// Update processes an input value, returning the filtered output
func (b *Biquad) Update(input float64) float64 {
	out := b.a0*input + b.z1
	b.z1 = input*b.a1 + b.z2 - b.b1*out
	b.z2 = input*b.a2 - b.b2*out
	return out
}
