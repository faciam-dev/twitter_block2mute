package model

import (
	"github.com/faciam_dev/twitter_block2mute/backend/entity"
	"gorm.io/gorm"
)

// for user_blocks
type UserBlock struct {
	gorm.Model
	UserID          uint   `gorm:"uniqueIndex:idx_user_id_target_twitter_id"`
	TargetTwitterID string `gorm:"uniqueIndex:idx_user_id_target_twitter_id;size:64;not null"`
	Flag            int
	User            User
}

func (u *UserBlock) FromDomain(blockEntity *entity.Block) {
	u.ID = blockEntity.ID
	u.TargetTwitterID = blockEntity.TargetTwitterID
	u.Flag = blockEntity.Flag
	u.UserID = blockEntity.UserID
	u.CreatedAt = blockEntity.CreatedAt
	u.UpdatedAt = blockEntity.UpdatedAt
}

func (u *UserBlock) ToDomain(i interface{}) {
	if s, ok := i.(*entity.Block); ok {
		s.ID = u.ID
		s.TargetTwitterID = u.TargetTwitterID
		s.Flag = u.Flag
		s.UserID = u.UserID
		s.CreatedAt = u.CreatedAt
		s.UpdatedAt = u.UpdatedAt
	}
}
