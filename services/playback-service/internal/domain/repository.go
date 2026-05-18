package domain

import "context"

type PlaybackRepository interface {
	StartPlayback(ctx context.Context, session *PlaybackSession) error
	UpdateProgress(ctx context.Context, session *PlaybackSession) error

	GetProgress(ctx context.Context, profileID, contentID, episodeID string) (*PlaybackSession, error)
	GetResumeData(ctx context.Context, profileID, contentID, episodeID string) (*PlaybackSession, error)
	GetContinueWatching(ctx context.Context, profileID string) ([]*PlaybackSession, error)

	MarkCompleted(ctx context.Context, profileID, contentID, episodeID string) error
	RemoveFromContinueWatching(ctx context.Context, profileID, contentID, episodeID string) error

	GetWatchHistory(ctx context.Context, profileID string) ([]*PlaybackSession, error)
	ClearWatchHistory(ctx context.Context, profileID string) error

	GetProfileStats(ctx context.Context, profileID string) (int32, int32, int32, error)

	GetPlaybackSession(ctx context.Context, sessionID string) (*PlaybackSession, error)
	EndPlaybackSession(ctx context.Context, sessionID string) error

	DeleteProgress(ctx context.Context, profileID, contentID, episodeID string) error
	GetLastWatchedContent(ctx context.Context, profileID string) (*PlaybackSession, error)
}