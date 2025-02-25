package rest

import (
	"database/sql"
	"errors"
	"html/template"
	"net/http"

	"github.com/go-chi/render"
	"github.com/gorilla/sessions"
	"github.com/nebiros/jobsity-chat-assignment/internal/rest/middlewareext"
	"github.com/nebiros/jobsity-chat-assignment/internal/rest/request"
	"github.com/nebiros/jobsity-chat-assignment/internal/rest/response"
	"github.com/nebiros/jobsity-chat-assignment/internal/static"
	"github.com/nebiros/jobsity-chat-assignment/jobsitychatassignment/chat"
	"github.com/nebiros/jobsity-chat-assignment/jobsitychatassignment/user"
)

type rootHandler struct {
	sessionStore sessions.Store
	userService  *user.Service
	chatService  *chat.Service
}

func newRootHandler(sessionStore sessions.Store, userService *user.Service, chatService *chat.Service) *rootHandler {
	return &rootHandler{
		sessionStore: sessionStore,
		userService:  userService,
		chatService:  chatService,
	}
}

func (h *rootHandler) chat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		tpl, err := template.ParseFS(static.FS, "chat.gohtml")
		if err != nil {
			response.WriteHTTPError(ctx, w, http.StatusInternalServerError, err)
			return
		}

		session := ctx.Value(middlewareext.SessionCtxKey{}).(*sessions.Session)

		data := struct {
			UserID   string
			Username string
		}{
			UserID:   session.Values["userId"].(string),
			Username: session.Values["username"].(string),
		}

		if err := tpl.Execute(w, data); err != nil {
			response.WriteHTTPError(ctx, w, http.StatusInternalServerError, err)
			return
		}
	}
}

func (h *rootHandler) ws() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if err := h.chatService.InitWebSocket("default", w, r); err != nil {
			response.WriteHTTPError(ctx, w, http.StatusInternalServerError, err)
			return
		}
	}
}

func (h *rootHandler) newUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		tpl, err := template.ParseFS(static.FS, "new_user.gohtml")
		if err != nil {
			response.WriteHTTPError(ctx, w, http.StatusInternalServerError, err)
			return
		}

		if err := tpl.Execute(w, nil); err != nil {
			response.WriteHTTPError(ctx, w, http.StatusInternalServerError, err)
			return
		}
	}
}

func (h *rootHandler) doNewUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var req request.CreateUser
		if err := render.Bind(r, &req); err != nil {
			response.WriteHTTPError(ctx, w, http.StatusBadRequest, err)
			return
		}

		newUser, err := req.ToUser()
		if err != nil {
			response.WriteHTTPError(ctx, w, http.StatusInternalServerError, err)
			return
		}

		if err := h.userService.CreateUser(ctx, newUser); err != nil {
			response.WriteHTTPError(ctx, w, http.StatusInternalServerError, err)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (h *rootHandler) logIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		session, err := h.sessionStore.Get(r, "user")
		if err != nil {
			response.WriteHTTPError(ctx, w, http.StatusInternalServerError, err)
			return
		}

		if session != nil {
			if currentUserID, ok := session.Values["userId"].(string); ok && currentUserID != "" {
				http.Redirect(w, r, "/users/chat", http.StatusFound)
				return
			}
		}

		tpl, err := template.ParseFS(static.FS, "log_in.gohtml")
		if err != nil {
			response.WriteHTTPError(ctx, w, http.StatusInternalServerError, err)
			return
		}

		if err := tpl.Execute(w, nil); err != nil {
			response.WriteHTTPError(ctx, w, http.StatusInternalServerError, err)
			return
		}
	}
}

func (h *rootHandler) doLogIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var req request.LogIn
		if err := render.Bind(r, &req); err != nil {
			response.WriteHTTPError(ctx, w, http.StatusBadRequest, err)
			return
		}

		loggedUser, err := h.userService.LogIn(ctx, req.Username, req.Password)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				response.WriteHTTPError(ctx, w, http.StatusNotFound, err)
				return
			}

			response.WriteHTTPError(ctx, w, http.StatusInternalServerError, err)
			return
		}

		session, err := h.sessionStore.Get(r, "user")
		if err != nil {
			response.WriteHTTPError(ctx, w, http.StatusInternalServerError, err)
			return
		}

		session.Values["userId"] = loggedUser.ID
		session.Values["username"] = loggedUser.Username

		if err := h.sessionStore.Save(r, w, session); err != nil {
			response.WriteHTTPError(ctx, w, http.StatusInternalServerError, err)
			return
		}

		http.Redirect(w, r, "/users/chat", http.StatusFound)
	}
}
