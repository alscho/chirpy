package main

import(
	"net/http"
	"encoding/json"
	"github.com/alscho/chirpy/internal/auth"
	"github.com/alscho/chirpy/internal/database"
	"github.com/google/uuid"
	"time"
	// "log"
)

const defaultExpirationTimeInSeconds = 3600

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email string `json:"email"`
	}
	
	type response struct {
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email string `json:"email"`
		Token string `json:"token"`
		RefreshToken string `json:"refresh_token"`
		IsChirpyRed bool `json:"is_chirpy_red"`
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

	user, err := cfg.db.GetHashedPasswordByEmail(r.Context(), params.Email)
	if err != nil || auth.CheckPasswordHash(params.Password, user.HashedPassword) != nil{
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.si)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create token", err)
		return
	}

	tokenUserID, err := auth.ValidateJWT(token, cfg.si)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no authorization, bad token", err)
		return
	}
	if tokenUserID != user.ID {
		respondWithError(w, http.StatusUnauthorized, "no authorization, wrong user", nil)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create refresh token", err)
		return
	}

	// log.Printf("email: %s\nuser: %s\ntoken: %s\nrefresh_token:%s", params.Email, user.ID, token, refreshToken)

	err = cfg.db.AddRefreshToken(r.Context(), database.AddRefreshTokenParams{
		Token: refreshToken, 
		UserID: user.ID,
		})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't add refresh token to database", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Token: token,
		RefreshToken: refreshToken,
		IsChirpyRed: user.IsChirpyRed,
	})
}


	