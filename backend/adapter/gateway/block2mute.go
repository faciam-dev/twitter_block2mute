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
	blockDbHandler handler.BlockDbHandler
	userDbHandler  handler.UserDbHandler
	muteDbHandler  handler.MuteDbHandler
	twitterHandler handler.TwitterHandler
	sessionHandler handler.SessionHandler
}

// NewBlockRepository はBlockRepositoryを返します．
func NewBlock2MuteRepository(
	dbHandler handler.BlockDbHandler,
	userDbHandler handler.UserDbHandler,
	muteDbHandler handler.MuteDbHandler,
	twitterHandler handler.TwitterHandler,
	sessionHandler handler.SessionHandler) port.Block2MuteRepository {
	return &Block2MuteRepository{
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
		return &block2Mute, err
	}

	// auth Twitter
	token := b.sessionHandler.Get("token")
	secret := b.sessionHandler.Get("secret")

	if token == nil || secret == nil {
		return &block2Mute, errors.New("session timeout or not found")
	}

	b.twitterHandler.UpdateTwitterApi(token.(string), secret.(string))

	// update blocks table
	convertedIDs := []string{}
	err := b.muteDbHandler.Transaction(func() error {
		// block2mute
		registedBlockEntities := []entity.Block{}
		if err := b.blockDbHandler.FindAllByUserID(&registedBlockEntities, userID); err != nil {
			log.Print(err)
			return err
		}
		sort.Slice(registedBlockEntities, func(i, j int) bool {
			return registedBlockEntities[i].TargetTwitterID <= registedBlockEntities[j].TargetTwitterID
		})

		registedMuteEntities := []entity.Mute{}
		if err := b.muteDbHandler.FindAllByUserID(&registedMuteEntities, userID); err != nil {
			log.Print(err)
			return err
		}
		sort.Slice(registedMuteEntities, func(i, j int) bool {
			return registedMuteEntities[i].TargetTwitterID <= registedMuteEntities[j].TargetTwitterID
		})

		// MuteでFlag=1のものはスキップし、Blockは変換できたものだけFlag=1をたてる。それ以外は処理しない。
		convertedBlockEntities := []entity.Block{}
		muteEntities := []entity.Mute{}
		for _, registedBlockEntity := range registedBlockEntities {

			// mute除外
			needleTwitterID := registedBlockEntity.TargetTwitterID
			idx := sort.Search(len(registedMuteEntities), func(i int) bool {
				return string(registedMuteEntities[i].TargetTwitterID) == needleTwitterID
			})

			if len(registedMuteEntities) > idx && registedMuteEntities[idx].Flag == 1 {
				continue
			}

			// block除外
			if registedBlockEntity.Flag == 1 {
				continue
			}

			// 移行処理
			// NOTE: 既にブロックを解除している場合はエラーを返さないようにする。
			if err := b.twitterHandler.DestroyBlock(user.TwitterID, registedBlockEntity.TargetTwitterID); err != nil {
				log.Printf("destroy block :%v %v", user.TwitterID, registedBlockEntity.TargetTwitterID)
				log.Print(err)
				continue
			}

			// NOTE: 既にミュートにしている場合はエラーを返さないようにする。
			if err := b.twitterHandler.CreateMute(user.TwitterID, registedBlockEntity.TargetTwitterID); err != nil {
				log.Printf("create mute :%v %v", user.TwitterID, registedBlockEntity.TargetTwitterID)
				log.Print(err)
				continue
			}
			registedBlockEntity.Flag = 1
			convertedBlockEntities = append(convertedBlockEntities, registedBlockEntity)
			mute := entity.Mute{
				TargetTwitterID: registedBlockEntity.TargetTwitterID,
				UserID:          registedBlockEntity.UserID,
				Flag:            1,
			}
			muteEntities = append(muteEntities, mute)
		}

		// 移行完了処理 blocks更新とmute更新
		if len(convertedBlockEntities) > 0 {
			if err := b.blockDbHandler.CreateNewBlocks(&convertedBlockEntities, "user_id", "twitter_id"); err != nil {
				log.Print(err)
				return err
			}
		}
		if len(muteEntities) > 0 {
			if err := b.muteDbHandler.CreateNew(&muteEntities, "user_id", "twitter_id"); err != nil {
				log.Print(err)
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
