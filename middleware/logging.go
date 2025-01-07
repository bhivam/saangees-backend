package middleware

import (
	"log"
	"net/http"
)

func Logging(next http.Handler, logger *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Printf(
			"RECIEVED REQUEST :: %v | %v | %v | %v\n",
			r.RequestURI,
			r.Method,
			r.RemoteAddr,
			r.UserAgent(),
		)
		next.ServeHTTP(w, r)
	})
}
