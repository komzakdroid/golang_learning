package middleware

import (
	"net/http"
	"os"
	"strings"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origins := os.Getenv("CORS_ALLOWED_ORIGINS")
		if origins == "" {
			origins = "*"
		}

		methods := os.Getenv("CORS_ALLOWED_METHODS")
		if methods == "" {
			methods = "GET,POST,PUT,DELETE,OPTIONS"
		}

		headers := os.Getenv("CORS_ALLOWED_HEADERS")
		if headers == "" {
			headers = "Content-Type,Authorization"
		}

		if origins == "*" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
			origin := r.Header.Get("Origin")
			allowedOrigins := strings.Split(origins, ",")
			for _, allowed := range allowedOrigins {
				if strings.TrimSpace(allowed) == origin {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", methods)
		w.Header().Set("Access-Control-Allow-Headers", headers)
		w.Header().Set("Access-Control-Max-Age", "3600")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
