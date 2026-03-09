package main

import (
	"encoding/json"
	"net/http"
	"strings"
)


func validateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct{
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	type MyResponse struct {
		CleanedBody string `json:"cleaned_body"`
	}
	cleanBody := cleanBody(params.Body)
	respondWithJSON(w, http.StatusOK, MyResponse{CleanedBody: cleanBody})

}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(code)
	respondWithJSON(w, code, map[string]string{"error": msg})
}
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	data, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to marshal response")
		return
	}
	w.Write(data)
}


func cleanBody(body string) string {
	body = strings.TrimSpace(body)
	for _, word := range strings.Split(body, " ") {
		if checkProfane(word) {
			body = strings.ReplaceAll(body, word, "****")
		}
	}
	return body
}

func checkProfane(body string) bool {
	bodyLower := strings.ToLower(body)
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	for _, word := range profaneWords {
		if strings.Contains(bodyLower, word) {
			return true
		}
	}
	return false
}