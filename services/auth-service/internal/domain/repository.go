package domain

import "context"

type AccountRepository interface {
	Create(ctx context.Context, account *Account) error
	GetByID(ctx context.Context, id string) (*Account, error)
	GetByEmail(ctx context.Context, email string) (*Account, error)
	Update(ctx context.Context, account *Account) error
	Delete(ctx context.Context, id string) error
	MarkVerified(ctx context.Context, id string) error
}

type ProfileRepository interface {
	Create(ctx context.Context, profile *Profile) error
	GetByAccountID(ctx context.Context, accountID string) ([]*Profile, error)
	GetByID(ctx context.Context, id string) (*Profile, error)
	Update(ctx context.Context, profile *Profile) error
	Delete(ctx context.Context, id string) error
}

type EmailVerificationRepository interface {
	Create(ctx context.Context, ev *EmailVerification) error
	GetByToken(ctx context.Context, token string) (*EmailVerification, error)
	Delete(ctx context.Context, id string) error
}

type PasswordResetRepository interface {
	Create(ctx context.Context, pr *PasswordReset) error
	GetByToken(ctx context.Context, token string) (*PasswordReset, error)
	Delete(ctx context.Context, id string) error
}

type RefreshTokenRepository interface {
	Create(ctx context.Context, rt *RefreshToken) error
	GetByToken(ctx context.Context, token string) (*RefreshToken, error)
	Delete(ctx context.Context, token string) error
	DeleteByAccountID(ctx context.Context, accountID string) error
}

type ProfileCache interface {
	SetProfiles(ctx context.Context, accountID string, profiles []*Profile) error
	GetProfiles(ctx context.Context, accountID string) ([]*Profile, error)
	InvalidateProfiles(ctx context.Context, accountID string) error
	SetProfile(ctx context.Context, profile *Profile) error
	GetProfile(ctx context.Context, profileID string) (*Profile, error)
	InvalidateProfile(ctx context.Context, profileID string) error
}
