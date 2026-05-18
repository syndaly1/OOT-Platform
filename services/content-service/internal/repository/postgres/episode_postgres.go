package postgres

import (
	"context"
	"database/sql"

	"content-service/internal/domain"
)

type EpisodePostgresRepository struct {
	db *sql.DB
}

func NewEpisodePostgresRepository(db *sql.DB) *EpisodePostgresRepository {
	return &EpisodePostgresRepository{db: db}
}

func (r *EpisodePostgresRepository) CreateEpisode(ctx context.Context, episode *domain.Episode) error {
	query := `
		INSERT INTO episodes (content_id, title, description, season_number, episode_number, duration_seconds)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		episode.ContentID,
		episode.Title,
		episode.Description,
		episode.SeasonNumber,
		episode.EpisodeNumber,
		episode.DurationSeconds,
	).Scan(&episode.ID, &episode.CreatedAt, &episode.UpdatedAt)
}

func (r *EpisodePostgresRepository) GetEpisodeByID(ctx context.Context, id string) (*domain.Episode, error) {
	query := `
		SELECT id, content_id, title, description, season_number, episode_number,
		       duration_seconds, skip_intro_start, skip_intro_end, created_at, updated_at
		FROM episodes
		WHERE id = $1
	`

	episode := &domain.Episode{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&episode.ID,
		&episode.ContentID,
		&episode.Title,
		&episode.Description,
		&episode.SeasonNumber,
		&episode.EpisodeNumber,
		&episode.DurationSeconds,
		&episode.SkipIntroStart,
		&episode.SkipIntroEnd,
		&episode.CreatedAt,
		&episode.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return episode, nil
}

func (r *EpisodePostgresRepository) ListEpisodesByContent(ctx context.Context, contentID string) ([]*domain.Episode, error) {
	query := `
		SELECT id, content_id, title, description, season_number, episode_number,
		       duration_seconds, skip_intro_start, skip_intro_end, created_at, updated_at
		FROM episodes
		WHERE content_id = $1
		ORDER BY season_number ASC, episode_number ASC
	`

	rows, err := r.db.QueryContext(ctx, query, contentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanEpisodes(rows)
}

func (r *EpisodePostgresRepository) UpdateEpisode(ctx context.Context, episode *domain.Episode) error {
	query := `
		UPDATE episodes
		SET title = $1,
		    description = $2,
		    season_number = $3,
		    episode_number = $4,
		    duration_seconds = $5,
		    updated_at = NOW()
		WHERE id = $6
		RETURNING updated_at
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		episode.Title,
		episode.Description,
		episode.SeasonNumber,
		episode.EpisodeNumber,
		episode.DurationSeconds,
		episode.ID,
	).Scan(&episode.UpdatedAt)
}

func (r *EpisodePostgresRepository) DeleteEpisode(ctx context.Context, id string) error {
	query := `DELETE FROM episodes WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *EpisodePostgresRepository) SetSkipIntroMetadata(
	ctx context.Context,
	episodeID string,
	skipIntroStart int,
	skipIntroEnd int,
) error {
	query := `
		UPDATE episodes
		SET skip_intro_start = $1,
		    skip_intro_end = $2,
		    updated_at = NOW()
		WHERE id = $3
	`

	_, err := r.db.ExecContext(ctx, query, skipIntroStart, skipIntroEnd, episodeID)
	return err
}

func scanEpisodes(rows *sql.Rows) ([]*domain.Episode, error) {
	var episodes []*domain.Episode

	for rows.Next() {
		episode := &domain.Episode{}

		err := rows.Scan(
			&episode.ID,
			&episode.ContentID,
			&episode.Title,
			&episode.Description,
			&episode.SeasonNumber,
			&episode.EpisodeNumber,
			&episode.DurationSeconds,
			&episode.SkipIntroStart,
			&episode.SkipIntroEnd,
			&episode.CreatedAt,
			&episode.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		episodes = append(episodes, episode)
	}

	return episodes, rows.Err()
}
