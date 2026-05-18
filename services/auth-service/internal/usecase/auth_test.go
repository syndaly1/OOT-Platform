package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"streaming-platform/auth-service/config"
	"streaming-platform/auth-service/internal/domain"
	"streaming-platform/auth-service/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type mockAccountRepo struct{ mock.Mock }

func (m *mockAccountRepo) Create(ctx context.Context, a *domain.Account) error {
	return m.Called(ctx, a).Error(0)
}
func (m *mockAccountRepo) GetByID(ctx context.Context, id string) (*domain.Account, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Account), args.Error(1)
}
func (m *mockAccountRepo) GetByEmail(ctx context.Context, email string) (*domain.Account, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Account), args.Error(1)
}
func (m *mockAccountRepo) Update(ctx context.Context, a *domain.Account) error {
	return m.Called(ctx, a).Error(0)
}
func (m *mockAccountRepo) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockAccountRepo) MarkVerified(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

type mockRefreshTokenRepo struct{ mock.Mock }

func (m *mockRefreshTokenRepo) Create(ctx context.Context, rt *domain.RefreshToken) error {
	return m.Called(ctx, rt).Error(0)
}
func (m *mockRefreshTokenRepo) GetByToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RefreshToken), args.Error(1)
}
func (m *mockRefreshTokenRepo) Delete(ctx context.Context, token string) error {
	return m.Called(ctx, token).Error(0)
}
func (m *mockRefreshTokenRepo) DeleteByAccountID(ctx context.Context, accountID string) error {
	return m.Called(ctx, accountID).Error(0)
}

type mockEmailVerRepo struct{ mock.Mock }

func (m *mockEmailVerRepo) Create(ctx context.Context, ev *domain.EmailVerification) error {
	return m.Called(ctx, ev).Error(0)
}
func (m *mockEmailVerRepo) GetByToken(ctx context.Context, token string) (*domain.EmailVerification, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.EmailVerification), args.Error(1)
}
func (m *mockEmailVerRepo) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

type mockPasswordResetRepo struct{ mock.Mock }

func (m *mockPasswordResetRepo) Create(ctx context.Context, pr *domain.PasswordReset) error {
	return m.Called(ctx, pr).Error(0)
}
func (m *mockPasswordResetRepo) GetByToken(ctx context.Context, token string) (*domain.PasswordReset, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PasswordReset), args.Error(1)
}
func (m *mockPasswordResetRepo) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

type mockNATS struct{ mock.Mock }

func (m *mockNATS) Publish(subject string, data []byte) error {
	return m.Called(subject, data).Error(0)
}

func newTestAuthUC(
	accountRepo *mockAccountRepo,
	rtRepo *mockRefreshTokenRepo,
	evRepo *mockEmailVerRepo,
	prRepo *mockPasswordResetRepo,
	nats *mockNATS,
) *usecase.AuthUsecase {
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:     "testsecret",
			AccessTTL:  15 * time.Minute,
			RefreshTTL: 7 * 24 * time.Hour,
		},
	}
	return usecase.NewAuthUsecase(accountRepo, rtRepo, evRepo, prRepo, nats, cfg)
}

func TestRegister_Success(t *testing.T) {
	accountRepo := new(mockAccountRepo)
	rtRepo := new(mockRefreshTokenRepo)
	evRepo := new(mockEmailVerRepo)
	prRepo := new(mockPasswordResetRepo)
	nats := new(mockNATS)

	accountRepo.On("GetByEmail", mock.Anything, "test@test.com").Return(nil, errors.New("not found"))
	accountRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Account")).Return(nil)
	evRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.EmailVerification")).Return(nil)
	nats.On("Publish", "user.registered", mock.Anything).Return(nil)

	uc := newTestAuthUC(accountRepo, rtRepo, evRepo, prRepo, nats)

	out, err := uc.Register(context.Background(), usecase.RegisterInput{
		Email:    "test@test.com",
		Username: "testuser",
		Password: "password123",
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, out.AccountID)
	assert.Contains(t, out.Message, "registered")
	accountRepo.AssertExpectations(t)
	evRepo.AssertExpectations(t)
	nats.AssertExpectations(t)
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	accountRepo := new(mockAccountRepo)
	rtRepo := new(mockRefreshTokenRepo)
	evRepo := new(mockEmailVerRepo)
	prRepo := new(mockPasswordResetRepo)
	nats := new(mockNATS)

	existing := &domain.Account{ID: "123", Email: "test@test.com"}
	accountRepo.On("GetByEmail", mock.Anything, "test@test.com").Return(existing, nil)

	uc := newTestAuthUC(accountRepo, rtRepo, evRepo, prRepo, nats)

	out, err := uc.Register(context.Background(), usecase.RegisterInput{
		Email:    "test@test.com",
		Username: "testuser",
		Password: "password123",
	})

	assert.Error(t, err)
	assert.Nil(t, out)
	assert.Equal(t, "email already in use", err.Error())
}

func TestLogin_Success(t *testing.T) {
	accountRepo := new(mockAccountRepo)
	rtRepo := new(mockRefreshTokenRepo)
	evRepo := new(mockEmailVerRepo)
	prRepo := new(mockPasswordResetRepo)
	nats := new(mockNATS)

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	require.NoError(t, err)

	account := &domain.Account{
		ID:         "acc-123",
		Email:      "test@test.com",
		Password:   string(hashedPwd),
		IsVerified: true,
	}

	accountRepo.On("GetByEmail", mock.Anything, "test@test.com").Return(account, nil)
	rtRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.RefreshToken")).Return(nil)

	uc := newTestAuthUC(accountRepo, rtRepo, evRepo, prRepo, nats)

	out, err := uc.Login(context.Background(), "test@test.com", "password123")

	require.NoError(t, err)
	require.NotNil(t, out)

	assert.NotEmpty(t, out.AccessToken)
	assert.NotEmpty(t, out.RefreshToken)
	assert.Equal(t, "acc-123", out.AccountID)

	accountRepo.AssertExpectations(t)
	rtRepo.AssertExpectations(t)
}

func TestLogin_NotVerified(t *testing.T) {
	accountRepo := new(mockAccountRepo)
	rtRepo := new(mockRefreshTokenRepo)
	evRepo := new(mockEmailVerRepo)
	prRepo := new(mockPasswordResetRepo)
	nats := new(mockNATS)

	account := &domain.Account{
		ID:         "acc-123",
		Email:      "test@test.com",
		Password:   "anything",
		IsVerified: false,
	}
	accountRepo.On("GetByEmail", mock.Anything, "test@test.com").Return(account, nil)

	uc := newTestAuthUC(accountRepo, rtRepo, evRepo, prRepo, nats)

	out, err := uc.Login(context.Background(), "test@test.com", "password123")

	assert.Error(t, err)
	assert.Nil(t, out)
	assert.Equal(t, "email not verified", err.Error())
}

func TestLogin_WrongPassword(t *testing.T) {
	accountRepo := new(mockAccountRepo)
	rtRepo := new(mockRefreshTokenRepo)
	evRepo := new(mockEmailVerRepo)
	prRepo := new(mockPasswordResetRepo)
	nats := new(mockNATS)

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	require.NoError(t, err)

	account := &domain.Account{
		ID:         "acc-123",
		Email:      "test@test.com",
		Password:   string(hashedPwd),
		IsVerified: true,
	}

	accountRepo.On("GetByEmail", mock.Anything, "test@test.com").Return(account, nil)

	uc := newTestAuthUC(accountRepo, rtRepo, evRepo, prRepo, nats)

	out, err := uc.Login(context.Background(), "test@test.com", "wrongpassword")

	require.Error(t, err)
	assert.Nil(t, out)
	assert.Equal(t, "invalid credentials", err.Error())

	accountRepo.AssertExpectations(t)
}

func TestVerifyEmail_Success(t *testing.T) {
	accountRepo := new(mockAccountRepo)
	rtRepo := new(mockRefreshTokenRepo)
	evRepo := new(mockEmailVerRepo)
	prRepo := new(mockPasswordResetRepo)
	nats := new(mockNATS)

	ev := &domain.EmailVerification{
		ID:        "ev-123",
		AccountID: "acc-123",
		Token:     "valid-token",
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}
	evRepo.On("GetByToken", mock.Anything, "valid-token").Return(ev, nil)
	accountRepo.On("MarkVerified", mock.Anything, "acc-123").Return(nil)
	evRepo.On("Delete", mock.Anything, "ev-123").Return(nil)

	uc := newTestAuthUC(accountRepo, rtRepo, evRepo, prRepo, nats)

	err := uc.VerifyEmail(context.Background(), "valid-token")

	assert.NoError(t, err)
	evRepo.AssertExpectations(t)
	accountRepo.AssertExpectations(t)
}

func TestVerifyEmail_ExpiredToken(t *testing.T) {
	accountRepo := new(mockAccountRepo)
	rtRepo := new(mockRefreshTokenRepo)
	evRepo := new(mockEmailVerRepo)
	prRepo := new(mockPasswordResetRepo)
	nats := new(mockNATS)

	ev := &domain.EmailVerification{
		ID:        "ev-123",
		AccountID: "acc-123",
		Token:     "expired-token",
		ExpiresAt: time.Now().Add(-1 * time.Hour), // уже истёк
	}
	evRepo.On("GetByToken", mock.Anything, "expired-token").Return(ev, nil)

	uc := newTestAuthUC(accountRepo, rtRepo, evRepo, prRepo, nats)

	err := uc.VerifyEmail(context.Background(), "expired-token")

	assert.Error(t, err)
	assert.Equal(t, "token expired", err.Error())
}

func TestLogout_Success(t *testing.T) {
	accountRepo := new(mockAccountRepo)
	rtRepo := new(mockRefreshTokenRepo)
	evRepo := new(mockEmailVerRepo)
	prRepo := new(mockPasswordResetRepo)
	nats := new(mockNATS)

	rtRepo.On("Delete", mock.Anything, "some-refresh-token").Return(nil)

	uc := newTestAuthUC(accountRepo, rtRepo, evRepo, prRepo, nats)

	err := uc.Logout(context.Background(), "some-refresh-token")

	assert.NoError(t, err)
	rtRepo.AssertExpectations(t)
}
