package model

import (
	"github.com/faciam_dev/twitter_block2mute/backend/entity"
	"gorm.io/gorm"
)

// for user_mutes
type UserMute struct {
	gorm.Model
	UserID          uint   `gorm:"uniqueIndex:idx_user_id_target_twitter_id"`
	TargetTwitterID string `gorm:"uniqueIndex:idx_user_id_target_twitter_id;size:64;not null"`
	Flag            int
	User            User
}

// for domain
type UserMuteModelForDomain struct {
	UserMute
	ModelForDomain[entity.Mute, UserMute]
}

func (u *UserMuteModelForDomain) FromDomain(entity *entity.Mute) UserMute {
	u.UserMute.ID = entity.ID
	u.UserMute.TargetTwitterID = entity.TargetTwitterID
	u.UserMute.Flag = entity.Flag
	u.UserMute.UserID = entity.UserID
	u.UserMute.CreatedAt = entity.CreatedAt
	u.UserMute.UpdatedAt = entity.UpdatedAt

	return u.UserMute
}

func (u *UserMuteModelForDomain) ToDomain(model UserMute, entity *entity.Mute) {
	entity.ID = model.ID
	entity.TargetTwitterID = model.TargetTwitterID
	entity.Flag = model.Flag
	entity.UserID = model.UserID
	entity.CreatedAt = model.CreatedAt
	entity.UpdatedAt = model.UpdatedAt
}

func (u *UserMuteModelForDomain) ToDomains(models []UserMute, entities *[]entity.Mute) {
	for _, model := range models {
		entity := &entity.Mute{}
		u.ToDomain(model, entity)
		*entities = append(*entities, *entity)
	}
}
