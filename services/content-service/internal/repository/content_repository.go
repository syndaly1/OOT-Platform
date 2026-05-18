package repository

import (
	"context"

	"content-service/internal/domain"
)

type ContentRepository interface {
	CreateContent(ctx context.Context, content *domain.Content) error
	GetContentByID(ctx context.Context, id string) (*domain.Content, error)
	ListContents(ctx context.Context, page, limit int32) ([]*domain.Content, error)
	UpdateContent(ctx context.Context, content *domain.Content) error
	DeleteContent(ctx context.Context, id string) error

	SearchContents(ctx context.Context, query string) ([]*domain.Content, error)
	GetContentsByGenre(ctx context.Context, genre string) ([]*domain.Content, error)
	GetTrendingContents(ctx context.Context, limit int32) ([]*domain.Content, error)
	GetFeaturedContents(ctx context.Context, limit int32) ([]*domain.Content, error)
}
