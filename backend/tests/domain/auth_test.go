package domain_test

import (
	"testing"

	"github.com/faciam_dev/twitter_block2mute/backend/entity"
)

func TestCreateAuthDomain(t *testing.T) {
	type args struct {
		Authenticated    int
		AuthUrl          string
		OAuthToken       string
		OAuthTokenSecret string
		Logout           int
	}
	tests := []struct {
		name string
		args args
		want entity.Auth
	}{
		{
			name: "success",
			args: args{
				Authenticated:    1,
				AuthUrl:          "http://localhost/auth/auth",
				OAuthToken:       "token",
				OAuthTokenSecret: "secret",
				Logout:           0,
			},
			want: entity.Auth{
				Authenticated:    1,
				AuthUrl:          "http://localhost/auth/auth",
				OAuthToken:       "token",
				OAuthTokenSecret: "secret",
				Logout:           0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := entity.Auth{
				Authenticated:    tt.args.Authenticated,
				AuthUrl:          tt.args.AuthUrl,
				OAuthToken:       tt.args.OAuthToken,
				OAuthTokenSecret: tt.args.OAuthTokenSecret,
				Logout:           tt.args.Logout,
			}
			if got != tt.want {
				t.Errorf("createAuthDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}
