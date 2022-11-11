package gateway

import (
	"errors"
	"log"
	"sort"

	"github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	"github.com/faciam_dev/twitter_block2mute/backend/entity"
	"github.com/faciam_dev/twitter_block2mute/backend/usecase/port"
)

type Block2MuteRepository struct {
	loggerHandler  handler.LoggerHandler
	blockDbHandler handler.BlockDbHandler
	userDbHandler  handler.UserDbHandler
	muteDbHandler  handler.MuteDbHandler
	twitterHandler handler.TwitterHandler
	sessionHandler handler.SessionHandler
}

// NewBlockRepository はBlockRepositoryを返します．
func NewBlock2MuteRepository(
	loggerHandler handler.LoggerHandler,
	dbHandler handler.BlockDbHandler,
	userDbHandler handler.UserDbHandler,
	muteDbHandler handler.MuteDbHandler,
	twitterHandler handler.TwitterHandler,
	sessionHandler handler.SessionHandler,
) port.Block2MuteRepository {
	return &Block2MuteRepository{
		loggerHandler:  loggerHandler,
		blockDbHandler: dbHandler,
		userDbHandler:  userDbHandler,
		muteDbHandler:  muteDbHandler,
		twitterHandler: twitterHandler,
		sessionHandler: sessionHandler,
	}
}

func (b *Block2MuteRepository) All(userID string) (*entity.Block2Mute, error) {
	block2Mute := entity.Block2Mute{}
	user := entity.User{}

	if err := b.userDbHandler.First(&user, userID); err != nil {
		b.loggerHandler.Errorf("user not found (user_id=%s)", userID)
		return &block2Mute, err
	}

	// auth Twitter
	token := b.sessionHandler.Get("token")
	secret := b.sessionHandler.Get("secret")

	if token == nil || secret == nil {
		return &block2Mute, errors.New("session timeout or not found")
	}

	b.twitterHandler.UpdateTwitterApi(token.(string), secret.(string))
	b.loggerHandler.Debugf("update twitter api (user_id=%s token=%s secret=%s)", userID, token.(string), secret.(string))

	// update blocks table
	convertedIDs := []string{}
	err := b.muteDbHandler.Transaction(func() error {
		// block2mute
		registedBlockEntities := []entity.Block{}
		if err := b.blockDbHandler.FindAllByUserID(&registedBlockEntities, userID); err != nil {
			b.loggerHandler.Errorw("blockDbHandler.FindAllByUserID() error.", "user_id", userID, "error", err)
			return err
		}
		sort.Slice(registedBlockEntities, func(i, j int) bool {
			return registedBlockEntities[i].GetTargetTwitterID() <= registedBlockEntities[j].GetTargetTwitterID()
		})
		b.loggerHandler.Debugf("user_id=%s Num_Of_blocks=%d", userID, len(registedBlockEntities))

		registedMuteEntities := []entity.Mute{}
		if err := b.muteDbHandler.FindAllByUserID(&registedMuteEntities, userID); err != nil {
			b.loggerHandler.Errorw("muteDbHandler.FindAllByUserID() error.", "user_id", userID, "error", err)
			return err
		}
		sort.Slice(registedMuteEntities, func(i, j int) bool {
			return registedMuteEntities[i].TargetTwitterID <= registedMuteEntities[j].TargetTwitterID
		})
		b.loggerHandler.Debugf("user_id=%s Num_Of_mutes=%d", userID, len(registedMuteEntities))

		// MuteでFlag=1のものはスキップし、Blockは変換できたものだけFlag=1をたてる。それ以外は処理しない。
		convertedBlockEntities := []entity.Block{}
		muteEntities := []entity.Mute{}
		for _, registedBlockEntity := range registedBlockEntities {

			// mute除外
			needleTwitterID := registedBlockEntity.GetTargetTwitterID()
			idx := sort.Search(len(registedMuteEntities), func(i int) bool {
				return string(registedMuteEntities[i].TargetTwitterID) == needleTwitterID
			})

			if len(registedMuteEntities) > idx && registedMuteEntities[idx].Flag == 1 {
				b.loggerHandler.Debugf("skip mute flag=1 user_id=%s target_twitter_id=%s", userID, registedMuteEntities[idx].TargetTwitterID)
				continue
			}

			// block除外
			if registedBlockEntity.GetFlag() == 1 {
				b.loggerHandler.Debugf("skip block flag=1 user_id=%s target_twitter_id=%s", userID, registedMuteEntities[idx].TargetTwitterID)
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
			mute := entity.Mute{
				TargetTwitterID: registedBlockEntity.GetTargetTwitterID(),
				UserID:          registedBlockEntity.GetUserID(),
				Flag:            1,
			}
			muteEntities = append(muteEntities, mute)
		}

		// 移行完了処理 blocks更新とmute更新
		if len(convertedBlockEntities) > 0 {
			if err := b.blockDbHandler.CreateNewBlocks(&convertedBlockEntities, "user_id", "twitter_id"); err != nil {
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
			if err := b.muteDbHandler.CreateNew(&muteEntities, "user_id", "twitter_id"); err != nil {
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
				if v.Flag == 1 {
					convertedIDs = append(convertedIDs, v.TargetTwitterID)
				}

			}
		}

		return nil
	})

	if err != nil {
		log.Print(err)
		return &block2Mute, err
	}

	block2Mute.NumberOfSuccess = uint(len(convertedIDs))
	block2Mute.SuccessTwitterIDs = convertedIDs

	return &block2Mute, nil
}
