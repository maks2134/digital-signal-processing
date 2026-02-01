package conv

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
