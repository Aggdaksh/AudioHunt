package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"Shazam-Vscode/backend/internal/audio/dsp"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func logDBDebug(hypothesisID, location, message string, data map[string]interface{}) {
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

// Connect establishes a connection to PostgreSQL.
func Connect(dsn string) error {
	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open DB: %w", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping DB: %w", err)
	}

	log.Println("Connected to database")
	return nil
}

// Close closes the database connection.
func Close() {
	if DB != nil {
		DB.Close()
	}
}

// Song represents a row in the songs table.
type Song struct {
	ID       int
	Title    string
	Artist   string
	FilePath string
}

// InsertSong adds a song and returns its ID.
func InsertSong(title, artist, filePath string) (int, error) {
	var id int
	err := DB.QueryRow(
		`INSERT INTO songs (title, artist, file_path) 
		 VALUES ($1, $2, $3) RETURNING id`,
		title, artist, filePath,
	).Scan(&id)

	return id, err
}

// InsertFingerprints bulk inserts fingerprints for a given song.
func InsertFingerprints(songID int, fingerprints []dsp.Fingerprint) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	stmt, err := tx.Prepare(`
		INSERT INTO fingerprints (hash, song_id, time_offset) 
		VALUES ($1, $2, $3)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, fp := range fingerprints {
		_, err = stmt.Exec(int64(fp.Hash), songID, fp.Offset)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// FingerprintMatch holds a database hit for a hash.
type FingerprintMatch struct {
	SongID     int
	TimeOffset int
	Title      string
	Artist     string
}

// GetMatchesForHash returns all matching fingerprints for a given hash.
func GetMatchesForHash(hash uint32) ([]FingerprintMatch, error) {
	rows, err := DB.Query(`
		SELECT f.song_id, f.time_offset, s.title, s.artist
		FROM fingerprints f
		JOIN songs s ON f.song_id = s.id
		WHERE f.hash = $1
	`, int64(hash))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []FingerprintMatch

	for rows.Next() {
		var m FingerprintMatch
		if err := rows.Scan(&m.SongID, &m.TimeOffset, &m.Title, &m.Artist); err != nil {
			continue
		}
		matches = append(matches, m)
	}
	// #region agent log
	logDBDebug("H2", "db.go:GetMatchesForHash", "Hash lookup completed", map[string]interface{}{"hash": hash, "matchCount": len(matches)})
	// #endregion

	return matches, nil
}
