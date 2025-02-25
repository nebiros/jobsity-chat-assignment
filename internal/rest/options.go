package rest

import (
	"database/sql"

	"github.com/gorilla/sessions"
)

type Config struct {
	DBClient     *sql.DB
	SessionStore sessions.Store
}

type Option func(*Config)

func WithDBClient(dbClient *sql.DB) Option {
	return func(c *Config) {
		c.DBClient = dbClient
	}
}

func WithSessionStore(sessionStore sessions.Store) Option {
	return func(c *Config) {
		c.SessionStore = sessionStore
	}
}
