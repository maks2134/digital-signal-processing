package fft

import (
	"gonum.org/v1/gonum/dsp/fourier"
)

func nextPowerOf2(n int) int {
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

func DFTLib(x []float64) []complex128 {
	N := len(x)
	fftSize := nextPowerOf2(N)

	xPadded := make([]float64, fftSize)
	copy(xPadded, x)

	fft := fourier.NewFFT(fftSize)
	X := fft.Coefficients(nil, xPadded)

	if N > len(X) {
		N = len(X)
	}
	if N <= 0 {
		return []complex128{}
	}
	return X[:N]
}

func IDFTLib(X []complex128) []float64 {
	N := len(X)
	if N == 0 {
		return []float64{}
	}

	fftSize := nextPowerOf2(N)
	coeffSize := fftSize/2 + 1

	if N <= coeffSize {
		XPadded := make([]complex128, coeffSize)
		copy(XPadded, X)

		fft := fourier.NewFFT(fftSize)
		result := make([]float64, fftSize)
		fft.Sequence(result, XPadded)

		if N > len(result) {
			N = len(result)
		}
		return result[:N]
	} else {
		fftSize = nextPowerOf2(N)
		coeffSize = fftSize/2 + 1
		if N > coeffSize {
			X = X[:coeffSize]
		}
		XPadded := make([]complex128, coeffSize)
		copy(XPadded, X)

		fft := fourier.NewFFT(fftSize)
		result := make([]float64, fftSize)
		fft.Sequence(result, XPadded)

		if N > len(result) {
			N = len(result)
		}
		return result[:N]
	}
}
