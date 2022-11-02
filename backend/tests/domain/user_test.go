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
		CreatedAt int64
		UpdatedAt int64
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
				CreatedAt: 1663673740,
				UpdatedAt: 1663673800,
			},
			want: entity.User{
				ID:        1,
				Name:      "test1",
				TwitterID: "12345678901234567890",
				CreatedAt: 1663673740,
				UpdatedAt: 1663673800,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := entity.User{
				ID:        tt.args.ID,
				Name:      tt.args.Name,
				TwitterID: tt.args.TwitterID,
				CreatedAt: tt.args.CreatedAt,
				UpdatedAt: tt.args.UpdatedAt,
			}
			if got != tt.want {
				t.Errorf("createUserDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}
