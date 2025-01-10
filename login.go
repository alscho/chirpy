package main

import(
	"net/http"
	"encoding/json"
	"github.com/alscho/chirpy/internal/auth"
	//"github.com/alscho/chirpy/internal/database"
	"github.com/google/uuid"
	"time"
)

const defaultExpirationTimeInSeconds = 3600

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email string `json:"email"`
		ExpiresInSeconds int `json:"expires_in_seconds"`
	}
	
	type response struct {
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email string `json:"email"`
		Token string `json:"token"`
	}
	
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	// initializing ExpiresInSeconds in case it's not set
	params.ExpiresInSeconds = defaultExpirationTimeInSeconds
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.ExpiresInSeconds <= 0 || params.ExpiresInSeconds > defaultExpirationTimeInSeconds {
		params.ExpiresInSeconds = defaultExpirationTimeInSeconds
	}

	expirationTime := time.Second * time.Duration(params.ExpiresInSeconds)

	if params.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Valid email is needed", nil)
		return
	}

	user, err := cfg.db.GetHashedPasswordByEmail(r.Context(), params.Email)
	if err != nil || auth.CheckPasswordHash(params.Password, user.HashedPassword) != nil{
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.si, expirationTime)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create token", err)
		return
	}

	tokenUserID, err := auth.ValidateJWT(token, cfg.si)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no authorization, bad token", err)
	}
	if tokenUserID != user.ID {
		respondWithError(w, http.StatusUnauthorized, "no authorization, wrong user", nil)
	}

	respondWithJSON(w, http.StatusOK, response{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Token: token,
	})
}


	