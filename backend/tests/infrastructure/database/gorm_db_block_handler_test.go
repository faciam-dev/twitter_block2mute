package database_test

import (
	"errors"
	"testing"

	"github.com/faciam_dev/twitter_block2mute/backend/entity"
)

func TestFindAllByTwitterID(t *testing.T) {
	type arg struct {
		value string
	}

	table := []struct {
		name         string
		arg          arg
		wantEntities []entity.Block
	}{
		{
			"success",
			arg{
				value: "1",
			},
			[]entity.Block{
				{
					UserID:          1,
					TargetTwitterID: "12345678901",
					Flag:            1,
				},
				{
					UserID:          1,
					TargetTwitterID: "12345678902",
					Flag:            1,
				},
				{
					UserID:          1,
					TargetTwitterID: "12345678903",
					Flag:            0,
				},
			},
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			blocks := []entity.Block{}
			if err := BlockDbHandler.FindAllByUserID(&blocks, tt.arg.value); err != nil {
				t.Error(err)
			}

			// 中身の比較
			for i, wantBlock := range tt.wantEntities {
				isError := true
				for _, gotBlock := range blocks {
					if gotBlock.UserID == wantBlock.UserID &&
						gotBlock.TargetTwitterID == wantBlock.TargetTwitterID &&
						gotBlock.Flag == wantBlock.Flag {
						isError = false
						break
					}
				}
				if isError {
					t.Errorf("FindAllByTwitterID(%v) got = %v, want = %v[%v]", tt.arg.value, wantBlock, blocks, i)
				}
			}
		})
	}
}

func TestCreateNewBlocks(t *testing.T) {
	type arg struct {
		column       string
		value        string
		createSource entity.Block
	}

	table := []struct {
		name      string
		arg       arg
		wantBlock entity.Block
	}{
		{
			"success_insert",
			arg{
				column: "id",
				value:  "4",
				createSource: entity.Block{
					UserID:          1,
					TargetTwitterID: "12345678904",
					Flag:            1,
				},
			},
			entity.Block{
				UserID:          1,
				TargetTwitterID: "12345678904",
				Flag:            1,
			},
		},
		{
			"success_update",
			arg{
				column: "id",
				value:  "3",
				createSource: entity.Block{
					UserID:          1,
					TargetTwitterID: "12345678903",
					Flag:            0,
				},
			},
			entity.Block{
				UserID:          1,
				TargetTwitterID: "12345678903",
				Flag:            0,
			},
		},
	}

	BlockDbHandler.Transaction(func() error {
		for _, tt := range table {
			t.Run(tt.name, func(t *testing.T) {
				// 作成・更新する
				blocks := []entity.Block{tt.arg.createSource}
				err := BlockDbHandler.CreateNewBlocks(&blocks, tt.arg.column, tt.arg.value)
				if err != nil {
					t.Error(err)
				}

				// 中身の比較
				if blocks[0].Flag != tt.wantBlock.Flag || blocks[0].UserID != tt.wantBlock.UserID || blocks[0].TargetTwitterID != tt.wantBlock.TargetTwitterID {
					t.Errorf("CreateNewBlocks()  = %v, want %v", blocks[0], tt.wantBlock)
				}
			})
		}

		return errors.New("rollback")
	})
}
