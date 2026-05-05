package audio

import (
	"Shazam-Vscode/backend/internal/audio/dsp"
)

// FingerprintConfig holds parameters for the fingerprinting pipeline.
type FingerprintConfig struct {
	// Spectrogram parameters
	WindowSize int
	HopSize    int

	// Peak picking
	PeakNeighborhood int
	PeakMinAmplitude float64

	// Hashing
	FanValue int
	MaxDelta int
}

// DefaultConfig returns sensible values used by Shazam-like systems.
func DefaultConfig() FingerprintConfig {
	return FingerprintConfig{
		WindowSize:       2048,
		HopSize:          512,
		PeakNeighborhood: 10,
		PeakMinAmplitude: 20.0,
		FanValue:         15,
		MaxDelta:         200,
	}
}

// ExtractFingerprints reads an audio file and returns its fingerprints.
// It assumes the file is WAV; for other formats you'd convert first.
func ExtractFingerprints(filePath string, cfg FingerprintConfig) ([]dsp.Fingerprint, error) {
	// Decode
	samples, sr, err := dsp.DecodeWav(filePath)
	if err != nil {
		return nil, err
	}

	// (Optional: resample if needed; we skip for simplicity)

	// Spectrogram
	spec := dsp.GenerateSpectrogram(samples, sr, cfg.WindowSize, cfg.HopSize)
	if spec == nil {
		return nil, nil
	}

	// Peaks
	peaks := dsp.FindPeaks(spec, cfg.PeakNeighborhood, cfg.PeakMinAmplitude)

	// Fingerprints
	fingerprints := dsp.GenerateFingerprints(peaks, cfg.FanValue, cfg.MaxDelta)

	return fingerprints, nil
}
