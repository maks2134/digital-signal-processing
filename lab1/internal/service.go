package internal

import (
	"lab1/pkg/conv"
	"lab1/pkg/corr"
	"lab1/pkg/fft"
	"lab1/pkg/signal"
	"math"
)

type Service interface {
	GenerateSignals() (x, y signal.Signal)
	Convolution(x, y signal.Signal) signal.Signal
	Correlation(x, y signal.Signal) signal.Signal
	DFT(sig signal.Signal) []complex128
	IDFT(spec []complex128, sampleRate float64) signal.Signal
	DFTLib(sig signal.Signal) []complex128
	IDFTLib(spec []complex128, sampleRate float64) signal.Signal
}

type service struct{}

func NewService() Service {
	return &service{}
}

func (s *service) GenerateSignals() (x, y signal.Signal) {
	common := signal.HarmonicParams{
		Amplitudes: []float64{0.8, 0.5, 0.3},
		BaseFreq:   330,
		Harmonics:  []float64{1, 2, 3},
		Duration:   0.05,
		SampleRate: 44100,
	}
	px := common
	px.Phi = 0
	py := common
	py.Phi = math.Pi / 2
	return signal.GenerateHarmonicSignal(px), signal.GenerateHarmonicSignal(py)
}

func (s *service) Convolution(x, y signal.Signal) signal.Signal {
	out := conv.Convolve(x.Samples, y.Samples)
	return signal.Signal{
		Samples:    out,
		SampleRate: x.SampleRate,
	}
}

func (s *service) Correlation(x, y signal.Signal) signal.Signal {
	out := corr.Correlate(x.Samples, y.Samples)
	return signal.Signal{
		Samples:    out,
		SampleRate: x.SampleRate,
	}
}

func (s *service) DFT(sig signal.Signal) []complex128 {
	return fft.DFT(sig.Samples)
}

func (s *service) IDFT(spec []complex128, sampleRate float64) signal.Signal {
	samples := fft.IDFT(spec)
	return signal.Signal{
		Samples:    samples,
		SampleRate: sampleRate,
	}
}

func (s *service) DFTLib(sig signal.Signal) []complex128 {
	return fft.DFTLib(sig.Samples)
}

func (s *service) IDFTLib(spec []complex128, sampleRate float64) signal.Signal {
	samples := fft.IDFTLib(spec)
	return signal.Signal{
		Samples:    samples,
		SampleRate: sampleRate,
	}
}
