package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	if code > 499 {
		log.Printf("Responding with 5xx error: %v", msg)
	}
	type errorReponse struct {
		Error string `json:"error"`
	}

	respondWithJson(w, code, errorReponse{Error: msg})

}

func respondWithJson(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling json: %v", err)
	}

	w.WriteHeader(code)
	w.Write(data)
}
