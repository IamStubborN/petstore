package mware

import (
	"net/http"
	"time"
)

func ResponseDefaultHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-AllowMethods-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-AllowMethods-Methods", "GET, POST, DELETE, PUT")
		w.Header().Set("Access-Control-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Date", time.Now().UTC().String())
		next.ServeHTTP(w, r)
	})
}
