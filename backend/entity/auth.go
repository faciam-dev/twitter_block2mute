package entity

type Auth struct {
	authenticated int
	authUrl       string
	logout        int
}

// Authのファクトリ関数
func NewAuth(authUrl string) *Auth {
	auth := Auth{}

	auth.authenticated = 0
	auth.logout = 0
	auth.authUrl = authUrl

	return &auth
}

// getter
// get logout
func (a *Auth) GetLogout() int {
	return a.logout
}

// get logout
func (a *Auth) GetAuthenticated() int {
	return a.authenticated
}

// get auth url
func (a *Auth) GetAuthUrl() string {
	return a.authUrl
}

// methods
// 認証成功時の処理
func (a *Auth) SuccessAuthenticated() *Auth {
	a.authenticated = 1
	return a
}

func (a *Auth) SuccessLogout() *Auth {
	a.logout = 1
	return a
}
