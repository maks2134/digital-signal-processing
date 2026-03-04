package internal

import (
	"lab1/pkg/conv"
	"lab1/pkg/corr"
	"lab1/pkg/fft"
	"lab1/pkg/signal"
	"lab2/pkg/filter"
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
	ConvolutionLib(x, y signal.Signal) signal.Signal
	CorrelationLib(x, y signal.Signal) signal.Signal

	// Filter methods
	GenerateSignalWithNoise(N int) signal.Signal
	ApplyMovingAverageFilter(sig signal.Signal, M int) signal.Signal
	ApplyFIRBandStopFilter(sig signal.Signal, f1, f2 float64, M int) (signal.Signal, []float64, []float64, []float64)
	ApplyIIRBandPassFilter(sig signal.Signal, f0, bw float64) (signal.Signal, []float64, []float64)
	GetFilterFrequencyResponse(coeffs interface{}, fs float64, nFreqs int) ([]float64, []float64)
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
		BaseFreq:   330,
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

func (s *service) ConvolutionLib(x, y signal.Signal) signal.Signal {
	out := conv.ConvolveLib(x.Samples, y.Samples)
	return signal.Signal{
		Samples:    out,
		SampleRate: x.SampleRate,
	}
}

func (s *service) CorrelationLib(x, y signal.Signal) signal.Signal {
	out := corr.CorrelateLib(x.Samples, y.Samples)
	return signal.Signal{
		Samples:    out,
		SampleRate: x.SampleRate,
	}
}

// GenerateSignalWithNoise - генерирует чистый сигнал 330 Гц с помехами в полосе [653,667] Гц
func (s *service) GenerateSignalWithNoise(N int) signal.Signal {
	if N <= 0 {
		N = 2560
	}
	sampleRate := float64(N) * 10
	duration := 1.0

	// Основной сигнал 330 Гц
	params := signal.HarmonicParams{
		Amplitudes: []float64{1.0, 0.5, 0.3},
		BaseFreq:   330,
		Harmonics:  []float64{1, 2, 3},
		Duration:   duration,
		SampleRate: sampleRate,
		Phi:        0,
	}

	sig := signal.GenerateHarmonicSignal(params)

	// Добавляем помехи (узкополосный шум в диапазоне [653,667] Гц)
	nSamples := len(sig.Samples)
	for n := 0; n < nSamples; n++ {
		t := float64(n) / sampleRate
		// Помеха на частоте 660 Гц (центр диапазона [653,667])
		noise := 0.3 * math.Sin(2*math.Pi*660*t)
		sig.Samples[n] += noise
	}

	if len(sig.Samples) > N {
		sig.Samples = sig.Samples[:N]
	}
	return sig
}

// ApplyMovingAverageFilter - применяет однородный фильтр
func (s *service) ApplyMovingAverageFilter(sig signal.Signal, M int) signal.Signal {
	filtered := filter.MovingAverageFilter(sig.Samples, M)
	return signal.Signal{
		Samples:    filtered,
		SampleRate: sig.SampleRate,
	}
}

// ApplyFIRBandStopFilter - применяет КИХ режекторный фильтр
func (s *service) ApplyFIRBandStopFilter(sig signal.Signal, f1, f2 float64, M int) (signal.Signal, []float64, []float64, []float64) {
	coeffs := filter.DesignFIRBandStop(f1, f2, sig.SampleRate, M)
	filtered := filter.FIRFilter(sig.Samples, coeffs)

	// Вычисляем частотную характеристику
	freqs, magnitude := filter.ComputeFrequencyResponse(coeffs, sig.SampleRate, 1024)

	return signal.Signal{
		Samples:    filtered,
		SampleRate: sig.SampleRate,
	}, coeffs, freqs, magnitude
}

// ApplyIIRBandPassFilter - применяет БИХ полосовой фильтр
func (s *service) ApplyIIRBandPassFilter(sig signal.Signal, f0, bw float64) (signal.Signal, []float64, []float64) {
	coeffs := filter.DesignBandPassIIR(f0, bw, sig.SampleRate)
	filtered := filter.IIRFilter(sig.Samples, coeffs)

	// Вычисляем частотную характеристику
	freqs, magnitude := filter.ComputeFrequencyResponse(coeffs, sig.SampleRate, 1024)

	return signal.Signal{
		Samples:    filtered,
		SampleRate: sig.SampleRate,
	}, freqs, magnitude
}

// GetFilterFrequencyResponse - получает частотную характеристику для display
func (s *service) GetFilterFrequencyResponse(coeffs interface{}, fs float64, nFreqs int) ([]float64, []float64) {
	return filter.ComputeFrequencyResponse(coeffs, fs, nFreqs)
}
