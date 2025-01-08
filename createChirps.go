package main

import(
	"net/http"
	"encoding/json"
	"strings"
	"github.com/google/uuid"
	"time"
	"github.com/alscho/chirpy/internal/database"
	//"log"
)

type Chirp struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body string `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type parameters struct {
	Body string `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirps(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	type response struct {
		Chirp
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	cleanedBody, _ := replaceProfanity(params.Body)

	//log.Printf("Trying to add chirp with body: '%s' and user_id: '%v'", cleanedBody, params.UserID)

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: cleanedBody,
		UserID: params.UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp or add chirp to database", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		Chirp: Chirp{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserID: chirp.UserID,
		},
	})
}

/*
func handlerValidateChirp (w http.ResponseWriter, r *http.Request){
	type respVal struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	cont := content{}
	err := decoder.Decode(&cont)
	
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(cont.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	cleanedBody, _ := replaceProfanity(cont.Body)

	respondWithJSON(w, http.StatusOK, respVal{
		CleanedBody: cleanedBody,
	})
}
*/

func replaceProfanity(in string) (out string, hasProfanity bool) {
	profanes := map[string]bool{
		"kerfuffle": true,
		"sharbert": true,
		"fornax": true,
	}

	hasProfanity = false

	words := strings.Split(in, " ")

	for i, word := range words {
		if _, exists := profanes[strings.ToLower(word)]; exists {
			words[i] = "****"
			hasProfanity = true
		}
	}

	out = strings.Join(words, " ")

	return out, hasProfanity
}