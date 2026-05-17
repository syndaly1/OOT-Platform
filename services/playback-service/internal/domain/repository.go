package domain

import "context"

type PlaybackRepository interface {
	StartPlayback(ctx context.Context, session *PlaybackSession) error
	UpdateProgress(ctx context.Context, session *PlaybackSession) error
	GetProgress(ctx context.Context, profileID, contentID string) (*PlaybackSession, error)
	MarkCompleted(ctx context.Context, profileID, contentID string) error
	GetContinueWatching(ctx context.Context, profileID string) ([]*PlaybackSession, error)
}