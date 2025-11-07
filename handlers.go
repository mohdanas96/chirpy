package main

import (
	"fmt"
	"log"
	"net/http"
)

func handlerReadiness(w http.ResponseWriter, req *http.Request) {
	body := []byte("OK")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	_, err := w.Write(body)
	if err != nil {
		log.Fatal("something went wrong while writing body :: %v", err)
	}
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hitCount := cfg.fileServerHits.Load()
	body := fmt.Sprintf("Hits: %v", hitCount)
	w.Write([]byte(body))
}

func (cfg *apiConfig) handlerResetMetrics(w http.ResponseWriter, _ *http.Request) {
	cfg.fileServerHits.Store(0)
	w.Write([]byte("Successful"))
}
