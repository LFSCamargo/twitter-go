package auth

import (
	"context"
	"net/http"

	userModel "github.com/LFSCamargo/twitter-go/database/models/user"
	"github.com/LFSCamargo/twitter-go/graph/services/user"
)

var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

// Middleware - Is the authentication middleare to get the user from the token
func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bearToken := r.Header.Get("Authorization")

			// Allow unauthenticated users in
			if bearToken == "" {
				next.ServeHTTP(w, r)
				return
			}

			user, err := user.GetUserFromToken(bearToken)
			if err != nil {
				http.Error(w, "Invalid Token", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), userCtxKey, user)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// ForContext - finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *userModel.User {
	raw, _ := ctx.Value(userCtxKey).(*userModel.User)
	return raw
}
