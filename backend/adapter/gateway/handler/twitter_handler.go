package handler

type TwitterHandler interface {
	SetCredentials(string, string)
	AuthorizationURL() (string, error)
	GetCredentials(string, string) (TwitterCredentials, TwitterValues, error)
	GetRateLimits() error
}

type TwitterCredentials interface {
	GetToken() string
	GetSecret() string
}

type TwitterValues interface {
	GetTwitterID() string
	GetTwitterScreenName() string
}
