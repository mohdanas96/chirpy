package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mohdanas96/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	type reqParams struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	type response struct {
		Chirp
	}

	decoder := json.NewDecoder(req.Body)

	params := reqParams{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode body", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	cleanBody := getCleanedBody(params.Body, badWords)

	chirp, err := cfg.dbQ.CreateChirp(context.Background(), database.CreateChirpParams{UserID: params.UserID, Body: cleanBody})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while creating chirp", err)
		return
	}

	respondWithJson(w, http.StatusCreated, response{
		Chirp: Chirp{ID: chirp.ID, CreatedAt: chirp.CreatedAt, UpdatedAt: chirp.UpdatedAt, Body: chirp.Body, UserId: chirp.UserID},
	})
}

func (cfg *apiConfig) GetAllChirps(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	fmt.Println("Working")

	chirps, err := cfg.dbQ.GetAllChirps(context.Background())
	if err != nil {
		respondWithError(w, 501, "Error while retrieving chirps", err)
		return
	}

	fmt.Println(chirps)

	data := make([]Chirp, len(chirps))
	for i := range chirps {
		fmt.Println(len(chirps), "LENENNE")
		data[i].ID = chirps[i].ID
		data[i].Body = chirps[i].Body
		data[i].CreatedAt = chirps[i].CreatedAt
		data[i].UpdatedAt = chirps[i].UpdatedAt
		data[i].UserId = chirps[i].UserID
	}

	respondWithJson(w, 200, data)
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}
