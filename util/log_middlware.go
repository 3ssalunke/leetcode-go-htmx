package util

import (
	"log"
	"net/http"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request details before handling it
		log.Printf("Request: %s %s", r.Method, r.URL.Path)

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
