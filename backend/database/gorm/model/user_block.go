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

/*
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
*/
// for domain
type UserBlockModelForDomain struct {
	UserBlock
	ModelForDomain[entity.Block, UserBlock]
}

func (u *UserBlockModelForDomain) FromDomain(entity *entity.Block) UserBlock {
	u.UserBlock.ID = entity.ID
	u.UserBlock.TargetTwitterID = entity.TargetTwitterID
	u.UserBlock.Flag = entity.Flag
	u.UserBlock.UserID = entity.UserID
	u.UserBlock.CreatedAt = entity.CreatedAt
	u.UserBlock.UpdatedAt = entity.UpdatedAt

	return u.UserBlock
}

func (u *UserBlockModelForDomain) ToDomain(model UserBlock, entity *entity.Block) {
	entity.ID = model.ID
	entity.TargetTwitterID = model.TargetTwitterID
	entity.Flag = model.Flag
	entity.UserID = model.UserID
	entity.CreatedAt = model.CreatedAt
	entity.UpdatedAt = model.UpdatedAt
}

func (u *UserBlockModelForDomain) ToDomains(models []UserBlock, entities *[]entity.Block) {
	for _, model := range models {
		entity := &entity.Block{}
		u.ToDomain(model, entity)
		*entities = append(*entities, *entity)
	}
}
