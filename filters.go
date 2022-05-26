package pctl

import (
	"math"
)

// LPF is a digital discrete-time single pole / first order low pass filter.
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

// HPF is a digital discrete-time single pole / first order high pass filter.
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

// NewBigQuadXXXX code adapted from Nigel Redmon's C++ Biquad implementation
// see https://www.earlevel.com/main/2012/11/26/biquad-c-source-code/
type NewBiquadFunc func(float64, float64, float64, float64) *Biquad

func NewBiquadLowpass(Fs, f, Q, g float64) *Biquad {
	Fc := f / Fs
	K := math.Tan(math.Pi * Fc)
	norm := 1 / (1 + K/Q + K*K)
	a0 := K * K * norm
	a1 := 2 * a0
	a2 := a0
	b1 := 2 * (K*K - 1) * norm
	b2 := (1 - K/Q + K*K) * norm
	return NewBiquad(a0, a1, a2, b1, b2)
}

func NewBiquadHighpass(Fs, f, Q, g float64) *Biquad {
	Fc := f / Fs
	K := math.Tan(math.Pi * Fc)
	norm := 1 / (1 + K/Q + K*K)
	a0 := norm
	a1 := -2 * a0
	a2 := a0
	b1 := 2 * (K*K - 1) * norm
	b2 := (1 - K/Q + K*K) * norm
	return NewBiquad(a0, a1, a2, b1, b2)
}

func NewBiquadBandpass(Fs, f, Q, g float64) *Biquad {
	Fc := f / Fs
	K := math.Tan(math.Pi * Fc)
	norm := 1 / (1 + K/Q + K*K)
	a0 := K / Q * norm
	a1 := 0.
	a2 := -a0
	b1 := 2 * (K*K - 1) * norm
	b2 := (1 - K/Q + K*K) * norm
	return NewBiquad(a0, a1, a2, b1, b2)
}

func NewBiquadNotch(Fs, f, Q, g float64) *Biquad {
	Fc := f / Fs
	K := math.Tan(math.Pi * Fc)
	norm := 1 / (1 + K/Q + K*K)
	a0 := (1 + K*K) * norm
	a1 := 2 * (K*K - 1) * norm
	a2 := a0
	b1 := a1
	b2 := (1 - K/Q + K*K) * norm
	return NewBiquad(a0, a1, a2, b1, b2)
}

func NewBiquadPeak(Fs, f, Q, g float64) *Biquad {
	Fc := f / Fs
	V := math.Pow(10, math.Abs(g)/20)
	K := math.Tan(math.Pi * Fc)
	var norm, a0, a1, a2, b1, b2 float64
	if g >= 0 {
		norm = 1 / (1 + 1/Q*K + K*K)
		a0 = (1 + V/Q*K + K*K) * norm
		a1 = 2 * (K*K - 1) * norm
		a2 = (1 - V/Q*K + K*K) * norm
		b1 = a1
		b2 = (1 - 1/Q*K + K*K) * norm
	} else {
		norm = 1 / (1 + V/Q*K + K*K)
		a0 = (1 + 1/Q*K + K*K) * norm
		a1 = 2 * (K*K - 1) * norm
		a2 = (1 - 1/Q*K + K*K) * norm
		b1 = a1
		b2 = (1 - V/Q*K + K*K) * norm
	}
	return NewBiquad(a0, a1, a2, b1, b2)
}

func NewBiquadLowShelf(Fs, f, Q, g float64) *Biquad {
	Fc := f / Fs
	V := math.Pow(10, math.Abs(g)/20)
	K := math.Tan(math.Pi * Fc)
	var norm, a0, a1, a2, b1, b2 float64
	if g >= 0 {
		norm = 1 / (1 + math.Sqrt(2)*K + K*K)
		a0 = (1 + math.Sqrt(2*V)*K + V*K*K) * norm
		a1 = 2 * (V*K*K - 1) * norm
		a2 = (1 - math.Sqrt(2*V)*K + V*K*K) * norm
		b1 = 2 * (K*K - 1) * norm
		b2 = (1 - math.Sqrt(2)*K + K*K) * norm
	} else {
		norm = 1 / (1 + math.Sqrt(2*V)*K + V*K*K)
		a0 = (1 + math.Sqrt(2)*K + K*K) * norm
		a1 = 2 * (K*K - 1) * norm
		a2 = (1 - math.Sqrt(2)*K + K*K) * norm
		b1 = 2 * (V*K*K - 1) * norm
		b2 = (1 - math.Sqrt(2*V)*K + V*K*K) * norm
	}
	return NewBiquad(a0, a1, a2, b1, b2)
}

func NewBiquadHighShelf(Fs, f, Q, g float64) *Biquad {
	Fc := f / Fs
	V := math.Pow(10, math.Abs(g)/20)
	K := math.Tan(math.Pi * Fc)
	var norm, a0, a1, a2, b1, b2 float64
	if g >= 0 {
		norm = 1 / (1 + math.Sqrt(2)*K + K*K)
		a0 = (V + math.Sqrt(2*V)*K + K*K) * norm
		a1 = 2 * (K*K - V) * norm
		a2 = (V - math.Sqrt(2*V)*K + K*K) * norm
		b1 = 2 * (K*K - 1) * norm
		b2 = (1 - math.Sqrt(2)*K + K*K) * norm
	} else {
		norm = 1 / (V + math.Sqrt(2*V)*K + K*K)
		a0 = (1 + math.Sqrt(2)*K + K*K) * norm
		a1 = 2 * (K*K - 1) * norm
		a2 = (1 - math.Sqrt(2)*K + K*K) * norm
		b1 = 2 * (K*K - V) * norm
		b2 = (V - math.Sqrt(2*V)*K + K*K) * norm
	}
	return NewBiquad(a0, a1, a2, b1, b2)
}

// Biquad is a digital discrete-time Biquad filter.  It is implemented using the
// "type 2 transposed" method which accumulates the least floating point error.
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
type Biquad struct {
	a0 float64
	a1 float64
	a2 float64
	b1 float64
	b2 float64
	z1 float64
	z2 float64
}

// NewBiquad returns a new biquad filter
func NewBiquad(a0, a1, a2, b1, b2 float64) *Biquad {
	return &Biquad{
		a0: a0,
		a1: a1,
		a2: a2,
		b1: b1,
		b2: b2,
	}
}

// Update processes an input value, returning the filtered output
func (b *Biquad) Update(input float64) float64 {
	out := b.a0*input + b.z1
	b.z1 = input*b.a1 + b.z2 - b.b1*out
	b.z2 = input*b.a2 - b.b2*out
	return out
}

// vectorDot takes the dot product of two vectors, it does not know the
// difference between row and column vectors
func vectorDot(a, b []float64) float64 {
	var out float64
	for i := 0; i < len(a); i++ {
		out += a[i] * b[i]
	}
	return out
}

// vectorMatrixProductSumScale computes Ax + By for matrix A, column vec x, row vec B, scalar y.
// if out is nil, a fresh slice is allocated, else it is re-used and the same slice
// is returned, i.e. the final arg and return are the same slice.
//
// This function is equivalent to the numpy code A @ x + B * y
// for A (nxm), x (1xm), B (1xm), y scalar
func vectorMatrixProductSumScale(x []float64, A [][]float64, B []float64, y float64, out []float64) []float64 {
	n := len(x)
	m := len(A)
	if out == nil {
		out = make([]float64, m)
	}
	for i := 0; i < m; i++ {
		out[i] = 0
		for j := 0; j < n; j++ {
			out[i] += A[i][j] * x[j]
		}
		out[i] += (B[i] * y)
	}
	return out
}

// StateSpaceFilter is a filter which operates on the state space of a system
// and is amenable to MIMO systems.  This implementation only operates on SISO.
type StateSpaceFilter struct {
	// x is the state of the system, column vector
	x []float64

	// A matrix of the state system
	a [][]float64

	// B Column vector of the system
	b []float64

	// C row vector of the system
	c []float64

	// D constant of the system
	d float64

	// scratch may also be the state of the system.
	// this implementation is allocation-free, and the state ping-pongs between
	// x and scratch.  It begins in x, after the first Update() is in scratch,
	// then x, then scratch, [...]
	scratch []float64
}

// NewStateSpaceFilter returns a new state-space filter with the given A,B,C,D representation and initial condition
// nil may be passed as a null initial condition (equivalent to zeros)
func NewStateSpaceFilter(A [][]float64, B, C []float64, D float64, initCond []float64) *StateSpaceFilter {
	if initCond == nil {
		initCond = make([]float64, len(B))
	}
	scratch := make([]float64, len(B))
	return &StateSpaceFilter{
		x:       initCond,
		a:       A,
		b:       B,
		c:       C,
		d:       D,
		scratch: scratch}
}

// Update updates the state-space filter and returns the filtered input
func (s *StateSpaceFilter) Update(input float64) float64 {
	vectorMatrixProductSumScale(s.x, s.a, s.b, input, s.scratch)
	out := vectorDot(s.x, s.c) + s.d*input
	s.x, s.scratch = s.scratch, s.x
	// careful in the implementation, this is a non-allocating approach.
	// s.x is a distinct slice to s.scratch, of the same size
	// "juggle the pointers" after, pointing x to "x prime" (scratch)
	// and using "old x" as the scratch on the next iteration
	return out
}

// Reset zeros the filter's internal state
func (s *StateSpaceFilter) Reset() {
	for i := 0; i < len(s.x); i++ {
		s.x[i] = 0
		s.scratch[i] = 0
	}
}
