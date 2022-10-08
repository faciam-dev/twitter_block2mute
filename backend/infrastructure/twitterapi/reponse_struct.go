package twitterapi

// TwitterUser
type TwitterUser struct {
	ID          string
	ScreenName  string
	TwitterName string
}

func (t TwitterUser) GetTwitterID() string {
	return t.ID
}

func (t TwitterUser) GetTwitterScreenName() string {
	return t.ScreenName
}

func (t TwitterUser) GetTwitterName() string {
	return t.TwitterName
}
