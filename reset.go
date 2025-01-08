package main

import "net/http"

func (cfg *apiConfig) handlerResetHits(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Reset hit counter to 0"))
	cfg.fileserverHits.Store(0)
}