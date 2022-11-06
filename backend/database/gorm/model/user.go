package model

import (
	"github.com/faciam_dev/twitter_block2mute/backend/entity"
	"gorm.io/gorm"
)

// for users table
type User struct {
	gorm.Model
	Name        string `gorm:"size:255;not null"`
	AccountName string `gorm:"size:255;not null"`
	TwitterID   string `gorm:"size:64;not null"`
	UserBlock   []UserBlock
	UserMute    []UserMute
}

/*
func (u *User) FromUserDomain(userEntity *entity.User) {
	u.ID = userEntity.GetID()
	u.Name = userEntity.GetName()
	u.AccountName = userEntity.GetAccountName()
	u.TwitterID = userEntity.GetTwitterID()
}

func (u *User) ToUserDomain(i interface{}) {
	if s, ok := i.(*entity.User); ok {
		s.ID = u.ID
		s.Name = u.Name
		s.AccountName = u.AccountName
		s.TwitterID = u.TwitterID
	}
}
*/

// for domain
type UserModelForDomain struct {
	User
	ModelForDomain[entity.User, User]
}

func (u *UserModelForDomain) FromDomain(userEntity *entity.User) User {
	u.User.ID = userEntity.GetID()
	u.User.Name = userEntity.GetName()
	u.User.AccountName = userEntity.GetAccountName()
	u.User.TwitterID = userEntity.GetTwitterID()

	return u.User
}

func (u *UserModelForDomain) ToDomain(userModel User, userEntity *entity.User) {

	userEntity.Update(userModel.ID, userModel.Name, userModel.AccountName, userModel.TwitterID)
	/*
		userEntity.ID = userModel.ID
		userEntity.Name = userModel.Name
		userEntity.AccountName = userModel.AccountName
		userEntity.TwitterID = userModel.TwitterID
	*/
}

func (u *UserModelForDomain) ToDomains(models []User, entities *[]entity.User) {
	for _, model := range models {
		entity := &entity.User{}
		u.ToDomain(model, entity)
		*entities = append(*entities, *entity)
	}
}
