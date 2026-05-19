package postgres

import (
	"context"

	"streaming-platform/services/playback-service/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PlaybackRepository struct {
	db *pgxpool.Pool
}

func NewPlaybackRepository(db *pgxpool.Pool) *PlaybackRepository {
	return &PlaybackRepository{db: db}
}

func (r *PlaybackRepository) StartPlayback(ctx context.Context, session *domain.PlaybackSession) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	updateQuery := `
		UPDATE watch_progress
		SET current_position_sec=$1,
			duration_sec=$2,
			watched_percent=$3,
			completed=$4,
			updated_at=NOW()
		WHERE profile_id=$5 AND content_id=$6 AND episode_id=$7
	`

	tag, err := tx.Exec(ctx, updateQuery,
		session.CurrentPositionSec,
		session.DurationSec,
		session.WatchedPercent,
		session.Completed,
		session.ProfileID,
		session.ContentID,
		session.EpisodeID,
	)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		insertQuery := `
			INSERT INTO watch_progress (
				profile_id, content_id, episode_id,
				current_position_sec, duration_sec,
				watched_percent, completed
			)
			VALUES ($1,$2,$3,$4,$5,$6,$7)
		`

		if _, err := tx.Exec(ctx, insertQuery,
			session.ProfileID,
			session.ContentID,
			session.EpisodeID,
			session.CurrentPositionSec,
			session.DurationSec,
			session.WatchedPercent,
			session.Completed,
		); err != nil {
			return err
		}
	}

	if err := deleteDuplicateProgressRows(ctx, tx, session.ProfileID, session.ContentID, session.EpisodeID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *PlaybackRepository) UpdateProgress(ctx context.Context, session *domain.PlaybackSession) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	updateQuery := `
		UPDATE watch_progress
		SET current_position_sec=$1,
			duration_sec=$2,
			watched_percent=$3,
			completed=$4,
			updated_at=NOW()
		WHERE profile_id=$5 AND content_id=$6 AND episode_id=$7
	`

	tag, err := tx.Exec(ctx, updateQuery,
		session.CurrentPositionSec,
		session.DurationSec,
		session.WatchedPercent,
		session.Completed,
		session.ProfileID,
		session.ContentID,
		session.EpisodeID,
	)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		insertQuery := `
			INSERT INTO watch_progress (
				profile_id, content_id, episode_id,
				current_position_sec, duration_sec,
				watched_percent, completed
			)
			VALUES ($1,$2,$3,$4,$5,$6,$7)
		`

		if _, err := tx.Exec(ctx, insertQuery,
			session.ProfileID,
			session.ContentID,
			session.EpisodeID,
			session.CurrentPositionSec,
			session.DurationSec,
			session.WatchedPercent,
			session.Completed,
		); err != nil {
			return err
		}
	}

	if err := deleteDuplicateProgressRows(ctx, tx, session.ProfileID, session.ContentID, session.EpisodeID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *PlaybackRepository) GetProgress(ctx context.Context, profileID, contentID, episodeID string) (*domain.PlaybackSession, error) {
	query := `
		SELECT profile_id, content_id, episode_id,
		       current_position_sec, duration_sec,
		       watched_percent, completed
		FROM watch_progress
		WHERE profile_id=$1 AND content_id=$2 AND episode_id=$3
		ORDER BY updated_at DESC
		LIMIT 1
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
		FROM (
			SELECT DISTINCT ON (profile_id, content_id, episode_id)
			       profile_id, content_id, episode_id,
			       current_position_sec, duration_sec,
			       watched_percent, completed, updated_at
			FROM watch_progress
			WHERE profile_id=$1 AND completed=false AND current_position_sec > 0
			ORDER BY profile_id, content_id, episode_id, updated_at DESC
		) latest
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

	return sessions, rows.Err()
}

func (r *PlaybackRepository) MarkCompleted(ctx context.Context, profileID, contentID, episodeID string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
		UPDATE watch_progress
		SET completed=true,
			watched_percent=100,
			updated_at=NOW()
		WHERE profile_id=$1 AND content_id=$2 AND episode_id=$3
	`

	tag, err := tx.Exec(ctx, query, profileID, contentID, episodeID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		insertQuery := `
			INSERT INTO watch_progress (
				profile_id, content_id, episode_id,
				current_position_sec, duration_sec,
				watched_percent, completed
			)
			VALUES ($1,$2,$3,0,0,100,true)
		`

		if _, err := tx.Exec(ctx, insertQuery, profileID, contentID, episodeID); err != nil {
			return err
		}
	}

	if err := deleteDuplicateProgressRows(ctx, tx, profileID, contentID, episodeID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *PlaybackRepository) RemoveFromContinueWatching(ctx context.Context, profileID, contentID, episodeID string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
		UPDATE watch_progress
		SET completed=true,
			updated_at=NOW()
		WHERE profile_id=$1 AND content_id=$2 AND episode_id=$3
	`

	if _, err := tx.Exec(ctx, query, profileID, contentID, episodeID); err != nil {
		return err
	}

	if err := deleteDuplicateProgressRows(ctx, tx, profileID, contentID, episodeID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *PlaybackRepository) GetWatchHistory(ctx context.Context, profileID string) ([]*domain.PlaybackSession, error) {
	query := `
		SELECT profile_id, content_id, episode_id,
		       current_position_sec, duration_sec,
		       watched_percent, completed
		FROM (
			SELECT DISTINCT ON (profile_id, content_id, episode_id)
			       profile_id, content_id, episode_id,
			       current_position_sec, duration_sec,
			       watched_percent, completed, updated_at
			FROM watch_progress
			WHERE profile_id=$1
			ORDER BY profile_id, content_id, episode_id, updated_at DESC
		) latest
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

	return sessions, rows.Err()
}

func (r *PlaybackRepository) ClearWatchHistory(ctx context.Context, profileID string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
		DELETE FROM watch_progress
		WHERE profile_id=$1
	`

	if _, err := tx.Exec(ctx, query, profileID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *PlaybackRepository) GetProfileStats(ctx context.Context, profileID string) (int32, int32, int32, error) {
	query := `
		SELECT
			COALESCE(SUM(current_position_sec), 0),
			COUNT(*) FILTER (WHERE completed=true),
			COUNT(*) FILTER (WHERE completed=false)
		FROM (
			SELECT DISTINCT ON (profile_id, content_id, episode_id)
			       profile_id, content_id, episode_id,
			       current_position_sec, completed, updated_at
			FROM watch_progress
			WHERE profile_id=$1
			ORDER BY profile_id, content_id, episode_id, updated_at DESC
		) latest
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
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
		DELETE FROM watch_progress
		WHERE profile_id=$1 AND content_id=$2 AND episode_id=$3
	`

	if _, err := tx.Exec(ctx, query, profileID, contentID, episodeID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *PlaybackRepository) GetLastWatchedContent(ctx context.Context, profileID string) (*domain.PlaybackSession, error) {
	query := `
		SELECT profile_id, content_id, episode_id,
		       current_position_sec, duration_sec,
		       watched_percent, completed
		FROM (
			SELECT DISTINCT ON (profile_id, content_id, episode_id)
			       profile_id, content_id, episode_id,
			       current_position_sec, duration_sec,
			       watched_percent, completed, updated_at
			FROM watch_progress
			WHERE profile_id=$1
			ORDER BY profile_id, content_id, episode_id, updated_at DESC
		) latest
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

func deleteDuplicateProgressRows(ctx context.Context, tx pgx.Tx, profileID, contentID, episodeID string) error {
	query := `
		DELETE FROM watch_progress
		WHERE profile_id=$1
		  AND content_id=$2
		  AND episode_id=$3
		  AND ctid NOT IN (
			  SELECT ctid
			  FROM watch_progress
			  WHERE profile_id=$1
			    AND content_id=$2
			    AND episode_id=$3
			  ORDER BY updated_at DESC
			  LIMIT 1
		  )
	`

	_, err := tx.Exec(ctx, query, profileID, contentID, episodeID)
	return err
}
