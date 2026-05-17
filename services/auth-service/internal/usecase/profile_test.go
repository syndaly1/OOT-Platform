package usecase_test

import (
	"context"
	"errors"
	"testing"

	"streaming-platform/auth-service/internal/domain"
	"streaming-platform/auth-service/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockProfileRepo struct{ mock.Mock }

func (m *mockProfileRepo) Create(ctx context.Context, p *domain.Profile) error {
	return m.Called(ctx, p).Error(0)
}
func (m *mockProfileRepo) GetByAccountID(ctx context.Context, accountID string) ([]*domain.Profile, error) {
	args := m.Called(ctx, accountID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Profile), args.Error(1)
}
func (m *mockProfileRepo) GetByID(ctx context.Context, id string) (*domain.Profile, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Profile), args.Error(1)
}
func (m *mockProfileRepo) Update(ctx context.Context, p *domain.Profile) error {
	return m.Called(ctx, p).Error(0)
}
func (m *mockProfileRepo) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

type mockProfileCache struct{ mock.Mock }

func (m *mockProfileCache) SetProfiles(ctx context.Context, accountID string, profiles []*domain.Profile) error {
	return m.Called(ctx, accountID, profiles).Error(0)
}
func (m *mockProfileCache) GetProfiles(ctx context.Context, accountID string) ([]*domain.Profile, error) {
	args := m.Called(ctx, accountID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Profile), args.Error(1)
}
func (m *mockProfileCache) InvalidateProfiles(ctx context.Context, accountID string) error {
	return m.Called(ctx, accountID).Error(0)
}
func (m *mockProfileCache) SetProfile(ctx context.Context, profile *domain.Profile) error {
	return m.Called(ctx, profile).Error(0)
}
func (m *mockProfileCache) GetProfile(ctx context.Context, profileID string) (*domain.Profile, error) {
	args := m.Called(ctx, profileID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Profile), args.Error(1)
}
func (m *mockProfileCache) InvalidateProfile(ctx context.Context, profileID string) error {
	return m.Called(ctx, profileID).Error(0)
}

func TestCreateProfile_Success(t *testing.T) {
	profileRepo := new(mockProfileRepo)
	profileCache := new(mockProfileCache)

	profileRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Profile")).Return(nil)
	profileCache.On("InvalidateProfiles", mock.Anything, "acc-123").Return(nil)

	uc := usecase.NewProfileUsecase(profileRepo, profileCache)

	profile, err := uc.CreateProfile(context.Background(), "acc-123", "John", "en", "", domain.ProfileTypeAdult)

	assert.NoError(t, err)
	assert.NotNil(t, profile)
	assert.Equal(t, "John", profile.Name)
	assert.Equal(t, "acc-123", profile.AccountID)
}

func TestGetProfiles_CacheHit(t *testing.T) {
	profileRepo := new(mockProfileRepo)
	profileCache := new(mockProfileCache)

	cached := []*domain.Profile{
		{ID: "p-1", AccountID: "acc-123", Name: "John"},
	}
	profileCache.On("GetProfiles", mock.Anything, "acc-123").Return(cached, nil)

	uc := usecase.NewProfileUsecase(profileRepo, profileCache)

	profiles, err := uc.GetProfiles(context.Background(), "acc-123")

	assert.NoError(t, err)
	assert.Len(t, profiles, 1)
	profileRepo.AssertNotCalled(t, "GetByAccountID")
}

func TestGetProfiles_CacheMiss(t *testing.T) {
	profileRepo := new(mockProfileRepo)
	profileCache := new(mockProfileCache)

	dbProfiles := []*domain.Profile{
		{ID: "p-1", AccountID: "acc-123", Name: "John"},
	}
	profileCache.On("GetProfiles", mock.Anything, "acc-123").Return(nil, errors.New("cache miss"))
	profileRepo.On("GetByAccountID", mock.Anything, "acc-123").Return(dbProfiles, nil)
	profileCache.On("SetProfiles", mock.Anything, "acc-123", dbProfiles).Return(nil)

	uc := usecase.NewProfileUsecase(profileRepo, profileCache)

	profiles, err := uc.GetProfiles(context.Background(), "acc-123")

	assert.NoError(t, err)
	assert.Len(t, profiles, 1)
	profileRepo.AssertCalled(t, "GetByAccountID", mock.Anything, "acc-123")
}

func TestDeleteProfile_Success(t *testing.T) {
	profileRepo := new(mockProfileRepo)
	profileCache := new(mockProfileCache)

	profile := &domain.Profile{ID: "p-1", AccountID: "acc-123", Name: "John"}
	profileRepo.On("GetByID", mock.Anything, "p-1").Return(profile, nil)
	profileRepo.On("Delete", mock.Anything, "p-1").Return(nil)
	profileCache.On("InvalidateProfile", mock.Anything, "p-1").Return(nil)
	profileCache.On("InvalidateProfiles", mock.Anything, "acc-123").Return(nil)

	uc := usecase.NewProfileUsecase(profileRepo, profileCache)

	err := uc.DeleteProfile(context.Background(), "p-1")

	assert.NoError(t, err)
	profileRepo.AssertExpectations(t)
	profileCache.AssertExpectations(t)
}
