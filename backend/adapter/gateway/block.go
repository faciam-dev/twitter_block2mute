package gateway

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/entity"
	"github.com/faciam_dev/twitter_block2mute/backend/infrastructure/database"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/port"
)

type BlockRepository struct {
	loggerHandler  handler.LoggerHandler
	dbHandler      handler.DBHandler
	twitterHandler handler.TwitterHandler
	sessionHandler handler.SessionHandler
}

// NewBlockRepository はBlockRepositoryを返します．
func NewBlockRepository(
	loggerHandler handler.LoggerHandler,
	dbHandler handler.DBHandler,
	twitterHandler handler.TwitterHandler,
	sessionHandler handler.SessionHandler,
) port.BlockRepository {
	return &BlockRepository{
		loggerHandler:  loggerHandler,
		dbHandler:      dbHandler,
		twitterHandler: twitterHandler,
		sessionHandler: sessionHandler,
	}
}

// ユーザーを取得する
func (u *BlockRepository) GetUser(userID string) *entity.User {
	user := entity.User{}
	userDBHandler := database.NewUserDbHandler(u.dbHandler.Connect())

	if err := userDBHandler.First(&user, userID); err != nil {
		u.loggerHandler.Errorf("user not found (user_id=%s)", userID)
		return &user
	}

	return &user
}

// APIからBlocksを得る。
func (u *BlockRepository) GetBlocks(user *entity.User) (*entity.Blocks, int, error) {
	blocks := entity.Blocks{}

	// auth Twitter
	token := u.sessionHandler.Get("token")
	secret := u.sessionHandler.Get("secret")

	if token == nil || secret == nil {
		return &blocks, 0, errors.New("session timeout or not found")
	}

	u.twitterHandler.UpdateTwitterApi(token.(string), secret.(string))
	u.loggerHandler.Debugf("update twitter api (user_id=%d token=%s secret=%s)", user.GetID(), token.(string), secret.(string))

	// blocks
	twitterUserIds, err := u.twitterHandler.GetBlockedUser(user.GetTwitterID())

	// 0件以外は変換できるようにする。
	// ブロック数が多い場合の対策。
	if err != nil && len(twitterUserIds.GetTwitterIDs()) == 0 {
		u.loggerHandler.Warnw("twitter blocklist API error.", "twitter_id", user.GetTwitterID(), "error", err)
		return &blocks, 0, err
	}

	nowTime := time.Now()

	total := twitterUserIds.GetTotal()
	for _, twitterUserId := range twitterUserIds.GetTwitterIDs() {
		block := entity.NewBlock(
			user.GetID(),
			twitterUserId,
			0,
			nowTime,
			nowTime,
		)
		blocks = append(blocks, *block)
	}
	blocks.SortByTargetTwtitterID()

	u.loggerHandler.Debugf("user_id=%d Num_Of_blocks=%d", user.GetID(), len(blocks))

	return &blocks, total, nil
}

// update blocks table
func (u *BlockRepository) TxUpdateAndDeleteBlocks(user *entity.User, blocks *entity.Blocks) error {
	//err := u.blockDbHandler.Transaction(func() error {
	err := u.dbHandler.Transaction(func(tx handler.DbConnection) error {
		blockDbHandler := database.NewBlockDbHandler(tx)

		// 登録済みのエンティティを取得する
		registedBlocks := []entity.Block{}
		if err := blockDbHandler.FindAllByUserID(&registedBlocks, strconv.FormatUint(uint64(user.GetID()), 10)); err != nil {
			u.loggerHandler.Errorw("FindAllByUserID", "error", err)
			return err
		}
		u.loggerHandler.Debugf("registedBlocks: num=%d", len(registedBlocks))

		registedBlockEntities := entity.Blocks{}
		registedBlockEntities.ToBlocksEntity(&registedBlocks)
		u.loggerHandler.Debugf("registedBlockEntities: num=%d", len(registedBlockEntities))

		registedBlockEntities.SortByTargetTwtitterID()
		IDs := registedBlockEntities.FindAllIDsNotFoundWithTwitterID(blocks)
		if len(IDs) > 0 {
			if err := blockDbHandler.DeleteByIds(&[]entity.Block{}, IDs); err != nil {
				log.Print(err)
				return err
			}
		}

		// 既存レコードは一切更新せずに登録する
		blocks = blocks.GetNotConvertedBlocks()
		if len(*blocks) > 0 {
			if err := blockDbHandler.CreateNewBlocks(blocks.ToBlockEntities(), "user_id", "twitter_id"); err != nil {
				u.loggerHandler.Errorw(
					"fail to create new blocks.",
					"user_id",
					user.GetID(),
					"twitter_id",
					user.GetTwitterID(),
					"error",
					err,
				)
				return err
			}
		}

		return nil
	})

	if err == nil {
		return err
	}

	return nil
}
