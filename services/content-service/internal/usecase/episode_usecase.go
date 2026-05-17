package usecase

import (
	"context"

	"content-service/internal/domain"
	"content-service/internal/repository"
)

type EpisodeUsecase struct {
	episodeRepo repository.EpisodeRepository
}

func NewEpisodeUsecase(
	episodeRepo repository.EpisodeRepository,
) *EpisodeUsecase {
	return &EpisodeUsecase{
		episodeRepo: episodeRepo,
	}
}

func (u *EpisodeUsecase) CreateEpisode(
	ctx context.Context,
	episode *domain.Episode,
) error {
	return u.episodeRepo.CreateEpisode(ctx, episode)
}

func (u *EpisodeUsecase) GetEpisodeByID(
	ctx context.Context,
	id string,
) (*domain.Episode, error) {
	return u.episodeRepo.GetEpisodeByID(ctx, id)
}

func (u *EpisodeUsecase) ListEpisodesByContent(
	ctx context.Context,
	contentID string,
) ([]*domain.Episode, error) {
	return u.episodeRepo.ListEpisodesByContent(ctx, contentID)
}

func (u *EpisodeUsecase) UpdateEpisode(
	ctx context.Context,
	episode *domain.Episode,
) error {
	return u.episodeRepo.UpdateEpisode(ctx, episode)
}

func (u *EpisodeUsecase) DeleteEpisode(
	ctx context.Context,
	id string,
) error {
	return u.episodeRepo.DeleteEpisode(ctx, id)
}

func (u *EpisodeUsecase) SetSkipIntroMetadata(
	ctx context.Context,
	episodeID string,
	skipIntroStart int,
	skipIntroEnd int,
) error {
	return u.episodeRepo.SetSkipIntroMetadata(ctx, episodeID, skipIntroStart, skipIntroEnd)
}
