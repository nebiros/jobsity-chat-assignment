package rest

import (
	"log/slog"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nebiros/jobsity-chat-assignment/internal/rest/middlewareext"
	"github.com/nebiros/jobsity-chat-assignment/jobsitychatassignment/chat"
	"github.com/nebiros/jobsity-chat-assignment/jobsitychatassignment/user"
	"github.com/nebiros/jobsity-chat-assignment/pkg/httpext"
	slogchi "github.com/samber/slog-chi"
)

func MakeRoutes(options ...Option) (router chi.Router, err error) {
	config := &Config{
		DBClient:     nil,
		SessionStore: nil,
	}

	for _, o := range options {
		o(config)
	}

	router = chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(slogchi.New(slog.Default()))
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(30 * 2 * time.Second))
	router.Use(middleware.StripSlashes)

	httpClient, err := httpext.NewClient()
	if err != nil {
		return nil, err
	}

	rh := newRootHandler(
		config.SessionStore,
		user.NewService(user.NewRepository(config.DBClient)),
		chat.NewService(httpClient),
	)

	router.Route("/users", func(router chi.Router) {
		router.Get("/new", rh.newUser())
		router.Post("/new", rh.doNewUser())
		router.With(middlewareext.Session(config.SessionStore)).Get("/chat", rh.chat())
		router.With(middlewareext.Session(config.SessionStore)).Get("/ws", rh.ws())
	})

	router.Get("/", rh.logIn())
	router.Post("/", rh.doLogIn())

	return router, nil
}
