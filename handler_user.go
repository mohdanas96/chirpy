package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	type reqParams struct {
		Email string `json:"email"`
	}

	type response struct {
		User
	}

	params := reqParams{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while decoding body", err)
		return
	}

	user, err := cfg.dbQ.CreateUser(context.Background(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while creating user", err)
		return
	}

	respondWithJson(w, http.StatusCreated, response{
		User: User{ID: user.ID, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Email: user.Email},
	})
}

func (cfg *apiConfig) handlerResetUsers(w http.ResponseWriter, req *http.Request) {
	platform := cfg.platform
	if platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Forbidden", nil)
		return
	}

	cfg.fileServerHits.Store(0)

	err := cfg.dbQ.DeleteAllUsers(context.Background())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while reseting users", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0 and database reset to initial state."))
}
