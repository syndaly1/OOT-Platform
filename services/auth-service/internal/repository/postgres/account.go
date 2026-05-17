package postgres

import (
	"context"
	"database/sql"
	"streaming-platform/auth-service/internal/domain"
)

type accountRepo struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) domain.AccountRepository {
	return &accountRepo{db: db}
}

func (r *accountRepo) Create(ctx context.Context, a *domain.Account) error {
	query := `INSERT INTO accounts (id, email, username, password, is_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, query,
		a.ID, a.Email, a.Username, a.Password, a.IsVerified, a.CreatedAt, a.UpdatedAt)
	return err
}

func (r *accountRepo) GetByID(ctx context.Context, id string) (*domain.Account, error) {
	query := `SELECT id, email, username, password, is_verified, created_at, updated_at FROM accounts WHERE id = $1`
	a := &domain.Account{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&a.ID, &a.Email, &a.Username, &a.Password, &a.IsVerified, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (r *accountRepo) GetByEmail(ctx context.Context, email string) (*domain.Account, error) {
	query := `SELECT id, email, username, password, is_verified, created_at, updated_at FROM accounts WHERE email = $1`
	a := &domain.Account{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&a.ID, &a.Email, &a.Username, &a.Password, &a.IsVerified, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (r *accountRepo) Update(ctx context.Context, a *domain.Account) error {
	query := `UPDATE accounts SET username = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, a.Username, a.UpdatedAt, a.ID)
	return err
}

func (r *accountRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM accounts WHERE id = $1`, id)
	return err
}

func (r *accountRepo) MarkVerified(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE accounts SET is_verified = TRUE WHERE id = $1`, id)
	return err
}
