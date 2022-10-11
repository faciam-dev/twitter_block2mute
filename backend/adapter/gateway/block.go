package gateway

import (
	"errors"
	"log"
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

	// update blocks table
	registedBlockEntities := []entity.Block{}

	err = u.blockDbHandler.Transaction(func() error {
		// 登録済みのエンティティを取得する
		if err := u.blockDbHandler.FindAllByUserID(&registedBlockEntities, userID); err != nil {
			log.Print(err)
			return err
		}

		log.Print(registedBlockEntities)

		// blocksに登録されていないものを一括削除する
		deleteTargetBlocks := []entity.Block{}
		for _, registedBlockEntity := range registedBlockEntities {
			isFound := false
			for _, block := range blocks {
				if registedBlockEntity.TargetTwitterID == block.TargetTwitterID {
					isFound = true
					break
				}
			}
			if !isFound {
				deleteTargetBlocks = append(deleteTargetBlocks, registedBlockEntity)
			}
		}
		//err := u.dbHandler.Delete(deleteTargetBlocks)
		log.Print(deleteTargetBlocks)

		// 既存レコードは一切更新せずに登録する
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
