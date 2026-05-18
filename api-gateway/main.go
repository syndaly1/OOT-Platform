package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	pb "streaming-platform/proto/playback"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type gateway struct {
	playback pb.PlaybackServiceClient
	marshal  protojson.MarshalOptions
}

func main() {
	playbackAddr := getenv("PLAYBACK_SERVICE_URL", "localhost:50053")

	conn, err := grpc.NewClient(playbackAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("playback grpc connect:", err)
	}
	defer conn.Close()

	g := &gateway{
		playback: pb.NewPlaybackServiceClient(conn),
		marshal:  protojson.MarshalOptions{UseProtoNames: false},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /playback/start", g.startPlayback)
	mux.HandleFunc("PUT /playback/progress", g.updateProgress)
	mux.HandleFunc("GET /playback/continue-watching/{profileId}", g.getContinueWatching)
	mux.HandleFunc("GET /playback/history/{profileId}", g.getWatchHistory)
	mux.HandleFunc("DELETE /playback/continue-watching/{profileId}/{contentId}", g.removeFromContinueWatching)

	addr := ":" + getenv("PORT", "8080")
	log.Printf("api-gateway listening on %s", addr)
	if err := http.ListenAndServe(addr, withCORS(mux)); err != nil {
		log.Fatal("http serve:", err)
	}
}

type startPlaybackRequest struct {
	ProfileID   string `json:"profile_id"`
	ProfileId   string `json:"profileId"`
	ContentID   string `json:"content_id"`
	ContentId   string `json:"contentId"`
	EpisodeID   string `json:"episode_id"`
	EpisodeId   string `json:"episodeId"`
	DurationSec int32  `json:"duration_sec"`
	Device      string `json:"device"`
}

func (r startPlaybackRequest) toProto() *pb.StartPlaybackRequest {
	return &pb.StartPlaybackRequest{
		ProfileId:   first(r.ProfileID, r.ProfileId),
		ContentId:   first(r.ContentID, r.ContentId),
		EpisodeId:   first(r.EpisodeID, r.EpisodeId),
		DurationSec: r.DurationSec,
		Device:      r.Device,
	}
}

type updateProgressRequest struct {
	SessionID          string `json:"session_id"`
	SessionId          string `json:"sessionId"`
	ProfileID          string `json:"profile_id"`
	ProfileId          string `json:"profileId"`
	ContentID          string `json:"content_id"`
	ContentId          string `json:"contentId"`
	EpisodeID          string `json:"episode_id"`
	EpisodeId          string `json:"episodeId"`
	CurrentPositionSec int32  `json:"current_position_sec"`
	DurationSec        int32  `json:"duration_sec"`
}

func (r updateProgressRequest) toProto() *pb.UpdateProgressRequest {
	return &pb.UpdateProgressRequest{
		SessionId:          first(r.SessionID, r.SessionId),
		ProfileId:          first(r.ProfileID, r.ProfileId),
		ContentId:          first(r.ContentID, r.ContentId),
		EpisodeId:          first(r.EpisodeID, r.EpisodeId),
		CurrentPositionSec: r.CurrentPositionSec,
		DurationSec:        r.DurationSec,
	}
}

func (g *gateway) startPlayback(w http.ResponseWriter, r *http.Request) {
	var req startPlaybackRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	ctx, cancel := requestContext(r)
	defer cancel()

	resp, err := g.playback.StartPlayback(ctx, req.toProto())
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}

	g.writeProto(w, resp)
}

func (g *gateway) updateProgress(w http.ResponseWriter, r *http.Request) {
	var req updateProgressRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	ctx, cancel := requestContext(r)
	defer cancel()

	resp, err := g.playback.UpdateProgress(ctx, req.toProto())
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}

	g.writeProto(w, resp)
}

func (g *gateway) getContinueWatching(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := requestContext(r)
	defer cancel()

	resp, err := g.playback.GetContinueWatching(ctx, &pb.GetContinueWatchingRequest{
		ProfileId: r.PathValue("profileId"),
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}

	g.writeProto(w, resp)
}

func (g *gateway) getWatchHistory(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := requestContext(r)
	defer cancel()

	resp, err := g.playback.GetWatchHistory(ctx, &pb.GetWatchHistoryRequest{
		ProfileId: r.PathValue("profileId"),
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}

	g.writeProto(w, resp)
}

func (g *gateway) removeFromContinueWatching(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := requestContext(r)
	defer cancel()

	resp, err := g.playback.RemoveFromContinueWatching(ctx, &pb.RemoveFromContinueWatchingRequest{
		ProfileId: r.PathValue("profileId"),
		ContentId: r.PathValue("contentId"),
		EpisodeId: r.URL.Query().Get("episodeId"),
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}

	g.writeProto(w, resp)
}

func (g *gateway) writeProto(w http.ResponseWriter, msg proto.Message) {
	data, err := g.marshal.Marshal(msg)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return false
	}
	return true
}

func writeError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}

func requestContext(r *http.Request) (context.Context, context.CancelFunc) {
	return context.WithTimeout(r.Context(), 10*time.Second)
}

func getenv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func first(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
