package main

import(
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidateChirp (w http.ResponseWriter, r *http.Request){
	type content struct {
		Body string `json:"body"`
	}

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