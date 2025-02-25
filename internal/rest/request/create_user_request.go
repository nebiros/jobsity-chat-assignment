package request

import (
	"errors"
	"net/http"
	"strings"

	domain "github.com/nebiros/jobsity-chat-assignment/jobsitychatassignment"
	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
)

type CreateUser struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

func (in *CreateUser) Bind(r *http.Request) error {
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

func (in *CreateUser) ToUser() (domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.Password), 8)
	if err != nil {
		return domain.User{}, err
	}

	return domain.User{
		ID:             ulid.Make().String(),
		Username:       in.Username,
		HashedPassword: string(hashedPassword),
	}, nil
}
