package gateway

import (
	"errors"
	"strconv"

	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/entity"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/port"
)

type UserRepository struct {
	loggerHandler handler.LoggerHandler
	dbHandler     handler.UserDbHandler
}

// NewUserRepository はUserRepositoryを返します．
func NewUserRepository(
	loggerHandler handler.LoggerHandler,
	dbHandler handler.UserDbHandler,
) port.UserRepository {
	return &UserRepository{
		loggerHandler: loggerHandler,
		dbHandler:     dbHandler,
	}
}

// DBからid=userIDに該当するデータを取得する。
func (u *UserRepository) GetUserByID(userID string) (*entity.User, error) {
	user := entity.User{}
	if err := u.dbHandler.First(&user, userID); err != nil {
		u.loggerHandler.Errorf("user not found (user_id=%s)", userID)
		return &entity.User{}, err
	}
	if strconv.FormatUint(uint64(user.GetID()), 10) != userID {
		return &entity.User{}, errors.New("user is not found")
	}
	u.loggerHandler.Debugf("user found (user_id=%s)", userID)
	return &user, nil
}

// DBにユーザーを追加する。既に存在する場合はデータを上書き更新する。
func (u *UserRepository) UpsertByTwitterID(newUser *entity.User, twitterID string) (*entity.User, error) {
	if err := u.dbHandler.Upsert(newUser, "twitter_id", twitterID); err != nil {
		u.loggerHandler.Errorw("upsert error", "twitter_id", twitterID, "error", err)
		return newUser, err
	}
	u.loggerHandler.Debugf("upsert ok (twitter_id=%s)", twitterID)
	return newUser, nil
}
