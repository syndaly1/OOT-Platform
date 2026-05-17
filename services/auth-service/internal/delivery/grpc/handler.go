package grpc

import (
	"context"
	"streaming-platform/auth-service/internal/domain"
	"streaming-platform/auth-service/internal/usecase"
	pb "streaming-platform/proto/auth"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	authUC    *usecase.AuthUsecase
	accountUC *usecase.AccountUsecase
	profileUC *usecase.ProfileUsecase
}

func NewAuthHandler(
	authUC *usecase.AuthUsecase,
	accountUC *usecase.AccountUsecase,
	profileUC *usecase.ProfileUsecase,
) *AuthHandler {
	return &AuthHandler{
		authUC:    authUC,
		accountUC: accountUC,
		profileUC: profileUC,
	}
}

func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	out, err := h.authUC.Register(ctx, usecase.RegisterInput{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &pb.RegisterResponse{AccountId: out.AccountID, Message: out.Message}, nil
}

func (h *AuthHandler) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	if err := h.authUC.VerifyEmail(ctx, req.Token); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &pb.VerifyEmailResponse{Success: true, Message: "email verified"}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	out, err := h.authUC.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	return &pb.LoginResponse{
		AccessToken:  out.AccessToken,
		RefreshToken: out.RefreshToken,
		AccountId:    out.AccountID,
	}, nil
}

func (h *AuthHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	out, err := h.authUC.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	return &pb.RefreshTokenResponse{
		AccessToken:  out.AccessToken,
		RefreshToken: out.RefreshToken,
	}, nil
}

func (h *AuthHandler) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	if err := h.authUC.Logout(ctx, req.RefreshToken); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.LogoutResponse{Success: true}, nil
}

func (h *AuthHandler) ForgotPassword(ctx context.Context, req *pb.ForgotPasswordRequest) (*pb.ForgotPasswordResponse, error) {
	if err := h.authUC.ForgotPassword(ctx, req.Email); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.ForgotPasswordResponse{Message: "if email exists, reset link was sent"}, nil
}

func (h *AuthHandler) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	if err := h.authUC.ResetPassword(ctx, req.Token, req.NewPassword); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &pb.ResetPasswordResponse{Success: true, Message: "password reset successfully"}, nil
}

func (h *AuthHandler) GetAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	acc, err := h.accountUC.GetAccount(ctx, req.AccountId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "account not found")
	}
	return &pb.GetAccountResponse{
		AccountId:  acc.ID,
		Email:      acc.Email,
		Username:   acc.Username,
		IsVerified: acc.IsVerified,
		CreatedAt:  acc.CreatedAt.String(),
	}, nil
}

func (h *AuthHandler) UpdateAccount(ctx context.Context, req *pb.UpdateAccountRequest) (*pb.UpdateAccountResponse, error) {
	if err := h.accountUC.UpdateAccount(ctx, req.AccountId, req.Username); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.UpdateAccountResponse{Success: true}, nil
}

func (h *AuthHandler) DeleteAccount(ctx context.Context, req *pb.DeleteAccountRequest) (*pb.DeleteAccountResponse, error) {
	if err := h.accountUC.DeleteAccount(ctx, req.AccountId); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.DeleteAccountResponse{Success: true}, nil
}

func (h *AuthHandler) CreateProfile(ctx context.Context, req *pb.CreateProfileRequest) (*pb.CreateProfileResponse, error) {
	pt := domain.ProfileType(req.Type.String())
	profile, err := h.profileUC.CreateProfile(ctx, req.AccountId, req.Name, req.Language, req.Avatar, pt)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.CreateProfileResponse{ProfileId: profile.ID, Message: "profile created"}, nil
}

func (h *AuthHandler) GetProfiles(ctx context.Context, req *pb.GetProfilesRequest) (*pb.GetProfilesResponse, error) {
	profiles, err := h.profileUC.GetProfiles(ctx, req.AccountId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var pbProfiles []*pb.Profile
	for _, p := range profiles {
		pbProfiles = append(pbProfiles, toPBProfile(p))
	}
	return &pb.GetProfilesResponse{Profiles: pbProfiles}, nil
}

func (h *AuthHandler) GetProfileById(ctx context.Context, req *pb.GetProfileByIdRequest) (*pb.GetProfileByIdResponse, error) {
	profile, err := h.profileUC.GetProfileByID(ctx, req.ProfileId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "profile not found")
	}
	return &pb.GetProfileByIdResponse{Profile: toPBProfile(profile)}, nil
}

func (h *AuthHandler) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	if err := h.profileUC.UpdateProfile(ctx, req.ProfileId, req.Name, req.Language, req.Avatar); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.UpdateProfileResponse{Success: true}, nil
}

func (h *AuthHandler) DeleteProfile(ctx context.Context, req *pb.DeleteProfileRequest) (*pb.DeleteProfileResponse, error) {
	if err := h.profileUC.DeleteProfile(ctx, req.ProfileId); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.DeleteProfileResponse{Success: true}, nil
}

func toPBProfile(p *domain.Profile) *pb.Profile {
	return &pb.Profile{
		ProfileId: p.ID,
		AccountId: p.AccountID,
		Name:      p.Name,
		Language:  p.Language,
		Avatar:    p.Avatar,
		CreatedAt: p.CreatedAt.String(),
	}
}
