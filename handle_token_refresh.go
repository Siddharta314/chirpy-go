package main

import (
	"net/http"
	"time"

	"github.com/Siddharta314/chirpygo/internal/auth"
)

func (apiCfg *apiConfig) refreshHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	//check refresh in database
	refreshTokenData, err := apiCfg.db.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	//check if refresh token is expired
	if refreshTokenData.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	//check if refresh token is revoked
	if refreshTokenData.RevokedAt.Valid{
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	
	// send a new JWT
	type response struct {
		Token string `json:"token"`
	}
	token, err := auth.MakeJWT(refreshTokenData.UserID, apiCfg.jwtSecret, time.Duration(60*60)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating token")
		return
	}
	respondWithJSON(w, http.StatusOK, response{
		Token: token,
	})
}

func (apiCfg *apiConfig) revokeHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	err = apiCfg.db.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error revoking token")
		return
	}
	respondWithJSON(w, http.StatusNoContent, nil)
}