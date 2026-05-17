package postgres

import (
	"context"
	"database/sql"
	"streaming-platform/auth-service/internal/domain"
)

type profileRepo struct {
	db *sql.DB
}

func NewProfileRepository(db *sql.DB) domain.ProfileRepository {
	return &profileRepo{db: db}
}

func (r *profileRepo) Create(ctx context.Context, p *domain.Profile) error {
	query := `INSERT INTO profiles (id, account_id, name, type, language, avatar, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, query,
		p.ID, p.AccountID, p.Name, p.Type, p.Language, p.Avatar, p.CreatedAt, p.UpdatedAt)
	return err
}

func (r *profileRepo) GetByAccountID(ctx context.Context, accountID string) ([]*domain.Profile, error) {
	query := `SELECT id, account_id, name, type, language, avatar, created_at, updated_at FROM profiles WHERE account_id = $1`
	rows, err := r.db.QueryContext(ctx, query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []*domain.Profile
	for rows.Next() {
		p := &domain.Profile{}
		err := rows.Scan(&p.ID, &p.AccountID, &p.Name, &p.Type, &p.Language, &p.Avatar, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, p)
	}
	return profiles, nil
}

func (r *profileRepo) GetByID(ctx context.Context, id string) (*domain.Profile, error) {
	query := `SELECT id, account_id, name, type, language, avatar, created_at, updated_at FROM profiles WHERE id = $1`
	p := &domain.Profile{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.AccountID, &p.Name, &p.Type, &p.Language, &p.Avatar, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *profileRepo) Update(ctx context.Context, p *domain.Profile) error {
	query := `UPDATE profiles SET name = $1, language = $2, avatar = $3, updated_at = $4 WHERE id = $5`
	_, err := r.db.ExecContext(ctx, query, p.Name, p.Language, p.Avatar, p.UpdatedAt, p.ID)
	return err
}

func (r *profileRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM profiles WHERE id = $1`, id)
	return err
}
