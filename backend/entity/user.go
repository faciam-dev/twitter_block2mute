package entity

type User struct {
	id          uint
	name        string
	accountName string
	twitterID   string
	createdAt   int64
	updatedAt   int64
}

// Authのファクトリ関数、指定なし
func NewBlankUser() *User {
	return &User{}
}

// Authのファクトリ関数
func NewUser(name string, accountName string, twitterID string) *User {
	user := User{}

	user.name = name
	user.accountName = accountName
	user.twitterID = twitterID

	return &user
}

// getter
// get id
func (u *User) GetID() uint {
	return u.id
}

// get name
func (u *User) GetName() string {
	return u.name
}

// get AccountName
func (u *User) GetAccountName() string {
	return u.accountName
}

// get twitterID
func (u *User) GetTwitterID() string {
	return u.twitterID
}

// 情報をUpdateする。濫用注意。
func (u *User) Update(id uint, name string, accountName string, twitterID string) *User {
	u.id = id
	u.name = name
	u.accountName = accountName
	u.twitterID = twitterID

	return u
}
