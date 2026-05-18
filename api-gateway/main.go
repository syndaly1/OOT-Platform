package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	authpb "streaming-platform/proto/auth"
	playbackpb "streaming-platform/proto/playback"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type gateway struct {
	auth     authpb.AuthServiceClient
	playback playbackpb.PlaybackServiceClient
	marshal  protojson.MarshalOptions
}

func main() {
	authAddr := getenv("AUTH_SERVICE_URL", "localhost:50051")
	playbackAddr := getenv("PLAYBACK_SERVICE_URL", "localhost:50053")

	authConn, err := grpc.NewClient(authAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("auth grpc connect:", err)
	}
	defer authConn.Close()

	playbackConn, err := grpc.NewClient(playbackAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("playback grpc connect:", err)
	}
	defer playbackConn.Close()

	g := &gateway{
		auth:     authpb.NewAuthServiceClient(authConn),
		playback: playbackpb.NewPlaybackServiceClient(playbackConn),
		marshal: protojson.MarshalOptions{
			UseProtoNames: false,
		},
	}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/auth/register", g.register)
	mux.HandleFunc("POST /api/auth/login", g.login)
	mux.HandleFunc("POST /api/auth/refresh", g.refreshToken)
	mux.HandleFunc("POST /api/auth/logout", g.logout)
	mux.HandleFunc("POST /api/auth/verify-email", g.verifyEmail)
	mux.HandleFunc("POST /api/auth/forgot-password", g.forgotPassword)
	mux.HandleFunc("POST /api/auth/reset-password", g.resetPassword)

	mux.HandleFunc("GET /api/auth/account/{accountId}", g.getAccount)
	mux.HandleFunc("PUT /api/auth/account/{accountId}", g.updateAccount)
	mux.HandleFunc("DELETE /api/auth/account/{accountId}", g.deleteAccount)

	mux.HandleFunc("POST /api/profiles", g.createProfile)
	mux.HandleFunc("GET /api/profiles", g.getProfiles)
	mux.HandleFunc("GET /api/profiles/{profileId}", g.getProfileByID)
	mux.HandleFunc("PUT /api/profiles/{profileId}", g.updateProfile)
	mux.HandleFunc("DELETE /api/profiles/{profileId}", g.deleteProfile)

	mux.HandleFunc("POST /api/playback/start", g.startPlayback)
	mux.HandleFunc("PUT /api/playback/progress", g.updateProgress)
	mux.HandleFunc("GET /api/playback/continue-watching/{profileId}", g.getContinueWatching)
	mux.HandleFunc("GET /api/playback/history/{profileId}", g.getWatchHistory)
	mux.HandleFunc("DELETE /api/playback/continue-watching/{profileId}/{contentId}", g.removeFromContinueWatching)

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
