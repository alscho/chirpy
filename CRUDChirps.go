package main

import(
	"net/http"
	"encoding/json"
	"strings"
	"github.com/google/uuid"
	"time"
	"github.com/alscho/chirpy/internal/database"
	"github.com/alscho/chirpy/internal/auth"
	"sort"
	"errors"
	"fmt"
	// "log"
)

type Chirp struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body string `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	
	// handling the chirpID
	idString := r.PathValue("chirpID")
	if idString == "" {
		respondWithError(w, http.StatusBadRequest, "no valid chirp id", nil)
		return
	}
	chirpUUID, err := uuid.Parse(idString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "UUID not valid - parsing impossible", err)
		return
	}

	// handling the authentication by extracting the user_id from the given bearer token
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

	//log.Printf("Trying to delete chirp.%v as user.%v", chirpUUID, userID)

	// trying to delete the chirp with chirpUUID and userID
	deletionResult, err := cfg.db.DeleteChirpFromIDAndUserID(r.Context(), database.DeleteChirpFromIDAndUserIDParams{
		ID: chirpUUID,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "deletion attempt unsuccessful - unknown why", err)
		return
	}
	if deletionAmount, _ := deletionResult.RowsAffected(); deletionAmount != 1 {
		if _, err2 := cfg.db.ExistChirp(r.Context(), chirpUUID); err2 != nil {
			respondWithError(w, http.StatusNotFound, "chirp doesn't exist", err2)
			return
		}
		respondWithError(w, http.StatusForbidden, "authorization error: not the owner of the chirp", err)
		return		
	}

	type EmptyStruct struct {}

	respondWithJSON(w, http.StatusNoContent, EmptyStruct{})
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

func parseAuthorID(r *http.Request) (uuid.UUID, error){
	str := r.URL.Query().Get("author_id")
	if str == "" {
		return uuid.UUID{}, errors.New("unset or empty parameter: author_id")
	}
	authorID, err := uuid.Parse(str)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("UUID not valid - parsing impossible: %v", err)
	}
	return authorID, nil
}

func parseSortOrder(r *http.Request) (string, error){
	sortBy := r.URL.Query().Get("sort")
	if sortBy != "asc" && sortBy != "desc"{
		return "", errors.New("unset or invalid parameter: sort")
	}
	return sortBy, nil
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	// checks if author_id is set - possible filter option later
	authorID, _ := parseAuthorID(r)
	sortBy, _ := parseSortOrder(r)

	// retrieves chirps depending on: author_id set or not
	var chirps []database.Chirp
	var err error
	if (authorID != uuid.UUID{}) {
		chirps, err = cfg.db.GetChirpsByUserID(r.Context(), authorID)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "This user has no chirps", err)
			return
		}
	} else {
		chirps, err = cfg.db.GetChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't GET chirps", err)
			return
		}
	} 

	// sorts chirps depending on: sortBy (created_at) - since GetChirpsByUserID and GetChirps already sort by created_at (ascending) only desc is checked
	if sortBy == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})
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