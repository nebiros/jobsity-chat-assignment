package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func NewClient(dbFilePath string) (*sql.DB, error) {
	dbClient, err := sql.Open("sqlite", dbFilePath)
	if err != nil {
		return nil, err
	}

	if err := dbClient.Ping(); err != nil {
		return nil, err
	}

	if err := initDB(dbClient); err != nil {
		return nil, fmt.Errorf("unable to init db '%s': %w", dbFilePath, err)
	}

	return dbClient, nil
}

func initDB(dbClient *sql.DB) error {
	_, err := dbClient.ExecContext(
		context.Background(),
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT NOT NULL,
			username TEXT NOT NULL, 
			hashedPassword TEXT NOT NULL 
		)`,
	)
	if err != nil {
		return err
	}

	return nil
}
