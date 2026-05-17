package usecase

import (
	"context"
	"errors"
	"streaming-platform/auth-service/config"
	"streaming-platform/auth-service/internal/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	accountRepo       domain.AccountRepository
	refreshTokenRepo  domain.RefreshTokenRepository
	emailVerRepo      domain.EmailVerificationRepository
	passwordResetRepo domain.PasswordResetRepository
	nats              NATSPublisher
	cfg               *config.Config
}

type NATSPublisher interface {
	Publish(subject string, data []byte) error
}

func NewAuthUsecase(
	accountRepo domain.AccountRepository,
	refreshTokenRepo domain.RefreshTokenRepository,
	emailVerRepo domain.EmailVerificationRepository,
	passwordResetRepo domain.PasswordResetRepository,
	nats NATSPublisher,
	cfg *config.Config,
) *AuthUsecase {
	return &AuthUsecase{
		accountRepo:       accountRepo,
		refreshTokenRepo:  refreshTokenRepo,
		emailVerRepo:      emailVerRepo,
		passwordResetRepo: passwordResetRepo,
		nats:              nats,
		cfg:               cfg,
	}
}

// ── Register ──

type RegisterInput struct {
	Email    string
	Username string
	Password string
}

type RegisterOutput struct {
	AccountID string
	Message   string
}

func (uc *AuthUsecase) Register(ctx context.Context, in RegisterInput) (*RegisterOutput, error) {
	existing, _ := uc.accountRepo.GetByEmail(ctx, in.Email)
	if existing != nil {
		return nil, errors.New("email already in use")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	account := &domain.Account{
		ID:         uuid.NewString(),
		Email:      in.Email,
		Username:   in.Username,
		Password:   string(hashed),
		IsVerified: false,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := uc.accountRepo.Create(ctx, account); err != nil {
		return nil, err
	}

	// создаём токен верификации email
	evToken := uuid.NewString()
	ev := &domain.EmailVerification{
		ID:        uuid.NewString(),
		AccountID: account.ID,
		Token:     evToken,
		ExpiresAt: now.Add(24 * time.Hour),
		CreatedAt: now,
	}
	if err := uc.emailVerRepo.Create(ctx, ev); err != nil {
		return nil, err
	}

	payload := []byte(`{"account_id":"` + account.ID + `","email":"` + account.Email + `","token":"` + evToken + `"}`)
	_ = uc.nats.Publish("user.registered", payload)

	return &RegisterOutput{AccountID: account.ID, Message: "registered successfully, check your email"}, nil
}

func (uc *AuthUsecase) VerifyEmail(ctx context.Context, token string) error {
	ev, err := uc.emailVerRepo.GetByToken(ctx, token)
	if err != nil {
		return errors.New("invalid or expired token")
	}

	if time.Now().After(ev.ExpiresAt) {
		return errors.New("token expired")
	}

	if err := uc.accountRepo.MarkVerified(ctx, ev.AccountID); err != nil {
		return err
	}

	return uc.emailVerRepo.Delete(ctx, ev.ID)
}

type LoginOutput struct {
	AccessToken  string
	RefreshToken string
	AccountID    string
}

func (uc *AuthUsecase) Login(ctx context.Context, email, password string) (*LoginOutput, error) {
	account, err := uc.accountRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !account.IsVerified {
		return nil, errors.New("email not verified")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	accessToken, err := uc.generateAccessToken(account.ID)
	if err != nil {
		return nil, err
	}

	refreshToken := uuid.NewString()
	now := time.Now()
	rt := &domain.RefreshToken{
		ID:        uuid.NewString(),
		AccountID: account.ID,
		Token:     refreshToken,
		ExpiresAt: now.Add(uc.cfg.JWT.RefreshTTL),
		CreatedAt: now,
	}
	if err := uc.refreshTokenRepo.Create(ctx, rt); err != nil {
		return nil, err
	}

	return &LoginOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		AccountID:    account.ID,
	}, nil
}

func (uc *AuthUsecase) RefreshToken(ctx context.Context, token string) (*LoginOutput, error) {
	rt, err := uc.refreshTokenRepo.GetByToken(ctx, token)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if time.Now().After(rt.ExpiresAt) {
		return nil, errors.New("refresh token expired")
	}

	// удаляем старый, создаём новый
	if err := uc.refreshTokenRepo.Delete(ctx, token); err != nil {
		return nil, err
	}

	accessToken, err := uc.generateAccessToken(rt.AccountID)
	if err != nil {
		return nil, err
	}

	newRefresh := uuid.NewString()
	now := time.Now()
	newRT := &domain.RefreshToken{
		ID:        uuid.NewString(),
		AccountID: rt.AccountID,
		Token:     newRefresh,
		ExpiresAt: now.Add(uc.cfg.JWT.RefreshTTL),
		CreatedAt: now,
	}
	if err := uc.refreshTokenRepo.Create(ctx, newRT); err != nil {
		return nil, err
	}

	return &LoginOutput{
		AccessToken:  accessToken,
		RefreshToken: newRefresh,
		AccountID:    rt.AccountID,
	}, nil
}

func (uc *AuthUsecase) Logout(ctx context.Context, refreshToken string) error {
	return uc.refreshTokenRepo.Delete(ctx, refreshToken)
}

func (uc *AuthUsecase) ForgotPassword(ctx context.Context, email string) error {
	account, err := uc.accountRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil
	}

	now := time.Now()
	resetToken := uuid.NewString()
	pr := &domain.PasswordReset{
		ID:        uuid.NewString(),
		AccountID: account.ID,
		Token:     resetToken,
		ExpiresAt: now.Add(1 * time.Hour),
		CreatedAt: now,
	}
	if err := uc.passwordResetRepo.Create(ctx, pr); err != nil {
		return err
	}

	payload := []byte(`{"email":"` + email + `","token":"` + resetToken + `"}`)
	_ = uc.nats.Publish("password.reset.requested", payload)

	return nil
}

func (uc *AuthUsecase) ResetPassword(ctx context.Context, token, newPassword string) error {
	pr, err := uc.passwordResetRepo.GetByToken(ctx, token)
	if err != nil {
		return errors.New("invalid or expired token")
	}

	if time.Now().After(pr.ExpiresAt) {
		return errors.New("token expired")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	account, err := uc.accountRepo.GetByID(ctx, pr.AccountID)
	if err != nil {
		return err
	}

	account.Password = string(hashed)
	account.UpdatedAt = time.Now()

	if err := uc.accountRepo.Update(ctx, account); err != nil {
		return err
	}

	return uc.passwordResetRepo.Delete(ctx, pr.ID)
}

func (uc *AuthUsecase) generateAccessToken(accountID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": accountID,
		"exp": time.Now().Add(uc.cfg.JWT.AccessTTL).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.cfg.JWT.Secret))
}
