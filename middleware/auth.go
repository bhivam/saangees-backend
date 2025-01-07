package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/bhivam/saangees-backend/util"
)

type AuthKey struct{}

func GetAuthMiddlewareFunc(
	tokenMaker *util.JWTMaker,
	logger *log.Logger,
	admin bool,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") == "" {
				logger.Println("No Authorization header found")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			token := strings.Split(r.Header.Get("Authorization"), "Bearer ")
			if len(token) != 2 {
				logger.Println("Invalid token format")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			claims, err := verifyClaimsFromAuthHeader(r, tokenMaker)
			if err != nil {
				logger.Println("Error validating token :: ", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if admin && !claims.IsAdmin {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), AuthKey{}, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func verifyClaimsFromAuthHeader(
	r *http.Request,
	tokenMaker *util.JWTMaker,
) (*util.UserClaims, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("no Authorization header found")
	}

	fields := strings.Fields(authHeader)
	if len(fields) != 2 || fields[0] != "Bearer" {
		return nil, fmt.Errorf("invalid authorization header")
	}

	token := fields[1]

	claims, err := tokenMaker.VerifyToken(token)
	if err != nil {
		return nil, fmt.Errorf("error verifying token [%v] :: %w", authHeader, err)
	}

	return claims, nil
}
