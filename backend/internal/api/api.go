package api

import (
	"log"
	"net/http"

	"Shazam-Vscode/backend/internal/songs"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func NewRouter() http.Handler {
	r := mux.NewRouter().StrictSlash(true)

	// Logging middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			log.Printf(">> INCOMING: %s %s", req.Method, req.URL.Path)
			next.ServeHTTP(w, req)
		})
	})

	songs.RegisterRoutes(r)

	// Debug: print registered routes
	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()
		log.Printf("Registered route: %s %v", path, methods)
		return nil
	})

	// Catch-all for unmatched routes
	// Catch-all for debugging
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		if req.Method == "POST" && req.URL.Path == "/api/recognize" {
			log.Println(">>> Manually calling Recognize handler from catch-all")
			h := songs.NewHandler() // create handler instance
			h.Recognize(w, req)
			return
		}
		log.Printf("!!! UNMATCHED: %s %s", req.Method, req.URL.Path)
		http.NotFound(w, req)
	})

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true, // Note: some browsers don't like wildcard origin with true credentials, but rs/cors handles '*' by reflecting the origin.
	})

	return c.Handler(r)
}
