package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
)

func NewRecoveryMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			defer func() {
				if err := recover(); err != nil {
					fmt.Printf("Recovered from panic: %v\n", err)
					debug.PrintStack()

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)

					resp := AuthErrorResponse{
						Code:    "INTERNAL",
						Message: "Errore interno inatteso",
					}

					_ = json.NewEncoder(w).Encode(resp)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
