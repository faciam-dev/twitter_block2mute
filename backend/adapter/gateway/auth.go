package gateway

import (
	"strconv"

	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/entity"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/port"
)

type AuthRepository struct {
	contextHandler handler.ContextHandler
	twitterHandler handler.TwitterHandler
	sessionHandler handler.SessionHandler
	userDbHandler handler.UserDbHandler
}

// NewAuthRepository はAuthRepositoryを返します．
func NewAuthRepository(contextHandler handler.ContextHandler, twitterHandler handler.TwitterHandler, sessionHandler handler.SessionHandler, userDbHandler handler.UserDbHandler) port.AuthRepository {
	authRepository := &AuthRepository{
		contextHandler: contextHandler,
		twitterHandler: twitterHandler,
		sessionHandler: sessionHandler,
		userDbHandler: userDbHandler,
	}
	authRepository.sessionHandler.SetContextHandler(contextHandler)
	return authRepository
}

// セッションに入っている値から認証を実施し、認証状態を取得する。
func (a *AuthRepository) IsAuthenticated() (*entity.Auth, error) {
	token := a.sessionHandler.Get("token")
	secret := a.sessionHandler.Get("secret") 

	auth := entity.Auth{
		Authenticated: 0,
	}

	if token == nil || secret == nil {
		return &auth, nil
	}

	a.twitterHandler.SetCredentials(token.(string), secret.(string))
	err := a.twitterHandler.GetRateLimits()
	if err == nil {
		// 認証成功
		auth.Authenticated = 1;
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
		return &auth, err;
	}

	auth.AuthUrl = uri

	return &auth, nil
}

// 認証コールバック
func (a *AuthRepository) Callback(token string, secret string, twitterID string, twitterName string) (*entity.Auth, error) {

	auth := entity.Auth{
		Authenticated: 0,
	}

	credentials, err := a.twitterHandler.GetCredentials(token, secret)
 
    if err != nil {
		// TODO: ログ書き込み
        return &auth, err
    }

    // 認証成功後処理
	// DB
	user := entity.User{} 
	if err := a.userDbHandler.FindByTwitterID(user, twitterID);err != nil {
		// TODO: ログ書き込み
		return &auth, err
	}
	if user.ID == 0 {
		user.Name = twitterName
		user.TwitterID = twitterID
	}
	if err := a.userDbHandler.Upsert(user, "id", strconv.Itoa(user.ID));err != nil {
		// TODO: ログ書き込み
		return &auth, err
	}

	// Session
	a.twitterHandler.SetCredentials(credentials.GetToken(), credentials.GetSecret())
	//a.api.VerifyCredentials()
	auth.Authenticated = 1;

	a.sessionHandler.Set("token", credentials.GetToken())
    a.sessionHandler.Set("secret", credentials.GetSecret())
    a.sessionHandler.Save()

    return &auth, nil
}
