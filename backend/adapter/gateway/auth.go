package gateway

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/faciam_dev/twitter_block2mute/backend/entity"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/port"
	"github.com/garyburd/go-oauth/oauth"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AuthRepository struct {
	context *gin.Context
	api *anaconda.TwitterApi
	callbackUrl string
}

// NewUserRepository はUserRepositoryを返します．
func NewAuthRepository(ctx *gin.Context, api *anaconda.TwitterApi, callbackUrl string) port.AuthRepository {
	return &AuthRepository{
		context: ctx,
		api: api,
		callbackUrl: callbackUrl,
	}
}

// セッションに入っている値から認証を実施し、認証状態を取得する。
func (a *AuthRepository) IsAuthenticated() (*entity.Auth, error) {
	session := a.GetSession() 

	token := session.Get("token")
	secret := session.Get("secret") 

	auth := entity.Auth{
		Authenticated: 0,
	}

	if token == nil || secret == nil {
		return &auth, nil
	}

	//api = anaconda.NewTwitterApi(token, secret)
	a.api.Credentials.Token = token.(string)
	a.api.Credentials.Secret = secret.(string)
	//a.api.VerifyCredentials()

	//api := anaconda.NewTwitterApiWithCredentials(token.(string), secret.(string), consumerKey, consumerSecret)
	r := []string{}
	_, err:= a.api.GetRateLimits(r)

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

	uri, _, err := a.api.AuthorizationURL(a.callbackUrl)

	if err != nil {
		return &auth, err;
	}

	auth.AuthUrl = uri

	return &auth, nil
}

// 認証コールバック
func (a *AuthRepository) Callback() (*entity.Auth, error) {
    token := a.context.Query("oauth_token")
    secret := a.context.Query("oauth_verifier")

	auth := entity.Auth{
		Authenticated: 0,
	}
 
    credentials, _, err := a.api.GetCredentials(&oauth.Credentials{
        Token: token,
    }, secret)
    if err != nil {
		// TODO: ログ書き込み
        return &auth, err
    }

    // 認証成功後処理	
	//api = anaconda.NewTwitterApi(credentials.Token, credentials.Secret)
	a.api.Credentials.Token = credentials.Token
	a.api.Credentials.Secret = credentials.Secret
	//a.api.VerifyCredentials()
	auth.Authenticated = 1;

    session := a.GetSession() 
    session.Set("token", credentials.Token)
    session.Set("secret", credentials.Secret)
    session.Save()
 
    return &auth, nil
}

func (a *AuthRepository) GetSession() sessions.Session {
	session := sessions.Default(a.context)
	return session
}
