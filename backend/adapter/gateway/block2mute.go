package gateway

import (
	"errors"
	"strconv"
	"time"

	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/entity"
	"github.com/faciam_dev/twitter_block2mute/backend/infrastructure/database"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/port"
)

type Block2MuteRepository struct {
	loggerHandler  handler.LoggerHandler
	dbHandler      handler.DBHandler
	twitterHandler handler.TwitterHandler
	sessionHandler handler.SessionHandler
}

// NewBlockRepository はBlockRepositoryを返します．
func NewBlock2MuteRepository(
	loggerHandler handler.LoggerHandler,
	dbHandler handler.DBHandler,
	twitterHandler handler.TwitterHandler,
	sessionHandler handler.SessionHandler,
) port.Block2MuteRepository {
	return &Block2MuteRepository{
		loggerHandler:  loggerHandler,
		dbHandler:      dbHandler,
		twitterHandler: twitterHandler,
		sessionHandler: sessionHandler,
	}
}

// ユーザーを取得する
func (b *Block2MuteRepository) GetUser(userID string) *entity.User {
	user := entity.User{}
	userDBHandler := database.NewUserDBHandler(b.dbHandler.Connect())

	if err := userDBHandler.First(&user, userID); err != nil {
		b.loggerHandler.Errorf("user not found (user_id=%s)", userID)
		return &user
	}

	return &user
}

// auth Twitter
func (b *Block2MuteRepository) AuthTwitter() error {
	token := b.sessionHandler.Get("token")
	secret := b.sessionHandler.Get("secret")

	if token == nil || secret == nil {
		return errors.New("session timeout or not found")
	}

	b.twitterHandler.UpdateTwitterApi(token.(string), secret.(string))
	b.loggerHandler.Debugf("update twitter api (token=%s secret=%s)", token.(string), secret.(string))

	return nil
}

func (b *Block2MuteRepository) All(user *entity.User) (*entity.Block2Mute, error) {
	userID := strconv.FormatUint(uint64(user.GetID()), 10)

	// update blocks table
	convertedIDs := []string{}
	//	err := b.muteDBHandler.Transaction(func() error {
	err := b.dbHandler.Transaction(func(tx handler.DBConnection) error {
		blockDBHandler := database.NewBlockDBHandler(tx)
		muteDBHandler := database.NewMuteHandler(tx)

		//err := b.muteDBHandler.Transaction(func() error {
		// block2mute
		registedBlocks := []entity.Block{}
		if err := blockDBHandler.FindAllByUserID(&registedBlocks, userID); err != nil {
			b.loggerHandler.Errorw("blockDBHandler.FindAllByUserID() error.", "user_id", userID, "error", err)
			return err
		}
		registedBlockEntities := entity.Blocks{}
		registedBlockEntities.ToBlocksEntity(&registedBlocks)
		registedBlockEntities.SortByTargetTwtitterID()

		b.loggerHandler.Debugf("user_id=%s Num_Of_blocks=%d", userID, len(registedBlockEntities))

		registedMutes := []entity.Mute{}
		if err := muteDBHandler.FindAllByUserID(&registedMutes, userID); err != nil {
			b.loggerHandler.Errorw("muteDBHandler.FindAllByUserID() error.", "user_id", userID, "error", err)
			return err
		}
		registedMuteEntities := entity.Mutes{}
		registedMuteEntities.ToMutesEntity(&registedMutes)
		registedMuteEntities.SortByTargetTwtitterID()

		b.loggerHandler.Debugf("user_id=%s Num_Of_mutes=%d", userID, len(registedMuteEntities))

		// MuteでFlag=1のものはスキップし、Blockは変換できたものだけFlag=1をたてる。それ以外は処理しない。
		convertedBlockEntities := []entity.Block{}
		muteEntities := []entity.Mute{}
		for _, registedBlockEntity := range registedBlockEntities {

			// mute除外
			needleTwitterID := registedBlockEntity.GetTargetTwitterID()
			if registedMuteEntities.IsConvertedByTwitterID(needleTwitterID) {
				b.loggerHandler.Debugf("skip mute flag=1 user_id=%s target_twitter_id=%s", userID, needleTwitterID)
				continue
			}

			// block除外
			if registedBlockEntity.IsConverted() {
				b.loggerHandler.Debugf("skip block flag=1 user_id=%s target_twitter_id=%s", userID, needleTwitterID)
				continue
			}

			// 移行処理
			// NOTE: 既にブロックを解除している場合はエラーを返さないようにする。
			if err := b.twitterHandler.DestroyBlock(user.GetTwitterID(), registedBlockEntity.GetTargetTwitterID()); err != nil {
				b.loggerHandler.Warnw(
					"destroy block error.",
					"TwitterID",
					user.GetTwitterID(),
					"TargetTwitterID",
					registedBlockEntity.GetTargetTwitterID(),
					"error",
					err,
				)
				continue
			}

			// NOTE: 既にミュートにしている場合はエラーを返さないようにする。
			if err := b.twitterHandler.CreateMute(user.GetTwitterID(), registedBlockEntity.GetTargetTwitterID()); err != nil {
				b.loggerHandler.Warnw(
					"create mute error.",
					"TwitterID",
					user.GetTwitterID(),
					"TargetTwitterID",
					registedBlockEntity.GetTargetTwitterID(),
					"error",
					err,
				)
				continue
			}
			registedBlockEntity.Converted()
			convertedBlockEntities = append(convertedBlockEntities, registedBlockEntity)
			nowTime := time.Now()
			mute := entity.NewMute(
				registedBlockEntity.GetUserID(),
				registedBlockEntity.GetTargetTwitterID(),
				1,
				nowTime,
				nowTime,
			)
			muteEntities = append(muteEntities, *mute)
		}

		// 移行完了処理 blocks更新とmute更新
		if len(convertedBlockEntities) > 0 {
			if err := blockDBHandler.CreateNewBlocks(&convertedBlockEntities, "user_id", "twitter_id"); err != nil {
				b.loggerHandler.Errorw(
					"fail to create new blocks.",
					"user_id",
					userID,
					"twitter_id",
					user.GetTwitterID(),
					"error",
					err,
				)
				return err
			}
		}
		if len(muteEntities) > 0 {
			if err := muteDBHandler.CreateNew(&muteEntities, "user_id", "twitter_id"); err != nil {
				b.loggerHandler.Errorw(
					"fail to create new mutes.",
					"user_id",
					userID,
					"twitter_id",
					user.GetTwitterID(),
					"error",
					err,
				)
				return err
			}

			// 移行完了数を数える。
			for _, v := range muteEntities {
				if v.IsConverted() {
					convertedIDs = append(convertedIDs, v.GetTargetTwitterID())
				}
			}
		}

		return nil
	})

	if err != nil {
		block2Mute := entity.Block2Mute{}
		b.loggerHandler.Errorw("transaction fails", "error", err)
		return &block2Mute, err
	}

	block2Mute := entity.NewBlock2Mute(
		uint(len(convertedIDs)),
		convertedIDs,
	)

	return block2Mute, nil
}
