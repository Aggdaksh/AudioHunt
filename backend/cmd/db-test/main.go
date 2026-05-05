// package dbtest
package main

import (
	"log"

	"Shazam-Vscode/backend/internal/audio/dsp"
	"Shazam-Vscode/backend/internal/db"
)

func main() {
	dsn := "postgres://shazam_user:pass123@localhost:5434/shazam_clone?sslmode=disable"
	if err := db.Connect(dsn); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Example: fingerprint a song and insert it
	filePath := "test.wav"
	cfg := dsp.DefaultConfig()
	fingerprints, err := dsp.ExtractFingerprints(filePath, cfg)
	if err != nil {
		log.Fatal(err)
	}

	songID, err := db.InsertSong("Test Song", "Test Artist", filePath)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Inserted song with ID: %d", songID)

	err = db.InsertFingerprints(songID, fingerprints)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Inserted %d fingerprints", len(fingerprints))

	// Test query
	matches, err := db.GetMatchesForHash(fingerprints[0].Hash)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Found %d matches for first hash", len(matches))
}
