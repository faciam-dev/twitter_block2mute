package domain_test

import (
	"testing"
	"time"

	"github.com/faciam_dev/twitter_block2mute/backend/entity"
)

func TestCreateMuteDomain(t *testing.T) {
	type args struct {
		ID              uint
		TargetTwitterID string
		Flag            int
		CreatedAt       time.Time
		UpdatedAt       time.Time
	}
	timeNow := time.Now()
	tests := []struct {
		name string
		args args
		want entity.Mute
	}{
		{
			name: "success",
			args: args{
				ID:              1,
				TargetTwitterID: "1234567890",
				Flag:            1,
				CreatedAt:       timeNow,
				UpdatedAt:       timeNow,
			},
			want: entity.Mute{
				ID:              1,
				TargetTwitterID: "1234567890",
				Flag:            1,
				CreatedAt:       timeNow,
				UpdatedAt:       timeNow,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := entity.Mute{
				ID:              tt.args.ID,
				TargetTwitterID: tt.args.TargetTwitterID,
				Flag:            tt.args.Flag,
				CreatedAt:       tt.args.CreatedAt,
				UpdatedAt:       tt.args.UpdatedAt,
			}
			if got != tt.want {
				t.Errorf("createBlockDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}
