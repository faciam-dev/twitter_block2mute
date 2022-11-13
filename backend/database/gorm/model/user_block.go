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

// for domain
type UserBlockModelForDomain struct {
	UserBlock
	ModelForDomain[entity.Block, UserBlock]
}

func (u *UserBlockModelForDomain) FromDomain(entity *entity.Block) UserBlock {
	u.UserBlock.ID = entity.GetID()
	u.UserBlock.TargetTwitterID = entity.GetTargetTwitterID()
	u.UserBlock.Flag = entity.GetFlag()
	u.UserBlock.UserID = entity.GetUserID()
	u.UserBlock.CreatedAt = entity.GetCreatedAt()
	u.UserBlock.UpdatedAt = entity.GetUpdatedAt()

	return u.UserBlock
}

func (u *UserBlockModelForDomain) ToDomain(model UserBlock, entity *entity.Block) {
	entity.Update(
		model.ID,
		model.UserID,
		model.TargetTwitterID,
		model.Flag,
		model.CreatedAt,
		model.UpdatedAt,
	)
}

func (u *UserBlockModelForDomain) ToDomains(models []UserBlock, entities *[]entity.Block) {
	for _, model := range models {
		entity := &entity.Block{}
		u.ToDomain(model, entity)
		*entities = append(*entities, *entity)
	}
}
