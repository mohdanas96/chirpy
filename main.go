package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

func main() {
	filepathRoot := "."
	port := ":8080"

	apiConfig := apiConfig{fileServerHits: atomic.Int32{}}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiConfig.middlewareMetricsInc((http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))))
	mux.HandleFunc("/healthz", handlerReadiness)
	mux.HandleFunc("/metrics", apiConfig.handlerMetrics)
	mux.HandleFunc("/reset", apiConfig.handlerResetMetrics)

	server := &http.Server{Addr: port, Handler: mux}

	log.Printf("Serving files from %s on port %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
