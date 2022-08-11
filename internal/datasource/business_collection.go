package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/internal/entity"
	"gorm.io/gorm"
)

type BusinessCollectionDatasource interface {
	GetBusinessCollection(tx *gorm.DB, where *entity.BusinessCollection) (*entity.BusinessCollection, error)
	ListBusinessCollection(tx *gorm.DB, where *entity.BusinessCollection) (*[]entity.BusinessCollection, error)
}

type businessCollectionDatasource struct{}

func (i *businessCollectionDatasource) ListBusinessCollection(tx *gorm.DB, where *entity.BusinessCollection) (*[]entity.BusinessCollection, error) {
	var itemsCategory []entity.BusinessCollection
	result := tx.Where(where).Find(&itemsCategory)
	if result.Error != nil {
		return nil, result.Error
	}
	return &itemsCategory, nil
}

func (i *businessCollectionDatasource) GetBusinessCollection(tx *gorm.DB, where *entity.BusinessCollection) (*entity.BusinessCollection, error) {
	var businessBusinessCollection *entity.BusinessCollection
	result := tx.Where(where).Take(&businessBusinessCollection)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return businessBusinessCollection, nil
}
