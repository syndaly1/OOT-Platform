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
	return &PlaybackHandler{
		usecase: u,
	}
}

func (h *PlaybackHandler) StartPlayback(
	ctx context.Context,
	req *pb.StartPlaybackRequest,
) (*pb.StartPlaybackResponse, error) {

	session := &domain.PlaybackSession{
		ProfileID:   req.ProfileId,
		ContentID:   req.ContentId,
		EpisodeID:   req.EpisodeId,
		DurationSec: req.DurationSec,
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

func (h *PlaybackHandler) UpdateProgress(
	ctx context.Context,
	req *pb.UpdateProgressRequest,
) (*pb.UpdateProgressResponse, error) {

	session := &domain.PlaybackSession{
		ProfileID:          req.ProfileId,
		ContentID:          req.ContentId,
		EpisodeID:          req.EpisodeId,
		CurrentPositionSec: req.CurrentPositionSec,
		DurationSec:        req.DurationSec,
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

func (h *PlaybackHandler) GetProgress(
	ctx context.Context,
	req *pb.GetProgressRequest,
) (*pb.GetProgressResponse, error) {

	session, err := h.usecase.GetProgress(
		ctx,
		req.ProfileId,
		req.ContentId,
		req.EpisodeId,
	)
	if err != nil {
		return nil, err
	}

	return &pb.GetProgressResponse{
		Progress: &pb.PlaybackProgress{
			ProfileId:          session.ProfileID,
			ContentId:          session.ContentID,
			EpisodeId:          session.EpisodeID,
			CurrentPositionSec: session.CurrentPositionSec,
			DurationSec:        session.DurationSec,
			WatchedPercent:     session.WatchedPercent,
			Completed:          session.Completed,
		},
	}, nil
}

func (h *PlaybackHandler) GetContinueWatching(
	ctx context.Context,
	req *pb.GetContinueWatchingRequest,
) (*pb.GetContinueWatchingResponse, error) {

	sessions, err := h.usecase.GetContinueWatching(ctx, req.ProfileId)
	if err != nil {
		return nil, err
	}

	items := make([]*pb.PlaybackProgress, 0, len(sessions))

	for _, session := range sessions {
		items = append(items, &pb.PlaybackProgress{
			ProfileId:          session.ProfileID,
			ContentId:          session.ContentID,
			EpisodeId:          session.EpisodeID,
			CurrentPositionSec: session.CurrentPositionSec,
			DurationSec:        session.DurationSec,
			WatchedPercent:     session.WatchedPercent,
			Completed:          session.Completed,
		})
	}

	return &pb.GetContinueWatchingResponse{
		Items: items,
	}, nil
}

func (h *PlaybackHandler) MarkCompleted(
	ctx context.Context,
	req *pb.MarkCompletedRequest,
) (*pb.MarkCompletedResponse, error) {

	err := h.usecase.MarkCompleted(
		ctx,
		req.ProfileId,
		req.ContentId,
		req.EpisodeId,
	)
	if err != nil {
		return nil, err
	}

	return &pb.MarkCompletedResponse{
		Success: true,
		Message: "content marked as completed",
	}, nil
}