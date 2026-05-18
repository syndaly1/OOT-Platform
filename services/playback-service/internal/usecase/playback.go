package usecase

import (
	"context"

	"streaming-platform/services/playback-service/internal/domain"
)

type PlaybackUsecase struct {
	repo domain.PlaybackRepository
}

func NewPlaybackUsecase(repo domain.PlaybackRepository) *PlaybackUsecase {
	return &PlaybackUsecase{repo: repo}
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

	return u.repo.UpdateProgress(ctx, session)
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
	return u.repo.MarkCompleted(ctx, profileID, contentID, episodeID)
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