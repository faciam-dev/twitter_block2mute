package gateway

import (
	"fmt"
	"log"

	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/entity"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/port"
)

type AuthRepository struct {
	contextHandler handler.ContextHandler
	twitterHandler handler.TwitterHandler
	sessionHandler handler.SessionHandler
	userDbHandler  handler.UserDbHandler
}

// NewAuthRepository はAuthRepositoryを返します．
func NewAuthRepository(contextHandler handler.ContextHandler, twitterHandler handler.TwitterHandler, sessionHandler handler.SessionHandler, userDbHandler handler.UserDbHandler) port.AuthRepository {
	authRepository := &AuthRepository{
		contextHandler: contextHandler,
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
	} else {
		log.Printf("GetUser() error: %v", err.Error())
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
		return &auth, err
	}

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
		log.Printf("token=%v,secret=%v", token, secret)
		log.Println(err)
		return &auth, err
	}

	// 認証成功後処理
	// DB
	user := entity.User{}
	if err := a.userDbHandler.FindByTwitterID(&user, values.GetTwitterID()); err != nil {
		log.Println(err)
		return &auth, err
	}

	//log.Printf("user-> id:%v tid:%v", user.ID, user.TwitterID)
	if user.ID == 0 {
		user.Name = values.GetTwitterScreenName()
		user.AccountName = values.GetTwitterScreenName()
		user.TwitterID = values.GetTwitterID()
	}
	if err := a.userDbHandler.Upsert(&user, "twitter_id", user.TwitterID); err != nil {
		log.Println(err)
		return &auth, err
	}

	a.twitterHandler.UpdateTwitterApi(credentials.GetToken(), credentials.GetSecret())
	auth.Authenticated = 1

	// Session
	a.sessionHandler.Set("token", credentials.GetToken())
	a.sessionHandler.Set("secret", credentials.GetSecret())
	a.sessionHandler.Set("user_id", fmt.Sprintf("%d", user.ID))
	a.sessionHandler.Set("twitter_id", values.GetTwitterID())

	if err := a.sessionHandler.Save(); err != nil {
		log.Println(err)
		return &auth, err
	}

	return &auth, nil
}
