package handler

type TwitterHandler interface {
	UpdateTwitterApi(string, string)
	AuthorizationURL() (string, error)
	GetCredentials(string, string) (TwitterCredentials, TwitterValues, error)
	GetUser(string) (TwitterUser, error)
}

type TwitterCredentials interface {
	GetToken() string
	GetSecret() string
}

type TwitterValues interface {
	GetTwitterID() string
	GetTwitterScreenName() string
}

type TwitterUser interface {
	GetTwitterID() string
	GetTwitterScreenName() string
	GetTwitterName() string
}
