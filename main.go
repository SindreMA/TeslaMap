package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/sindrema/teslamap/internal/config"
	"github.com/sindrema/teslamap/internal/db"
	"github.com/sindrema/teslamap/internal/handler"
)

//go:embed all:frontend/dist
var frontendFS embed.FS

func main() {
	cfg := config.Load()
	pool := db.Connect(cfg)
	defer pool.Close()

	h := &handler.Handler{DB: pool}

	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("GET /api/cars", h.ListCars)
	mux.HandleFunc("GET /api/cars/{id}", h.GetCarPosition)
	mux.HandleFunc("GET /api/health", h.Health)

	// Default redirect
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, fmt.Sprintf("/car/%d", cfg.DefaultCarID), http.StatusFound)
			return
		}
		// Serve SPA
		serveSPA(w, r)
	})

	log.Printf("Listening on :%s (default car: %d)", cfg.Port, cfg.DefaultCarID)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatal(err)
	}
}

func serveSPA(w http.ResponseWriter, r *http.Request) {
	distFS, err := fs.Sub(frontendFS, "frontend/dist")
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "" {
		path = "index.html"
	}

	// Try to serve the file directly
	f, err := distFS.Open(path)
	if err == nil {
		f.Close()
		http.FileServerFS(distFS).ServeHTTP(w, r)
		return
	}

	// Fall back to index.html for SPA routing
	indexFile, err := fs.ReadFile(distFS, "index.html")
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(indexFile)
}

// Ensure frontend/dist exists at compile time for embed.
// During development, create a placeholder.
func init() {
	if _, err := os.Stat("frontend/dist"); os.IsNotExist(err) {
		// This only matters during development without a built frontend
		_ = os.MkdirAll("frontend/dist", 0o755)
	}
}
