package database

import (
	"errors"
	"log"
	"reflect"

	"github.com/faciam_dev/twitter_block2mute/backend/database/gorm/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormDbEntityHandler[E any, M any] struct {
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
		backUpDb := g.db
		g.db = tx
		err := fn()
		g.db = backUpDb
		return err
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

	if err := g.db.First(&model, ID).Error; err != nil {
		log.Print(err)
		return err
	}

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

	if err := g.db.Where(columnName+" = ?", searchValue).Find(&model).Limit(1).Error; err != nil {
		log.Print(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

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

	if err := g.db.Where(columnName+" = ?", searchValue).Find(&models).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	g.ModelForDomain.ToDomains(models, entities)

	return nil
}

// 条件に該当するレコードを新規作成または更新する
func (g *GormDbEntityHandler[E, M]) Upsert(domainEntitiesInterface interface{}, columnName string, searchValue string) error {
	domainEntity, err := g.InterfaceToEntity(domainEntitiesInterface)

	if err != nil {
		return err
	}
	model := g.ModelForDomain.FromDomain(domainEntity)

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
		return &[]E{}, errors.New("interface (" + reflect.TypeOf(interfaceSlice).String() + ") is not []entity.*" + reflect.TypeOf([]E{}).String())
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
