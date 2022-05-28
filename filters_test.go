package pctl

import (
	"math"
	"math/rand"
	"testing"
)

const biquadFilterCoefTol = 1e-8

func approxEqualAbs(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

// these tests verify physical properties of the filters, but not e.g.
// that the -3dB point is correct, or any internal state variables are correct

func TestLowPassFilterAsymptotic(t *testing.T) {
	lpf := NewLPF(1e6, 1e-3)
	// 1Mhz low-pass filter has corner frequency of a MHz, corresponds to
	// microseconds to reach nearly steady state
	process := lpf.Update(0) // seed the filter
	target := 1.
	for i := 0; i < 5; i++ {
		process = lpf.Update(target)
	}
	err := 1 - process
	if math.Abs(err) > 1e-5 {
		t.Errorf("process of %f has error of %f, expected to converge to 1 along step", process, err)
	}
}

func TestHighPassRejectsOscillation(t *testing.T) {
	hpf := NewHPF(1e6, 1e-3)
	// cutoff of 1Mhz means that inputs at ms should be rejected
	process := hpf.Update(0)
	target := 1.
	for i := 0; i < 5; i++ {
		process = hpf.Update(target)
	}
	err := 0 - process
	if math.Abs(err) > 1e-5 {
		t.Errorf("process of %f has error of %f, expected to resist motion in ms domain", process, err)
	}
}

func TestBiquadFilterAsymptotic(t *testing.T) {
	// Bq at sample rate 1kHz, 250Hz corner, Q=sqrt(2)/2, -6dB gain
	a0 := 0.2928920553392428
	a1 := 0.5857841106784856
	a2 := a0
	b1 := -1.3007020142696517e-16
	b2 := 0.17156822135697122
	bq := NewBiquad(a0, a1, a2, b1, b2)
	process := bq.Update(0)
	target := 1.
	for i := 0; i < 100; i++ {
		process = bq.Update(target)
	}
	err := target - process
	if math.Abs(err) > 1e-5 {
		t.Errorf("process of %f has error of %f, expected to converge to target=1", process, err)
	}
}

func TestStateSpaceReducesNoise(t *testing.T) {
	A := [][]float64{
		{2, -1},
		{1, 0},
	}
	B := []float64{5e-5, 0}
	C := []float64{4, 0.02}
	D := 5e-5
	filt := NewStateSpaceFilter(A, B, C, D, nil)
	var maxIn float64
	var maxOut float64
	for i := 0; i < 100; i++ {
		in := rand.Float64()
		if in > maxIn {
			maxIn = in
		}
		out := filt.Update(in)
		if out > maxOut {
			maxOut = out
		}
	}
	if maxOut >= maxIn {
		t.Errorf("state-space lowpass filter failed to reduce peak noise")
	}
}

func TestStateSpaceNonzeroOutput(t *testing.T) {
	A := [][]float64{
		{2, -1},
		{1, 0},
	}
	B := []float64{5e-5, 0}
	C := []float64{4, 0.02}
	D := 5e-5
	filt := NewStateSpaceFilter(A, B, C, D, nil)
	var maxOut float64
	for i := 0; i < 100; i++ {
		in := rand.Float64()
		out := filt.Update(in)
		if out > maxOut {
			maxOut = out
		}
	}
	if maxOut == 0 {
		t.Errorf("state-space lowpass filter returned zero where it should not.")
	}
}

func testBiquadvsEarLevel(t *testing.T, newF NewBiquadFunc, a0, a1, a2, b1, b2 float64) {
	// assumes below parameters (default for biquad calculator v3)
	// were used to compute a0..b2
	b := newF(44100, 100, 0.7071, 6)
	if !approxEqualAbs(b.a0, a0, biquadFilterCoefTol) {
		t.Errorf("a0 %f != %f", b.a0, a0)
	}
	if !approxEqualAbs(b.a1, a1, biquadFilterCoefTol) {
		t.Errorf("a1 %f != %f", b.a1, a1)
	}
	if !approxEqualAbs(b.a2, a2, biquadFilterCoefTol) {
		t.Errorf("a2 %f != %f", b.a2, a2)
	}
	if !approxEqualAbs(b.b1, b1, biquadFilterCoefTol) {
		t.Errorf("b1 %f != %f", b.b1, b1)
	}
	if !approxEqualAbs(b.b2, b2, biquadFilterCoefTol) {
		t.Errorf("b2 %f != %f", b.b2, b2)
	}
}

// these tests are regression against earlevel.com, where the C++ implementation
// was borrowed from

func TestLowPassBiquadCorrectCoefs(t *testing.T) {
	testBiquadvsEarLevel(t, NewBiquadLowpass,
		0.00005024141818873903,
		0.00010048283637747806,
		0.00005024141818873903,
		-1.979851353142371,
		0.9800523188151258)
}

func TestHighPassBiquadCorrectCoefs(t *testing.T) {
	testBiquadvsEarLevel(t, NewBiquadHighpass,
		0.9899759179893742,
		-1.9799518359787485,
		0.9899759179893742,
		-1.979851353142371,
		0.9800523188151258)
}

func TestBandPassBiquadCorrectCoefs(t *testing.T) {
	testBiquadvsEarLevel(t, NewBiquadBandpass,
		0.009973840592437116,
		0,
		-0.009973840592437116,
		-1.979851353142371,
		0.9800523188151258)
}

func TestNotchBiquadCorrectCoefs(t *testing.T) {
	testBiquadvsEarLevel(t, NewBiquadNotch,
		0.990026159407563,
		-1.979851353142371,
		0.990026159407563,
		-1.979851353142371,
		0.9800523188151258)
}

func TestPeakBiquadCorrectCoefs(t *testing.T) {
	testBiquadvsEarLevel(t, NewBiquadPeak,
		1.0099265876771597,
		-1.979851353142371,
		0.9701257311379663,
		-1.979851353142371,
		0.9800523188151258)
}

func TestLowShelfBiquadCorrectCoefs(t *testing.T) {
	testBiquadvsEarLevel(t, NewBiquadLowShelf,
		1.0041645480379269,
		-1.9797515357244457,
		0.9759879669583226,
		-1.9798515425143588,
		0.9800525082063363)
}

func TestHighShelfBiquadCorrectCoefs(t *testing.T) {
	testBiquadvsEarLevel(t, NewBiquadHighShelf,
		1.9894003627867007,
		-3.95042317880182,
		1.9612237817070965,
		-1.9798515425143588,
		0.9800525082063363)
}

func TestIdentityFIRFilterDoesNothing(t *testing.T) {
	f := NewFIRFilter([]float64{1, 0, 0, 0})
	input := []float64{3.14, 2.87, 1}
	for i := 0; i < 3; i++ {
		in := input[i]
		out := f.Update(in)
		if !approxEqualAbs(in, out, 1e-16) {
			t.Errorf("sample %d had input-output mismatch %f != %f", i, in, out)
		}
	}
}

func TestLagFIRFilterLags(t *testing.T) {
	f := NewFIRFilter([]float64{0, 1, 0, 0})
	input := []float64{3.14, 2.87, 1}
	f.Update(input[0])
	for i := 1; i < 3; i++ {
		in := input[i]
		out := f.Update(in)
		if !approxEqualAbs(input[i-1], out, 1e-16) {
			t.Errorf("sample %d had input-output mismatch %f != %f", i, in, out)
		}
	}
}
