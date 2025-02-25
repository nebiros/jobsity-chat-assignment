package user

import (
	"context"
	"database/sql"

	domain "github.com/nebiros/jobsity-chat-assignment/jobsitychatassignment"
	"golang.org/x/crypto/bcrypt"
)

type Repository struct {
	dbClient *sql.DB
}

func NewRepository(dbClient *sql.DB) *Repository {
	return &Repository{dbClient: dbClient}
}

func (r *Repository) Create(ctx context.Context, in domain.User) error {
	q := `INSERT INTO users (id, username, hashedPassword) VALUES (?, ?, ?)`

	stmt, err := r.dbClient.PrepareContext(ctx, q)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, in.ID, in.Username, in.HashedPassword)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) LogIn(ctx context.Context, username, password string) (domain.User, error) {
	q := `SELECT id, username, hashedPassword FROM users WHERE username = ?`

	stmt, err := r.dbClient.PrepareContext(ctx, q)
	if err != nil {
		return domain.User{}, err
	}

	defer stmt.Close()

	var user domain.User
	if err := stmt.QueryRowContext(ctx, username).Scan(&user.ID, &user.Username, &user.HashedPassword); err != nil {
		return domain.User{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return domain.User{}, err
	}

	return user, nil
}
