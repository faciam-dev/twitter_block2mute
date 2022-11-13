package entity

import (
	"sort"
	"time"
)

type Block struct {
	id              uint
	userID          uint
	targetTwitterID string
	flag            int
	createdAt       time.Time
	updatedAt       time.Time
}

type Blocks []Block

func NewBlock(
	userID uint,
	targetTwitterID string,
	flag int,
	createdAt time.Time,
	updatedAt time.Time,
) *Block {
	block := Block{}
	block.userID = userID
	block.targetTwitterID = targetTwitterID
	block.flag = flag

	return &block
}

// getter
func (b *Block) GetID() uint {
	return b.id
}

func (b *Block) GetUserID() uint {
	return b.userID
}

func (b *Block) GetTargetTwitterID() string {
	return b.targetTwitterID
}

func (b *Block) GetFlag() int {
	return b.flag
}

func (b *Block) GetCreatedAt() time.Time {
	return b.createdAt
}

func (b *Block) GetUpdatedAt() time.Time {
	return b.updatedAt
}

// ステータス変更など
// Update
func (b *Block) Update(
	id uint,
	userID uint,
	targetTwitterID string,
	flag int,
	createdAt time.Time,
	updatedAt time.Time,

) {
	b.id = id
	b.userID = userID
	b.targetTwitterID = targetTwitterID
	b.flag = flag
	b.createdAt = createdAt
	b.updatedAt = updatedAt
}

// Flag
func (b *Block) Converted() {
	b.flag = 1
}

func (b *Block) NotConverted() {
	b.flag = 0
}

func (b *Block) IsConverted() bool {
	return b.flag != 0
}

// ソート
func (bs *Blocks) SortByTargetTwtitterID() {
	sort.Slice(*bs, func(i, j int) bool {
		return (*bs)[i].GetTargetTwitterID() <= (*bs)[j].GetTargetTwitterID()
	})
}

// 検索
// TargetTwitterIDによる検索
func (bs *Blocks) FindByTargetTwitterID(needle string) int {
	return sort.Search(len(*bs), func(i int) bool {
		return string((*bs)[i].GetTargetTwitterID()) == needle
	})
}

// 対象のBlocksに存在しないIDsを検索する。
func (bs *Blocks) FindAllIDsNotFoundWithTwitterID(targetBlocks *Blocks) []uint {
	IDs := []uint{}
	for _, b := range *bs {
		idx := targetBlocks.FindByTargetTwitterID(b.GetTargetTwitterID())

		// targetBlocksの長さが0なら登録されていない。
		// または、idxがtargetBlocks総数以上の場合は登録されていない
		if len(*targetBlocks) == 0 || len(*targetBlocks) <= idx {
			IDs = append(IDs, b.GetID())
		}
	}

	return IDs
}

// フィルタ性質のあるゲッタ
// 変換済みではないblocksを得る
func (bs *Blocks) GetNotConvertedBlocks() *Blocks {
	blocks := &Blocks{}
	for _, v := range *bs {
		if v.flag == 0 {
			*blocks = append(*blocks, v)
		}
	}
	return blocks
}

// コンバート
// Blocks への変換
func (bs *Blocks) ToBlocksEntity(blocks *[]Block) {
	for _, v := range *blocks {
		*bs = append(*bs, v)
	}
}

// *[]Block への変換
func (bs *Blocks) ToBlockEntities() *[]Block {
	blocks := []Block{}
	for _, v := range *bs {
		blocks = append(blocks, v)
	}
	return &blocks
}
