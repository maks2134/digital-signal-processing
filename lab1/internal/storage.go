package internal

type SignalDTO struct {
	Samples    []float64 `json:"samples"`
	SampleRate float64   `json:"sampleRate"`
}

type SpectrumDTO struct {
	Freqs []float64 `json:"freqs"`
	Mag   []float64 `json:"mag"`
}

type AnalysisResult struct {
	X            SignalDTO   `json:"x"`
	Y            SignalDTO   `json:"y"`
	Conv         SignalDTO   `json:"conv"`
	Corr         SignalDTO   `json:"corr"`
	SpectrumX    SpectrumDTO `json:"spectrumX"`
	SpectrumY    SpectrumDTO `json:"spectrumY"`
	SpectrumConv SpectrumDTO `json:"spectrumConv"`
	SpectrumCorr SpectrumDTO `json:"spectrumCorr"`
	IDFTX        SignalDTO   `json:"idftX"`
	IDFTY        SignalDTO   `json:"idftY"`
	IDFTConv     SignalDTO   `json:"idftConv"`
}
