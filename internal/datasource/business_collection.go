package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/internal/entity"
	"gorm.io/gorm"
)

type BusinessCollectionDatasource interface {
	GetBusinessCollection(tx *gorm.DB, where *entity.BusinessCollection, fields *[]string) (*entity.BusinessCollection, error)
	ListBusinessCollection(tx *gorm.DB, where *entity.BusinessCollection, fields *[]string) (*[]entity.BusinessCollection, error)
}

type businessCollectionDatasource struct{}

func (i *businessCollectionDatasource) ListBusinessCollection(tx *gorm.DB, where *entity.BusinessCollection, fields *[]string) (*[]entity.BusinessCollection, error) {
	var itemsCategory []entity.BusinessCollection
	selectFields := &[]string{"*"}
	if fields != nil {
		selectFields = fields
	}
	result := tx.Where(where).Select(*selectFields).Find(&itemsCategory)
	if result.Error != nil {
		return nil, result.Error
	}
	return &itemsCategory, nil
}

func (i *businessCollectionDatasource) GetBusinessCollection(tx *gorm.DB, where *entity.BusinessCollection, fields *[]string) (*entity.BusinessCollection, error) {
	var businessBusinessCollection *entity.BusinessCollection
	selectFields := &[]string{"*"}
	if fields != nil {
		selectFields = fields
	}
	result := tx.Where(where).Select(*selectFields).Take(&businessBusinessCollection)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return businessBusinessCollection, nil
}
