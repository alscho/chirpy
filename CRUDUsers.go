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
	IsChirpyRed bool `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handlerUpdateUsers(w http.ResponseWriter, r *http.Request){
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no valid token found", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.si)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no authorization, bad token", err)
		return
	}

	type parameters struct {
		Password string `json:"password"`
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)	
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	if params.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Valid email is needed", nil)
		return
	}
	if params.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Valid password is needed", nil)
		return
	}
	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user, password invalid", err)
		return
	}

	// sets password and email, returns nothing
	cfg.db.UpdatePasswordEmailFromToken(r.Context(), database.UpdatePasswordEmailFromTokenParams{
		ID: userID,
		Email: params.Email,
		HashedPassword: hash,
	})

	// gets whole user back by email
	user, err := cfg.db.GetHashedPasswordByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't fetch user by email", err)
		return
	}

	type response struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string `json:"email"`
		IsChirpyRed bool `json:"is_chirpy_red"`
	}

	respondWithJSON(w, http.StatusOK, response{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})


}


func (cfg *apiConfig) handlerCreateUsers(w http.ResponseWriter, r *http.Request) {
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
	if params.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Valid email is needed", nil)
		return
	}

	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't create user, password invalid", err)
		return
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
			IsChirpyRed: user.IsChirpyRed,
		},
	})
}