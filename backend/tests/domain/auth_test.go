package domain_test

import (
	"testing"

	"github.com/faciam_dev/twitter_block2mute/backend/entity"
)

func TestCreateAuthDomain(t *testing.T) {
	type args struct {
		Authenticated int
		AuthUrl       string
		Logout        int
	}

	tests := []struct {
		name string
		args args
		want *entity.Auth
	}{

		{
			name: "success_1",
			args: args{
				Authenticated: 1,
				AuthUrl:       "http://localhost/auth/auth",
				Logout:        0,
			},
			want: entity.NewAuth("http://localhost/auth/auth").SuccessAuthenticated(),
		},
		{
			name: "success_2",
			args: args{
				Authenticated: 0,
				AuthUrl:       "http://localhost/auth/auth",
				Logout:        1,
			},
			want: entity.NewAuth("http://localhost/auth/auth").SuccessLogout(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := entity.NewAuth(tt.args.AuthUrl)
			if tt.args.Authenticated == 1 {
				got.SuccessAuthenticated()
			}
			if tt.args.Logout == 1 {
				got.SuccessLogout()
			}

			if *got != *tt.want {
				t.Errorf("createAuthDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}
