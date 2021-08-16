package repository

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

const TokenQueryTimeout = time.Duration(10) * time.Second

type TokenRepository interface {
	IsTokenValid(ctx context.Context, token string) (bool, error)
	InsertToken(ctx context.Context, token, created string) error
}

type tokenRepository struct {
	Db *sqlx.DB
}

func NewTokenRepository(db *sqlx.DB) TokenRepository {
	return &tokenRepository{Db: db}
}

func (tr *tokenRepository) IsTokenValid(ctx context.Context, token string) (bool, error) {
	timoutContext, cancel := context.WithTimeout(ctx, TokenQueryTimeout)
	defer cancel()

	var count int
	if err := tr.Db.GetContext(timoutContext, &count, "select count(*) from token where token = ?", token); err != nil {
		return false, fmt.Errorf("Failed to validate token: %v", err)
	}

	if count < 1 {
		return false, nil
	}

	return true, nil
}

func (tr *tokenRepository) InsertToken(ctx context.Context, token, created string) error {
	timoutContext, cancel := context.WithTimeout(ctx, TokenQueryTimeout)
	defer cancel()

	result, err := tr.Db.ExecContext(timoutContext, "insert into token (token, created) values (?, ?)", token, created)
	if err != nil {
		return fmt.Errorf("Failed to insert token: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Failed fetch affected row count: %v", err)
	}

	if rows < 1 {
		return fmt.Errorf("No rows inserted into database")
	}

	return nil
}
