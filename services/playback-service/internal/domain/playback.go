package domain

type PlaybackSession struct {
	ProfileID string
	ContentID string
	EpisodeID string

	CurrentPositionSec int32
	DurationSec        int32

	WatchedPercent float32

	Completed bool
}