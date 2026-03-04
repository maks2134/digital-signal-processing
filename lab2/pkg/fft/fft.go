package fft

import (
	"math"
	"math/cmplx"
)

func DFT(x []float64) []complex128 {
	N := len(x)
	X := make([]complex128, N)
	if N == 0 {
		return X
	}

	twoPi := 2 * math.Pi
	for k := 0; k < N; k++ {
		var sum complex128
		for n := 0; n < N; n++ {
			angle := -twoPi * float64(k*n) / float64(N)
			w := cmplx.Exp(complex(0, angle))
			sum += complex(x[n], 0) * w
		}
		X[k] = sum
	}
	return X
}

func IDFT(X []complex128) []float64 {
	N := len(X)
	x := make([]float64, N)
	if N == 0 {
		return x
	}

	twoPi := 2 * math.Pi
	for n := 0; n < N; n++ {
		var sum complex128
		for k := 0; k < N; k++ {
			angle := twoPi * float64(k*n) / float64(N)
			w := cmplx.Exp(complex(0, angle))
			sum += X[k] * w
		}
		x[n] = real(sum) / float64(N)
	}
	return x
}

func FFT(x []float64) []complex128 {
	N := len(x)
	if N == 0 {
		return []complex128{}
	}
	size := 1
	for size < N {
		size <<= 1
	}
	if size != N {
		xPadded := make([]float64, size)
		copy(xPadded, x)
		x = xPadded
		N = size
	}

	if N == 1 {
		return []complex128{complex(x[0], 0)}
	}

	even := make([]float64, N/2)
	odd := make([]float64, N/2)
	for i := 0; i < N/2; i++ {
		even[i] = x[2*i]
		odd[i] = x[2*i+1]
	}

	evenFFT := FFT(even)
	oddFFT := FFT(odd)

	Y := make([]complex128, N)
	twoPi := 2 * math.Pi
	for k := 0; k < N/2; k++ {
		angle := -twoPi * float64(k) / float64(N) //степень eшки
		twiddle := cmplx.Exp(complex(0, angle))   //Wn
		t := oddFFT[k] * twiddle                  // w * bj нечет
		Y[k] = evenFFT[k] + t
		Y[k+N/2] = evenFFT[k] - t
	}

	return Y
}

func IFFT(X []complex128) []float64 {
	N := len(X)
	if N == 0 {
		return []float64{}
	}
	size := 1
	for size < N {
		size <<= 1
	}
	if size != N {
		XPadded := make([]complex128, size)
		copy(XPadded, X)
		X = XPadded
		N = size
	}

	if N == 1 {
		return []float64{real(X[0])}
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
