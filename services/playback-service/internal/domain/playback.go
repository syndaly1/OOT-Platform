package domain

type PlaybackSession struct {
	SessionID string
	ProfileID string
	ContentID string
	EpisodeID string
	Device    string

	CurrentPositionSec int32
	DurationSec        int32
	WatchedPercent     float32
	Completed          bool

	Status string
}