package main

import "net/http"

func (cfg *apiConfig) handlerResetHits(w http.ResponseWriter, r *http.Request){
	if cfg.pf != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Reset all users - database empty"))
	cfg.db.ResetUsers(r.Context())
}