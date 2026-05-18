package domain

import "context"

type PlaybackEventPublisher interface {
	PublishProgressUpdated(ctx context.Context, session *PlaybackSession) error
	PublishCompleted(ctx context.Context, profileID, contentID, episodeID string) error
}
