package entity

import "time"

type Mute struct {
	ID              uint
	UserID          uint
	TargetTwitterID string
	Flag            int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
