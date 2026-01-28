package internal

import "lab1/pkg/lab1"

type Service interface {
	GenerateSignals() (x, y lab1.Signal)
	Convolution(x, y lab1.Signal) lab1.Signal
	Correlation(x, y lab1.Signal) lab1.Signal
	DFT(sig lab1.Signal) []lab1.Complex
	IDFT(spec []lab1.Complex, sampleRate float64) lab1.Signal
}

type service struct{}

func NewService() Service {
	return &service{}
}

func (s *service) GenerateSignals() (x, y lab1.Signal) {
	common := lab1.HarmonicParams{
		Amplitudes: []float64{0.8, 0.5, 0.3},
		BaseFreq:   330,
		Harmonics:  []float64{1, 2, 3},
		Duration:   0.05,
		SampleRate: 44100,
	}
	px := common
	px.Phi = 0
	py := common
	py.Phi = lab1.PiOver2()
	return lab1.GenerateHarmonicSignal(px), lab1.GenerateHarmonicSignal(py)
}

func (s *service) Convolution(x, y lab1.Signal) lab1.Signal {
	out := lab1.Convolve(x.Samples, y.Samples)
	return lab1.Signal{
		Samples:    out,
		SampleRate: x.SampleRate,
	}
}

func (s *service) Correlation(x, y lab1.Signal) lab1.Signal {
	out := lab1.Correlate(x.Samples, y.Samples)
	return lab1.Signal{
		Samples:    out,
		SampleRate: x.SampleRate,
	}
}

func (s *service) DFT(sig lab1.Signal) []lab1.Complex {
	return lab1.DFT(sig.Samples)
}

func (s *service) IDFT(spec []lab1.Complex, sampleRate float64) lab1.Signal {
	samples := lab1.IDFT(spec)
	return lab1.Signal{
		Samples:    samples,
		SampleRate: sampleRate,
	}
}
