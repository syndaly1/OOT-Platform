package domain

import "time"

type PlaybackSession struct {
	ID                 string
	ProfileID          string
	ContentID          string
	EpisodeID          string
	CurrentPositionSec int32
	DurationSec        int32
	WatchedPercent     float32
	Completed          bool
	StartedAt          time.Time
	UpdatedAt          time.Time
}