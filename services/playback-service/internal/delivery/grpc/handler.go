package grpc

import (
	"context"

	pb "streaming-platform/proto/playback"
	"streaming-platform/services/playback-service/internal/domain"
	"streaming-platform/services/playback-service/internal/usecase"
)

type PlaybackHandler struct {
	pb.UnimplementedPlaybackServiceServer
	usecase *usecase.PlaybackUsecase
}

func NewPlaybackHandler(u *usecase.PlaybackUsecase) *PlaybackHandler {
	return &PlaybackHandler{usecase: u}
}

func (h *PlaybackHandler) StartPlayback(ctx context.Context, req *pb.StartPlaybackRequest) (*pb.StartPlaybackResponse, error) {
	session := &domain.PlaybackSession{
		ProfileID:   req.ProfileId,
		ContentID:   req.ContentId,
		EpisodeID:   req.EpisodeId,
		DurationSec: req.DurationSec,
		Device:      req.Device,
	}

	err := h.usecase.StartPlayback(ctx, session)
	if err != nil {
		return nil, err
	}

	return &pb.StartPlaybackResponse{
		SessionId: "session-created",
		Message:   "playback started successfully",
	}, nil
}

func (h *PlaybackHandler) UpdateProgress(ctx context.Context, req *pb.UpdateProgressRequest) (*pb.UpdateProgressResponse, error) {
	session := &domain.PlaybackSession{
		SessionID:           req.SessionId,
		ProfileID:           req.ProfileId,
		ContentID:           req.ContentId,
		EpisodeID:           req.EpisodeId,
		CurrentPositionSec:  req.CurrentPositionSec,
		DurationSec:         req.DurationSec,
	}

	err := h.usecase.UpdateProgress(ctx, session)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateProgressResponse{
		Success:        true,
		WatchedPercent: session.WatchedPercent,
		Completed:      session.Completed,
		Message:        "progress updated",
	}, nil
}

func (h *PlaybackHandler) GetProgress(ctx context.Context, req *pb.GetProgressRequest) (*pb.GetProgressResponse, error) {
	session, err := h.usecase.GetProgress(ctx, req.ProfileId, req.ContentId, req.EpisodeId)
	if err != nil {
		return nil, err
	}

	return &pb.GetProgressResponse{Progress: toPlaybackProgress(session)}, nil
}

func (h *PlaybackHandler) GetResumeData(ctx context.Context, req *pb.GetResumeDataRequest) (*pb.GetResumeDataResponse, error) {
	session, err := h.usecase.GetResumeData(ctx, req.ProfileId, req.ContentId, req.EpisodeId)
	if err != nil {
		return nil, err
	}

	shouldShowContinue := session.CurrentPositionSec > 0 && !session.Completed

	return &pb.GetResumeDataResponse{
		ProfileId:          session.ProfileID,
		ContentId:          session.ContentID,
		EpisodeId:          session.EpisodeID,
		LastPositionSec:    session.CurrentPositionSec,
		Completed:          session.Completed,
		ShouldShowContinue: shouldShowContinue,
		IntroStartSec:      0,
		IntroEndSec:        82,
	}, nil
}

func (h *PlaybackHandler) GetContinueWatching(ctx context.Context, req *pb.GetContinueWatchingRequest) (*pb.GetContinueWatchingResponse, error) {
	sessions, err := h.usecase.GetContinueWatching(ctx, req.ProfileId)
	if err != nil {
		return nil, err
	}

	items := make([]*pb.PlaybackProgress, 0, len(sessions))
	for _, session := range sessions {
		items = append(items, toPlaybackProgress(session))
	}

	return &pb.GetContinueWatchingResponse{Items: items}, nil
}

func (h *PlaybackHandler) RemoveFromContinueWatching(ctx context.Context, req *pb.RemoveFromContinueWatchingRequest) (*pb.RemoveFromContinueWatchingResponse, error) {
	err := h.usecase.RemoveFromContinueWatching(ctx, req.ProfileId, req.ContentId, req.EpisodeId)
	if err != nil {
		return nil, err
	}

	return &pb.RemoveFromContinueWatchingResponse{Success: true}, nil
}

func (h *PlaybackHandler) MarkCompleted(ctx context.Context, req *pb.MarkCompletedRequest) (*pb.MarkCompletedResponse, error) {
	err := h.usecase.MarkCompleted(ctx, req.ProfileId, req.ContentId, req.EpisodeId)
	if err != nil {
		return nil, err
	}

	return &pb.MarkCompletedResponse{
		Success: true,
		Message: "content marked as completed",
	}, nil
}

func (h *PlaybackHandler) EndPlaybackSession(ctx context.Context, req *pb.EndPlaybackSessionRequest) (*pb.EndPlaybackSessionResponse, error) {
	err := h.usecase.EndPlaybackSession(ctx, req.SessionId)
	if err != nil {
		return nil, err
	}

	return &pb.EndPlaybackSessionResponse{Success: true}, nil
}

func (h *PlaybackHandler) GetPlaybackSession(ctx context.Context, req *pb.GetPlaybackSessionRequest) (*pb.GetPlaybackSessionResponse, error) {
	session, err := h.usecase.GetPlaybackSession(ctx, req.SessionId)
	if err != nil {
		return nil, err
	}

	return &pb.GetPlaybackSessionResponse{
		Session: &pb.PlaybackSession{
			SessionId: session.SessionID,
			ProfileId: session.ProfileID,
			ContentId: session.ContentID,
			EpisodeId: session.EpisodeID,
			Device:    session.Device,
			Status:    pb.PlaybackStatus_PLAYING,
		},
	}, nil
}

func (h *PlaybackHandler) GetWatchHistory(ctx context.Context, req *pb.GetWatchHistoryRequest) (*pb.GetWatchHistoryResponse, error) {
	sessions, err := h.usecase.GetWatchHistory(ctx, req.ProfileId)
	if err != nil {
		return nil, err
	}

	history := make([]*pb.WatchHistoryItem, 0, len(sessions))

	for _, session := range sessions {
		history = append(history, &pb.WatchHistoryItem{
			ProfileId:       session.ProfileID,
			ContentId:       session.ContentID,
			EpisodeId:       session.EpisodeID,
			TotalWatchedSec: session.CurrentPositionSec,
			Completed:       session.Completed,
		})
	}

	return &pb.GetWatchHistoryResponse{History: history}, nil
}

func (h *PlaybackHandler) ClearWatchHistory(ctx context.Context, req *pb.ClearWatchHistoryRequest) (*pb.ClearWatchHistoryResponse, error) {
	err := h.usecase.ClearWatchHistory(ctx, req.ProfileId)
	if err != nil {
		return nil, err
	}

	return &pb.ClearWatchHistoryResponse{Success: true}, nil
}

func (h *PlaybackHandler) GetProfileStats(ctx context.Context, req *pb.GetProfileStatsRequest) (*pb.GetProfileStatsResponse, error) {
	total, completed, inProgress, err := h.usecase.GetProfileStats(ctx, req.ProfileId)
	if err != nil {
		return nil, err
	}

	return &pb.GetProfileStatsResponse{
		ProfileId:        req.ProfileId,
		TotalWatchedSec: total,
		CompletedCount:  completed,
		InProgressCount: inProgress,
	}, nil
}

func (h *PlaybackHandler) SkipIntro(ctx context.Context, req *pb.SkipIntroRequest) (*pb.SkipIntroResponse, error) {
	newPosition := h.usecase.SkipIntro(ctx, req.IntroEndSec)

	return &pb.SkipIntroResponse{
		Success:        true,
		NewPositionSec: newPosition,
		Message:        "intro skipped successfully",
	}, nil
}

func (h *PlaybackHandler) GetLastWatchedContent(ctx context.Context, req *pb.GetLastWatchedContentRequest) (*pb.GetLastWatchedContentResponse, error) {
	session, err := h.usecase.GetLastWatchedContent(ctx, req.ProfileId)
	if err != nil {
		return nil, err
	}

	return &pb.GetLastWatchedContentResponse{
		Progress: toPlaybackProgress(session),
	}, nil
}

func (h *PlaybackHandler) DeleteProgress(ctx context.Context, req *pb.DeleteProgressRequest) (*pb.DeleteProgressResponse, error) {
	err := h.usecase.DeleteProgress(ctx, req.ProfileId, req.ContentId, req.EpisodeId)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteProgressResponse{Success: true}, nil
}

func toPlaybackProgress(session *domain.PlaybackSession) *pb.PlaybackProgress {
	return &pb.PlaybackProgress{
		ProfileId:          session.ProfileID,
		ContentId:          session.ContentID,
		EpisodeId:          session.EpisodeID,
		CurrentPositionSec: session.CurrentPositionSec,
		DurationSec:        session.DurationSec,
		WatchedPercent:     session.WatchedPercent,
		Completed:          session.Completed,
	}
}