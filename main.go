package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mhmdanas10/chirpy/internal/database"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	dbQ            *database.Queries
	platform       string
}

func main() {
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")

	if platform == "" {
		log.Fatal("Platform must be set")
	}

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		fmt.Printf("Cannot open connection to DB: %v", err)
		return
	}

	dbQueries := database.New(db)

	filepathRoot := "."
	port := ":8080"

	apiConfig := apiConfig{fileServerHits: atomic.Int32{}, dbQ: dbQueries, platform: platform}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiConfig.middlewareMetricsInc((http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiConfig.handlerMetrics)
	mux.HandleFunc("POST /admin/reset_metrics", apiConfig.handlerResetMetrics)
	mux.HandleFunc("POST /admin/reset", apiConfig.handlerResetUsers)
	mux.HandleFunc("POST /api/users", apiConfig.handlerCreateUser)
	mux.HandleFunc("POST /api/chirps", apiConfig.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiConfig.GetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiConfig.GetChirp)

	server := &http.Server{Addr: port, Handler: mux}

	log.Printf("Serving files from %s on port %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
