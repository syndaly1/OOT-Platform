package repository

import (
	"context"

	"content-service/internal/domain"
)

type EpisodeRepository interface {
	CreateEpisode(ctx context.Context, episode *domain.Episode) error
	GetEpisodeByID(ctx context.Context, id string) (*domain.Episode, error)
	ListEpisodesByContent(ctx context.Context, contentID string) ([]*domain.Episode, error)
	UpdateEpisode(ctx context.Context, episode *domain.Episode) error
	DeleteEpisode(ctx context.Context, id string) error

	SetSkipIntroMetadata(ctx context.Context, episodeID string, skipIntroStart int, skipIntroEnd int) error
}
