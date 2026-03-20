package main

import (
	"database/sql"
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
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

type LoginResponse struct {
    User
    Token string `json:"token"`
	RefreshToken string `json:"refresh_token"`
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
	respondWithJSON(w, http.StatusCreated, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}

func (apiCfg *apiConfig) updateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}
	//get user with token
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID, err := auth.ValidateJWT(token, apiCfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	//decode req body
	decoder := json.NewDecoder(r.Body)
	var params parameters
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error parsing request body")
		return
	}
	hashed_password, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing password")
		return
	}

	//call database to update user
	user, err := apiCfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID: userID,
		Email: params.Email,
		HashedPassword: hashed_password,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user")
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}

func (apiCfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	var params parameters
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
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

	token, err := auth.MakeJWT(user.ID, apiCfg.jwtSecret, time.Duration(60*60)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating token")
		return
	}
	refreshTokenStr:= auth.MakeRefreshToken()
    if refreshTokenStr == "" {
        respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token")
        return
    }
	_, err = apiCfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
        Token:     refreshTokenStr,
        UserID:    user.ID,
        ExpiresAt: time.Now().AddDate(0, 0, 60),
    })
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't save refresh token")
        return
    }
	respondWithJSON(w, http.StatusOK, LoginResponse{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token: token,
		RefreshToken: refreshTokenStr,
	})
}


func (apiCfg *apiConfig) pokaWebhookHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct{
		Event string `json:"event"`
		Data struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}
	var params parameters
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	_, err = apiCfg.db.UpgradeUserRed(r.Context(), params.Data.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
            respondWithError(w, http.StatusNotFound, "User not found")
            return
        }
		respondWithError(w, http.StatusInternalServerError, "Error updating user plan")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
