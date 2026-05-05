package dsp

import (
	"encoding/json"
	"os"
	"time"
)

func logFingerprintDebug(hypothesisID, location, message string, data map[string]interface{}) {
	payload := map[string]interface{}{
		"sessionId":    "042fea",
		"runId":        "pre-fix",
		"hypothesisId": hypothesisID,
		"location":     location,
		"message":      message,
		"data":         data,
		"timestamp":    time.Now().UnixMilli(),
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return
	}
	f, err := os.OpenFile("/Users/dakshaggarwal/Desktop/Shazam-Vscode copy/.cursor/debug-042fea.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	_, _ = f.Write(append(b, '\n'))
}

// Fingerprint is a hash + the absolute time offset (in frames) of the anchor peak.
type Fingerprint struct {
	Hash   uint32
	Offset int
}

// GenerateFingerprints creates combinatorial hashes from a list of peaks.
// Uses the classic Shazam technique:
//   - Sort peaks by time.
//   - For each peak (anchor), pair it with up to fanValue next peaks (targets).
//   - For each pair, compute a 32-bit hash from anchor frequency, target frequency,
//     and time delta.
//
// Parameters:
//   - peaks: slice of Peak structs
//   - fanValue: how many target peaks to pair with each anchor (e.g., 15)
//   - maxDelta: maximum allowed time difference in frames (e.g., 200)
//
// Returns a slice of Fingerprints.
func GenerateFingerprints(peaks []Peak, fanValue, maxDelta int) []Fingerprint {
	if len(peaks) < 2 {
		return nil
	}
	skippedForDelta := 0

	// Sort peaks by time (ascending)
	sorted := make([]Peak, len(peaks))
	copy(sorted, peaks)
	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].Time < sorted[i].Time {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	fingerprints := []Fingerprint{}

	for i := 0; i < len(sorted); i++ {
		anchor := sorted[i]
		// Pair with next fanValue peaks
		for j := 1; j <= fanValue && i+j < len(sorted); j++ {
			target := sorted[i+j]
			delta := target.Time - anchor.Time
			if delta <= 0 || delta > maxDelta {
				continue
			}
			if delta > 255 {
				skippedForDelta++
				continue
			}

			// 32-bit hash:
			// Bits 31-20: anchor frequency (12 bits, enough for up to 4096 bins)
			// Bits 19-8:  target frequency  (12 bits)
			// Bits 7-0:   time delta        (8 bits, up to 255 frames)
			hash := (uint32(anchor.FreqBin) << 20) |
				(uint32(target.FreqBin) << 8) |
				uint32(delta)

			fingerprints = append(fingerprints, Fingerprint{
				Hash:   hash,
				Offset: anchor.Time,
			})
		}
	}
	// #region agent log
	logFingerprintDebug("H3", "fingerprint.go:GenerateFingerprints", "Fingerprint generation stats", map[string]interface{}{"peakCount": len(peaks), "fingerprintCount": len(fingerprints), "skippedDeltaOutOf8Bit": skippedForDelta, "maxDelta": maxDelta})
	// #endregion
	return fingerprints
}
