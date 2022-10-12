package gateway

import (
	"errors"
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/entity"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/port"
)

type BlockRepository struct {
	blockDbHandler handler.BlockDbHandler
	userDbHandler  handler.UserDbHandler
	twitterHandler handler.TwitterHandler
	sessionHandler handler.SessionHandler
}

// NewBlockRepository はBlockRepositoryを返します．
func NewBlockRepository(dbHandler handler.BlockDbHandler, userDbHandler handler.UserDbHandler, twitterHandler handler.TwitterHandler, sessionHandler handler.SessionHandler) port.BlockRepository {
	return &BlockRepository{
		blockDbHandler: dbHandler,
		userDbHandler:  userDbHandler,
		twitterHandler: twitterHandler,
		sessionHandler: sessionHandler,
	}
}

// DBからid=userIDに該当するデータを取得する。
func (u *BlockRepository) GetUserIDs(userID string) (*[]entity.Block, int, error) {

	blocks := []entity.Block{}
	user := entity.User{}

	if err := u.userDbHandler.First(&user, userID); err != nil {
		return &blocks, 0, err
	}

	// auth Twitter
	token := u.sessionHandler.Get("token")
	secret := u.sessionHandler.Get("secret")

	if token == nil || secret == nil {
		return &blocks, 0, errors.New("session timeout or not found")
	}

	u.twitterHandler.UpdateTwitterApi(token.(string), secret.(string))

	// blocks
	twitterUserIds, err := u.twitterHandler.GetBlockedUser(user.TwitterID)
	if err != nil {
		return &blocks, 0, err
	}

	convertedSUserID64, _ := strconv.ParseInt(userID, 10, 64)
	nowTime := time.Now()

	total := twitterUserIds.GetTotal()
	for _, twitterUserId := range twitterUserIds.GetTwitterIDs() {
		block := entity.Block{
			UserID:          uint(convertedSUserID64),
			TargetTwitterID: twitterUserId,
			Flag:            0,
			CreatedAt:       nowTime,
			UpdatedAt:       nowTime,
		}
		blocks = append(blocks, block)
	}

	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i].TargetTwitterID <= blocks[j].TargetTwitterID
	})

	// update blocks table
	err = u.blockDbHandler.Transaction(func() error {
		// 登録済みのエンティティを取得する
		registedBlockEntities := []entity.Block{}
		if err := u.blockDbHandler.FindAllByUserID(&registedBlockEntities, userID); err != nil {
			log.Print(err)
			return err
		}

		sort.Slice(registedBlockEntities, func(i, j int) bool {
			return registedBlockEntities[i].TargetTwitterID <= registedBlockEntities[j].TargetTwitterID
		})

		//log.Print(registedBlockEntities)

		// api問い合わせ結果のblocksに登録されていないものを一括削除する
		IDs := []uint{}
		for _, registedBlockEntity := range registedBlockEntities {
			needle := registedBlockEntity.TargetTwitterID
			idx := sort.Search(len(blocks), func(i int) bool {
				return string(blocks[i].TargetTwitterID) >= needle
			})

			if blocks[idx].TargetTwitterID != needle {
				IDs = append(IDs, registedBlockEntity.ID)
			}
		}

		if err := u.blockDbHandler.DeleteByIds(&[]entity.Block{}, IDs); err != nil {
			log.Print(err)
			return err
		}

		// 既存レコードは一切更新せずに登録する
		// 登録済みエンティティでフラグが1であるものは除外する（Upsertでうまくいかないので）
		for n, block := range blocks {
			needle := block.TargetTwitterID
			idx := sort.Search(len(registedBlockEntities), func(i int) bool {
				return string(registedBlockEntities[i].TargetTwitterID) >= needle
			})

			if registedBlockEntities[idx].TargetTwitterID == needle {
				blocks[n].Flag = registedBlockEntities[idx].Flag
			}
		}
		if err := u.blockDbHandler.CreateNewBlocks(&blocks, "user_id", "twitter_id"); err != nil {
			log.Print(err)
			return err
		}

		return nil
	})

	if err == nil {
		return &blocks, total, err
	}

	return &blocks, total, nil
}
