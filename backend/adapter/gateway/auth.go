package gateway

import (
	"fmt"

	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/entity"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/port"
)

type AuthRepository struct {
	contextHandler handler.ContextHandler
	loggerHandler  handler.LoggerHandler
	twitterHandler handler.TwitterHandler
	sessionHandler handler.SessionHandler
	userDBHandler  handler.UserDBHandler
}

// NewAuthRepository はAuthRepositoryを返します．
func NewAuthRepository(
	contextHandler handler.ContextHandler,
	loggerHandler handler.LoggerHandler,
	twitterHandler handler.TwitterHandler,
	sessionHandler handler.SessionHandler,
	userDBHandler handler.UserDBHandler,
) port.AuthRepository {
	authRepository := &AuthRepository{
		contextHandler: contextHandler,
		loggerHandler:  loggerHandler,
		twitterHandler: twitterHandler,
		sessionHandler: sessionHandler,
		userDBHandler:  userDBHandler,
	}
	authRepository.sessionHandler.SetContextHandler(contextHandler)
	return authRepository
}

// セッションに入っている値から認証を実施する。
func (a *AuthRepository) IsAuthenticated() error {
	token := a.sessionHandler.Get("token")
	secret := a.sessionHandler.Get("secret")
	twitterID := a.sessionHandler.Get("twitter_id")

	if token == nil || secret == nil || twitterID == nil {
		return nil
	}

	a.loggerHandler.Debugf("Try IsAuthenticated(%s): token=%s secret=%s", twitterID.(string), token.(string), secret.(string))

	a.UpdateTwitterApi(token.(string), secret.(string))
	_, err := a.twitterHandler.GetUser(twitterID.(string))

	a.loggerHandler.Debugf("IsAuthenticated(%s) OK: token=%s secret=%s", twitterID.(string), token.(string), secret.(string))

	return err
}

// 認証用URLを得る
func (a *AuthRepository) GetAuthUrl() (string, error) {
	uri, err := a.twitterHandler.AuthorizationURL()

	if err != nil {
		a.loggerHandler.Errorw("Auth() error", "error", err.Error())
		return uri, err
	}

	a.loggerHandler.Debugf("Auth(): url=%s", uri)

	return uri, nil
}

// コールバックパラメータによるユーザー認証
func (a *AuthRepository) AuthByCallbackParams(token string, secret string) (handler.TwitterCredentials, handler.TwitterValues, error) {
	credentials, values, err := a.twitterHandler.GetCredentials(token, secret)

	if err != nil {
		a.loggerHandler.Errorf("token=%s,secret=%s", token, secret)
		a.loggerHandler.Errorw("error", "error", err)
		return credentials, values, err
	}

	a.loggerHandler.Debugf("Callback() OK: token=%s secret=%s", token, secret)

	return credentials, values, nil
}

// TwitterIDによるユーザー取得
func (a *AuthRepository) FindUserByTwitterID(twitterID string) (*entity.User, error) {
	user := entity.User{}

	if err := a.userDBHandler.FindByTwitterID(&user, twitterID); err != nil {
		a.loggerHandler.Errorw("Callback() FindByTwitterID: error", "error", err)
		return &user, err
	}

	a.loggerHandler.Debugf("user-> id:%s tid:%s", user.GetID(), twitterID)

	return &user, nil
}

// userのUpsert処理
func (a *AuthRepository) UpsertUser(user *entity.User) error {
	if err := a.userDBHandler.Upsert(user, "twitter_id", user.GetTwitterID()); err != nil {
		a.loggerHandler.Errorw("Callback() Upsert", "error", err)
		return err
	}

	a.loggerHandler.Debugf("user(ID:%s TwitterID:%s) Upsert OK", user.GetID(), user.GetTwitterID())

	return nil
}

// TwitterのAPIをtokenとsecretで更新しなおす
func (a *AuthRepository) UpdateTwitterApi(token string, secret string) {
	a.twitterHandler.UpdateTwitterApi(token, secret)
}

// セッションを更新する
func (a *AuthRepository) UpdateSession(token string, secret string, userID int, twitterID string) error {
	a.sessionHandler.Set("token", token)
	a.sessionHandler.Set("secret", secret)
	a.sessionHandler.Set("user_id", fmt.Sprintf("%d", userID))
	a.sessionHandler.Set("twitter_id", twitterID)

	if err := a.sessionHandler.Save(); err != nil {
		a.loggerHandler.Errorw("Session error", "error", err)
		return err
	}

	a.loggerHandler.Debug("Session Set OK")

	return nil
}

// ログアウトのためにセッションをクリアする。
func (a *AuthRepository) Logout() error {
	a.sessionHandler.Clear()

	if err := a.sessionHandler.Save(); err != nil {
		a.loggerHandler.Errorw("Session error", "error", err)
		return err
	}
	a.loggerHandler.Debug("Session Clear OK")

	return nil
}
