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
			profile_id, content_id, episode_id,
			current_position_sec, duration_sec,
			watched_percent, completed
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
			completed=$4,
			updated_at=NOW()
		WHERE profile_id=$5 AND content_id=$6 AND episode_id=$7
	`

	_, err := r.db.Exec(ctx, query,
		session.CurrentPositionSec,
		session.DurationSec,
		session.WatchedPercent,
		session.Completed,
		session.ProfileID,
		session.ContentID,
		session.EpisodeID,
	)

	return err
}

func (r *PlaybackRepository) GetProgress(ctx context.Context, profileID, contentID, episodeID string) (*domain.PlaybackSession, error) {
	query := `
		SELECT profile_id, content_id, episode_id,
		       current_position_sec, duration_sec,
		       watched_percent, completed
		FROM watch_progress
		WHERE profile_id=$1 AND content_id=$2 AND episode_id=$3
	`

	var session domain.PlaybackSession

	err := r.db.QueryRow(ctx, query, profileID, contentID, episodeID).Scan(
		&session.ProfileID,
		&session.ContentID,
		&session.EpisodeID,
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

func (r *PlaybackRepository) GetResumeData(ctx context.Context, profileID, contentID, episodeID string) (*domain.PlaybackSession, error) {
	return r.GetProgress(ctx, profileID, contentID, episodeID)
}

func (r *PlaybackRepository) GetContinueWatching(ctx context.Context, profileID string) ([]*domain.PlaybackSession, error) {
	query := `
		SELECT profile_id, content_id, episode_id,
		       current_position_sec, duration_sec,
		       watched_percent, completed
		FROM watch_progress
		WHERE profile_id=$1 AND completed=false AND current_position_sec > 0
		ORDER BY updated_at DESC
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
			&session.EpisodeID,
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

func (r *PlaybackRepository) MarkCompleted(ctx context.Context, profileID, contentID, episodeID string) error {
	query := `
		UPDATE watch_progress
		SET completed=true,
			watched_percent=100,
			updated_at=NOW()
		WHERE profile_id=$1 AND content_id=$2 AND episode_id=$3
	`

	_, err := r.db.Exec(ctx, query, profileID, contentID, episodeID)
	return err
}

func (r *PlaybackRepository) RemoveFromContinueWatching(ctx context.Context, profileID, contentID, episodeID string) error {
	query := `
		UPDATE watch_progress
		SET completed=true,
			updated_at=NOW()
		WHERE profile_id=$1 AND content_id=$2 AND episode_id=$3
	`

	_, err := r.db.Exec(ctx, query, profileID, contentID, episodeID)
	return err
}

func (r *PlaybackRepository) GetWatchHistory(ctx context.Context, profileID string) ([]*domain.PlaybackSession, error) {
	query := `
		SELECT profile_id, content_id, episode_id,
		       current_position_sec, duration_sec,
		       watched_percent, completed
		FROM watch_progress
		WHERE profile_id=$1
		ORDER BY updated_at DESC
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
			&session.EpisodeID,
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

func (r *PlaybackRepository) ClearWatchHistory(ctx context.Context, profileID string) error {
	query := `
		DELETE FROM watch_progress
		WHERE profile_id=$1
	`

	_, err := r.db.Exec(ctx, query, profileID)
	return err
}

func (r *PlaybackRepository) GetProfileStats(ctx context.Context, profileID string) (int32, int32, int32, error) {
	query := `
		SELECT
			COALESCE(SUM(current_position_sec), 0),
			COUNT(*) FILTER (WHERE completed=true),
			COUNT(*) FILTER (WHERE completed=false)
		FROM watch_progress
		WHERE profile_id=$1
	`

	var totalWatchedSec int32
	var completedCount int32
	var inProgressCount int32

	err := r.db.QueryRow(ctx, query, profileID).Scan(
		&totalWatchedSec,
		&completedCount,
		&inProgressCount,
	)

	return totalWatchedSec, completedCount, inProgressCount, err
}

func (r *PlaybackRepository) GetPlaybackSession(ctx context.Context, sessionID string) (*domain.PlaybackSession, error) {
	return &domain.PlaybackSession{
		SessionID: sessionID,
		Status:    "PLAYING",
	}, nil
}

func (r *PlaybackRepository) EndPlaybackSession(ctx context.Context, sessionID string) error {
	return nil
}

func (r *PlaybackRepository) DeleteProgress(ctx context.Context, profileID, contentID, episodeID string) error {
	query := `
		DELETE FROM watch_progress
		WHERE profile_id=$1 AND content_id=$2 AND episode_id=$3
	`

	_, err := r.db.Exec(ctx, query, profileID, contentID, episodeID)
	return err
}

func (r *PlaybackRepository) GetLastWatchedContent(ctx context.Context, profileID string) (*domain.PlaybackSession, error) {
	query := `
		SELECT profile_id, content_id, episode_id,
		       current_position_sec, duration_sec,
		       watched_percent, completed
		FROM watch_progress
		WHERE profile_id=$1
		ORDER BY updated_at DESC
		LIMIT 1
	`

	var session domain.PlaybackSession

	err := r.db.QueryRow(ctx, query, profileID).Scan(
		&session.ProfileID,
		&session.ContentID,
		&session.EpisodeID,
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