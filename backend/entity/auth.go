package entity

type Auth struct {
	Authenticated    int
	AuthUrl          string
	OAuthToken       string
	OAuthTokenSecret string
	Logout           int
}
