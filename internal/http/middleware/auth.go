package middleware

import (
	"encoding/json"
	"net/http"
)

type AuthErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewAuthMiddleware(expectedToken string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("X-Auth-Token")

			if token == "" || token != expectedToken {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)

				resp := AuthErrorResponse{
					Code:    "UNAUTHORIZED",
					Message: "Token mancante o non valido",
				}

				_ = json.NewEncoder(w).Encode(resp)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
