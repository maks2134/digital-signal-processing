package filter

import (
	"math"
)

func MovingAverageFilter(signal []float64, M int) []float64 {
	if M <= 0 || len(signal) == 0 {
		return signal
	}
	result := make([]float64, len(signal))
	for n := 0; n < len(signal); n++ {
		sum := 0.0
		count := 0
		for m := 0; m < M; m++ {
			if n-m >= 0 {
				sum += signal[n-m]
				count++
			}
		}
		if count > 0 {
			result[n] = sum / float64(count)
		}
	}
	return result
}

// Хеннинга
func HannWindow(M int) []float64 {
	window := make([]float64, M)
	for n := 0; n < M; n++ {
		window[n] = 0.5 * (1 - math.Cos(2*math.Pi*float64(n)/float64(M-1)))
	}
	return window
}

// ких режекторный
func DesignFIRBandStop(f1, f2, fs float64, M int) []float64 {
	if M%2 == 0 {
		M++
	}

	w1 := 2 * math.Pi * f1 / fs
	w2 := 2 * math.Pi * f2 / fs
	_ = w1
	_ = w2

	N := M / 2
	h_lp := make([]float64, M)

	fc1 := f1 / fs
	for n := 0; n < M; n++ {
		if n == N {
			h_lp[n] = 2 * fc1
		} else {
			h_lp[n] = math.Sin(2*math.Pi*fc1*float64(n-N)) / (math.Pi * float64(n-N))
		}
	}

	fc2 := f2 / fs
	h_hp := make([]float64, M)
	for n := 0; n < M; n++ {
		if n == N {
			h_hp[n] = 1 - 2*fc2
		} else {
			h_hp[n] = -math.Sin(2*math.Pi*fc2*float64(n-N)) / (math.Pi * float64(n-N))
		}
	}

	window := HannWindow(M)
	h_bandstop := make([]float64, M)

	for n := 0; n < M; n++ {
		h_bandstop[n] = (h_lp[n] + h_hp[n]) * window[n]
	}

	return h_bandstop
}

// ких
func FIRFilter(signal, coeffs []float64) []float64 {
	if len(signal) == 0 || len(coeffs) == 0 {
		return signal
	}

	result := make([]float64, len(signal))

	for n := 0; n < len(signal); n++ {
		sum := 0.0
		for k := 0; k < len(coeffs); k++ {
			if n-k >= 0 {
				sum += coeffs[k] * signal[n-k]
			}
		}
		result[n] = sum
	}

	return result
}

type BiquadCoefficients struct {
	B0, B1, B2 float64
	A1, A2     float64
}

func DesignBandPassIIR(f0, bw, fs float64) BiquadCoefficients {
	w0 := 2 * math.Pi * f0 / fs
	_ = 2 * math.Pi * bw / fs

	Q := f0 / bw

	sinW0 := math.Sin(w0)
	cosW0 := math.Cos(w0)
	alpha := sinW0 / (2 * Q)

	b0 := alpha
	b1 := 0.0
	b2 := -alpha
	a0 := 1 + alpha
	a1 := -2 * cosW0
	a2 := 1 - alpha

	return BiquadCoefficients{
		B0: b0 / a0,
		B1: b1 / a0,
		B2: b2 / a0,
		A1: a1 / a0,
		A2: a2 / a0,
	}
}

// бих
func IIRFilter(signal []float64, coeffs BiquadCoefficients) []float64 {
	if len(signal) == 0 {
		return signal
	}

	result := make([]float64, len(signal))

	x_prev1 := 0.0
	x_prev2 := 0.0
	y_prev1 := 0.0
	y_prev2 := 0.0

	for n := 0; n < len(signal); n++ {
		y := coeffs.B0*signal[n] + coeffs.B1*x_prev1 + coeffs.B2*x_prev2 -
			coeffs.A1*y_prev1 - coeffs.A2*y_prev2

		result[n] = y

		x_prev2 = x_prev1
		x_prev1 = signal[n]
		y_prev2 = y_prev1
		y_prev1 = y
	}

	return result
}

func ComputeFrequencyResponse(coeffs interface{}, fs float64, nFreqs int) ([]float64, []float64) {
	freqs := make([]float64, nFreqs)
	magnitude := make([]float64, nFreqs)

	switch c := coeffs.(type) {
	case []float64:
		for k := 0; k < nFreqs; k++ {
			freq := float64(k) * fs / float64(2*nFreqs)
			freqs[k] = freq
			w := 2 * math.Pi * freq / fs

			real := 0.0
			imag := 0.0

			for n := 0; n < len(c); n++ {
				real += c[n] * math.Cos(float64(n)*w)
				imag -= c[n] * math.Sin(float64(n)*w)
			}

			magnitude[k] = math.Sqrt(real*real + imag*imag)
		}

	case BiquadCoefficients:
		for k := 0; k < nFreqs; k++ {
			freq := float64(k) * fs / float64(2*nFreqs)
			if freq > fs/2 {
				freq = fs / 2
			}
			freqs[k] = freq
			w := 2 * math.Pi * freq / fs

			numReal := c.B0 + c.B1*math.Cos(w) + c.B2*math.Cos(2*w)
			numImag := -c.B1*math.Sin(w) - c.B2*math.Sin(2*w)

			denReal := 1 + c.A1*math.Cos(w) + c.A2*math.Cos(2*w)
			denImag := -c.A1*math.Sin(w) - c.A2*math.Sin(2*w)

			numMag := math.Sqrt(numReal*numReal + numImag*numImag)
			denMag := math.Sqrt(denReal*denReal + denImag*denImag)

			if denMag > 1e-10 {
				magnitude[k] = numMag / denMag
			}
		}
	}

	return freqs, magnitude
}
