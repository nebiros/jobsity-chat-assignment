package middlewareext

import (
	"context"
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/nebiros/jobsity-chat-assignment/internal/rest/response"
)

type (
	SessionCtxKey struct{}
)

func Session(sessionStore sessions.Store) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			session, err := sessionStore.Get(r, "user")
			if err != nil {
				http.SetCookie(w, &http.Cookie{Name: "user", MaxAge: -1, Path: "/"})
				return
			}

			if session != nil {
				if _, ok := session.Values["userId"]; !ok {
					response.WriteHTTPError(ctx, w, http.StatusForbidden, errors.New("not logged in"))
					return
				}
			}

			r = r.WithContext(context.WithValue(ctx, SessionCtxKey{}, session))
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
