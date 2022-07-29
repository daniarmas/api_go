package repository

import (
	"github.com/daniarmas/api_go/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BusinessPaymentMethodRepository interface {
	ListBusinessPaymentMethodWithEnabled(tx *gorm.DB, where *entity.BusinessPaymentMethod) (*[]entity.BusinessPaymentMethodWithEnabled, error)
	ListBusinessPaymentMethod(tx *gorm.DB, where *entity.BusinessPaymentMethod) (*[]entity.BusinessPaymentMethod, error)
	CreateBusinessPaymentMethod(tx *gorm.DB, data *entity.BusinessPaymentMethod) (*entity.BusinessPaymentMethod, error)
	UpdateBusinessPaymentMethod(tx *gorm.DB, where *entity.BusinessPaymentMethod, data *entity.BusinessPaymentMethod) (*entity.BusinessPaymentMethod, error)
	GetBusinessPaymentMethod(tx *gorm.DB, where *entity.BusinessPaymentMethod) (*entity.BusinessPaymentMethod, error)
	DeleteBusinessPaymentMethod(tx *gorm.DB, where *entity.BusinessPaymentMethod, ids *[]uuid.UUID) (*[]entity.BusinessPaymentMethod, error)
}

type businessPaymentMethodRepository struct{}

func (i *businessPaymentMethodRepository) ListBusinessPaymentMethodWithEnabled(tx *gorm.DB, where *entity.BusinessPaymentMethod) (*[]entity.BusinessPaymentMethodWithEnabled, error) {
	result, err := Datasource.NewBusinessPaymentMethodDatasource().ListBusinessPaymentMethodWithEnabled(tx, where)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *businessPaymentMethodRepository) ListBusinessPaymentMethod(tx *gorm.DB, where *entity.BusinessPaymentMethod) (*[]entity.BusinessPaymentMethod, error) {
	result, err := Datasource.NewBusinessPaymentMethodDatasource().ListBusinessPaymentMethod(tx, where)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *businessPaymentMethodRepository) DeleteBusinessPaymentMethod(tx *gorm.DB, where *entity.BusinessPaymentMethod, ids *[]uuid.UUID) (*[]entity.BusinessPaymentMethod, error) {
	res, err := Datasource.NewBusinessPaymentMethodDatasource().DeleteBusinessPaymentMethod(tx, where, ids)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *businessPaymentMethodRepository) GetBusinessPaymentMethod(tx *gorm.DB, where *entity.BusinessPaymentMethod) (*entity.BusinessPaymentMethod, error) {
	res, err := Datasource.NewBusinessPaymentMethodDatasource().GetBusinessPaymentMethod(tx, where)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *businessPaymentMethodRepository) CreateBusinessPaymentMethod(tx *gorm.DB, data *entity.BusinessPaymentMethod) (*entity.BusinessPaymentMethod, error) {
	res, err := Datasource.NewBusinessPaymentMethodDatasource().CreateBusinessPaymentMethod(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *businessPaymentMethodRepository) UpdateBusinessPaymentMethod(tx *gorm.DB, where *entity.BusinessPaymentMethod, data *entity.BusinessPaymentMethod) (*entity.BusinessPaymentMethod, error) {
	result, err := Datasource.NewBusinessPaymentMethodDatasource().UpdateBusinessPaymentMethod(tx, where, data)
	if err != nil {
		return nil, err
	}
	return result, nil
}
