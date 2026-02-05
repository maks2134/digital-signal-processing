package conv

func Convolve(x, y []float64) []float64 {
	N := len(x)
	M := len(y)
	out := make([]float64, N+M-1)
	for n := 0; n < N+M-1; n++ {
		var sum float64
		for m := 0; m < N; m++ {
			h := n - m
			if h >= 0 && h < M {
				sum += x[m] * y[h]
			}
		}
		out[n] = sum
	}
	return out
}
