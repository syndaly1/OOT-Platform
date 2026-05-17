package usecase

import (
	"context"

	"content-service/internal/domain"
	"content-service/internal/repository"
)

type ContentUsecase struct {
	contentRepo repository.ContentRepository
}

func NewContentUsecase(contentRepo repository.ContentRepository) *ContentUsecase {
	return &ContentUsecase{
		contentRepo: contentRepo,
	}
}

func (u *ContentUsecase) CreateContent(
	ctx context.Context,
	content *domain.Content,
) error {
	return u.contentRepo.CreateContent(ctx, content)
}

func (u *ContentUsecase) GetContentByID(
	ctx context.Context,
	id string,
) (*domain.Content, error) {
	return u.contentRepo.GetContentByID(ctx, id)
}

func (u *ContentUsecase) ListContents(
	ctx context.Context,
	page,
	limit int32,
) ([]*domain.Content, error) {
	return u.contentRepo.ListContents(ctx, page, limit)
}

func (u *ContentUsecase) UpdateContent(
	ctx context.Context,
	content *domain.Content,
) error {
	return u.contentRepo.UpdateContent(ctx, content)
}

func (u *ContentUsecase) DeleteContent(
	ctx context.Context,
	id string,
) error {
	return u.contentRepo.DeleteContent(ctx, id)
}

func (u *ContentUsecase) SearchContents(
	ctx context.Context,
	query string,
) ([]*domain.Content, error) {
	return u.contentRepo.SearchContents(ctx, query)
}

func (u *ContentUsecase) GetContentsByGenre(
	ctx context.Context,
	genre string,
) ([]*domain.Content, error) {
	return u.contentRepo.GetContentsByGenre(ctx, genre)
}

func (u *ContentUsecase) GetTrendingContents(
	ctx context.Context,
	limit int32,
) ([]*domain.Content, error) {
	return u.contentRepo.GetTrendingContents(ctx, limit)
}

func (u *ContentUsecase) GetFeaturedContents(
	ctx context.Context,
	limit int32,
) ([]*domain.Content, error) {
	return u.contentRepo.GetFeaturedContents(ctx, limit)
}
