package infrastructure

import (
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/controller"
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway"
	"github.com/faciam_dev/twitter_block2mute/backend/adapter/presenter"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/interactor"
	"github.com/gin-gonic/gin"
)

type Routing struct {
    DB *DB
    Gin *gin.Engine
    Port string
}

func NewRouting(db *DB) *Routing {
    c := NewConfig()
    r := &Routing{
        DB: db,
        Gin: gin.Default(),
        Port: c.Routing.Port,
    }
    r.setRouting()
    return r
}

func (r *Routing) setRouting() {
	userController := controller.User{
		OutputFactory: presenter.NewUserOutputPort,
		InputFactory:  interactor.NewUserInputPort,
		RepoFactory:   gateway.NewUserRepository,
		Conn:          r.DB.Connection,
	}
    r.Gin.GET("/users/:id", func (c *gin.Context) {userController.GetUserByID(c) })
}

func (r *Routing) Run() {
    r.Gin.Run(r.Port)
}