package database

import (
	"errors"
	"log"

	"github.com/faciam_dev/twitter_block2mute/backend/database/gorm/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormDbEntityHandler[E any, M any] struct {
	//GormCommonDbHandler
	db *gorm.DB
	model.ModelForDomain[E, M]
}

func NewDbEntityHandler[E any, M any](gormDbHandler GormDbHandler) GormDbEntityHandler[E, M] {
	gormDbEntityHandler := new(GormDbEntityHandler[E, M])
	gormDbEntityHandler.db = gormDbHandler.Conn

	return *gormDbEntityHandler
}

// トランザクションfunc
func (g *GormDbEntityHandler[E, M]) Transaction(fn func() error) error {
	return g.db.Transaction(func(tx *gorm.DB) error {
		return fn()
	})
}

// トランザクション開始
func (g *GormDbEntityHandler[E, M]) Begin() {
	g.db = g.db.Begin()
}

// トランザクションコミット
func (g *GormDbEntityHandler[E, M]) Commit() {
	g.db = g.db.Commit()
}

// トランザクションロールバック
func (g *GormDbEntityHandler[E, M]) Rollback() {
	g.db = g.db.Rollback()
}

// レコード1件をIDで取得
func (g *GormDbEntityHandler[E, M]) First(domainEntityInterface interface{}, ID string) error {

	domainEntity, err := g.InterfaceToEntity(domainEntityInterface)
	if err != nil {
		return err
	}

	model := g.ModelForDomain.FromDomain(domainEntity)
	//model := modelForDomain.FromDomain(domainEntity)
	/*
		model, err := g.entityToModel(modelForDomain, domainEntity)
		if err != nil {
			log.Print(err)
			return err
		}
	*/

	if err := g.db.First(&model, ID).Error; err != nil {
		log.Print(err)
		return err
	}
	/*
		if err = g.GormCommonDbHandler.First(model, ID); err != nil {
			return err
		}
	*/

	g.ModelForDomain.ToDomain(model, domainEntity)

	return nil
}

// レコードを作成
func (g *GormDbEntityHandler[E, M]) Create(domainEntityInterface interface{}) error {
	domainEntity, err := g.InterfaceToEntity(domainEntityInterface)
	if err != nil {
		return err
	}

	model := g.ModelForDomain.FromDomain(domainEntity)

	if err := g.db.Create(&model).Error; err != nil {
		return err
	}

	g.ModelForDomain.ToDomain(model, domainEntity)

	return nil
}

// レコードをID基準で更新
func (g *GormDbEntityHandler[E, M]) Update(domainEntityInterface interface{}, ID string) error {
	domainEntity, err := g.InterfaceToEntity(domainEntityInterface)
	if err != nil {
		return err
	}

	model := g.ModelForDomain.FromDomain(domainEntity)

	if err := g.db.Model(&model).Where("id = ?", ID).Updates(model).Error; err != nil {
		return err
	}

	g.ModelForDomain.ToDomain(model, domainEntity)

	return nil
}

// レコード1件を検索。※columnNameにユーザーからの入力値を絶対に使わないこと。
func (g *GormDbEntityHandler[E, M]) Find(domainEntityInterface interface{}, columnName string, searchValue string) error {
	domainEntity, err := g.InterfaceToEntity(domainEntityInterface)
	if err != nil {
		return err
	}

	model := g.ModelForDomain.FromDomain(domainEntity)
	//model := modelForDomain.FromDomain(domainEntity)
	/*
		model, err := g.entityToModel(modelForDomain, domainEntity)
		if err != nil {
			log.Print(err)
			return err
		}
	*/

	//log.Print(model)

	if err := g.db.Where(columnName+" = ?", searchValue).Find(&model).Limit(1).Error; err != nil {
		log.Print(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	/*
		if err = g.GormCommonDbHandler.First(model, ID); err != nil {
			return err
		}
	*/

	g.ModelForDomain.ToDomain(model, domainEntity)

	return nil

}

// レコード全件を検索。※columnNameにユーザーからの入力値を絶対に使わないこと。
func (g *GormDbEntityHandler[E, M]) FindAll(domainEntitiesInterface interface{}, columnName string, searchValue string) error {

	entities, err := g.InterfaceToEntities(domainEntitiesInterface)
	if err != nil {
		return err
	}

	models := []M{}
	for _, entity := range *entities {
		model := g.ModelForDomain.FromDomain(&entity)
		models = append(models, model)
	}

	/*
		models, err := g.entitiesToModels(domainEntitiesInterface)
		if err != nil {
			return err
		}
	*/

	if err := g.db.Where(columnName+" = ?", searchValue).Find(&models).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	g.ModelForDomain.ToDomains(models, entities)

	/*
		resultEntities := []E{}
		//log.Print(models)
		for _, model := range models {
			entity := new(E)
			g.ModelForDomain.ToDomain(model, entity)
			//entities = append(entities, *entity)
			resultEntities = append(resultEntities, *entity)
		}

		entities = &resultEntities

		domainEntitiesInterface = entities

		log.Print(domainEntitiesInterface)
	*/

	return nil
}

// 条件に該当するレコードを新規作成または更新する
func (g *GormDbEntityHandler[E, M]) Upsert(domainEntitiesInterface interface{}, columnName string, searchValue string) error {
	domainEntity, err := g.InterfaceToEntity(domainEntitiesInterface)

	if err != nil {
		return err
	}
	model := g.ModelForDomain.FromDomain(domainEntity)

	/*
		models := []M{}
		for _, entity := range *entities {
			model := g.ModelForDomain.FromDomain(&entity)
			models = append(models, model)
		}
	*/

	err = g.db.Clauses(clause.OnConflict{
		Where:     clause.Where{Exprs: []clause.Expression{clause.Eq{Column: columnName, Value: searchValue}}},
		UpdateAll: true,
	}).Create(&model).Error

	if err != nil {
		return err
	}

	g.ModelForDomain.ToDomain(model, domainEntity)

	return nil
}

// entityをmodelにする。
/*
func (g *GormDbEntityHandler[T, E, MD, M]) entityToModel(modelForDomain MD, domainEntity interface{}) (M, error) {
	switch casted := domainEntity.(type) {
	case *E:
		model := modelForDomain.FromDomain(casted)
		return model, nil
	default:
		model := modelForDomain.Blank()
		return model, errors.New("argment type is not entity")
	}
}
*/

// プライマリキーに対応するモデルのレコードの削除処理
func (g *GormDbEntityHandler[E, M]) Delete(domainEntityInterface interface{}, searchValue string) error {

	domainEntity, err := g.InterfaceToEntity(domainEntityInterface)
	if err != nil {
		return err
	}

	model := g.ModelForDomain.FromDomain(domainEntity)

	if err := g.db.Delete(&model, searchValue).Error; err != nil {
		return err
	}

	return nil
}

// 複数プライマリキーに対応するモデルのレコードの削除処理
func (g *GormDbEntityHandler[E, M]) DeleteByIds(domainEntitiesInterface interface{}, IDs []uint) error {

	entities, err := g.InterfaceToEntities(domainEntitiesInterface)
	if err != nil {
		return err
	}

	models := []M{}
	for _, entity := range *entities {
		model := g.ModelForDomain.FromDomain(&entity)
		models = append(models, model)
	}

	if err := g.db.Delete(&models, IDs).Error; err != nil {
		return err
	}

	return nil
}

// エンティティに対応するモデルの全レコード削除処理
func (g *GormDbEntityHandler[E, M]) DeleteAll(domainEntityInterface interface{}, columnName string, searchValue string) error {

	domainEntity, err := g.InterfaceToEntity(domainEntityInterface)
	if err != nil {
		return err
	}

	model := g.ModelForDomain.FromDomain(domainEntity)

	if err := g.db.Where(columnName+" = ?", searchValue).Delete(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	return nil
}

// エンティティ変換：interfaceをentityにする
func (g *GormDbEntityHandler[E, M]) InterfaceToEntity(interf interface{}) (*E, error) {
	switch casted := interf.(type) {
	case *E:
		return casted, nil
	default:
		return new(E), errors.New("interface is not entity.*")
	}
}

// エンティティスライス変換：interfaceを[]entityにする。
func (g *GormDbEntityHandler[E, M]) InterfaceToEntities(interfaceSlice interface{}) (*[]E, error) {
	switch casted := interfaceSlice.(type) {
	case *[]E:
		return casted, nil
	default:
		return &[]E{}, errors.New("interface is not []entity.*")
	}
}

// モデルスライス変換：interface([]entity)を[]modelにする。
func (g *GormDbEntityHandler[E, M]) EntitiesToModels(entities []E) []M {
	models := []M{}

	for _, entity := range entities {
		model := g.ModelForDomain.FromDomain(&entity)
		models = append(models, model)
	}

	return models
}
