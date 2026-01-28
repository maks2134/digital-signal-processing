package lab1

func Convolve(x, y []float64) []float64 {
	N := len(x)
	M := len(y)
	out := make([]float64, N+M-1)

	for n := 0; n < N+M-1; n++ {
		var sum float64
		for k := 0; k < N; k++ {
			j := n - k
			if j >= 0 && j < M {
				sum += x[k] * y[j]
			}
		}
		out[n] = sum
	}
	return out
}

func Correlate(x, y []float64) []float64 {
	N := len(x)
	M := len(y)
	out := make([]float64, N+M-1)
	for l := -(M - 1); l <= N-1; l++ {
		var sum float64
		for n := 0; n < N; n++ {
			j := n + l
			if j >= 0 && j < M {
				sum += x[n] * y[j]
			}
		}
		out[l+M-1] = sum
	}
	return out
}
