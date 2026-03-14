package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Siddharta314/chirpygo/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (apiCfg *apiConfig) createChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct{
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	decoder := json.NewDecoder(r.Body)
	var params parameters
	err := decoder.Decode(&params)
	if err != nil {
		http.Error(w, "Error parsing request body", http.StatusInternalServerError)
		return
	}
	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	chirp, err := apiCfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   params.Body,
		UserID: params.UserID,
	})
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, 201, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}


func (apiCfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := apiCfg.db.GetChirps(r.Context())
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Error getting chirps")
        return
    }

    chirps := []Chirp{} 
    for _, dbChirp := range dbChirps {
        chirps = append(chirps, Chirp{
            ID:        dbChirp.ID,
            CreatedAt: dbChirp.CreatedAt,
            UpdatedAt: dbChirp.UpdatedAt,
            Body:      dbChirp.Body,
            UserID:    dbChirp.UserID,
        })
    }

    respondWithJSON(w, 200, chirps)
}


func (apiCfg *apiConfig) getChirpByID(w http.ResponseWriter, r *http.Request) {
	chipID := r.PathValue("chirpID")
	if chipID == ""{
		respondWithError(w, http.StatusBadRequest, "Chirp ID is required")
		return
	}
	dbChirp, err := apiCfg.db.GetChirpByID(r.Context(), uuid.MustParse(chipID))
    if err != nil {
        respondWithError(w, http.StatusNotFound, "Chirp not found")
        return
    }
	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}
	respondWithJSON(w, 200, chirp)
}