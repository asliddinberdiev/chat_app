package user

import (
	"context"
	"database/sql"
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type repository struct {
	db DBTX
}

func NewRepository(db DBTX) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, user *User) (*User, error) {
	query := "INSERT INTO users(id, username, email, password) VALUES ($1, $2, $3, $4) RETURNING id"
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Username, user.Email, user.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	input := User{}

	query := "SELECT id, username, email, password FROM users WHERE email = $1"
	err := r.db.QueryRowContext(ctx, query, email).Scan(&input.ID, &input.Username, &input.Email, &input.Password)

	if err != nil {
		return nil, err
	}

	return &input, nil
}
