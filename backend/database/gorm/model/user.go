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

func (u *User) FromUserDomain(userEntity *entity.User) {
	u.ID = userEntity.ID
	u.Name = userEntity.Name
	u.AccountName = userEntity.AccountName
	u.TwitterID = userEntity.TwitterID
}

func (u *User) ToUserDomain(i interface{}) {
	if s, ok := i.(*entity.User); ok {
		s.ID = u.ID
		s.Name = u.Name
		s.AccountName = u.AccountName
		s.TwitterID = u.TwitterID
	}
}

// for domain
type UserModelForDomain struct {
	User
	ModelForDomain[entity.User, User]
}

func (u *UserModelForDomain) FromDomain(userEntity *entity.User) User {
	u.User.ID = userEntity.ID
	u.User.Name = userEntity.Name
	u.User.AccountName = userEntity.AccountName
	u.User.TwitterID = userEntity.TwitterID

	return u.User
}

func (u *UserModelForDomain) ToDomain(userModel User, userEntity *entity.User) {
	userEntity.ID = userModel.ID
	userEntity.Name = userModel.Name
	userEntity.AccountName = userModel.AccountName
	userEntity.TwitterID = userModel.TwitterID
}

func (u *UserModelForDomain) ToDomains(models []User, entities *[]entity.User) {
	for _, model := range models {
		entity := &entity.User{}
		u.ToDomain(model, entity)
		*entities = append(*entities, *entity)
	}
}

/*
	func (u *UserModelForDomain[E, M]) ToDomain(userModel User, userEntity interface{}) {
		switch userEntity.(type) {
		case E:
			log.Print("cast")
				casted.ID = userModel.ID
				casted.Name = userModel.Name
				casted.AccountName = userModel.AccountName
				casted.TwitterID = userModel.TwitterID
		default:
			log.Print("nottt")
		}
	}
*/
/*
func (u *UserModelForDomain[E, M]) Blank() M {
	return M(u.User)
}
*/
