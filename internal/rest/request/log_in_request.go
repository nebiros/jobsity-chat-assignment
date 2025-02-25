package request

import (
	"errors"
	"net/http"
	"strings"
)

type LogIn struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

func (in *LogIn) Bind(r *http.Request) error {
	username := strings.TrimSpace(in.Username)
	if username == "" {
		return errors.New("username seems empty")
	}

	password := strings.TrimSpace(in.Password)
	if password == "" {
		return errors.New("password seems empty")
	}

	return nil
}
