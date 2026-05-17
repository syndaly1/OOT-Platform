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

	return &pb.UpdateProgressResponse{
		Success:        true,
		WatchedPercent: 45.5,
		Completed:      false,
		Message:        "progress updated",
	}, nil
}

func (h *PlaybackHandler) GetProgress(
	ctx context.Context,
	req *pb.GetProgressRequest,
) (*pb.GetProgressResponse, error) {

	return &pb.GetProgressResponse{
		Progress: &pb.PlaybackProgress{
			ProfileId:          req.ProfileId,
			ContentId:          req.ContentId,
			EpisodeId:          req.EpisodeId,
			CurrentPositionSec: 320,
			DurationSec:        1200,
			WatchedPercent:     26.6,
			Completed:          false,
			UpdatedAt:          "2026-05-17",
		},
	}, nil
}