package main

import (
	"context"
	"lab1/pkg/signal"
	"lab2/internal"
	"math/cmplx"
	"os"
	"path/filepath"
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
	phase := make([]float64, N)
	df := sampleRate / float64(N)
	for k := 0; k < N; k++ {
		freqs[k] = float64(k) * df
		mag[k] = cmplx.Abs(X[k])
		phase[k] = cmplx.Phase(X[k])
	}
	return internal.SpectrumDTO{
		Freqs: freqs,
		Mag:   mag,
		Phase: phase,
	}
}

// FilterResult DTO for filter analysis results
type FilterResult struct {
	Filtered      internal.SignalDTO   `json:"filtered"`
	Freqs         []float64            `json:"freqs"`
	Magnitude     []float64            `json:"magnitude"`
	InputSignal   internal.SignalDTO   `json:"inputSignal"`
	InputSpectrum internal.SpectrumDTO `json:"inputSpectrum"`
}

// ApplyFilters applies all three filters and returns the results
func (a *App) ApplyFilters(N int) map[string]FilterResult {
	// Generate signal with noise
	inputSig := a.dspServ.GenerateSignalWithNoise(N)

	// Calculate input signal spectrum
	spec := a.dspServ.DFTLib(inputSig)
	inputSpectrum := toSpectrumDTO(spec, inputSig.SampleRate)

	// Apply Moving Average Filter (M=15)
	maFiltered := a.dspServ.ApplyMovingAverageFilter(inputSig, 15)
	_ = a.dspServ.DFTLib(maFiltered) // for potential future use

	// Apply FIR Band Stop Filter (653-667 Hz, M=201)
	firFiltered, _, firFreqs, firMag := a.dspServ.ApplyFIRBandStopFilter(inputSig, 653, 667, 201)
	_ = a.dspServ.DFTLib(firFiltered) // for potential future use

	// Apply IIR Band Pass Filter (f0=400 Hz, BW=80 Hz)
	iirFiltered, iirFreqs, iirMag := a.dspServ.ApplyIIRBandPassFilter(inputSig, 400, 80)
	_ = a.dspServ.DFTLib(iirFiltered) // for potential future use

	return map[string]FilterResult{
		"movingAverage": {
			Filtered:      toSignalDTO(maFiltered),
			InputSignal:   toSignalDTO(inputSig),
			InputSpectrum: inputSpectrum,
			Freqs:         []float64{}, // MA doesn't have freq response visualization
			Magnitude:     []float64{},
		},
		"firBandStop": {
			Filtered:      toSignalDTO(firFiltered),
			InputSignal:   toSignalDTO(inputSig),
			InputSpectrum: inputSpectrum,
			Freqs:         firFreqs,
			Magnitude:     firMag,
		},
		"iirBandPass": {
			Filtered:      toSignalDTO(iirFiltered),
			InputSignal:   toSignalDTO(inputSig),
			InputSpectrum: inputSpectrum,
			Freqs:         iirFreqs,
			Magnitude:     iirMag,
		},
	}
}

// SaveWAVFile saves signal samples to a WAV file
func (a *App) saveWAVFile(filename string, samples []float64, sampleRate int) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Normalize samples to int16 range
	var maxSample float64
	for _, s := range samples {
		if s < 0 {
			s = -s
		}
		if s > maxSample {
			maxSample = s
		}
	}

	if maxSample == 0 {
		maxSample = 1
	}

	// WAV header parameters
	numChannels := 1
	bitsPerSample := 16
	blockAlign := numChannels * bitsPerSample / 8
	byteRate := sampleRate * blockAlign

	// Calculate sizes
	subchunk2Size := len(samples) * blockAlign
	chunkSize := 36 + subchunk2Size

	// Write RIFF header
	file.WriteString("RIFF")
	var header [4]byte
	header[0] = byte(chunkSize)
	header[1] = byte(chunkSize >> 8)
	header[2] = byte(chunkSize >> 16)
	header[3] = byte(chunkSize >> 24)
	file.Write(header[:])
	file.WriteString("WAVE")

	// Write fmt subchunk
	file.WriteString("fmt ")
	subchunk1size := uint32(16)
	file.Write([]byte{byte(subchunk1size), byte(subchunk1size >> 8), byte(subchunk1size >> 16), byte(subchunk1size >> 24)})

	audioFormat := uint16(1)
	file.Write([]byte{byte(audioFormat), byte(audioFormat >> 8)})

	channels := uint16(numChannels)
	file.Write([]byte{byte(channels), byte(channels >> 8)})

	sampleRateVal := uint32(sampleRate)
	file.Write([]byte{byte(sampleRateVal), byte(sampleRateVal >> 8), byte(sampleRateVal >> 16), byte(sampleRateVal >> 24)})

	byteRateVal := uint32(byteRate)
	file.Write([]byte{byte(byteRateVal), byte(byteRateVal >> 8), byte(byteRateVal >> 16), byte(byteRateVal >> 24)})

	blockAlignVal := uint16(blockAlign)
	file.Write([]byte{byte(blockAlignVal), byte(blockAlignVal >> 8)})

	bitsVal := uint16(bitsPerSample)
	file.Write([]byte{byte(bitsVal), byte(bitsVal >> 8)})

	// Write data subchunk
	file.WriteString("data")
	dataSizeVal := uint32(subchunk2Size)
	file.Write([]byte{byte(dataSizeVal), byte(dataSizeVal >> 8), byte(dataSizeVal >> 16), byte(dataSizeVal >> 24)})

	// Write samples
	for _, sample := range samples {
		// Normalize to int16 range
		val := int16((sample / maxSample) * 32767)
		file.Write([]byte{byte(val), byte(val >> 8)})
	}

	return nil
}

// ExportFilterResults exports all filter results to WAV files
func (a *App) ExportFilterResults(N int) map[string]string {
	// Create output directory if it doesn't exist
	outputDir := "output"
	os.MkdirAll(outputDir, 0755)

	// Generate signal with noise
	inputSig := a.dspServ.GenerateSignalWithNoise(N)

	// Apply filters
	maFiltered := a.dspServ.ApplyMovingAverageFilter(inputSig, 15)
	firFiltered, _, _, _ := a.dspServ.ApplyFIRBandStopFilter(inputSig, 653, 667, 201)
	iirFiltered, _, _ := a.dspServ.ApplyIIRBandPassFilter(inputSig, 400, 80)

	// Get sample rate as integer
	sampleRate := int(inputSig.SampleRate)

	// Save to WAV files
	results := make(map[string]string)

	inputPath := filepath.Join(outputDir, "input_with_noise.wav")
	a.saveWAVFile(inputPath, inputSig.Samples, sampleRate)
	results["input"] = inputPath

	maPath := filepath.Join(outputDir, "moving_average_filtered.wav")
	a.saveWAVFile(maPath, maFiltered.Samples, sampleRate)
	results["moving_average"] = maPath

	firPath := filepath.Join(outputDir, "fir_bandstop_filtered.wav")
	a.saveWAVFile(firPath, firFiltered.Samples, sampleRate)
	results["fir_bandstop"] = firPath

	iirPath := filepath.Join(outputDir, "iir_bandpass_filtered.wav")
	a.saveWAVFile(iirPath, iirFiltered.Samples, sampleRate)
	results["iir_bandpass"] = iirPath

	return results
}
