package handler

type TwitterHandler interface {
	SetCredentials(string, string)
	AuthorizationURL() (string, error)
	GetCredentials(string, string) (TwitterCredentials, error)
	GetRateLimits() error
}

type TwitterCredentials interface {
	GetToken() string
	GetSecret() string
}