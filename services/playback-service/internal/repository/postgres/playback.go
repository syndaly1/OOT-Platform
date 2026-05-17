package postgres

import (
	"context"

	"streaming-platform/services/playback-service/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PlaybackRepository struct {
	db *pgxpool.Pool
}

func NewPlaybackRepository(db *pgxpool.Pool) *PlaybackRepository {
	return &PlaybackRepository{db: db}
}

func (r *PlaybackRepository) StartPlayback(ctx context.Context, session *domain.PlaybackSession) error {
	query := `
		INSERT INTO watch_progress (
			profile_id,
			content_id,
			episode_id,
			current_position_sec,
			duration_sec,
			watched_percent,
			completed
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
	`

	_, err := r.db.Exec(ctx, query,
		session.ProfileID,
		session.ContentID,
		session.EpisodeID,
		session.CurrentPositionSec,
		session.DurationSec,
		session.WatchedPercent,
		session.Completed,
	)

	return err
}

func (r *PlaybackRepository) UpdateProgress(ctx context.Context, session *domain.PlaybackSession) error {
	query := `
		UPDATE watch_progress
		SET current_position_sec=$1,
			duration_sec=$2,
			watched_percent=$3,
			completed=$4
		WHERE profile_id=$5 AND content_id=$6
	`
	_, err := r.db.Exec(ctx, query,
		session.CurrentPositionSec,
		session.DurationSec,
		session.WatchedPercent,
		session.Completed,
		session.ProfileID,
		session.ContentID,
	)

	return err
}

func (r *PlaybackRepository) GetProgress(ctx context.Context, profileID, contentID string) (*domain.PlaybackSession, error) {
	query := `
		SELECT profile_id, content_id, current_position_sec, duration_sec, watched_percent, completed
		FROM watch_progress
		WHERE profile_id=$1 AND content_id=$2
	`

	var session domain.PlaybackSession

	err := r.db.QueryRow(ctx, query, profileID, contentID).Scan(
		&session.ProfileID,
		&session.ContentID,
		&session.CurrentPositionSec,
		&session.DurationSec,
		&session.WatchedPercent,
		&session.Completed,
	)

	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *PlaybackRepository) MarkCompleted(ctx context.Context, profileID, contentID string) error {
	query := `
		UPDATE watch_progress
		SET completed=true
		WHERE profile_id=$1 AND content_id=$2
	`

	_, err := r.db.Exec(ctx, query, profileID, contentID)
	return err
}

func (r *PlaybackRepository) GetContinueWatching(ctx context.Context, profileID string) ([]*domain.PlaybackSession, error) {
	query := `
		SELECT profile_id, content_id, current_position_sec, duration_sec, watched_percent, completed
		FROM watch_progress
		WHERE profile_id=$1 AND completed=false
	`

	rows, err := r.db.Query(ctx, query, profileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*domain.PlaybackSession

	for rows.Next() {
		var session domain.PlaybackSession

		err := rows.Scan(
			&session.ProfileID,
			&session.ContentID,
			&session.CurrentPositionSec,
			&session.DurationSec,
			&session.WatchedPercent,
			&session.Completed,
		)

		if err != nil {
			return nil, err
		}

		sessions = append(sessions, &session)
	}

	return sessions, nil
}