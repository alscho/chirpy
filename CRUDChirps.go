package main

import(
	"net/http"
	"encoding/json"
	"strings"
	"github.com/google/uuid"
	"time"
	"github.com/alscho/chirpy/internal/database"
	"github.com/alscho/chirpy/internal/auth"
	//"log"
)

type Chirp struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body string `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	idString := r.PathValue("chirpID")

	if idString == "" {
		respondWithError(w, http.StatusBadRequest, "no valid chirp id", nil)
		return
	}

	//log.Printf("Is this the chirp id: '%s' ?", idString)
	
	chirpUUID, err := uuid.Parse(idString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "UUID not valid - parsing impossible", err)
		return
	}

	type response struct {
		Chirp
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirpUUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't GET chirp with id: '"+idString+"'.", err)
		return
	}

	//log.Printf("This is the chirp struct: %v", chirp)

	respondWithJSON(w, http.StatusOK, response{
		Chirp: Chirp{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserID: chirp.UserID,
		},
	})
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't GET chirps", err)
		return
	}

	resp := make([]Chirp, 0)

	for _, chirp := range chirps {
		resp = append(resp, Chirp{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserID: chirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, resp)
}

func (cfg *apiConfig) handlerCreateChirps(w http.ResponseWriter, r *http.Request) {
	
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "problem getting the token", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.si)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "token bad", err)
		return
	}

	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)

	// log.Printf("parameters struct: %v", params)

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

	// log.Printf("Trying to add chirp with body: '%s' and user_id: '%v'", cleanedBody, params.UserID)

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: cleanedBody,
		UserID: userID,
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