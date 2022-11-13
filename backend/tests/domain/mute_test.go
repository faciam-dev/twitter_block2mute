package domain_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/faciam_dev/twitter_block2mute/backend/entity"
)

type MuteArgs struct {
	ID              uint
	TargetTwitterID string
	Flag            int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Mutesを生成する
func CreateMutes(args []MuteArgs) *entity.Mutes {
	mutes := entity.Mutes{}
	for _, arg := range args {
		mutes = append(
			mutes,
			*entity.NewMute(
				arg.ID,
				arg.TargetTwitterID,
				arg.Flag,
				arg.CreatedAt,
				arg.UpdatedAt,
			),
		)
	}
	return &mutes
}

func TestCreateMuteDomain(t *testing.T) {
	type args struct {
		ID              uint
		TargetTwitterID string
		Flag            int
		CreatedAt       time.Time
		UpdatedAt       time.Time
	}
	nowTime := time.Now()
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
				CreatedAt:       nowTime,
				UpdatedAt:       nowTime,
			},
			want: *entity.NewMute(1, "1234567890", 1, nowTime, nowTime),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := *entity.NewMute(
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

func TestMuteSortByTargetTwtitterID(t *testing.T) {
	timeNow := time.Now()
	tests := []struct {
		name string
		args []MuteArgs
		want entity.Mutes
	}{
		{
			name: "success",
			args: []MuteArgs{
				{
					ID:              1,
					TargetTwitterID: "1234567891",
					Flag:            1,
					CreatedAt:       timeNow,
					UpdatedAt:       timeNow,
				},
				{
					ID:              2,
					TargetTwitterID: "1234567890",
					Flag:            0,
					CreatedAt:       timeNow,
					UpdatedAt:       timeNow,
				},
				{
					ID:              3,
					TargetTwitterID: "1234567893",
					Flag:            0,
					CreatedAt:       timeNow,
					UpdatedAt:       timeNow,
				},
			},
			want: entity.Mutes{
				*entity.NewMute(
					2, "1234567890", 0, timeNow, timeNow,
				),
				*entity.NewMute(
					1, "1234567891", 1, timeNow, timeNow,
				),
				*entity.NewMute(
					3, "1234567893", 0, timeNow, timeNow,
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateMutes(tt.args)

			got.SortByTargetTwtitterID()

			if reflect.DeepEqual(got, tt.want) {
				t.Errorf("SortByTargetTwtitterID() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestMuteFindByTargetTwitterID(t *testing.T) {
	timeNow := time.Now()
	tests := []struct {
		name   string
		args   []MuteArgs
		needle string
		want   bool
	}{
		{
			name: "found",
			args: []MuteArgs{
				{
					ID:              1,
					TargetTwitterID: "1234567891",
					Flag:            1,
					CreatedAt:       timeNow,
					UpdatedAt:       timeNow,
				},
				{
					ID:              2,
					TargetTwitterID: "1234567890",
					Flag:            1,
					CreatedAt:       timeNow,
					UpdatedAt:       timeNow,
				},
				{
					ID:              3,
					TargetTwitterID: "1234567893",
					Flag:            0,
					CreatedAt:       timeNow,
					UpdatedAt:       timeNow,
				},
			},
			needle: "1234567890",
			want:   true,
		},
		{
			name: "not_found",
			args: []MuteArgs{
				{
					ID:              1,
					TargetTwitterID: "1234567891",
					Flag:            1,
					CreatedAt:       timeNow,
					UpdatedAt:       timeNow,
				},
				{
					ID:              2,
					TargetTwitterID: "1234567890",
					Flag:            0,
					CreatedAt:       timeNow,
					UpdatedAt:       timeNow,
				},
				{
					ID:              3,
					TargetTwitterID: "1234567893",
					Flag:            0,
					CreatedAt:       timeNow,
					UpdatedAt:       timeNow,
				},
			},
			needle: "1234567894",
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blocks := CreateMutes(tt.args)

			got := blocks.IsConvertedByTwitterID(tt.needle)

			if got != tt.want {
				t.Errorf("SortByTargetTwtitterID() = %v, want %v", got, tt.want)
			}

		})
	}
}
