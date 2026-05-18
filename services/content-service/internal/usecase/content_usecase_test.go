package usecase

import (
 "context"
 "errors"
 "testing"
 "time"

 "content-service/internal/domain"

 "github.com/stretchr/testify/assert"
 "github.com/stretchr/testify/require"
)

type mockContentRepository struct {
 createContentFn        func(ctx context.Context, content *domain.Content) error
 getContentByIDFn      func(ctx context.Context, id string) (*domain.Content, error)
 listContentsFn        func(ctx context.Context, page, limit int32) ([]*domain.Content, error)
 updateContentFn       func(ctx context.Context, content *domain.Content) error
 deleteContentFn       func(ctx context.Context, id string) error
 searchContentsFn      func(ctx context.Context, query string) ([]*domain.Content, error)
 getContentsByGenreFn  func(ctx context.Context, genre string) ([]*domain.Content, error)
 getTrendingContentsFn func(ctx context.Context, limit int32) ([]*domain.Content, error)
 getFeaturedContentsFn func(ctx context.Context, limit int32) ([]*domain.Content, error)
}

func (m *mockContentRepository) CreateContent(ctx context.Context, content *domain.Content) error {
 return m.createContentFn(ctx, content)
}

func (m *mockContentRepository) GetContentByID(ctx context.Context, id string) (*domain.Content, error) {
 return m.getContentByIDFn(ctx, id)
}

func (m *mockContentRepository) ListContents(ctx context.Context, page, limit int32) ([]*domain.Content, error) {
 return m.listContentsFn(ctx, page, limit)
}

func (m *mockContentRepository) UpdateContent(ctx context.Context, content *domain.Content) error {
 return m.updateContentFn(ctx, content)
}

func (m *mockContentRepository) DeleteContent(ctx context.Context, id string) error {
 return m.deleteContentFn(ctx, id)
}

func (m *mockContentRepository) SearchContents(ctx context.Context, query string) ([]*domain.Content, error) {
 return m.searchContentsFn(ctx, query)
}

func (m *mockContentRepository) GetContentsByGenre(ctx context.Context, genre string) ([]*domain.Content, error) {
 return m.getContentsByGenreFn(ctx, genre)
}

func (m *mockContentRepository) GetTrendingContents(ctx context.Context, limit int32) ([]*domain.Content, error) {
 return m.getTrendingContentsFn(ctx, limit)
}

func (m *mockContentRepository) GetFeaturedContents(ctx context.Context, limit int32) ([]*domain.Content, error) {
 return m.getFeaturedContentsFn(ctx, limit)
}

type publishedEvent struct {
 subject string
 data    any
}

type mockPublisher struct {
 events []publishedEvent
}

func (m *mockPublisher) Publish(subject string, data any) {
 m.events = append(m.events, publishedEvent{
  subject: subject,
  data:    data,
 })
}

func sampleContent() *domain.Content {
 now := time.Now()

 return &domain.Content{
  ID:          "content-1",
  Title:       "Night Has Come",
  Description: "A mystery thriller.",
  PosterURL:   "/images/kdrama-1.jpg",
  Type:        "tv_series",
  Genre:       "thriller",
  ReleaseYear: 2023,
  Rating:      8.2,
  IsFeatured:  true,
  IsTrending:  true,
  CreatedAt:   now,
  UpdatedAt:   now,
 }
}

func TestCreateContentSuccessPublishesEvent(t *testing.T) {
 ctx := context.Background()
 content := sampleContent()

 repo := &mockContentRepository{
  createContentFn: func(ctx context.Context, got *domain.Content) error {
   assert.Equal(t, content.ID, got.ID)
   assert.Equal(t, content.Title, got.Title)
   assert.Equal(t, content.PosterURL, got.PosterURL)
   return nil
  },
 }

 publisher := &mockPublisher{}
 uc := NewContentUsecase(repo, publisher)

 err := uc.CreateContent(ctx, content)

 require.NoError(t, err)
 require.Len(t, publisher.events, 1)
 assert.Equal(t, "content.created", publisher.events[0].subject)
 assert.Equal(t, content, publisher.events[0].data)
}

func TestCreateContentRepositoryErrorDoesNotPublishEvent(t *testing.T) {
 ctx := context.Background()
 expectedErr := errors.New("database error")

 repo := &mockContentRepository{
  createContentFn: func(ctx context.Context, content *domain.Content) error {
   return expectedErr
  },
 }

 publisher := &mockPublisher{}
 uc := NewContentUsecase(repo, publisher)

 err := uc.CreateContent(ctx, sampleContent())

 require.Error(t, err)
 assert.Equal(t, expectedErr, err)
 assert.Empty(t, publisher.events)
}

func TestGetContentByIDSuccess(t *testing.T) {
 ctx := context.Background()
 content := sampleContent()

 repo := &mockContentRepository{
  getContentByIDFn: func(ctx context.Context, id string) (*domain.Content, error) {
   assert.Equal(t, "content-1", id)
   return content, nil
  },
 }

 uc := NewContentUsecase(repo, nil)

 got, err := uc.GetContentByID(ctx, "content-1")

 require.NoError(t, err)
 require.NotNil(t, got)
 assert.Equal(t, content.ID, got.ID)
 assert.Equal(t, content.Title, got.Title)
 assert.Equal(t, content.PosterURL, got.PosterURL)
}

func TestListContentsSuccess(t *testing.T) {
 ctx := context.Background()
 contents := []*domain.Content{
  sampleContent(),
  {
   ID:          "content-2",
   Title:       "The First Frost",
   Description: "A romantic drama.",
   PosterURL:   "/images/cdrama-1.jpg",
   Type:        "tv_series",
   Genre:       "romance",
   ReleaseYear: 2025,
   Rating:      8.5,
   IsFeatured:  false,
   IsTrending:  true,
  },
 }

 repo := &mockContentRepository{
  listContentsFn: func(ctx context.Context, page, limit int32) ([]*domain.Content, error) {
   assert.Equal(t, int32(1), page)
   assert.Equal(t, int32(10), limit)
   return contents, nil
  },
 }

 uc := NewContentUsecase(repo, nil)

 got, err := uc.ListContents(ctx, 1, 10)

 require.NoError(t, err)
 require.Len(t, got, 2)
 assert.Equal(t, "content-1", got[0].ID)
 assert.Equal(t, "content-2", got[1].ID)
}

func TestUpdateContentSuccessPublishesEvent(t *testing.T) {
 ctx := context.Background()
 content := sampleContent()
 content.Title = "Updated Title"

 repo := &mockContentRepository{
  updateContentFn: func(ctx context.Context, got *domain.Content) error {
   assert.Equal(t, "content-1", got.ID)
   assert.Equal(t, "Updated Title", got.Title)
   return nil
  },
 }

 publisher := &mockPublisher{}
 uc := NewContentUsecase(repo, publisher)

 err := uc.UpdateContent(ctx, content)

 require.NoError(t, err)
 require.Len(t, publisher.events, 1)
 assert.Equal(t, "content.updated", publisher.events[0].subject)
 assert.Equal(t, content, publisher.events[0].data)
}

func TestUpdateContentRepositoryErrorDoesNotPublishEvent(t *testing.T) {
 ctx := context.Background()
 expectedErr := errors.New("update failed")

 repo := &mockContentRepository{
  updateContentFn: func(ctx context.Context, content *domain.Content) error {
   return expectedErr
  },
 }

 publisher := &mockPublisher{}
 uc := NewContentUsecase(repo, publisher)

 err := uc.UpdateContent(ctx, sampleContent())

 require.Error(t, err)
 assert.Equal(t, expectedErr, err)
 assert.Empty(t, publisher.events)
}

func TestDeleteContentSuccess(t *testing.T) {
 ctx := context.Background()

 repo := &mockContentRepository{
  deleteContentFn: func(ctx context.Context, id string) error {
   assert.Equal(t, "content-1", id)
   return nil
  },
 }

 uc := NewContentUsecase(repo, nil)

 err := uc.DeleteContent(ctx, "content-1")

 require.NoError(t, err)
}

func TestSearchContentsSuccess(t *testing.T) {
 ctx := context.Background()
 contents := []*domain.Content{sampleContent()}

 repo := &mockContentRepository{
  searchContentsFn: func(ctx context.Context, query string) ([]*domain.Content, error) {
   assert.Equal(t, "night", query)
   return contents, nil
  },
 }

 uc := NewContentUsecase(repo, nil)

 got, err := uc.SearchContents(ctx, "night")

 require.NoError(t, err)
 require.Len(t, got, 1)
 assert.Equal(t, "Night Has Come", got[0].Title)
}

func TestGetContentsByGenreSuccess(t *testing.T) {
 ctx := context.Background()
 contents := []*domain.Content{sampleContent()}

 repo := &mockContentRepository{
  getContentsByGenreFn: func(ctx context.Context, genre string) ([]*domain.Content, error) {
   assert.Equal(t, "thriller", genre)
   return contents, nil
  },
 }

 uc := NewContentUsecase(repo, nil)

 got, err := uc.GetContentsByGenre(ctx, "thriller")

 require.NoError(t, err)
 require.Len(t, got, 1)
 assert.Equal(t, "thriller", got[0].Genre)
}

func TestGetTrendingContentsSuccess(t *testing.T) {
 ctx := context.Background()
 contents := []*domain.Content{sampleContent()}

 repo := &mockContentRepository{
  getTrendingContentsFn: func(ctx context.Context, limit int32) ([]*domain.Content, error) {
   assert.Equal(t, int32(5), limit)
   return contents, nil
  },
 }

 uc := NewContentUsecase(repo, nil)

 got, err := uc.GetTrendingContents(ctx, 5)

 require.NoError(t, err)
 require.Len(t, got, 1)
 assert.True(t, got[0].IsTrending)
}

func TestGetFeaturedContentsSuccess(t *testing.T) {
 ctx := context.Background()
 contents := []*domain.Content{sampleContent()}

 repo := &mockContentRepository{
  getFeaturedContentsFn: func(ctx context.Context, limit int32) ([]*domain.Content, error) {
   assert.Equal(t, int32(3), limit)
   return contents, nil
  },
 }

 uc := NewContentUsecase(repo, nil)

 got, err := uc.GetFeaturedContents(ctx, 3)

 require.NoError(t, err)
 require.Len(t, got, 1)
 assert.True(t, got[0].IsFeatured)
}
