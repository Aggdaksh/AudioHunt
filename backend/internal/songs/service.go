package songs

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"Shazam-Vscode/backend/internal/audio"
	"Shazam-Vscode/backend/internal/audio/dsp"
	"Shazam-Vscode/backend/internal/db"
)

type Service struct{}

func logDebug(hypothesisID, location, message string, data map[string]interface{}) {
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

func NewService() *Service {
	return &Service{}
}

func (s *Service) AddSong(title, artist, tempFilePath string) (int, error) {
	cfg := audio.DefaultConfig()
	fingerprints, err := audio.ExtractFingerprints(tempFilePath, cfg)
	if err != nil {
		return 0, fmt.Errorf("fingerprint extraction failed: %w", err)
	}

	songID, err := db.InsertSong(title, artist, tempFilePath)
	if err != nil {
		return 0, fmt.Errorf("db insert song failed: %w", err)
	}

	err = db.InsertFingerprints(songID, fingerprints)
	if err != nil {
		return 0, fmt.Errorf("db insert fingerprints failed: %w", err)
	}

	return songID, nil
}

func (s *Service) buildHistogram(queryFPs []dsp.Fingerprint) map[int]map[int]int {
	histogram := make(map[int]map[int]int)
	for _, fp := range queryFPs {
		matches, err := db.GetMatchesForHash(fp.Hash)
		if err != nil {
			continue
		}
		for _, m := range matches {
			delta := m.TimeOffset - fp.Offset
			if _, ok := histogram[m.SongID]; !ok {
				histogram[m.SongID] = make(map[int]int)
			}
			histogram[m.SongID][delta]++
		}
	}
	return histogram
}

func findBestMatch(histogram map[int]map[int]int) (int, int) {
	bestSongID := -1
	maxCount := 0
	for songID, deltas := range histogram {
		for _, count := range deltas {
			if count > maxCount {
				maxCount = count
				bestSongID = songID
			}
		}
	}
	return bestSongID, maxCount
}

func (s *Service) Recognize(snippetPath string) (*db.Song, int, error) {
	cfg := audio.DefaultConfig()
	queryFPs, err := audio.ExtractFingerprints(snippetPath, cfg)
	if err != nil {
		return nil, 0, fmt.Errorf("fingerprint extraction failed: %w", err)
	}
	if len(queryFPs) == 0 {
		return nil, 0, fmt.Errorf("no fingerprints extracted — audio too short or too quiet")
	}
	// #region agent log
	logDebug("H2", "service.go:Recognize", "Query fingerprints extracted", map[string]interface{}{"queryFingerprintCount": len(queryFPs)})
	// #endregion

	histogram := s.buildHistogram(queryFPs)
	if len(histogram) == 0 {
		return nil, 0, fmt.Errorf("no matching fingerprints found in catalog")
	}

	bestSongID, maxCount := findBestMatch(histogram)

	if bestSongID == -1 {
		// #region agent log
		logDebug("H2", "service.go:Recognize", "No histogram winner", map[string]interface{}{"histogramSongCount": len(histogram)})
		// #endregion
		return nil, 0, fmt.Errorf("no match found")
	}
	// #region agent log
	logDebug("H2", "service.go:Recognize", "Best match selected", map[string]interface{}{"bestSongID": bestSongID, "maxVotes": maxCount})
	// #endregion

	var song db.Song
	err = db.DB.QueryRow("SELECT id, title, artist, file_path FROM songs WHERE id = $1", bestSongID).
		Scan(&song.ID, &song.Title, &song.Artist, &song.FilePath)
	if err != nil {
		return nil, 0, err
	}

	confidence := maxCount
	if confidence < 2 {
		return nil, 0, fmt.Errorf("match confidence too low")
	}
	// #region agent log
	logDebug("H2", "service.go:Recognize", "Confidence computed", map[string]interface{}{"confidence": confidence})
	// #endregion

	return &song, confidence, nil
}