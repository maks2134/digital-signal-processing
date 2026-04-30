package main

import (
	"context"
	"lab1/internal"
	"lab1/pkg/signal"
	"math/cmplx"
)

type App struct {
	ctx     context.Context
	dspServ internal.Service
}

func NewApp() *App {
	return &App{
		dspServ: internal.NewService(),
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) Analyze(N int) internal.AnalysisResult {
	x, y := a.dspServ.GenerateSignals(N)
	conv := a.dspServ.Convolution(x, y)
	corr := a.dspServ.Correlation(x, y)
	specX := a.dspServ.DFT(x)
	specY := a.dspServ.DFT(y)
	specConv := a.dspServ.DFT(conv)
	specCorr := a.dspServ.DFT(corr)

	idftX := a.dspServ.IDFT(specX, x.SampleRate)
	idftY := a.dspServ.IDFT(specY, y.SampleRate)
	idftConv := a.dspServ.IDFT(specConv, conv.SampleRate)

	return internal.AnalysisResult{
		X:            toSignalDTO(x),
		Y:            toSignalDTO(y),
		Conv:         toSignalDTO(conv),
		Corr:         toSignalDTO(corr),
		SpectrumX:    toSpectrumDTO(specX, x.SampleRate),
		SpectrumY:    toSpectrumDTO(specY, y.SampleRate),
		SpectrumConv: toSpectrumDTO(specConv, conv.SampleRate),
		SpectrumCorr: toSpectrumDTO(specCorr, corr.SampleRate),
		IDFTX:        toSignalDTO(idftX),
		IDFTY:        toSignalDTO(idftY),
		IDFTConv:     toSignalDTO(idftConv),
	}
}

func (a *App) AnalyzeLib(N int) internal.AnalysisResult {
	x, y := a.dspServ.GenerateSignals(N)
	conv := a.dspServ.ConvolutionLib(x, y)
	corr := a.dspServ.CorrelationLib(x, y)
	specX := a.dspServ.DFTLib(x)
	specY := a.dspServ.DFTLib(y)
	specConv := a.dspServ.DFTLib(conv)
	specCorr := a.dspServ.DFTLib(corr)

	idftX := a.dspServ.IDFTLib(specX, x.SampleRate)
	idftY := a.dspServ.IDFTLib(specY, y.SampleRate)
	idftConv := a.dspServ.IDFTLib(specConv, conv.SampleRate)

	return internal.AnalysisResult{
		X:            toSignalDTO(x),
		Y:            toSignalDTO(y),
		Conv:         toSignalDTO(conv),
		Corr:         toSignalDTO(corr),
		SpectrumX:    toSpectrumDTO(specX, x.SampleRate),
		SpectrumY:    toSpectrumDTO(specY, y.SampleRate),
		SpectrumConv: toSpectrumDTO(specConv, conv.SampleRate),
		SpectrumCorr: toSpectrumDTO(specCorr, corr.SampleRate),
		IDFTX:        toSignalDTO(idftX),
		IDFTY:        toSignalDTO(idftY),
		IDFTConv:     toSignalDTO(idftConv),
	}
}

func toSignalDTO(s signal.Signal) internal.SignalDTO {
	return internal.SignalDTO{
		Samples:    s.Samples,
		SampleRate: s.SampleRate,
	}
}

func toSpectrumDTO(X []complex128, sampleRate float64) internal.SpectrumDTO {
	N := len(X)
	freqs := make([]float64, N)
	mag := make([]float64, N)
	df := sampleRate / float64(N)
	for k := 0; k < N; k++ {
		freqs[k] = float64(k) * df
		mag[k] = cmplx.Abs(X[k])
	}
	return internal.SpectrumDTO{
		Freqs: freqs,
		Mag:   mag,
	}
}
