package usecase

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"streaming-platform/services/playback-service/internal/domain"
)

const (
	SubjectProgressUpdated = "playback.progress.updated"
	SubjectCompleted       = "playback.completed"
)

type EventPublisher interface {
	Publish(subject string, data []byte) error
}

type PlaybackUsecase struct {
	repo      domain.PlaybackRepository
	publisher EventPublisher
}

func NewPlaybackUsecase(repo domain.PlaybackRepository, publisher EventPublisher) *PlaybackUsecase {
	return &PlaybackUsecase{
		repo:      repo,
		publisher: publisher,
	}
}

func (u *PlaybackUsecase) StartPlayback(ctx context.Context, session *domain.PlaybackSession) error {
	session.CurrentPositionSec = 0
	session.WatchedPercent = 0
	session.Completed = false
	session.Status = "PLAYING"

	return u.repo.StartPlayback(ctx, session)
}

func (u *PlaybackUsecase) UpdateProgress(ctx context.Context, session *domain.PlaybackSession) error {
	if session.DurationSec > 0 {
		session.WatchedPercent =
			float32(session.CurrentPositionSec) /
				float32(session.DurationSec) * 100
	}

	if session.WatchedPercent >= 90 {
		session.Completed = true
		session.Status = "COMPLETED"
	} else {
		session.Completed = false
		session.Status = "PLAYING"
	}

	if err := u.repo.UpdateProgress(ctx, session); err != nil {
		return err
	}

	u.publishProgressUpdated(session)
	return nil
}

func (u *PlaybackUsecase) GetProgress(ctx context.Context, profileID, contentID, episodeID string) (*domain.PlaybackSession, error) {
	return u.repo.GetProgress(ctx, profileID, contentID, episodeID)
}

func (u *PlaybackUsecase) GetResumeData(ctx context.Context, profileID, contentID, episodeID string) (*domain.PlaybackSession, error) {
	return u.repo.GetResumeData(ctx, profileID, contentID, episodeID)
}

func (u *PlaybackUsecase) GetContinueWatching(ctx context.Context, profileID string) ([]*domain.PlaybackSession, error) {
	return u.repo.GetContinueWatching(ctx, profileID)
}

func (u *PlaybackUsecase) MarkCompleted(ctx context.Context, profileID, contentID, episodeID string) error {
	if err := u.repo.MarkCompleted(ctx, profileID, contentID, episodeID); err != nil {
		return err
	}

	u.publishCompleted(profileID, contentID, episodeID)
	return nil
}

func (u *PlaybackUsecase) RemoveFromContinueWatching(ctx context.Context, profileID, contentID, episodeID string) error {
	return u.repo.RemoveFromContinueWatching(ctx, profileID, contentID, episodeID)
}

func (u *PlaybackUsecase) GetWatchHistory(ctx context.Context, profileID string) ([]*domain.PlaybackSession, error) {
	return u.repo.GetWatchHistory(ctx, profileID)
}

func (u *PlaybackUsecase) ClearWatchHistory(ctx context.Context, profileID string) error {
	return u.repo.ClearWatchHistory(ctx, profileID)
}

func (u *PlaybackUsecase) GetProfileStats(ctx context.Context, profileID string) (int32, int32, int32, error) {
	return u.repo.GetProfileStats(ctx, profileID)
}

func (u *PlaybackUsecase) GetPlaybackSession(ctx context.Context, sessionID string) (*domain.PlaybackSession, error) {
	return u.repo.GetPlaybackSession(ctx, sessionID)
}

func (u *PlaybackUsecase) EndPlaybackSession(ctx context.Context, sessionID string) error {
	return u.repo.EndPlaybackSession(ctx, sessionID)
}

func (u *PlaybackUsecase) SkipIntro(ctx context.Context, introEndSec int32) int32 {
	return introEndSec
}

func (u *PlaybackUsecase) DeleteProgress(ctx context.Context, profileID, contentID, episodeID string) error {
	return u.repo.DeleteProgress(ctx, profileID, contentID, episodeID)
}

func (u *PlaybackUsecase) GetLastWatchedContent(ctx context.Context, profileID string) (*domain.PlaybackSession, error) {
	return u.repo.GetLastWatchedContent(ctx, profileID)
}

type progressUpdatedEvent struct {
	ProfileID string `json:"profile_id"`
	ContentID string `json:"content_id"`
	EpisodeID string `json:"episode_id"`
	UpdatedAt string `json:"updated_at"`
}

type completedEvent struct {
	ProfileID   string `json:"profile_id"`
	ContentID   string `json:"content_id"`
	EpisodeID   string `json:"episode_id"`
	CompletedAt string `json:"completed_at"`
}

func (u *PlaybackUsecase) publishProgressUpdated(session *domain.PlaybackSession) {
	event := progressUpdatedEvent{
		ProfileID: session.ProfileID,
		ContentID: session.ContentID,
		EpisodeID: session.EpisodeID,
		UpdatedAt: time.Now().UTC().Format(time.RFC3339Nano),
	}

	u.publishEvent(SubjectProgressUpdated, event)
}

func (u *PlaybackUsecase) publishCompleted(profileID, contentID, episodeID string) {
	event := completedEvent{
		ProfileID:   profileID,
		ContentID:   contentID,
		EpisodeID:   episodeID,
		CompletedAt: time.Now().UTC().Format(time.RFC3339Nano),
	}

	u.publishEvent(SubjectCompleted, event)
}

func (u *PlaybackUsecase) publishEvent(subject string, event any) {
	if u.publisher == nil {
		return
	}

	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("warning: failed to marshal playback event %s: %v", subject, err)
		return
	}

	if err := u.publisher.Publish(subject, data); err != nil {
		log.Printf("warning: failed to publish playback event %s: %v", subject, err)
	}
}
