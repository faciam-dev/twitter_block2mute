package gateway

import (
	"errors"

	"github.com/faciam_dev/twitter_block2mute/backend/entity"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/port"
	"gorm.io/gorm"
)

type UserRepository struct {
	conn *gorm.DB
}

// NewUserRepository はUserRepositoryを返します．
func NewUserRepository(conn *gorm.DB) port.UserRepository {
	return &UserRepository{
		conn: conn,
	}
}

// GetUserByID はDBからデータを取得します．
func (u *UserRepository) GetUserByID(userID string) (*entity.User, error) {
	conn := u.GetDBConn()
	user := entity.User{}

    conn.First(&user, userID)
    if user.ID <= 0 {
        return &entity.User{}, errors.New("user is not found")
    }
    return &user, nil


	/*
	db := interactor.DB.Connect()
    foundUser, err := interactor.User.FindById(db, userID)
	row := conn.QueryRowContext(ctx, "SELECT * FROM `user` WHERE id=?", userID)
	user := entity.User{}
	err := row.Scan(&user.ID, &user.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("User Not Found. UserID = %s", userID)
		}
		log.Println(err)
		return nil, errors.New("Internal Server Error. adapter/gateway/GetUserByID")
	}
	return &user, nil
	*/
}

// GetDBConn はconnectionを取得します．
func (u *UserRepository) GetDBConn() *gorm.DB {
	return u.conn
}