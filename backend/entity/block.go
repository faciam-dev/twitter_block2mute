package entity

import "time"

type Block struct {
	ID              uint
	UserID          uint
	TargetTwitterID string
	Flag            int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
