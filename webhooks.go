package main

import(
	"net/http"
	"github.com/google/uuid"
	"encoding/json"
	"github.com/alscho/chirpy/internal/auth"
	// "log"
)

func (cfg *apiConfig) handlerPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no valid apiKey found", err)
		return
	}
	if apiKey != cfg.pk {
		respondWithError(w, http.StatusUnauthorized, "not the correct apiKey from the payments party", nil)
	}
	
	type Data struct {
		UserID uuid.UUID `json:"user_id"`
	}
	type parameters struct {
		Event string `json:"event"`
		Data Data `json:"data"`
	}


	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode request.Body", err)
		return
	}

	type EmptyStruct struct {}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, EmptyStruct{})
		return
	}

	// log.Printf("these are the parameters: %v", params)
	// log.Printf("trying to upgrade user with params.Data.UserID: %v", params.Data.UserID)

	upgradeResult, err := cfg.db.UpgradeUserToChirpyRed(r.Context(), params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Wasn't able to speak to database", err)
		return
	}
	if upgradeAmount, _ := upgradeResult.RowsAffected(); upgradeAmount != 1 {
		respondWithError(w, http.StatusNotFound, "Didn't find user", nil)
		return		
	}

	respondWithJSON(w, http.StatusNoContent, EmptyStruct{})
}