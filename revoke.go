package main

import(
	"net/http"
	"log"
	"io"
	"github.com/alscho/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
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

	err = cfg.db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no valid refresh token", err)
		return
	}

	type EmptyStruct struct {}

	respondWithJSON(w, http.StatusNoContent, EmptyStruct{})
}