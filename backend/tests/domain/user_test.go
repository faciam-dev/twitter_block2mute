package domain_test

import (
	"testing"

	"github.com/faciam_dev/twitter_block2mute/backend/entity"
)

func TestCreateUserDomain(t *testing.T) {
	type args struct {
		ID        uint
		Name      string
		TwitterID string
	}
	tests := []struct {
		name string
		args args
		want entity.User
	}{
		{
			name: "success",
			args: args{
				ID:        1,
				Name:      "test1",
				TwitterID: "12345678901234567890",
			},
			want: *entity.NewBlankUser().Update(
				1,
				"test1",
				"test1",
				"12345678901234567890",
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := *entity.NewBlankUser().Update(
				tt.args.ID,
				tt.args.Name,
				tt.args.Name,
				tt.args.TwitterID,
			)

			if got != tt.want {
				t.Errorf("createUserDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}
