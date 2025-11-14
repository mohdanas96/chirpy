package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mhmdanas10/chirpy/internal/auth"
	"github.com/mhmdanas10/chirpy/internal/database"
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
		Email    string `json:"email"`
		Password string `json:"password"`
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

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while hashing password", err)
	}

	createUserParams := database.CreateUserParams{Email: params.Email, HashedPassword: hashedPassword}

	user, err := cfg.dbQ.CreateUser(context.Background(), createUserParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while creating user", err)
		return
	}

	respondWithJson(w, http.StatusCreated, response{
		User: User{ID: user.ID, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Email: user.Email},
	})
}

func (cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	type reqParams struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		User
	}

	params := reqParams{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid body", err)
		return
	}

	user, err := cfg.dbQ.GetUserByEmail(context.Background(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	valid, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while comparing password", err)
		return
	}

	if !valid {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	data := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	respondWithJson(w, 200, data)
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
