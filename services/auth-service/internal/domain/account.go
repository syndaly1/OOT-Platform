package domain

import "time"

type Account struct {
	ID         string
	Email      string
	Username   string
	Password   string
	IsVerified bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type ProfileType string

const (
	ProfileTypeAdult ProfileType = "adult"
	ProfileTypeTeen  ProfileType = "teen"
	ProfileTypeKid   ProfileType = "kid"
)

type Profile struct {
	ID        string
	AccountID string
	Name      string
	Type      ProfileType
	Language  string
	Avatar    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type EmailVerification struct {
	ID        string
	AccountID string
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type PasswordReset struct {
	ID        string
	AccountID string
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type RefreshToken struct {
	ID        string
	AccountID string
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}
