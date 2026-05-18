package postgres

import (
	"context"
	"database/sql"

	"content-service/internal/domain"
)

type ContentPostgresRepository struct {
	db *sql.DB
}

func NewContentPostgresRepository(db *sql.DB) *ContentPostgresRepository {
	return &ContentPostgresRepository{db: db}
}

func (r *ContentPostgresRepository) CreateContent(ctx context.Context, content *domain.Content) error {
	query := `
		INSERT INTO contents (title, description, type, genre, release_year, rating, is_featured, is_trending)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		content.Title,
		content.Description,
		content.Type,
		content.Genre,
		content.ReleaseYear,
		content.Rating,
		content.IsFeatured,
		content.IsTrending,
	).Scan(&content.ID, &content.CreatedAt, &content.UpdatedAt)
}

func (r *ContentPostgresRepository) GetContentByID(ctx context.Context, id string) (*domain.Content, error) {
	query := `
		SELECT id, title, description, type, genre, release_year, rating, is_featured, is_trending, created_at, updated_at
		FROM contents
		WHERE id = $1
	`

	content := &domain.Content{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&content.ID,
		&content.Title,
		&content.Description,
		&content.Type,
		&content.Genre,
		&content.ReleaseYear,
		&content.Rating,
		&content.IsFeatured,
		&content.IsTrending,
		&content.CreatedAt,
		&content.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return content, nil
}

func (r *ContentPostgresRepository) ListContents(ctx context.Context, page, limit int32) ([]*domain.Content, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	query := `
		SELECT id, title, description, type, genre, release_year, rating, is_featured, is_trending, created_at, updated_at
		FROM contents
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanContents(rows)
}

func (r *ContentPostgresRepository) UpdateContent(ctx context.Context, content *domain.Content) error {
	query := `
		UPDATE contents
		SET title = $1,
		    description = $2,
		    type = $3,
		    genre = $4,
		    release_year = $5,
		    rating = $6,
		    is_featured = $7,
		    is_trending = $8,
		    updated_at = NOW()
		WHERE id = $9
		RETURNING updated_at
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		content.Title,
		content.Description,
		content.Type,
		content.Genre,
		content.ReleaseYear,
		content.Rating,
		content.IsFeatured,
		content.IsTrending,
		content.ID,
	).Scan(&content.UpdatedAt)
}

func (r *ContentPostgresRepository) DeleteContent(ctx context.Context, id string) error {
	query := `DELETE FROM contents WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *ContentPostgresRepository) SearchContents(ctx context.Context, queryText string) ([]*domain.Content, error) {
	query := `
		SELECT id, title, description, type, genre, release_year, rating, is_featured, is_trending, created_at, updated_at
		FROM contents
		WHERE title ILIKE '%' || $1 || '%'
		   OR description ILIKE '%' || $1 || '%'
		   OR genre ILIKE '%' || $1 || '%'
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, queryText)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanContents(rows)
}

func (r *ContentPostgresRepository) GetContentsByGenre(ctx context.Context, genre string) ([]*domain.Content, error) {
	query := `
		SELECT id, title, description, type, genre, release_year, rating, is_featured, is_trending, created_at, updated_at
		FROM contents
		WHERE genre ILIKE $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, genre)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanContents(rows)
}

func (r *ContentPostgresRepository) GetTrendingContents(ctx context.Context, limit int32) ([]*domain.Content, error) {
	if limit <= 0 {
		limit = 10
	}

	query := `
		SELECT id, title, description, type, genre, release_year, rating, is_featured, is_trending, created_at, updated_at
		FROM contents
		WHERE is_trending = TRUE
		ORDER BY rating DESC, created_at DESC
		LIMIT $1
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanContents(rows)
}

func (r *ContentPostgresRepository) GetFeaturedContents(ctx context.Context, limit int32) ([]*domain.Content, error) {
	if limit <= 0 {
		limit = 10
	}

	query := `
		SELECT id, title, description, type, genre, release_year, rating, is_featured, is_trending, created_at, updated_at
		FROM contents
		WHERE is_featured = TRUE
		ORDER BY rating DESC, created_at DESC
		LIMIT $1
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanContents(rows)
}

func scanContents(rows *sql.Rows) ([]*domain.Content, error) {
	var contents []*domain.Content

	for rows.Next() {
		content := &domain.Content{}

		err := rows.Scan(
			&content.ID,
			&content.Title,
			&content.Description,
			&content.Type,
			&content.Genre,
			&content.ReleaseYear,
			&content.Rating,
			&content.IsFeatured,
			&content.IsTrending,
			&content.CreatedAt,
			&content.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		contents = append(contents, content)
	}

	return contents, rows.Err()
}
