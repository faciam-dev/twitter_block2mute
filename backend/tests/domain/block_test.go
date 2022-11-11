package domain_test

import (
	"testing"
	"time"

	"github.com/faciam_dev/twitter_block2mute/backend/entity"
)

func TestCreateBlockDomain(t *testing.T) {
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
		want entity.Block
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
			want: *entity.NewBlock(
				1, "1234567890", 1, timeNow, timeNow,
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := *entity.NewBlock(
				tt.args.ID,
				tt.args.TargetTwitterID,
				tt.args.Flag,
				tt.args.CreatedAt,
				tt.args.UpdatedAt,
			)
			if got != tt.want {
				t.Errorf("createBlockDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}
