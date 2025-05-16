package middleware

import (
	"net/http"
)

type contextKey string

const UserIDKey contextKey = "userID"

func Corsmiddleware(app http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS,PUT,DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization,X-Refresh-Token")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		app.ServeHTTP(w, r)
	})
}
