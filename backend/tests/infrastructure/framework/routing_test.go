package framework_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/faciam_dev/twitter_block2mute/backend/config"
	"github.com/gin-gonic/gin"
)

func TestGetUserSelf(t *testing.T) {

	table := []struct {
		name       string
		method     string
		id         int
		wantStatus int
		wantBody   string
	}{
		{
			"success_200",
			"POST",
			1,
			200,
			"{\"id\":\"1\",\"name\":\"test\"}",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {

			config := config.NewConfig(".env.test")

			router := NewTestRouting(config)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			requestBody := bytes.NewBufferString("")

			c.Request, _ = http.NewRequest(tt.method, "/user/user/self", requestBody)

			router.Gin.ServeHTTP(w, c.Request)

			// statusの比較
			if w.Code != tt.wantStatus {
				t.Errorf("GetUserByID() status = %v, want %v", w.Code, tt.wantStatus)
			}

			// 中身の比較
			if w.Code == 200 && w.Body.String() != tt.wantBody {
				t.Errorf("GetUserByID() body = %v, want %v", w.Body, tt.wantBody)
			}
		})
	}
}

func TestIsAuth(t *testing.T) {

	table := []struct {
		name       string
		method     string
		wantStatus int
		wantBody   string
	}{
		{
			"success_200",
			"POST",
			200,
			"{\"result\":1}",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {

			config := config.NewConfig(".env.test")

			router := NewTestRouting(config)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			requestBody := bytes.NewBufferString("")
			c.Request, _ = http.NewRequest(tt.method, "/auth/is_auth", requestBody)

			router.Gin.ServeHTTP(w, c.Request)

			// statusの比較
			if w.Code != tt.wantStatus {
				t.Errorf("Auth() status = %v, want %v", w.Code, tt.wantStatus)
			}

			// 中身の比較
			if w.Code == 200 && w.Body.String() != tt.wantBody {
				t.Errorf("Auth() body = %v, want %v", w.Body, tt.wantBody)
			}
		})
	}
}

func TestAuth(t *testing.T) {

	table := []struct {
		name       string
		method     string
		wantStatus int
		wantBody   string
	}{
		{
			"success_200",
			"POST",
			200,
			"{\"to_url\":\"http://test/test\"}",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {

			config := config.NewConfig(".env.test")

			router := NewTestRouting(config)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			requestBody := bytes.NewBufferString("")
			c.Request, _ = http.NewRequest(tt.method, "/auth/auth", requestBody)

			router.Gin.ServeHTTP(w, c.Request)

			// statusの比較
			if w.Code != tt.wantStatus {
				t.Errorf("Auth() status = %v, want %v", w.Code, tt.wantStatus)
			}

			// 中身の比較
			if w.Code == 200 && w.Body.String() != tt.wantBody {
				t.Errorf("Auth() body = %v, want %v", w.Body, tt.wantBody)
			}
		})
	}
}

func TestCallback(t *testing.T) {

	type arg struct {
		oauthToken    string
		oauthVerifier string
		userId        string
		screenName    string
	}

	table := []struct {
		name       string
		method     string
		arg        arg
		wantStatus int
		wantBody   string
	}{
		{
			"success_200",
			"GET",
			arg{
				oauthToken:    "token",
				oauthVerifier: "secret",
				userId:        "1234567890",
				screenName:    "test",
			},
			200,
			"{}",
		},
		{
			"success_500",
			"GET",
			arg{
				oauthToken:    "errortoken",
				oauthVerifier: "errorsecret",
				userId:        "1234567890",
				screenName:    "test",
			},
			500,
			"{\"error\":{}}",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {

			config := config.NewConfig(".env.test")

			router := NewTestRouting(config)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(tt.method, "/auth/auth_callback", nil)

			queryString := c.Request.URL.Query()
			queryString.Add("oauth_token", string(tt.arg.oauthToken))
			queryString.Add("oauth_verifier", string(tt.arg.oauthVerifier))
			queryString.Add("user_id", string(tt.arg.oauthVerifier))
			queryString.Add("screen_name", string(tt.arg.oauthVerifier))
			c.Request.URL.RawQuery = queryString.Encode()

			router.Gin.ServeHTTP(w, c.Request)

			// statusの比較
			if w.Code != tt.wantStatus {
				t.Errorf("Auth() status = %v, want %v", w.Code, tt.wantStatus)
			}

			// 中身の比較
			if w.Body.String() != tt.wantBody {
				t.Errorf("Auth() body = %v, want %v", w.Body, tt.wantBody)
			}
		})
	}
}

func TestLogout(t *testing.T) {

	table := []struct {
		name       string
		method     string
		wantStatus int
		wantBody   string
	}{
		{
			"success_200",
			"POST",
			200,
			"{\"result\":1}",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {

			config := config.NewConfig(".env.test")

			router := NewTestRouting(config)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			requestBody := bytes.NewBufferString("")
			c.Request, _ = http.NewRequest(tt.method, "/auth/logout", requestBody)

			router.Gin.ServeHTTP(w, c.Request)

			// statusの比較
			if w.Code != tt.wantStatus {
				t.Errorf("Auth() status = %v, want %v", w.Code, tt.wantStatus)
			}

			// 中身の比較
			if w.Code == 200 && w.Body.String() != tt.wantBody {
				t.Errorf("Auth() body = %v, want %v", w.Body, tt.wantBody)
			}
		})
	}
}

func TestGetBlockByID(t *testing.T) {

	table := []struct {
		name       string
		method     string
		id         int
		wantStatus int
		wantBody   string
	}{
		{
			"success_200",
			"POST",
			1,
			200,
			"{\"ids\":[\"1\",\"2\",\"3\"],\"total\":3}",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {

			config := config.NewConfig(".env.test")

			router := NewTestRouting(config)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			requestBody := bytes.NewBufferString("")
			c.Request, _ = http.NewRequest(tt.method, "/blocks/ids", requestBody)

			router.Gin.ServeHTTP(w, c.Request)

			// statusの比較
			if w.Code != tt.wantStatus {
				t.Errorf("GetBlockByID() status = %v, want %v", w.Code, tt.wantStatus)
			}

			// 中身の比較
			if w.Code == 200 && w.Body.String() != tt.wantBody {
				t.Errorf("GetBlockByID() body = %v, want %v", w.Body, tt.wantBody)
			}
		})
	}
}

func TestAll(t *testing.T) {

	table := []struct {
		name       string
		method     string
		id         int
		wantStatus int
		wantBody   string
	}{
		{
			"success_200",
			"POST",
			1,
			200,
			"{\"num_success\":3,\"success_ids\":[\"1\",\"2\",\"3\"]}",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {

			config := config.NewConfig(".env.test")

			router := NewTestRouting(config)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			requestBody := bytes.NewBufferString("")
			c.Request, _ = http.NewRequest(tt.method, "/block2mute/all", requestBody)

			router.Gin.ServeHTTP(w, c.Request)

			// statusの比較
			if w.Code != tt.wantStatus {
				t.Errorf("All() status = %v, want %v", w.Code, tt.wantStatus)
			}

			// 中身の比較
			if w.Code == 200 && w.Body.String() != tt.wantBody {
				t.Errorf("All() body = %v, want %v", w.Body, tt.wantBody)
			}
		})
	}
}
