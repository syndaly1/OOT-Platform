package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"streaming-platform/services/playback-service/internal/domain"
)

type mockPlaybackRepository struct {
	startPlaybackFn              func(ctx context.Context, session *domain.PlaybackSession) error
	updateProgressFn             func(ctx context.Context, session *domain.PlaybackSession) error
	getProgressFn                func(ctx context.Context, profileID, contentID, episodeID string) (*domain.PlaybackSession, error)
	getResumeDataFn              func(ctx context.Context, profileID, contentID, episodeID string) (*domain.PlaybackSession, error)
	getContinueWatchingFn        func(ctx context.Context, profileID string) ([]*domain.PlaybackSession, error)
	markCompletedFn              func(ctx context.Context, profileID, contentID, episodeID string) error
	removeFromContinueWatchingFn func(ctx context.Context, profileID, contentID, episodeID string) error
	getWatchHistoryFn            func(ctx context.Context, profileID string) ([]*domain.PlaybackSession, error)
	clearWatchHistoryFn          func(ctx context.Context, profileID string) error
	getProfileStatsFn            func(ctx context.Context, profileID string) (int32, int32, int32, error)
	getPlaybackSessionFn         func(ctx context.Context, sessionID string) (*domain.PlaybackSession, error)
	endPlaybackSessionFn         func(ctx context.Context, sessionID string) error
	deleteProgressFn             func(ctx context.Context, profileID, contentID, episodeID string) error
	getLastWatchedContentFn      func(ctx context.Context, profileID string) (*domain.PlaybackSession, error)
}

func (m *mockPlaybackRepository) StartPlayback(ctx context.Context, session *domain.PlaybackSession) error {
	return m.startPlaybackFn(ctx, session)
}

func (m *mockPlaybackRepository) UpdateProgress(ctx context.Context, session *domain.PlaybackSession) error {
	return m.updateProgressFn(ctx, session)
}

func (m *mockPlaybackRepository) GetProgress(ctx context.Context, profileID, contentID, episodeID string) (*domain.PlaybackSession, error) {
	return m.getProgressFn(ctx, profileID, contentID, episodeID)
}

func (m *mockPlaybackRepository) GetResumeData(ctx context.Context, profileID, contentID, episodeID string) (*domain.PlaybackSession, error) {
	return m.getResumeDataFn(ctx, profileID, contentID, episodeID)
}

func (m *mockPlaybackRepository) GetContinueWatching(ctx context.Context, profileID string) ([]*domain.PlaybackSession, error) {
	return m.getContinueWatchingFn(ctx, profileID)
}

func (m *mockPlaybackRepository) MarkCompleted(ctx context.Context, profileID, contentID, episodeID string) error {
	return m.markCompletedFn(ctx, profileID, contentID, episodeID)
}

func (m *mockPlaybackRepository) RemoveFromContinueWatching(ctx context.Context, profileID, contentID, episodeID string) error {
	return m.removeFromContinueWatchingFn(ctx, profileID, contentID, episodeID)
}

func (m *mockPlaybackRepository) GetWatchHistory(ctx context.Context, profileID string) ([]*domain.PlaybackSession, error) {
	return m.getWatchHistoryFn(ctx, profileID)
}

func (m *mockPlaybackRepository) ClearWatchHistory(ctx context.Context, profileID string) error {
	return m.clearWatchHistoryFn(ctx, profileID)
}

func (m *mockPlaybackRepository) GetProfileStats(ctx context.Context, profileID string) (int32, int32, int32, error) {
	return m.getProfileStatsFn(ctx, profileID)
}

func (m *mockPlaybackRepository) GetPlaybackSession(ctx context.Context, sessionID string) (*domain.PlaybackSession, error) {
	return m.getPlaybackSessionFn(ctx, sessionID)
}

func (m *mockPlaybackRepository) EndPlaybackSession(ctx context.Context, sessionID string) error {
	return m.endPlaybackSessionFn(ctx, sessionID)
}

func (m *mockPlaybackRepository) DeleteProgress(ctx context.Context, profileID, contentID, episodeID string) error {
	return m.deleteProgressFn(ctx, profileID, contentID, episodeID)
}

func (m *mockPlaybackRepository) GetLastWatchedContent(ctx context.Context, profileID string) (*domain.PlaybackSession, error) {
	return m.getLastWatchedContentFn(ctx, profileID)
}

type publishedPlaybackEvent struct {
	subject string
	data    []byte
}

type mockPlaybackPublisher struct {
	events []publishedPlaybackEvent
}

func (m *mockPlaybackPublisher) Publish(subject string, data []byte) error {
	m.events = append(m.events, publishedPlaybackEvent{
		subject: subject,
		data:    data,
	})
	return nil
}

func samplePlaybackSession() *domain.PlaybackSession {
	return &domain.PlaybackSession{
		SessionID:          "session-1",
		ProfileID:          "profile-1",
		ContentID:          "content-1",
		EpisodeID:          "episode-1",
		Device:             "web",
		CurrentPositionSec: 600,
		DurationSec:        3600,
	}
}

func TestStartPlaybackSetsInitialState(t *testing.T) {
	ctx := context.Background()
	session := samplePlaybackSession()
	session.CurrentPositionSec = 999
	session.WatchedPercent = 50
	session.Completed = true
	session.Status = "COMPLETED"

	repo := &mockPlaybackRepository{
		startPlaybackFn: func(ctx context.Context, got *domain.PlaybackSession) error {
			if got.CurrentPositionSec != 0 {
				t.Fatalf("expected CurrentPositionSec 0, got %d", got.CurrentPositionSec)
			}
			if got.WatchedPercent != 0 {
				t.Fatalf("expected WatchedPercent 0, got %f", got.WatchedPercent)
			}
			if got.Completed {
				t.Fatal("expected Completed false")
			}
			if got.Status != "PLAYING" {
				t.Fatalf("expected Status PLAYING, got %s", got.Status)
			}
			return nil
		},
	}

	uc := NewPlaybackUsecase(repo, nil)

	if err := uc.StartPlayback(ctx, session); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateProgressCalculatesPercentAndPublishesEvent(t *testing.T) {
	ctx := context.Background()
	session := samplePlaybackSession()

	repo := &mockPlaybackRepository{
		updateProgressFn: func(ctx context.Context, got *domain.PlaybackSession) error {
			if got.WatchedPercent != float32(16.666668) {
				t.Fatalf("expected watched percent 16.666668, got %f", got.WatchedPercent)
			}
			if got.Completed {
				t.Fatal("expected Completed false")
			}
			if got.Status != "PLAYING" {
				t.Fatalf("expected Status PLAYING, got %s", got.Status)
			}
			return nil
		},
	}

	publisher := &mockPlaybackPublisher{}
	uc := NewPlaybackUsecase(repo, publisher)

	if err := uc.UpdateProgress(ctx, session); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(publisher.events) != 1 {
		t.Fatalf("expected 1 published event, got %d", len(publisher.events))
	}

	if publisher.events[0].subject != SubjectProgressUpdated {
		t.Fatalf("expected subject %s, got %s", SubjectProgressUpdated, publisher.events[0].subject)
	}

	var payload map[string]string
	if err := json.Unmarshal(publisher.events[0].data, &payload); err != nil {
		t.Fatalf("failed to unmarshal event payload: %v", err)
	}

	if payload["profile_id"] != "profile-1" {
		t.Fatalf("expected profile_id profile-1, got %s", payload["profile_id"])
	}
	if payload["content_id"] != "content-1" {
		t.Fatalf("expected content_id content-1, got %s", payload["content_id"])
	}
	if payload["episode_id"] != "episode-1" {
		t.Fatalf("expected episode_id episode-1, got %s", payload["episode_id"])
	}
}

func TestUpdateProgressMarksCompletedAtNinetyPercent(t *testing.T) {
	ctx := context.Background()
	session := samplePlaybackSession()
	session.CurrentPositionSec = 3240
	session.DurationSec = 3600

	repo := &mockPlaybackRepository{
		updateProgressFn: func(ctx context.Context, got *domain.PlaybackSession) error {
			if got.WatchedPercent != 90 {
				t.Fatalf("expected watched percent 90, got %f", got.WatchedPercent)
			}
			if !got.Completed {
				t.Fatal("expected Completed true")
			}
			if got.Status != "COMPLETED" {
				t.Fatalf("expected Status COMPLETED, got %s", got.Status)
			}
			return nil
		},
	}

	publisher := &mockPlaybackPublisher{}
	uc := NewPlaybackUsecase(repo, publisher)

	if err := uc.UpdateProgress(ctx, session); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(publisher.events) != 1 {
		t.Fatalf("expected 1 published event, got %d", len(publisher.events))
	}
}

func TestUpdateProgressRepositoryErrorDoesNotPublishEvent(t *testing.T) {
	ctx := context.Background()
	expectedErr := errors.New("database error")

	repo := &mockPlaybackRepository{
		updateProgressFn: func(ctx context.Context, session *domain.PlaybackSession) error {
			return expectedErr
		},
	}

	publisher := &mockPlaybackPublisher{}
	uc := NewPlaybackUsecase(repo, publisher)

	err := uc.UpdateProgress(ctx, samplePlaybackSession())

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}

	if len(publisher.events) != 0 {
		t.Fatalf("expected no published events, got %d", len(publisher.events))
	}
}

func TestMarkCompletedPublishesCompletedEvent(t *testing.T) {
	ctx := context.Background()

	repo := &mockPlaybackRepository{
		markCompletedFn: func(ctx context.Context, profileID, contentID, episodeID string) error {
			if profileID != "profile-1" {
				t.Fatalf("expected profile-1, got %s", profileID)
			}
			if contentID != "content-1" {
				t.Fatalf("expected content-1, got %s", contentID)
			}
			if episodeID != "episode-1" {
				t.Fatalf("expected episode-1, got %s", episodeID)
			}
			return nil
		},
	}

	publisher := &mockPlaybackPublisher{}
	uc := NewPlaybackUsecase(repo, publisher)

	if err := uc.MarkCompleted(ctx, "profile-1", "content-1", "episode-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(publisher.events) != 1 {
		t.Fatalf("expected 1 published event, got %d", len(publisher.events))
	}

	if publisher.events[0].subject != SubjectCompleted {
		t.Fatalf("expected subject %s, got %s", SubjectCompleted, publisher.events[0].subject)
	}
}

func TestGetContinueWatchingReturnsSessions(t *testing.T) {
	ctx := context.Background()
	expected := []*domain.PlaybackSession{
		samplePlaybackSession(),
	}

	repo := &mockPlaybackRepository{
		getContinueWatchingFn: func(ctx context.Context, profileID string) ([]*domain.PlaybackSession, error) {
			if profileID != "profile-1" {
				t.Fatalf("expected profile-1, got %s", profileID)
			}
			return expected, nil
		},
	}

	uc := NewPlaybackUsecase(repo, nil)

	got, err := uc.GetContinueWatching(ctx, "profile-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 session, got %d", len(got))
	}

	if got[0].ContentID != "content-1" {
		t.Fatalf("expected content-1, got %s", got[0].ContentID)
	}
}
