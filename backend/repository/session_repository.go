package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"time"
)

const SessionQueryTimeout = time.Duration(10) * time.Second

type Session struct {
	Id      string `db:"id"`
	User    string `db:"user"`
	Created string `db:"created"`
	Stopped string `db:"stopped"`
	Name    string `db:"name"`
}

type SessionRepository interface {
	StartSession(ctx context.Context, token, createdAt string) (string, error)
	StopSession(ctx context.Context, id, stoppedAt string) error
	GetSession(ctx context.Context, id string) (*Session, error)
	ListSessions(ctx context.Context, token, from, to string) ([]*Session, error)
	SetName(ctx context.Context, id, name string) error
}

type sessionRepository struct {
	Db *sqlx.DB
}

func NewSessionRepository(db *sqlx.DB) SessionRepository {
	return &sessionRepository{Db: db}
}

func (sr *sessionRepository) StartSession(ctx context.Context, token, createdAt string) (string, error) {
	timoutContext, cancel := context.WithTimeout(ctx, SessionQueryTimeout)
	defer cancel()

	id := uuid.New().String()

	result, err := sr.Db.ExecContext(timoutContext, "insert into session (id, user, created) values (?, ?, ?)", id, token, createdAt)
	if err != nil {
		return "", fmt.Errorf("Failed to insert session: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("Failed fetch affected row count: %v", err)
	}

	if rows < 1 {
		return "", fmt.Errorf("No rows inserted into database")
	}

	return id, nil
}

func (sr *sessionRepository) StopSession(ctx context.Context, id, stoppedAt string) error {
	timoutContext, cancel := context.WithTimeout(ctx, SessionQueryTimeout)
	defer cancel()

	result, err := sr.Db.ExecContext(timoutContext, "update session set stopped = ? where id = ?", stoppedAt, id)
	if err != nil {
		return fmt.Errorf("Failed to update session: %v", err)
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

func (sr *sessionRepository) GetSession(ctx context.Context, id string) (*Session, error) {
	timoutContext, cancel := context.WithTimeout(ctx, SessionQueryTimeout)
	defer cancel()

	var session Session
	if err := sr.Db.GetContext(timoutContext, &session, "select * from session where id = ?", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("Failed to get session: %v", err)
	}

	return &session, nil
}

func (sr *sessionRepository) ListSessions(ctx context.Context, token, from, to string) ([]*Session, error) {
	timoutContext, cancel := context.WithTimeout(ctx, SessionQueryTimeout)
	defer cancel()

	var session []*Session
	var err error
	if from == "" || to == "" {
		err = sr.Db.SelectContext(timoutContext, &session, "select * from session where user = ?", token)
	} else {
		// Maybe this would be cleaner with a NULL stopped field, but it's good enough for now
		err = sr.Db.SelectContext(timoutContext, &session,
			"select * from session where user = ? and stopped >= ? and created <= ?"+ // Everything between the dates
				" union "+
				"select * from session where user = ? and created between ? and ? and stopped = '0000-00-00 00:00:00'"+ // Or the creation is between the dates and the session isn't yet stopped
				"order by created desc",
			token, from, to, token, from, to)
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("Failed to get session: %v", err)
	}

	return session, nil
}

func (sr *sessionRepository) SetName(ctx context.Context, id, name string) error {
	timoutContext, cancel := context.WithTimeout(ctx, SessionQueryTimeout)
	defer cancel()

	result, err := sr.Db.ExecContext(timoutContext, "update session set name = ? where id = ?", name, id)
	if err != nil {
		return fmt.Errorf("Failed to update session: %v", err)
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
