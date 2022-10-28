package framework_test

import (
	"errors"
	"net/http"

	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/config"
	"github.com/faciam_dev/twitter_block2mute/backend/entity"
	"github.com/faciam_dev/twitter_block2mute/backend/infrastructure/framework"
	"github.com/gin-gonic/gin"
)

// テスト用ルーティング処理
// コントローラーmockの作成
// user
type User struct {
}

func (u *User) GetUserByID(contextHandler handler.ContextHandler) {
	table := []struct {
		id   string
		name string
	}{
		{
			"1",
			"test",
		},
	}

	userID := contextHandler.Param("id")
	outputFlag := false
	for _, tt := range table {
		if tt.id == userID {
			contextHandler.JSON(http.StatusOK, map[string]interface{}{
				"id":   userID,
				"name": tt.name,
			})
			outputFlag = true
			break
		}
	}
	if !outputFlag {
		contextHandler.JSON(http.StatusNotFound, map[string]interface{}{})
	}

}

// auth
type Auth struct {
}

func (a *Auth) IsAuth(contextHandler handler.ContextHandler) {
	auth := entity.Auth{
		Authenticated: 1,
	}

	contextHandler.JSON(http.StatusOK, map[string]interface{}{
		"result": auth.Authenticated,
	})
}

func (a *Auth) Auth(contextHandler handler.ContextHandler) {
	auth := entity.Auth{
		AuthUrl: "http://test/test",
	}

	contextHandler.JSON(http.StatusOK, map[string]interface{}{
		"to_url": auth.AuthUrl,
	})
}

func (a *Auth) Callback(contextHandler handler.ContextHandler) {
	token := contextHandler.Query("oauth_token")
	secret := contextHandler.Query("oauth_verifier")
	// twitterID := contextHandler.Query("user_id")
	// twitterName := contextHandler.Query("screen_name")

	// エラーがある場合
	if token == "errortoken" && secret == "errorsecret" {
		contextHandler.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": errors.New("callback error"),
		})
		return
	}

	// エラーなく表示
	contextHandler.JSON(http.StatusOK, map[string]interface{}{})
}

func (a *Auth) Logout(contextHandler handler.ContextHandler) {
	auth := entity.Auth{
		Logout: 1,
	}

	contextHandler.JSON(http.StatusOK, map[string]interface{}{
		"result": auth.Logout,
	})
}

// block
type Block struct {
}

func (b *Block) GetBlockByID(contextHandler handler.ContextHandler) {
	table := []struct {
		IDs   []string
		Total int
	}{
		{
			[]string{"1", "2", "3"},
			3,
		},
	}

	for _, tt := range table {
		contextHandler.JSON(http.StatusOK, map[string]interface{}{
			"ids":   tt.IDs,
			"total": tt.Total,
		})
	}

}

// block
type Block2Mute struct {
}

func (b *Block2Mute) All(contextHandler handler.ContextHandler) {
	table := []struct {
		NumSuccess int
		SuccessIDs []string
	}{
		{
			3,
			[]string{"1", "2", "3"},
		},
	}

	for _, tt := range table {
		contextHandler.JSON(http.StatusOK, map[string]interface{}{
			"num_success": tt.NumSuccess,
			"success_ids": tt.SuccessIDs,
		})
	}

}

// ルーティングの作成
type TestRouting struct {
	Gin  *gin.Engine
	Port string
}

func NewTestRouting(config *config.Config) *TestRouting {
	t := &TestRouting{
		Gin:  gin.Default(),
		Port: config.Routing.Port,
	}
	t.setTestRouting()
	return t
}

func (t *TestRouting) setTestRouting() {
	userController := User{}

	authController := Auth{}

	blockController := Block{}

	block2MuteController := Block2Mute{}

	// ルーティング割当
	// user
	t.Gin.POST("/user/user/:id", func(c *gin.Context) {
		userController.GetUserByID(framework.NewGinContextHandler(c))
	})

	// auth
	t.Gin.POST("/auth/auth", func(c *gin.Context) {
		authController.Auth(framework.NewGinContextHandler(c))
	})
	t.Gin.POST("/auth/is_auth", func(c *gin.Context) {
		authController.IsAuth(framework.NewGinContextHandler(c))
	})
	t.Gin.GET("/auth/auth_callback", func(c *gin.Context) {
		authController.Callback(framework.NewGinContextHandler(c))
	})
	t.Gin.POST("/auth/logout", func(c *gin.Context) {
		authController.Logout(framework.NewGinContextHandler(c))
	})

	// block
	t.Gin.POST("/blocks/ids", func(c *gin.Context) {
		blockController.GetBlockByID(framework.NewGinContextHandler(c))
	})

	// block2mute
	t.Gin.POST("/block2mute/all", func(c *gin.Context) {
		block2MuteController.All(framework.NewGinContextHandler(c))
	})

}
