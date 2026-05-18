package nats

import (
	"context"
	"encoding/json"

	"streaming-platform/services/playback-service/internal/domain"

	natsgo "github.com/nats-io/nats.go"
)

const (
	SubjectProgressUpdated = "playback.progress.updated"
	SubjectCompleted       = "playback.completed"
)

type Publisher struct {
	conn *natsgo.Conn
}

func NewPublisher(conn *natsgo.Conn) *Publisher {
	return &Publisher{conn: conn}
}

type progressUpdatedEvent struct {
	SessionID          string  `json:"session_id,omitempty"`
	ProfileID          string  `json:"profile_id"`
	ContentID          string  `json:"content_id"`
	EpisodeID          string  `json:"episode_id,omitempty"`
	CurrentPositionSec int32   `json:"current_position_sec"`
	DurationSec        int32   `json:"duration_sec"`
	WatchedPercent     float32 `json:"watched_percent"`
	Completed          bool    `json:"completed"`
	Status             string  `json:"status"`
}

type completedEvent struct {
	ProfileID string `json:"profile_id"`
	ContentID string `json:"content_id"`
	EpisodeID string `json:"episode_id,omitempty"`
}

func (p *Publisher) PublishProgressUpdated(ctx context.Context, session *domain.PlaybackSession) error {
	event := progressUpdatedEvent{
		SessionID:          session.SessionID,
		ProfileID:          session.ProfileID,
		ContentID:          session.ContentID,
		EpisodeID:          session.EpisodeID,
		CurrentPositionSec: session.CurrentPositionSec,
		DurationSec:        session.DurationSec,
		WatchedPercent:     session.WatchedPercent,
		Completed:          session.Completed,
		Status:             session.Status,
	}

	return p.publish(ctx, SubjectProgressUpdated, event)
}

func (p *Publisher) PublishCompleted(ctx context.Context, profileID, contentID, episodeID string) error {
	event := completedEvent{
		ProfileID: profileID,
		ContentID: contentID,
		EpisodeID: episodeID,
	}

	return p.publish(ctx, SubjectCompleted, event)
}

func (p *Publisher) publish(ctx context.Context, subject string, event any) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	done := make(chan error, 1)
	go func() {
		done <- p.conn.Publish(subject, data)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}
