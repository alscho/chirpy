package main

import(
	"fmt"
	"net/http"
	"os"
)

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request){
	
	contents, err := os.ReadFile("template_hits.html")
	if err != nil {
		fmt.Println("File reading error: %v", err)
		return
	}
	
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(string(contents), cfg.fileserverHits.Load())))
	}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
    })
}