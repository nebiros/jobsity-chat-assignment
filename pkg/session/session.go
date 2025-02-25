package session

import (
	"github.com/gorilla/sessions"
)

func NewCookieStore(sessionKey string) *sessions.CookieStore {
	return sessions.NewCookieStore([]byte(sessionKey))
}
