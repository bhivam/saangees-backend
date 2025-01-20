package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/bhivam/saangees-backend/data"
	"github.com/bhivam/saangees-backend/util"
)

func GetAuthMiddlewareFunc(
	userStore data.UserStore,
	logger *log.Logger,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Vary", "Authorization")

			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				ctx := context.WithValue(r.Context(), util.UserContextKey{}, data.AnonymousUser)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			tokenParts := strings.Split(r.Header.Get("Authorization"), " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				logger.Println("Invalid token format")
				http.Error(w, "Bad Auth Header", http.StatusBadRequest)
				return
			}

			token := tokenParts[1]

			// TODO validate token

			user, err := userStore.GetByToken(data.ScopeAuthentication, token)
			if err != nil {
				logger.Println("Error getting user :: ", err)
				http.Error(w, "Error getting user from token", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), util.UserContextKey{}, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
