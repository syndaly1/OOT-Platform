package postgres

import (
	"context"
	"database/sql"
	"streaming-platform/auth-service/internal/domain"
)

type refreshTokenRepo struct {
	db *sql.DB
}

func NewRefreshTokenRepository(db *sql.DB) domain.RefreshTokenRepository {
	return &refreshTokenRepo{db: db}
}

func (r *refreshTokenRepo) Create(ctx context.Context, rt *domain.RefreshToken) error {
	query := `INSERT INTO refresh_tokens (id, account_id, token, expires_at, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, rt.ID, rt.AccountID, rt.Token, rt.ExpiresAt, rt.CreatedAt)
	return err
}

func (r *refreshTokenRepo) GetByToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	query := `SELECT id, account_id, token, expires_at, created_at FROM refresh_tokens WHERE token = $1`
	rt := &domain.RefreshToken{}
	err := r.db.QueryRowContext(ctx, query, token).Scan(&rt.ID, &rt.AccountID, &rt.Token, &rt.ExpiresAt, &rt.CreatedAt)
	if err != nil {
		return nil, err
	}
	return rt, nil
}

func (r *refreshTokenRepo) Delete(ctx context.Context, token string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM refresh_tokens WHERE token = $1`, token)
	return err
}

func (r *refreshTokenRepo) DeleteByAccountID(ctx context.Context, accountID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM refresh_tokens WHERE account_id = $1`, accountID)
	return err
}

type emailVerificationRepo struct {
	db *sql.DB
}

func NewEmailVerificationRepository(db *sql.DB) domain.EmailVerificationRepository {
	return &emailVerificationRepo{db: db}
}

func (r *emailVerificationRepo) Create(ctx context.Context, ev *domain.EmailVerification) error {
	query := `INSERT INTO email_verifications (id, account_id, token, expires_at, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, ev.ID, ev.AccountID, ev.Token, ev.ExpiresAt, ev.CreatedAt)
	return err
}

func (r *emailVerificationRepo) GetByToken(ctx context.Context, token string) (*domain.EmailVerification, error) {
	query := `SELECT id, account_id, token, expires_at, created_at FROM email_verifications WHERE token = $1`
	ev := &domain.EmailVerification{}
	err := r.db.QueryRowContext(ctx, query, token).Scan(&ev.ID, &ev.AccountID, &ev.Token, &ev.ExpiresAt, &ev.CreatedAt)
	if err != nil {
		return nil, err
	}
	return ev, nil
}

func (r *emailVerificationRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM email_verifications WHERE id = $1`, id)
	return err
}

type passwordResetRepo struct {
	db *sql.DB
}

func NewPasswordResetRepository(db *sql.DB) domain.PasswordResetRepository {
	return &passwordResetRepo{db: db}
}

func (r *passwordResetRepo) Create(ctx context.Context, pr *domain.PasswordReset) error {
	query := `INSERT INTO password_resets (id, account_id, token, expires_at, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, pr.ID, pr.AccountID, pr.Token, pr.ExpiresAt, pr.CreatedAt)
	return err
}

func (r *passwordResetRepo) GetByToken(ctx context.Context, token string) (*domain.PasswordReset, error) {
	query := `SELECT id, account_id, token, expires_at, created_at FROM password_resets WHERE token = $1`
	pr := &domain.PasswordReset{}
	err := r.db.QueryRowContext(ctx, query, token).Scan(&pr.ID, &pr.AccountID, &pr.Token, &pr.ExpiresAt, &pr.CreatedAt)
	if err != nil {
		return nil, err
	}
	return pr, nil
}

func (r *passwordResetRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM password_resets WHERE id = $1`, id)
	return err
}
