package main

import(
	"net/http"
	"github.com/alscho/chirpy/internal/auth"
	"io"
	"log"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "no valid refresh token found", err)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
    	respondWithError(w, http.StatusInternalServerError, "Failed to read request body", err)
    	return
	}

	if len(bodyBytes) > 0 {
    	// If you wish to log the potential malicious content:
    	log.Printf("Unexpected body: %s", string(bodyBytes))
    
    	respondWithError(w, http.StatusBadRequest, "Failed to read request body", err)
    	return
	}

	userID, err := cfg.db.GetUserFromValidRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no valid refresh token", err)
		return
	}

	token, err := auth.MakeJWT(userID, cfg.si)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create token", err)
		return
	}

	_, err = auth.ValidateJWT(token, cfg.si)
	// not sure if safeguard like this is needed here, since I look up the userID myself. Maybe just Validate to be sure that nothing broke...
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "no authorization, bad token - token generation broken?", err)
		return
	}

	type response struct {
		Token string `json:"token"`
	}
	
	respondWithJSON(w, http.StatusOK, response{
		Token: token,
	})

}