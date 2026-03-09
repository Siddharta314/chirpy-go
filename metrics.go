package main

import (
	"fmt"
	"net/http"
)

func (apiCfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`<html>
		<body>
			<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
		</html>`, apiCfg.fileserverHits.Load())))
}

func (apiCfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	apiCfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
}