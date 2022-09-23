package gateway

import (
	"errors"

	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/entity"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/port"
)

type UserRepository struct {
	dbHandler handler.UserDbHandler
}

// NewUserRepository はUserRepositoryを返します．
func NewUserRepository(dbHandler handler.UserDbHandler) port.UserRepository {
	return &UserRepository{
		dbHandler: dbHandler,
	}
}

// DBからid=userIDに該当するデータを取得します．
func (u *UserRepository) GetUserByID(userID string) (*entity.User, error) {
	user := entity.User{}

	if err := u.dbHandler.First(&user, userID); err != nil {
		return &entity.User{}, err
	}
	if user.ID <= 0 {
		return &entity.User{}, errors.New("user is not found")
	}
	return &user, nil
}

// DBにユーザーを追加する。既に存在する場合はデータを上書き更新するする。
func (u *UserRepository) UpsertByTwitterID(newUser *entity.User, twitterID string) (*entity.User, error) {
	if err := u.dbHandler.Upsert(newUser, "twitter_id", twitterID); err != nil {
		return newUser, err
	}
	return newUser, nil
}
