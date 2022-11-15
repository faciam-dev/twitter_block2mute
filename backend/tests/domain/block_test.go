package domain_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/faciam_dev/twitter_block2mute/backend/entity"
)

type BlockArgs struct {
	ID              uint
	TargetTwitterID string
	Flag            int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Blocksを生成する
func CreateBlocks(args []BlockArgs) *entity.Blocks {
	blocks := entity.Blocks{}
	for _, arg := range args {
		blocks = append(
			blocks,
			*entity.NewBlock(
				arg.ID,
				arg.TargetTwitterID,
				arg.Flag,
				arg.CreatedAt,
				arg.UpdatedAt,
			),
		)
	}
	return &blocks
}

func TestCreateBlockDomain(t *testing.T) {
	timeNow := time.Now()
	tests := []struct {
		name string
		args BlockArgs
		want entity.Block
	}{
		{
			name: "success",
			args: BlockArgs{
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

func TestBlocksSortByTargetTwtitterID(t *testing.T) {
	timeNow := time.Now()
	tests := []struct {
		name string
		args []BlockArgs
		want entity.Blocks
	}{
		{
			name: "success",
			args: []BlockArgs{
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
			want: entity.Blocks{
				*entity.NewBlock(
					2, "1234567890", 0, timeNow, timeNow,
				),
				*entity.NewBlock(
					1, "1234567891", 1, timeNow, timeNow,
				),
				*entity.NewBlock(
					3, "1234567893", 0, timeNow, timeNow,
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateBlocks(tt.args)

			got.SortByTargetTwtitterID()

			if reflect.DeepEqual(got, tt.want) {
				t.Errorf("SortByTargetTwtitterID() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestBlocksFindByTargetTwitterID(t *testing.T) {
	timeNow := time.Now()
	tests := []struct {
		name   string
		args   []BlockArgs
		needle string
		want   int
	}{
		{
			name: "found",
			args: []BlockArgs{
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
			needle: "1234567890",
			want:   1,
		},
		{
			name: "not_found",
			args: []BlockArgs{
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
			want:   3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blocks := CreateBlocks(tt.args)

			got := blocks.FindByTargetTwitterID(tt.needle)

			if got != tt.want {
				t.Errorf("SortByTargetTwtitterID() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestBlocksFindAllIDsNotFoundWithTwitterID(t *testing.T) {
	timeNow := time.Now()
	tests := []struct {
		name             string
		myBlocksArgs     []BlockArgs
		targetBlocksArgs []BlockArgs
		want             entity.Blocks
	}{
		{
			name: "found",
			myBlocksArgs: []BlockArgs{
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
			targetBlocksArgs: []BlockArgs{
				{
					ID:              1,
					TargetTwitterID: "1234567891",
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
			want: entity.Blocks{
				*entity.NewBlock(
					1, "1234567891", 1, timeNow, timeNow,
				),
				*entity.NewBlock(
					3, "1234567893", 0, timeNow, timeNow,
				),
			},
		},
		{
			name: "not_found",
			myBlocksArgs: []BlockArgs{
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
			targetBlocksArgs: []BlockArgs{},
			want: entity.Blocks{
				*entity.NewBlock(
					2, "1234567890", 0, timeNow, timeNow,
				),
				*entity.NewBlock(
					1, "1234567891", 1, timeNow, timeNow,
				),
				*entity.NewBlock(
					3, "1234567893", 0, timeNow, timeNow,
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			myBlocks := CreateBlocks(tt.myBlocksArgs)
			targetBlocks := CreateBlocks(tt.targetBlocksArgs)

			got := myBlocks.FindAllIDsNotFoundWithTwitterID(targetBlocks)

			if reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindAllIDsNotFoundWithTwitterID() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestGetNotConvertedBlocks(t *testing.T) {
	timeNow := time.Now()
	tests := []struct {
		name string
		args []BlockArgs
		want entity.Blocks
	}{
		{
			name: "found",
			args: []BlockArgs{
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
			want: entity.Blocks{
				*entity.NewBlock(
					2, "1234567890", 0, timeNow, timeNow,
				),
				*entity.NewBlock(
					3, "1234567893", 0, timeNow, timeNow,
				),
			},
		},
		{
			name: "not_found",
			args: []BlockArgs{
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
					Flag:            1,
					CreatedAt:       timeNow,
					UpdatedAt:       timeNow,
				},
			},
			want: entity.Blocks{
				*entity.NewBlock(
					2, "1234567890", 1, timeNow, timeNow,
				),
				*entity.NewBlock(
					1, "1234567891", 1, timeNow, timeNow,
				),
				*entity.NewBlock(
					3, "1234567893", 1, timeNow, timeNow,
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blocks := CreateBlocks(tt.args)

			got := blocks.GetNotConvertedBlocks()

			if reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindAllIDsNotFoundWithTwitterID() = %v, want %v", got, tt.want)
			}

		})
	}
}
