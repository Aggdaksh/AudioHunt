package songs

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

type Handler struct {
	service *Service
}

func NewHandler() *Handler {
	return &Handler{service: NewService()}
}

func RegisterRoutes(r *mux.Router) {
	h := NewHandler()
	r.HandleFunc("/api/songs", h.AddSong).Methods("POST")
	r.HandleFunc("/api/recognize", h.Recognize).Methods("POST")

}

func (h *Handler) AddSong(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 50<<20)

	if err := r.ParseMultipartForm(50 << 20); err != nil {
		http.Error(w, "File too large or form parse error", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("audio_file")
	if err != nil {
		http.Error(w, "Missing audio_file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	title := r.FormValue("title")
	artist := r.FormValue("artist")
	if title == "" || artist == "" {
		http.Error(w, "title and artist are required", http.StatusBadRequest)
		return
	}

	uploadDir := "./uploads"
	os.MkdirAll(uploadDir, os.ModePerm)
	filePath := filepath.Join(uploadDir, header.Filename)

	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	io.Copy(dst, file)

	songID, err := h.service.AddSong(title, artist, filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      songID,
		"message": "Song added successfully",
	})
}

func (h *Handler) Recognize(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("audio_snippet")
	if err != nil {
		http.Error(w, "Missing audio_snippet", http.StatusBadRequest)
		return
	}
	defer file.Close()

	tmpFile, err := os.CreateTemp("", "snippet-*.wav")
	if err != nil {
		http.Error(w, "Failed to create temp file", http.StatusInternalServerError)
		return
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	io.Copy(tmpFile, file)

	song, confidence, err := h.service.Recognize(tmpFile.Name())
	if err != nil {
		log.Printf("!!! RECOGNIZE ERROR: %v", err)
		if err.Error() == "no fingerprints extracted — audio too short or too quiet" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err.Error() == "match confidence too low" || err.Error() == "no match found" || err.Error() == "no matching fingerprints found in catalog" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"song_id":    song.ID,
		"title":      song.Title,
		"artist":     song.Artist,
		"confidence": confidence,
	})
}
