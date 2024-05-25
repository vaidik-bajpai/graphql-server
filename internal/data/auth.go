package data

import (
	"context"
	"net/http"

	"github.com/vaidik-bajpai/hackernews/internal/jwt"
)

var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")

		// Allow unauthenticated users in
		if header == "" {
			next.ServeHTTP(w, r)
			return
		}

		//validate jwt token
		tokenStr := header
		username, err := jwt.ParseToken(tokenStr)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusForbidden)
			return
		}

		// create user and check if user exists in db
		user := User{Username: username}
		id, err := GetUserIdByUsername(username)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		user.ID = id
		// put it in context
		ctx := context.WithValue(r.Context(), userCtxKey, &user)

		// and call the next with our new context
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *User {
	raw, _ := ctx.Value(userCtxKey).(*User)
	return raw
}
