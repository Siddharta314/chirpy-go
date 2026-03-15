package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Siddharta314/chirpygo/internal/auth"
	"github.com/Siddharta314/chirpygo/internal/database"
	"github.com/google/uuid"
)

type User struct {
    ID        uuid.UUID `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Email     string    `json:"email"`
}

type LoginResponse struct {
    User
    Token string `json:"token"`
}

func (apiCfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	var params parameters
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error parsing request body")
		return
	}
	hashed_password, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing password")
		return
	}
	user, err := apiCfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email: params.Email,
		HashedPassword: hashed_password,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}
	respondWithJSON(w, 201, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}


func (apiCfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct{
		Email string `json:"email"`
		Password string `json:"password"`
		ExpiresIn int `json:"expires_in_seconds"`
	}
	decoder := json.NewDecoder(r.Body)
	var params parameters
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if params.ExpiresIn == 0 {
		params.ExpiresIn = 60 * 60  // 1 hour
	}

	user, err := apiCfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
	correct, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !correct{
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	token, err := auth.MakeJWT(user.ID, apiCfg.jwtSecret, time.Duration(params.ExpiresIn)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating token")
		return
	}
	respondWithJSON(w, http.StatusOK, LoginResponse{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token: token,
	})
}