package internal

import (
	"lab1/pkg/conv"
	"lab1/pkg/corr"
	"lab1/pkg/fft"
	"lab1/pkg/signal"
	"math"
)

type Service interface {
	GenerateSignals(N int) (x, y signal.Signal)
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

func (s *service) GenerateSignals(N int) (x, y signal.Signal) {
	if N <= 0 {
		N = 256
	}
	sampleRate := float64(N) * 10
	duration := 1.0

	common := signal.HarmonicParams{
		Amplitudes: []float64{0.8, 0.5, 0.3},
		BaseFreq:   50,
		Harmonics:  []float64{1, 2, 3},
		Duration:   duration,
		SampleRate: sampleRate,
	}
	px := common
	px.Phi = 0
	py := common
	py.Phi = math.Pi / 2
	xSig := signal.GenerateHarmonicSignal(px)
	ySig := signal.GenerateHarmonicSignal(py)

	if len(xSig.Samples) > N {
		xSig.Samples = xSig.Samples[:N]
		ySig.Samples = ySig.Samples[:N]
	}
	return xSig, ySig
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
	return fft.FFT(sig.Samples)
}

func (s *service) IDFT(spec []complex128, sampleRate float64) signal.Signal {
	samples := fft.IFFT(spec)
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
