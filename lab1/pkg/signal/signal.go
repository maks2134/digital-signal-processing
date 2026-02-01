package signal

import "math"

type Signal struct {
	Samples    []float64
	SampleRate float64
}

type HarmonicParams struct {
	Amplitudes []float64
	BaseFreq   float64
	Harmonics  []float64
	Phi        float64
	Duration   float64
	SampleRate float64
}

func GenerateHarmonicSignal(p HarmonicParams) Signal {
	nSamples := int(p.Duration * p.SampleRate)
	s := make([]float64, nSamples)
	for n := 0; n < nSamples; n++ {
		t := float64(n) / p.SampleRate
		var v float64
		for i := range p.Amplitudes {
			A := p.Amplitudes[i]
			h := p.Harmonics[i]
			v += A * math.Sin(2*math.Pi*h*p.BaseFreq*t+p.Phi)
		}
		s[n] = v
	}
	return Signal{
		Samples:    s,
		SampleRate: p.SampleRate,
	}
}
