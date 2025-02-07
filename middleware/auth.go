package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/bhivam/saangees-backend/data"
	"github.com/bhivam/saangees-backend/util"
)

func GetAuthMiddlewareFunc(
	userStore data.UserStore,
	logger *log.Logger,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// get token from cookie
			cookie, err := r.Cookie("token")
			if err == http.ErrNoCookie {
				// if no cookie, pass anonymous user
				ctx := context.WithValue(r.Context(), util.UserContextKey{}, data.AnonymousUser)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			} else if err != nil {
				logger.Println("Error getting cookie :: ", err)
				http.Error(w, "Error getting cookie", http.StatusInternalServerError)
				return
			}

			// get token from cookie
			token := cookie.Value

			if token == "" {
				ctx := context.WithValue(r.Context(), util.UserContextKey{}, data.AnonymousUser)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

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
