package gateway

import (
	"errors"

	"github.com/faciam_dev/twitter_block2mute/backend/entity"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/port"
)

type UserRepository struct {
	dbHandler DbHandler
}

// NewUserRepository はUserRepositoryを返します．
func NewUserRepository(dbHandler DbHandler) port.UserRepository {
	return &UserRepository{
		dbHandler: dbHandler,
	}
}

// GetUserByID はDBからデータを取得します．
func (u *UserRepository) GetUserByID(userID string) (*entity.User, error) {
	user := entity.User{}

    if err := u.dbHandler.First(&user, userID);err != nil {
		return &entity.User{},err
	}
    if user.ID <= 0 {
        return &entity.User{}, errors.New("user is not found")
    }
    return &user, nil
}