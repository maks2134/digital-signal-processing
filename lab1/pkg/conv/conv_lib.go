package conv

import "gonum.org/v1/gonum/dsp/fourier"

func nextPowerOf2Conv(n int) int {
	if n <= 0 {
		return 1
	}
	if n&(n-1) == 0 {
		return n
	}
	power := 1
	for power < n {
		power <<= 1
	}
	return power
}

func ConvolveLib(x, y []float64) []float64 {
	N := len(x)
	M := len(y)
	if N == 0 || M == 0 {
		return []float64{}
	}

	L := N + M - 1
	fftSize := nextPowerOf2Conv(L)

	xPadded := make([]float64, fftSize)
	yPadded := make([]float64, fftSize)
	copy(xPadded, x)
	copy(yPadded, y)

	fft := fourier.NewFFT(fftSize)

	X := fft.Coefficients(nil, xPadded)
	Y := fft.Coefficients(nil, yPadded)

	coeffLen := len(X)
	if len(Y) < coeffLen {
		coeffLen = len(Y)
	}

	Z := make([]complex128, coeffLen)
	for i := 0; i < coeffLen; i++ {
		Z[i] = X[i] * Y[i]
	}

	outFull := make([]float64, fftSize)
	fft.Sequence(outFull, Z)

	return outFull[:L]
}
