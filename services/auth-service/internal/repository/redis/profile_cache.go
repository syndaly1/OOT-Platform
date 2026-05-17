package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"streaming-platform/auth-service/internal/domain"

	"github.com/redis/go-redis/v9"
)

type ProfileCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewProfileCache(client *redis.Client) *ProfileCache {
	return &ProfileCache{
		client: client,
		ttl:    10 * time.Minute,
	}
}

func (c *ProfileCache) SetProfiles(ctx context.Context, accountID string, profiles []*domain.Profile) error {
	key := fmt.Sprintf("profiles:account:%s", accountID)
	data, err := json.Marshal(profiles)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, c.ttl).Err()
}

func (c *ProfileCache) GetProfiles(ctx context.Context, accountID string) ([]*domain.Profile, error) {
	key := fmt.Sprintf("profiles:account:%s", accountID)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	var profiles []*domain.Profile
	if err := json.Unmarshal(data, &profiles); err != nil {
		return nil, err
	}
	return profiles, nil
}

func (c *ProfileCache) InvalidateProfiles(ctx context.Context, accountID string) error {
	key := fmt.Sprintf("profiles:account:%s", accountID)
	return c.client.Del(ctx, key).Err()
}

func (c *ProfileCache) SetProfile(ctx context.Context, profile *domain.Profile) error {
	key := fmt.Sprintf("profile:%s", profile.ID)
	data, err := json.Marshal(profile)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, c.ttl).Err()
}

func (c *ProfileCache) GetProfile(ctx context.Context, profileID string) (*domain.Profile, error) {
	key := fmt.Sprintf("profile:%s", profileID)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	var profile domain.Profile
	if err := json.Unmarshal(data, &profile); err != nil {
		return nil, err
	}
	return &profile, nil
}

func (c *ProfileCache) InvalidateProfile(ctx context.Context, profileID string) error {
	key := fmt.Sprintf("profile:%s", profileID)
	return c.client.Del(ctx, key).Err()
}
