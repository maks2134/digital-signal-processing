package fft

import (
	"math"
	"math/cmplx"
)

func DFT(x []float64) []complex128 {
	N := len(x)
	X := make([]complex128, N)
	twoPi := 2 * math.Pi
	for k := 0; k < N; k++ {
		var re, im float64
		for n := 0; n < N; n++ {
			angle := -twoPi * float64(k*n) / float64(N)
			re += x[n] * math.Cos(angle)
			im += x[n] * math.Sin(angle)
		}
		X[k] = complex(re, im)
	}
	return X
}

func FFT(x []float64) []complex128 {
	N := len(x)
	if N <= 1 {
		result := make([]complex128, N)
		for i := range x {
			result[i] = complex(x[i], 0)
		}
		return result
	}

	if N&(N-1) != 0 {
		return DFT(x)
	}

	even := make([]float64, N/2)
	odd := make([]float64, N/2)
	for i := 0; i < N/2; i++ {
		even[i] = x[2*i]
		odd[i] = x[2*i+1]
	}

	evenFFT := FFT(even)
	oddFFT := FFT(odd)

	X := make([]complex128, N)
	twoPi := 2 * math.Pi
	for k := 0; k < N/2; k++ {
		angle := -twoPi * float64(k) / float64(N)
		twiddle := cmplx.Exp(complex(0, angle))
		t := oddFFT[k] * twiddle
		X[k] = evenFFT[k] + t
		X[k+N/2] = evenFFT[k] - t
	}

	return X
}

func IDFT(X []complex128) []float64 {
	N := len(X)
	x := make([]float64, N)
	twoPi := 2 * math.Pi

	for n := 0; n < N; n++ {
		var re float64
		for k := 0; k < N; k++ {
			angle := twoPi * float64(k*n) / float64(N)
			c := math.Cos(angle)
			s := math.Sin(angle)
			re += real(X[k])*c - imag(X[k])*s
		}
		x[n] = re / float64(N)
	}
	return x
}

func IFFT(X []complex128) []float64 {
	N := len(X)
	if N <= 1 {
		result := make([]float64, N)
		for i := range X {
			result[i] = real(X[i]) / float64(N)
		}
		return result
	}

	if N&(N-1) != 0 {
		return IDFT(X)
	}

	even := make([]complex128, N/2)
	odd := make([]complex128, N/2)
	for i := 0; i < N/2; i++ {
		even[i] = X[2*i]
		odd[i] = X[2*i+1]
	}

	evenIFFT := IFFTRecursive(even)
	oddIFFT := IFFTRecursive(odd)

	x := make([]float64, N)
	twoPi := 2 * math.Pi
	for k := 0; k < N/2; k++ {
		angle := twoPi * float64(k) / float64(N)
		twiddle := cmplx.Exp(complex(0, angle))
		t := oddIFFT[k] * twiddle
		combined := evenIFFT[k] + t
		x[k] = real(combined) / float64(N)
		x[k+N/2] = real(evenIFFT[k]-t) / float64(N)
	}

	return x
}

func IFFTRecursive(X []complex128) []complex128 {
	N := len(X)
	if N <= 1 {
		return X
	}

	even := make([]complex128, N/2)
	odd := make([]complex128, N/2)
	for i := 0; i < N/2; i++ {
		even[i] = X[2*i]
		odd[i] = X[2*i+1]
	}

	evenIFFT := IFFTRecursive(even)
	oddIFFT := IFFTRecursive(odd)

	result := make([]complex128, N)
	twoPi := 2 * math.Pi
	for k := 0; k < N/2; k++ {
		angle := twoPi * float64(k) / float64(N)
		twiddle := cmplx.Exp(complex(0, angle))
		t := oddIFFT[k] * twiddle
		result[k] = evenIFFT[k] + t
		result[k+N/2] = evenIFFT[k] - t
	}

	return result
}
