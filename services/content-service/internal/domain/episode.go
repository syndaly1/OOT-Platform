package domain

import "time"

type Episode struct {
	ID              string
	ContentID       string
	Title           string
	Description     string
	SeasonNumber    int
	EpisodeNumber   int
	DurationSeconds int
	SkipIntroStart  int
	SkipIntroEnd    int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
