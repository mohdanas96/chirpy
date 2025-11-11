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
	"github.com/mohdanas96/chirpy/internal/database"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	dbQ            *database.Queries
}

func main() {
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		fmt.Printf("Cannot open connection to DB: %v", err)
		return
	}

	dbQueries := database.New(db)

	filepathRoot := "."
	port := ":8080"

	apiConfig := apiConfig{fileServerHits: atomic.Int32{}, dbQ: dbQueries}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiConfig.middlewareMetricsInc((http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiConfig.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiConfig.handlerResetMetrics)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	server := &http.Server{Addr: port, Handler: mux}

	log.Printf("Serving files from %s on port %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
