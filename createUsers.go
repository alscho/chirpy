package main

import (
	"net/http"
	"encoding/json"
	"github.com/google/uuid"
	"time"
	"github.com/alscho/chirpy/internal/auth"
	"github.com/alscho/chirpy/internal/database"
	// "log"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string `json:"email"`
}

func (cfg *apiConfig) handlerCreateUsers(w http.ResponseWriter, r *http.Request) {
	
	// log.Printf("Trying to create new user...")

	type parameters struct {
		Password string `json:"password"`
		Email string `json:"email"`
	}

	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// log.Printf("Trying to create new user with email '%s'.", params.Email)

	if params.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Valid email is needed", nil)
	}

	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't create user, password invalid", err)
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email: params.Email,
		HashedPassword: hash,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}
	
	respondWithJSON(w, http.StatusCreated, response{
		User: User{
			ID: user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email: user.Email,
		},
	})
}