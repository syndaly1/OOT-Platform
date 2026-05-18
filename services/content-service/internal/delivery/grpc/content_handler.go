package grpc

import (
	"context"
	"time"

	"content-service/internal/domain"
	"content-service/internal/usecase"
	"content-service/pkg/contentpb"
)

type ContentHandler struct {
	contentpb.UnimplementedContentServiceServer
	contentUsecase *usecase.ContentUsecase
	episodeUsecase *usecase.EpisodeUsecase
}

func NewContentHandler(
	contentUsecase *usecase.ContentUsecase,
	episodeUsecase *usecase.EpisodeUsecase,
) *ContentHandler {
	return &ContentHandler{
		contentUsecase: contentUsecase,
		episodeUsecase: episodeUsecase,
	}
}

func (h *ContentHandler) CreateContent(ctx context.Context, req *contentpb.CreateContentRequest) (*contentpb.ContentResponse, error) {
	content := &domain.Content{
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
		Genre:       req.Genre,
		ReleaseYear: int(req.ReleaseYear),
		Rating:      req.Rating,
		IsFeatured:  req.IsFeatured,
		IsTrending:  req.IsTrending,
	}

	if err := h.contentUsecase.CreateContent(ctx, content); err != nil {
		return nil, err
	}

	return &contentpb.ContentResponse{Content: toProtoContent(content)}, nil
}

func (h *ContentHandler) GetContent(ctx context.Context, req *contentpb.GetContentRequest) (*contentpb.ContentResponse, error) {
	content, err := h.contentUsecase.GetContentByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &contentpb.ContentResponse{Content: toProtoContent(content)}, nil
}

func (h *ContentHandler) ListContent(ctx context.Context, req *contentpb.ListContentRequest) (*contentpb.ListContentResponse, error) {
	contents, err := h.contentUsecase.ListContents(ctx, req.Page, req.Limit)
	if err != nil {
		return nil, err
	}

	return &contentpb.ListContentResponse{Contents: toProtoContents(contents)}, nil
}

func (h *ContentHandler) UpdateContent(ctx context.Context, req *contentpb.UpdateContentRequest) (*contentpb.ContentResponse, error) {
	content := &domain.Content{
		ID:          req.Id,
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
		Genre:       req.Genre,
		ReleaseYear: int(req.ReleaseYear),
		Rating:      req.Rating,
		IsFeatured:  req.IsFeatured,
		IsTrending:  req.IsTrending,
	}

	if err := h.contentUsecase.UpdateContent(ctx, content); err != nil {
		return nil, err
	}

	updatedContent, err := h.contentUsecase.GetContentByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &contentpb.ContentResponse{Content: toProtoContent(updatedContent)}, nil
}

func (h *ContentHandler) DeleteContent(ctx context.Context, req *contentpb.DeleteContentRequest) (*contentpb.DeleteContentResponse, error) {
	if err := h.contentUsecase.DeleteContent(ctx, req.Id); err != nil {
		return nil, err
	}

	return &contentpb.DeleteContentResponse{Message: "content deleted successfully"}, nil
}

func (h *ContentHandler) SearchContent(ctx context.Context, req *contentpb.SearchContentRequest) (*contentpb.ListContentResponse, error) {
	contents, err := h.contentUsecase.SearchContents(ctx, req.Query)
	if err != nil {
		return nil, err
	}

	return &contentpb.ListContentResponse{Contents: toProtoContents(contents)}, nil
}

func (h *ContentHandler) GetContentByGenre(ctx context.Context, req *contentpb.GetContentByGenreRequest) (*contentpb.ListContentResponse, error) {
	contents, err := h.contentUsecase.GetContentsByGenre(ctx, req.Genre)
	if err != nil {
		return nil, err
	}

	return &contentpb.ListContentResponse{Contents: toProtoContents(contents)}, nil
}

func (h *ContentHandler) GetTrendingContent(ctx context.Context, req *contentpb.GetTrendingContentRequest) (*contentpb.ListContentResponse, error) {
	contents, err := h.contentUsecase.GetTrendingContents(ctx, req.Limit)
	if err != nil {
		return nil, err
	}

	return &contentpb.ListContentResponse{Contents: toProtoContents(contents)}, nil
}

func (h *ContentHandler) GetFeaturedContent(ctx context.Context, req *contentpb.GetFeaturedContentRequest) (*contentpb.ListContentResponse, error) {
	contents, err := h.contentUsecase.GetFeaturedContents(ctx, req.Limit)
	if err != nil {
		return nil, err
	}

	return &contentpb.ListContentResponse{Contents: toProtoContents(contents)}, nil
}

func (h *ContentHandler) CreateEpisode(ctx context.Context, req *contentpb.CreateEpisodeRequest) (*contentpb.EpisodeResponse, error) {
	episode := &domain.Episode{
		ContentID:       req.ContentId,
		Title:           req.Title,
		Description:     req.Description,
		SeasonNumber:    int(req.SeasonNumber),
		EpisodeNumber:   int(req.EpisodeNumber),
		DurationSeconds: int(req.DurationSeconds),
	}

	if err := h.episodeUsecase.CreateEpisode(ctx, episode); err != nil {
		return nil, err
	}

	return &contentpb.EpisodeResponse{Episode: toProtoEpisode(episode)}, nil
}

func (h *ContentHandler) GetEpisode(ctx context.Context, req *contentpb.GetEpisodeRequest) (*contentpb.EpisodeResponse, error) {
	episode, err := h.episodeUsecase.GetEpisodeByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &contentpb.EpisodeResponse{Episode: toProtoEpisode(episode)}, nil
}

func (h *ContentHandler) ListEpisodesByContent(ctx context.Context, req *contentpb.ListEpisodesByContentRequest) (*contentpb.ListEpisodesResponse, error) {
	episodes, err := h.episodeUsecase.ListEpisodesByContent(ctx, req.ContentId)
	if err != nil {
		return nil, err
	}

	return &contentpb.ListEpisodesResponse{Episodes: toProtoEpisodes(episodes)}, nil
}

func (h *ContentHandler) UpdateEpisode(ctx context.Context, req *contentpb.UpdateEpisodeRequest) (*contentpb.EpisodeResponse, error) {
	episode := &domain.Episode{
		ID:              req.Id,
		Title:           req.Title,
		Description:     req.Description,
		SeasonNumber:    int(req.SeasonNumber),
		EpisodeNumber:   int(req.EpisodeNumber),
		DurationSeconds: int(req.DurationSeconds),
	}

	if err := h.episodeUsecase.UpdateEpisode(ctx, episode); err != nil {
		return nil, err
	}

	updatedEpisode, err := h.episodeUsecase.GetEpisodeByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &contentpb.EpisodeResponse{Episode: toProtoEpisode(updatedEpisode)}, nil
}

func (h *ContentHandler) DeleteEpisode(ctx context.Context, req *contentpb.DeleteEpisodeRequest) (*contentpb.DeleteEpisodeResponse, error) {
	if err := h.episodeUsecase.DeleteEpisode(ctx, req.Id); err != nil {
		return nil, err
	}

	return &contentpb.DeleteEpisodeResponse{Message: "episode deleted successfully"}, nil
}

func (h *ContentHandler) SetSkipIntroMetadata(ctx context.Context, req *contentpb.SetSkipIntroMetadataRequest) (*contentpb.EpisodeResponse, error) {
	if err := h.episodeUsecase.SetSkipIntroMetadata(ctx, req.EpisodeId, int(req.SkipIntroStart), int(req.SkipIntroEnd)); err != nil {
		return nil, err
	}

	episode, err := h.episodeUsecase.GetEpisodeByID(ctx, req.EpisodeId)
	if err != nil {
		return nil, err
	}

	return &contentpb.EpisodeResponse{Episode: toProtoEpisode(episode)}, nil
}

func toProtoContent(content *domain.Content) *contentpb.Content {
	return &contentpb.Content{
		Id:          content.ID,
		Title:       content.Title,
		Description: content.Description,
		Type:        content.Type,
		Genre:       content.Genre,
		ReleaseYear: int32(content.ReleaseYear),
		Rating:      content.Rating,
		IsFeatured:  content.IsFeatured,
		IsTrending:  content.IsTrending,
		CreatedAt:   formatTime(content.CreatedAt),
		UpdatedAt:   formatTime(content.UpdatedAt),
	}
}

func toProtoContents(contents []*domain.Content) []*contentpb.Content {
	result := make([]*contentpb.Content, 0, len(contents))
	for _, content := range contents {
		result = append(result, toProtoContent(content))
	}
	return result
}

func toProtoEpisode(episode *domain.Episode) *contentpb.Episode {
	return &contentpb.Episode{
		Id:              episode.ID,
		ContentId:       episode.ContentID,
		Title:           episode.Title,
		Description:     episode.Description,
		SeasonNumber:    int32(episode.SeasonNumber),
		EpisodeNumber:   int32(episode.EpisodeNumber),
		DurationSeconds: int32(episode.DurationSeconds),
		SkipIntroStart:  int32(episode.SkipIntroStart),
		SkipIntroEnd:    int32(episode.SkipIntroEnd),
		CreatedAt:       formatTime(episode.CreatedAt),
		UpdatedAt:       formatTime(episode.UpdatedAt),
	}
}

func toProtoEpisodes(episodes []*domain.Episode) []*contentpb.Episode {
	result := make([]*contentpb.Episode, 0, len(episodes))
	for _, episode := range episodes {
		result = append(result, toProtoEpisode(episode))
	}
	return result
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}
