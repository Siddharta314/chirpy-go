package main

import (
	"encoding/json"
	"net/http"
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

	respondWithJSON(w, http.StatusOK, map[string]bool{"valid": true})

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