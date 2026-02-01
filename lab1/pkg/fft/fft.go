package fft

import (
	"math"
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
