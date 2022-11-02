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
	userDbHandler  handler.UserDbHandler
}

// NewAuthRepository はAuthRepositoryを返します．
func NewAuthRepository(
	contextHandler handler.ContextHandler,
	loggerHandler handler.LoggerHandler,
	twitterHandler handler.TwitterHandler,
	sessionHandler handler.SessionHandler,
	userDbHandler handler.UserDbHandler,
) port.AuthRepository {
	authRepository := &AuthRepository{
		contextHandler: contextHandler,
		loggerHandler:  loggerHandler,
		twitterHandler: twitterHandler,
		sessionHandler: sessionHandler,
		userDbHandler:  userDbHandler,
	}
	authRepository.sessionHandler.SetContextHandler(contextHandler)
	return authRepository
}

// セッションに入っている値から認証を実施し、認証状態を取得する。
func (a *AuthRepository) IsAuthenticated() (*entity.Auth, error) {
	token := a.sessionHandler.Get("token")
	secret := a.sessionHandler.Get("secret")
	twitterID := a.sessionHandler.Get("twitter_id")

	auth := entity.Auth{
		Authenticated: 0,
	}

	if token == nil || secret == nil {
		return &auth, nil
	}

	a.twitterHandler.UpdateTwitterApi(token.(string), secret.(string))
	_, err := a.twitterHandler.GetUser(twitterID.(string))
	if err == nil {
		// 認証成功
		auth.Authenticated = 1
		a.loggerHandler.Debugf("IsAuthenticated(%s) OK: token=%s secret=%s", twitterID.(string), token.(string), secret.(string))
	} else {
		a.loggerHandler.Errorw("IsAuthenticated() error", "error", err.Error())
		//log.Printf("GetUser() error: %v", err.Error())
	}

	return &auth, nil
}

// 認証を実施する（認証用URLを返す）
func (a *AuthRepository) Auth() (*entity.Auth, error) {
	auth := entity.Auth{
		Authenticated: 0,
	}

	uri, err := a.twitterHandler.AuthorizationURL()

	if err != nil {
		a.loggerHandler.Errorw("Auth() error", "error", err.Error())
		return &auth, err
	}

	a.loggerHandler.Debugf("Auth(): url=%s", uri)

	auth.AuthUrl = uri

	return &auth, nil
}

// 認証コールバック
func (a *AuthRepository) Callback(token string, secret string, twitterID string, twitterName string) (*entity.Auth, error) {

	auth := entity.Auth{
		Authenticated: 0,
	}

	credentials, values, err := a.twitterHandler.GetCredentials(token, secret)

	if err != nil {
		a.loggerHandler.Errorf("token=%s,secret=%s", token, secret)
		a.loggerHandler.Errorw("error", "error", err)
		return &auth, err
	}

	a.loggerHandler.Debugf("Callback() OK: token=%s secret=%s", token, secret)

	// 認証成功後処理
	// DB
	user := entity.User{}
	if err := a.userDbHandler.FindByTwitterID(&user, values.GetTwitterID()); err != nil {
		a.loggerHandler.Errorw("Callback() FindByTwitterID: error", "error", err)
		return &auth, err
	}

	a.loggerHandler.Debugf("user-> id:%s tid:%s", user.ID, user.TwitterID)
	//log.Printf("user-> id:%v tid:%v", user.ID, user.TwitterID)
	if user.ID == 0 {
		user.Name = values.GetTwitterScreenName()
		user.AccountName = values.GetTwitterScreenName()
		user.TwitterID = values.GetTwitterID()
	}
	if err := a.userDbHandler.Upsert(&user, "twitter_id", user.TwitterID); err != nil {
		a.loggerHandler.Errorw("Callback() Upsert", "error", err)
		return &auth, err
	}
	a.loggerHandler.Debugf("user(ID:%s TwitterID:%s) Upsert OK", user.ID, user.TwitterID)

	a.twitterHandler.UpdateTwitterApi(credentials.GetToken(), credentials.GetSecret())
	auth.Authenticated = 1

	// Session
	a.sessionHandler.Set("token", credentials.GetToken())
	a.sessionHandler.Set("secret", credentials.GetSecret())
	a.sessionHandler.Set("user_id", fmt.Sprintf("%d", user.ID))
	a.sessionHandler.Set("twitter_id", values.GetTwitterID())

	if err := a.sessionHandler.Save(); err != nil {
		a.loggerHandler.Errorw("Session error", "error", err)
		return &auth, err
	}

	a.loggerHandler.Debug("Session Set OK")

	return &auth, nil
}

func (a *AuthRepository) Logout() (*entity.Auth, error) {
	auth := entity.Auth{
		Logout: 0,
	}

	a.sessionHandler.Clear()

	auth.Logout = 1

	if err := a.sessionHandler.Save(); err != nil {
		a.loggerHandler.Errorw("Session error", "error", err)
		return &auth, err
	}

	a.loggerHandler.Debug("Session Clear OK")

	return &auth, nil
}
