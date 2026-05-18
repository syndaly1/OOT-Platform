package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	authpb "streaming-platform/proto/auth"
	contentpb "streaming-platform/proto/content"
	playbackpb "streaming-platform/proto/playback"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type gateway struct {
	auth     authpb.AuthServiceClient
	content  contentpb.ContentServiceClient
	playback playbackpb.PlaybackServiceClient
	marshal  protojson.MarshalOptions
}

func main() {
	authAddr := getenv("AUTH_SERVICE_URL", "localhost:50051")
	contentAddr := getenv("CONTENT_SERVICE_URL", "localhost:50052")
	playbackAddr := getenv("PLAYBACK_SERVICE_URL", "localhost:50053")

	authConn, err := grpc.NewClient(authAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("auth grpc connect:", err)
	}
	defer authConn.Close()

	contentConn, err := grpc.NewClient(contentAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("content grpc connect:", err)
	}
	defer contentConn.Close()

	playbackConn, err := grpc.NewClient(playbackAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("playback grpc connect:", err)
	}
	defer playbackConn.Close()

	g := &gateway{
		auth:     authpb.NewAuthServiceClient(authConn),
		content:  contentpb.NewContentServiceClient(contentConn),
		playback: playbackpb.NewPlaybackServiceClient(playbackConn),
		marshal: protojson.MarshalOptions{
			UseProtoNames: false,
		},
	}

	mux := http.NewServeMux()

	// Auth routes
	mux.HandleFunc("POST /api/auth/register", g.register)
	mux.HandleFunc("POST /api/auth/login", g.login)
	mux.HandleFunc("POST /api/auth/refresh", g.refreshToken)
	mux.HandleFunc("POST /api/auth/logout", g.logout)
	mux.HandleFunc("POST /api/auth/verify-email", g.verifyEmail)
	mux.HandleFunc("POST /api/auth/forgot-password", g.forgotPassword)
	mux.HandleFunc("POST /api/auth/reset-password", g.resetPassword)

	// Account routes
	mux.HandleFunc("GET /api/auth/account/{accountId}", g.getAccount)
	mux.HandleFunc("PUT /api/auth/account/{accountId}", g.updateAccount)
	mux.HandleFunc("DELETE /api/auth/account/{accountId}", g.deleteAccount)

	// Profile routes
	mux.HandleFunc("POST /api/profiles", g.createProfile)
	mux.HandleFunc("GET /api/profiles", g.getProfiles)
	mux.HandleFunc("GET /api/profiles/{profileId}", g.getProfileByID)
	mux.HandleFunc("PUT /api/profiles/{profileId}", g.updateProfile)
	mux.HandleFunc("DELETE /api/profiles/{profileId}", g.deleteProfile)

	// Content routes
	mux.HandleFunc("POST /api/content", g.createContent)
	mux.HandleFunc("GET /api/content", g.listContent)
	mux.HandleFunc("GET /api/content/search", g.searchContent)
	mux.HandleFunc("GET /api/content/genre/{genre}", g.getContentByGenre)
	mux.HandleFunc("GET /api/content/trending", g.getTrendingContent)
	mux.HandleFunc("GET /api/content/featured", g.getFeaturedContent)
	mux.HandleFunc("GET /api/content/{contentId}", g.getContent)
	mux.HandleFunc("PUT /api/content/{contentId}", g.updateContent)
	mux.HandleFunc("DELETE /api/content/{contentId}", g.deleteContent)

	// Episode routes
	mux.HandleFunc("POST /api/content/{contentId}/episodes", g.createEpisode)
	mux.HandleFunc("GET /api/content/{contentId}/episodes", g.listEpisodesByContent)
	mux.HandleFunc("GET /api/episodes/{episodeId}", g.getEpisode)
	mux.HandleFunc("PUT /api/episodes/{episodeId}", g.updateEpisode)
	mux.HandleFunc("DELETE /api/episodes/{episodeId}", g.deleteEpisode)
	mux.HandleFunc("PUT /api/episodes/{episodeId}/skip-intro", g.setSkipIntroMetadata)

	// Playback routes
	mux.HandleFunc("POST /api/playback/start", g.startPlayback)
	mux.HandleFunc("PUT /api/playback/progress", g.updateProgress)
	mux.HandleFunc("GET /api/playback/continue-watching/{profileId}", g.getContinueWatching)
	mux.HandleFunc("GET /api/playback/history/{profileId}", g.getWatchHistory)
	mux.HandleFunc("DELETE /api/playback/continue-watching/{profileId}/{contentId}", g.removeFromContinueWatching)

	// Old playback routes, kept for compatibility
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

// AUTH REQUESTS

type registerRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type refreshTokenRequest struct {
	RefreshToken      string `json:"refresh_token"`
	RefreshTokenCamel string `json:"refreshToken"`
}

type logoutRequest struct {
	RefreshToken      string `json:"refresh_token"`
	RefreshTokenCamel string `json:"refreshToken"`
}

type verifyEmailRequest struct {
	Token string `json:"token"`
}

type forgotPasswordRequest struct {
	Email string `json:"email"`
}

type resetPasswordRequest struct {
	Token            string `json:"token"`
	NewPassword      string `json:"new_password"`
	NewPasswordCamel string `json:"newPassword"`
}

type updateAccountRequest struct {
	Username string `json:"username"`
}

type createProfileRequest struct {
	AccountID string `json:"account_id"`
	AccountId string `json:"accountId"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Language  string `json:"language"`
	Avatar    string `json:"avatar"`
}

type updateProfileRequest struct {
	Name     string `json:"name"`
	Language string `json:"language"`
	Avatar   string `json:"avatar"`
}

// CONTENT REQUESTS

type createContentRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	PosterURL   string  `json:"poster_url"`
	PosterUrl   string  `json:"posterUrl"`
	Type        string  `json:"type"`
	Genre       string  `json:"genre"`
	ReleaseYear int32   `json:"release_year"`
	Rating      float64 `json:"rating"`
	IsFeatured  bool    `json:"is_featured"`
	IsTrending  bool    `json:"is_trending"`
}

type updateContentRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	PosterURL   string  `json:"poster_url"`
	PosterUrl   string  `json:"posterUrl"`
	Type        string  `json:"type"`
	Genre       string  `json:"genre"`
	ReleaseYear int32   `json:"release_year"`
	Rating      float64 `json:"rating"`
	IsFeatured  bool    `json:"is_featured"`
	IsTrending  bool    `json:"is_trending"`
}

type createEpisodeRequest struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	SeasonNumber    int32  `json:"season_number"`
	EpisodeNumber   int32  `json:"episode_number"`
	DurationSeconds int32  `json:"duration_seconds"`
}

type updateEpisodeRequest struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	SeasonNumber    int32  `json:"season_number"`
	EpisodeNumber   int32  `json:"episode_number"`
	DurationSeconds int32  `json:"duration_seconds"`
}

type setSkipIntroRequest struct {
	SkipIntroStart int32 `json:"skip_intro_start"`
	SkipIntroEnd   int32 `json:"skip_intro_end"`
}

// PLAYBACK REQUESTS

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

func (r startPlaybackRequest) toProto() *playbackpb.StartPlaybackRequest {
	return &playbackpb.StartPlaybackRequest{
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

func (r updateProgressRequest) toProto() *playbackpb.UpdateProgressRequest {
	return &playbackpb.UpdateProgressRequest{
		SessionId:          first(r.SessionID, r.SessionId),
		ProfileId:          first(r.ProfileID, r.ProfileId),
		ContentId:          first(r.ContentID, r.ContentId),
		EpisodeId:          first(r.EpisodeID, r.EpisodeId),
		CurrentPositionSec: r.CurrentPositionSec,
		DurationSec:        r.DurationSec,
	}
}

// AUTH HANDLERS

func (g *gateway) register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	ctx, cancel := requestContext(r)
	defer cancel()

	resp, err := g.auth.Register(ctx, &authpb.RegisterRequest{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}

	g.writeProto(w, resp)
}

func (g *gateway) login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	ctx, cancel := requestContext(r)
	defer cancel()

	resp, err := g.auth.Login(ctx, &authpb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		writeError(w, http.StatusUnauthorized, err)
		return
	}

	g.writeProto(w, resp)
}

func (g *gateway) refreshToken(w http.ResponseWriter, r *http.Request) {
	var req refreshTokenRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	ctx, cancel := requestContext(r)
	defer cancel()

	resp, err := g.auth.RefreshToken(ctx, &authpb.RefreshTokenRequest{
		RefreshToken: first(req.RefreshToken, req.RefreshTokenCamel),
	})
	if err != nil {
		writeError(w, http.StatusUnauthorized, err)
		return
	}

	g.writeProto(w, resp)
}

func (g *gateway) logout(w http.ResponseWriter, r *http.Request) {
	var req logoutRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	ctx, cancel := requestContext(r)
	defer cancel()

	resp, err := g.auth.Logout(ctx, &authpb.LogoutRequest{
		RefreshToken: first(req.RefreshToken, req.RefreshTokenCamel),
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}

	g.writeProto(w, resp)
}

func (g *gateway) verifyEmail(w http.ResponseWriter, r *http.Request) {
	var req verifyEmailRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	ctx, cancel := requestContext(r)
	defer cancel()

	resp, err := g.auth.VerifyEmail(ctx, &authpb.VerifyEmailRequest{
		Token: req.Token,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	g.writeProto(w, resp)
}

func (g *gateway) forgotPassword(w http.ResponseWriter, r *http.Request) {
	var req forgotPasswordRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	ctx, cancel := requestContext(r)
	defer cancel()

	resp, err := g.auth.ForgotPassword(ctx, &authpb.ForgotPasswordRequest{
		Email: req.Email,
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}

	g.writeProto(w, resp)
}

func (g *gateway) resetPassword(w http.ResponseWriter, r *http.Request) {
	var req resetPasswordRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	ctx, cancel := requestContext(r)
	defer cancel()

	resp, err := g.auth.ResetPassword(ctx, &authpb.ResetPasswordRequest{
		Token:       req.Token,
		NewPassword: first(req.NewPassword, req.NewPasswordCamel),
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	g.writeProto(w, resp)
}

func (g *gateway) getAccount(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := requestContext(r)
	defer cancel()

	resp, err := g.auth.GetAccount(ctx, &authpb.GetAccountRequest{
		AccountId: r.PathValue("accountId"),
	})
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}

	g.writeProto(w, resp)
}

func (g *gateway) updateAccount(w http.ResponseWriter, r *http.Request) {
	var req updateAccountRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	ctx, cancel := requestContext(r)
	defer cancel()

	resp, err := g.auth.UpdateAccount(ctx, &authpb.UpdateAccountRequest{
		AccountId: r.PathValue("accountId"),
		Username:  req.Username,
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}

	g.writeProto(w, resp)
}

func (g *gateway) deleteAccount(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := requestContext(r)
	defer cancel()

	resp, err := g.auth.DeleteAccount(ctx, &authpb.DeleteAccountRequest{
		AccountId: r.PathValue("accountId"),
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}

	g.writeProto(w, resp)
}

func (g *gateway) createProfile(w http.ResponseWriter, r *http.Request) {
	var req createProfileRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	ctx, cancel := requestContext(r)
	defer cancel()

	resp, err := g.auth.CreateProfile(ctx, &authpb.CreateProfileRequest{
		AccountId: first(req.AccountID, req.AccountId),
		Name:      req.Name,
		Type:      profileTypeFromString(req.Type),
		Language:  req.Language,
		Avatar:    req.Avatar,
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}

	g.writeProto(w, resp)
}

func (g *gateway) getProfiles(w http.ResponseWriter, r *http.Request) {
	accountID := first(r.URL.Query().Get("accountId"), r.URL.Query().Get("account_id"))
	if accountID == "" {
		writeErrorMessage(w, http.StatusBadRequest, "accountId is required")
		return
	}

	ctx, cancel := requestContext(r)
	defer cancel()

	resp, err := g.auth.GetProfiles(ctx, &authpb.GetProfilesRequest{
		AccountId: accountID,
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}

	g.writeProto(w, resp)
}

func (g *gateway) getProfileByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := requestContext(r)
	defer cancel()

	resp, err := g.auth.GetProfileById(ctx, &authpb.GetProfileByIdRequest{
		ProfileId: r.PathValue("profileId"),
	})
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}

	g.writeProto(w, resp)
}

func (g *gateway) updateProfile(w http.ResponseWriter, r *http.Request) {
	var req updateProfileRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	ctx, cancel := requestContext(r)
	defer cancel()

	resp, err := g.auth.UpdateProfile(ctx, &authpb.UpdateProfileRequest{
		ProfileId: r.PathValue("profileId"),
		Name:      req.Name,
		Language:  req.Language,
		Avatar:    req.Avatar,
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}

	g.writeProto(w, resp)
}

func (g *gateway) deleteProfile(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := requestContext(r)
	defer cancel()

	resp, err := g.auth.DeleteProfile(ctx, &authpb.DeleteProfileRequest{
		ProfileId: r.PathValue("profileId"),
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}

	g.writeProto(w, resp)
}

// CONTENT HANDLERS

func (g *gateway) createContent(w http.ResponseWriter, r *http.Request) {
	var req createContentRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	ctx, cancel := requestContext(r)
	defer cancel()

	resp, err := g.content.CreateContent(ctx, &contentpb.CreateContentRequest{
		Title:       req.Title,
		Description: req.Description,
		PosterUrl:   first(req.PosterURL, req.PosterUrl),
		Type:        req.Type,
		Genre:       req.Genre,
		ReleaseYear: req.ReleaseYear,
		Rating:      req.Rating,
		IsFeatured:  req.IsFeatured,
		IsTrending:  req.IsTrending,
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}

	g.writeProto(w, resp)
}

func (g *gateway) getContent(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := requestContext(r)
	defer cancel()

	resp, err := g.content.GetContent(ctx, &contentpb.GetContentRequest{
		Id: r.PathValue("contentId"),
	})
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}

	g.writeProto(w, resp)
}

func (g *gateway) listContent(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := requestContext(r)
	defer cancel()

	resp, err := g.content.ListContent(ctx, &contentpb.ListContentRequest{
		Page:  int32(queryInt(r, "page", 1)),
		Limit: int32(queryInt(r, "limit", 20)),
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}

	g.writeProto(w, resp)
}

func (g *gateway) updateContent(w http.ResponseWriter, r *http.Request) {
	var req updateContentRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	ctx, cancel := requestContext(r)
	defer cancel()

	resp, err := g.content.UpdateContent(ctx, &contentpb.UpdateContentRequest{
		Id:          r.PathValue("contentId"),
		Title:       req.Title,
		Description: req.Description,
		PosterUrl:   first(req.PosterURL, req.PosterUrl),
		Type:        req.Type,
		Genre:       req.Genre,
		ReleaseYear: req.ReleaseYear,
		Rating:      req.Rating,
		IsFeatured:  req.IsFeatured,
		IsTrending:  req.IsTrending,
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}

	g.writeProto(w, resp)
}

func (g *gateway) deleteContent(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := requestContext(r)
	defer cancel()

	resp, err := g.content.DeleteContent(ctx, &contentpb.DeleteContentRequest{
		Id: r.PathValue("contentId"),
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}

	g.writeProto(w, resp)
}

func (g *gateway) searchContent(w http.ResponseWriter, r *http.Request) {
	query := first(r.URL.Query().Get("q"), r.URL.Query().Get("query"))
	if query == "" {
		writeErrorMessage(w, http.StatusBadRequest, "q is required")
		return
	}

	ctx, cancel := requestContext(r)
	defer cancel()

	resp, err := g.content.SearchContent(ctx, &contentpb.SearchContentRequest{
		Query: query,
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}

	g.writeProto(w, resp)
}

func (g *gateway) getContentByGenre(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := requestContext(r)
	defer cancel()

	resp, err := g.content.GetContentByGenre(ctx, &contentpb.GetContentByGenreRequest{
		Genre: r.PathValue("genre"),
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}

	g.writeProto(w, resp)
}

func (g *gateway) getTrendingContent(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := requestContext(r)
	defer cancel()

	resp, err := g.content.GetTrendingContent(ctx, &contentpb.GetTrendingContentRequest{
		Limit: int32(queryInt(r, "limit", 10)),
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}

	g.writeProto(w, resp)
}

func (g *gateway) getFeaturedContent(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := requestContext(r)
	defer cancel()

	resp, err := g.content.GetFeaturedContent(ctx, &contentpb.GetFeaturedContentRequest{
		Limit: int32(queryInt(r, "limit", 10)),
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}

	g.writeProto(w, resp)
}

// EPISODE HANDLERS

func (g *gateway) createEpisode(w http.ResponseWriter, r *http.Request) {
	var req createEpisodeRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	ctx, cancel := requestContext(r)
	defer cancel()

	resp, err := g.content.CreateEpisode(ctx, &contentpb.CreateEpisodeRequest{
		ContentId:       r.PathValue("contentId"),
		Title:           req.Title,
		Description:     req.Description,
		SeasonNumber:    req.SeasonNumber,
		EpisodeNumber:   req.EpisodeNumber,
		DurationSeconds: req.DurationSeconds,
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}

	g.writeProto(w, resp)
}

func (g *gateway) getEpisode(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := requestContext(r)
	defer cancel()

	resp, err := g.content.GetEpisode(ctx, &contentpb.GetEpisodeRequest{
		Id: r.PathValue("episodeId"),
	})
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}

	g.writeProto(w, resp)
}

func (g *gateway) listEpisodesByContent(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := requestContext(r)
	defer cancel()

	resp, err := g.content.ListEpisodesByContent(ctx, &contentpb.ListEpisodesByContentRequest{
		ContentId: r.PathValue("contentId"),
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}

	g.writeProto(w, resp)
}

func (g *gateway) updateEpisode(w http.ResponseWriter, r *http.Request) {
	var req updateEpisodeRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	ctx, cancel := requestContext(r)
	defer cancel()

	resp, err := g.content.UpdateEpisode(ctx, &contentpb.UpdateEpisodeRequest{
		Id:              r.PathValue("episodeId"),
		Title:           req.Title,
		Description:     req.Description,
		SeasonNumber:    req.SeasonNumber,
		EpisodeNumber:   req.EpisodeNumber,
		DurationSeconds: req.DurationSeconds,
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}

	g.writeProto(w, resp)
}

func (g *gateway) deleteEpisode(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := requestContext(r)
	defer cancel()

	resp, err := g.content.DeleteEpisode(ctx, &contentpb.DeleteEpisodeRequest{
		Id: r.PathValue("episodeId"),
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}

	g.writeProto(w, resp)
}

func (g *gateway) setSkipIntroMetadata(w http.ResponseWriter, r *http.Request) {
	var req setSkipIntroRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	ctx, cancel := requestContext(r)
	defer cancel()

	resp, err := g.content.SetSkipIntroMetadata(ctx, &contentpb.SetSkipIntroMetadataRequest{
		EpisodeId:      r.PathValue("episodeId"),
		SkipIntroStart: req.SkipIntroStart,
		SkipIntroEnd:   req.SkipIntroEnd,
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}

	g.writeProto(w, resp)
}

// PLAYBACK HANDLERS

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

	resp, err := g.playback.GetContinueWatching(ctx, &playbackpb.GetContinueWatchingRequest{
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

	resp, err := g.playback.GetWatchHistory(ctx, &playbackpb.GetWatchHistoryRequest{
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

	resp, err := g.playback.RemoveFromContinueWatching(ctx, &playbackpb.RemoveFromContinueWatchingRequest{
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

// HELPERS

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
	writeErrorMessage(w, status, err.Error())
}

func writeErrorMessage(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
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
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func queryInt(r *http.Request, key string, fallback int) int {
	raw := strings.TrimSpace(r.URL.Query().Get(key))
	if raw == "" {
		return fallback
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}

	return value
}

func profileTypeFromString(value string) authpb.ProfileType {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "teen":
		return authpb.ProfileType_TEEN
	case "kid":
		return authpb.ProfileType_KID
	case "adult":
		return authpb.ProfileType_ADULT
	default:
		return authpb.ProfileType_ADULT
	}
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
