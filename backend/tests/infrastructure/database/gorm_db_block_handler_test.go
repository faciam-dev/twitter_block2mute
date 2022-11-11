package database_test

import (
	"errors"
	"testing"
	"time"

	"github.com/faciam_dev/twitter_block2mute/backend/entity"
)

func TestFindAllByTwitterID(t *testing.T) {
	type arg struct {
		value string
	}
	nowTime := time.Now()
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
				*entity.NewBlock(1, "12345678901", 1, nowTime, nowTime),
				*entity.NewBlock(1, "12345678902", 1, nowTime, nowTime),
				*entity.NewBlock(1, "12345678903", 0, nowTime, nowTime),
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
					if gotBlock.GetUserID() == wantBlock.GetUserID() &&
						gotBlock.GetTargetTwitterID() == wantBlock.GetTargetTwitterID() &&
						gotBlock.GetFlag() == wantBlock.GetFlag() {
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
	nowTime := time.Now()
	table := []struct {
		name      string
		arg       arg
		wantBlock entity.Block
	}{
		{
			"success_insert",
			arg{
				column:       "id",
				value:        "4",
				createSource: *entity.NewBlock(1, "12345678904", 1, nowTime, nowTime),
			},
			*entity.NewBlock(1, "12345678904", 1, nowTime, nowTime),
		},
		{
			"success_update",
			arg{
				column:       "id",
				value:        "3",
				createSource: *entity.NewBlock(1, "12345678903", 0, nowTime, nowTime),
			},
			*entity.NewBlock(1, "12345678903", 0, nowTime, nowTime),
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
				if blocks[0].GetFlag() != tt.wantBlock.GetFlag() || blocks[0].GetUserID() != tt.wantBlock.GetUserID() || blocks[0].GetTargetTwitterID() != tt.wantBlock.GetTargetTwitterID() {
					t.Errorf("CreateNewBlocks()  = %v, want %v", blocks[0], tt.wantBlock)
				}
			})
		}

		return errors.New("rollback")
	})
}
