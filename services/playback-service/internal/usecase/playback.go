package usecase

import (
	"context"

	"streaming-platform/services/playback-service/internal/domain"
)

type PlaybackUsecase struct {
	repo domain.PlaybackRepository
}

func NewPlaybackUsecase(repo domain.PlaybackRepository) *PlaybackUsecase {
	return &PlaybackUsecase{
		repo: repo,
	}
}

func (u *PlaybackUsecase) StartPlayback(ctx context.Context, session *domain.PlaybackSession) error {
	session.CurrentPositionSec = 0
	session.WatchedPercent = 0
	session.Completed = false

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
	}

	return u.repo.UpdateProgress(ctx, session)
}

func (u *PlaybackUsecase) GetProgress(
	ctx context.Context,
	profileID string,
	contentID string,
	episodeID string,
) (*domain.PlaybackSession, error) {
	return u.repo.GetProgress(ctx, profileID, contentID, episodeID)
}

func (u *PlaybackUsecase) GetContinueWatching(
	ctx context.Context,
	profileID string,
) ([]*domain.PlaybackSession, error) {
	return u.repo.GetContinueWatching(ctx, profileID)
}

func (u *PlaybackUsecase) MarkCompleted(
	ctx context.Context,
	profileID string,
	contentID string,
	episodeID string,
) error {
	return u.repo.MarkCompleted(ctx, profileID, contentID, episodeID)
}

func (u *PlaybackUsecase) GetResumeData(
	ctx context.Context,
	profileID string,
	contentID string,
	episodeID string,
) (*domain.PlaybackSession, error) {

	return u.repo.GetResumeData(ctx, profileID, contentID, episodeID)
}