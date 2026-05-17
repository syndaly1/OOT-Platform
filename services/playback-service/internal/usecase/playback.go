package usecase

import (
	"context"

	"streaming-platform/services/playback-service/internal/domain"
)

type PlaybackUsecase struct {
	repo domain.PlaybackRepository
}

func NewPlaybackUsecase(
	repo domain.PlaybackRepository,
) *PlaybackUsecase {

	return &PlaybackUsecase{
		repo: repo,
	}
}

func (u *PlaybackUsecase) StartPlayback(
	ctx context.Context,
	session *domain.PlaybackSession,
) error {

	session.CurrentPositionSec = 0
	session.WatchedPercent = 0
	session.Completed = false

	return u.repo.StartPlayback(ctx, session)
}