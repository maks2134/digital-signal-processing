package main

import (
	"context"
	"fmt"
	"lab1/internal"

	"lab1/pkg/lab1"
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

func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) Analyze() internal.AnalysisResult {
	x, y := a.dspServ.GenerateSignals()
	conv := a.dspServ.Convolution(x, y)
	corr := a.dspServ.Correlation(x, y)
	specX := a.dspServ.DFT(x)
	specY := a.dspServ.DFT(y)
	specConv := a.dspServ.DFT(conv)
	return internal.AnalysisResult{
		X:            toSignalDTO(x),
		Y:            toSignalDTO(y),
		Conv:         toSignalDTO(conv),
		Corr:         toSignalDTO(corr),
		SpectrumX:    toSpectrumDTO(specX, x.SampleRate),
		SpectrumY:    toSpectrumDTO(specY, y.SampleRate),
		SpectrumConv: toSpectrumDTO(specConv, conv.SampleRate),
	}
}

func toSignalDTO(s lab1.Signal) internal.SignalDTO {
	return internal.SignalDTO{
		Samples:    s.Samples,
		SampleRate: s.SampleRate,
	}
}

func toSpectrumDTO(X []lab1.Complex, sampleRate float64) internal.SpectrumDTO {
	N := len(X)
	freqs := make([]float64, N)
	mag := make([]float64, N)
	phase := make([]float64, N)
	df := sampleRate / float64(N)
	for k := 0; k < N; k++ {
		freqs[k] = float64(k) * df
		mag[k] = X[k].Abs()
		phase[k] = X[k].Phase()
	}
	return internal.SpectrumDTO{
		Freqs: freqs,
		Mag:   mag,
		Phase: phase,
	}
}
