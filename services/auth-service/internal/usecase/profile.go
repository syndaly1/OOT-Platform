package usecase

import (
	"context"
	"log"
	"time"

	"streaming-platform/auth-service/internal/domain"

	"github.com/google/uuid"
)

type ProfileUsecase struct {
	profileRepo  domain.ProfileRepository
	profileCache domain.ProfileCache
}

func NewProfileUsecase(profileRepo domain.ProfileRepository, profileCache domain.ProfileCache) *ProfileUsecase {
	return &ProfileUsecase{
		profileRepo:  profileRepo,
		profileCache: profileCache,
	}
}

func (uc *ProfileUsecase) CreateProfile(ctx context.Context, accountID, name, language, avatar string, profileType domain.ProfileType) (*domain.Profile, error) {
	now := time.Now()
	profile := &domain.Profile{
		ID:        uuid.NewString(),
		AccountID: accountID,
		Name:      name,
		Type:      profileType,
		Language:  language,
		Avatar:    avatar,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := uc.profileRepo.Create(ctx, profile); err != nil {
		return nil, err
	}
	_ = uc.profileCache.InvalidateProfiles(ctx, accountID)
	return profile, nil
}

func (uc *ProfileUsecase) GetProfiles(ctx context.Context, accountID string) ([]*domain.Profile, error) {
	cached, err := uc.profileCache.GetProfiles(ctx, accountID)
	if err == nil && cached != nil {
		log.Printf("profiles cache hit for account %s", accountID)
		return cached, nil
	}

	profiles, err := uc.profileRepo.GetByAccountID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	_ = uc.profileCache.SetProfiles(ctx, accountID, profiles)
	return profiles, nil
}

func (uc *ProfileUsecase) GetProfileByID(ctx context.Context, id string) (*domain.Profile, error) {
	cached, err := uc.profileCache.GetProfile(ctx, id)
	if err == nil && cached != nil {
		return cached, nil
	}

	profile, err := uc.profileRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	_ = uc.profileCache.SetProfile(ctx, profile)
	return profile, nil
}

func (uc *ProfileUsecase) UpdateProfile(ctx context.Context, id, name, language, avatar string) error {
	profile, err := uc.profileRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	profile.Name = name
	profile.Language = language
	profile.Avatar = avatar
	profile.UpdatedAt = time.Now()

	if err := uc.profileRepo.Update(ctx, profile); err != nil {
		return err
	}

	_ = uc.profileCache.InvalidateProfile(ctx, id)
	_ = uc.profileCache.InvalidateProfiles(ctx, profile.AccountID)
	return nil
}

func (uc *ProfileUsecase) DeleteProfile(ctx context.Context, id string) error {
	profile, err := uc.profileRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := uc.profileRepo.Delete(ctx, id); err != nil {
		return err
	}

	_ = uc.profileCache.InvalidateProfile(ctx, id)
	_ = uc.profileCache.InvalidateProfiles(ctx, profile.AccountID)
	return nil
}
