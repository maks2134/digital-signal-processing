package corr

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
