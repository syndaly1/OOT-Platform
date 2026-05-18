package domain

import "time"

type Content struct {
	ID          string
	Title       string
	Description string
	PosterURL   string
	Type        string
	Genre       string
	ReleaseYear int
	Rating      float64
	IsFeatured  bool
	IsTrending  bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
