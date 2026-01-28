package lab1

import "math"

// Complex is a minimal complex number implementation to avoid external deps.
type Complex struct {
	Re float64
	Im float64
}

func (c Complex) Abs() float64 {
	return math.Hypot(c.Re, c.Im)
}

func (c Complex) Phase() float64 {
	return math.Atan2(c.Im, c.Re)
}

func DFT(x []float64) []Complex {
	N := len(x)
	X := make([]Complex, N)
	twoPi := 2 * math.Pi
	for k := 0; k < N; k++ {
		var re, im float64
		for n := 0; n < N; n++ {
			angle := -twoPi * float64(k*n) / float64(N)
			re += x[n] * math.Cos(angle)
			im += x[n] * math.Sin(angle)
		}
		X[k] = Complex{Re: re, Im: im}
	}
	return X
}

func IDFT(X []Complex) []float64 {
	N := len(X)
	x := make([]float64, N)
	twoPi := 2 * math.Pi

	for n := 0; n < N; n++ {
		var re, im float64
		for k := 0; k < N; k++ {
			angle := twoPi * float64(k*n) / float64(N)
			c := math.Cos(angle)
			s := math.Sin(angle)
			re += X[k].Re*c - X[k].Im*s
			im += X[k].Re*s + X[k].Im*c
		}
		x[n] = re / float64(N)
		_ = im
	}
	return x
}
