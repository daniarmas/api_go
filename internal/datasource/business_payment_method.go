package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BusinessPaymentMethodDatasource interface {
	ListBusinessPaymentMethodWithEnabled(tx *gorm.DB, where *entity.BusinessPaymentMethod) (*[]entity.BusinessPaymentMethodWithEnabled, error)
	ListBusinessPaymentMethod(tx *gorm.DB, where *entity.BusinessPaymentMethod) (*[]entity.BusinessPaymentMethod, error)
	CreateBusinessPaymentMethod(tx *gorm.DB, data *entity.BusinessPaymentMethod) (*entity.BusinessPaymentMethod, error)
	UpdateBusinessPaymentMethod(tx *gorm.DB, where *entity.BusinessPaymentMethod, data *entity.BusinessPaymentMethod) (*entity.BusinessPaymentMethod, error)
	GetBusinessPaymentMethod(tx *gorm.DB, where *entity.BusinessPaymentMethod) (*entity.BusinessPaymentMethod, error)
	DeleteBusinessPaymentMethod(tx *gorm.DB, where *entity.BusinessPaymentMethod, ids *[]uuid.UUID) (*[]entity.BusinessPaymentMethod, error)
}

type businessPaymentMethodDatasource struct{}

func (i *businessPaymentMethodDatasource) ListBusinessPaymentMethod(tx *gorm.DB, where *entity.BusinessPaymentMethod) (*[]entity.BusinessPaymentMethod, error) {
	var res []entity.BusinessPaymentMethod
	result := tx.Where(where).Find(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return &res, nil
}

func (i *businessPaymentMethodDatasource) ListBusinessPaymentMethodWithEnabled(tx *gorm.DB, where *entity.BusinessPaymentMethod) (*[]entity.BusinessPaymentMethodWithEnabled, error) {
	var res []entity.BusinessPaymentMethodWithEnabled
	result := tx.Model(&entity.BusinessPaymentMethod{}).Select(`business_payment_method.id, business_payment_method.address, business_payment_method.type, business_payment_method.business_id, business_payment_method.payment_method_id, business_payment_method.create_time, business_payment_method.update_time, payment_method.enabled`).Where(where).Joins("left join payment_method on payment_method.id = business_payment_method.payment_method_id").Find(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return &res, nil
}

func (i *businessPaymentMethodDatasource) CreateBusinessPaymentMethod(tx *gorm.DB, data *entity.BusinessPaymentMethod) (*entity.BusinessPaymentMethod, error) {
	result := tx.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

func (i *businessPaymentMethodDatasource) UpdateBusinessPaymentMethod(tx *gorm.DB, where *entity.BusinessPaymentMethod, data *entity.BusinessPaymentMethod) (*entity.BusinessPaymentMethod, error) {
	result := tx.Clauses(clause.Returning{}).Where(where).Updates(&data)
	if result.RowsAffected == 0 {
		return nil, errors.New("record not found")
	} else if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return data, nil
}

func (i *businessPaymentMethodDatasource) GetBusinessPaymentMethod(tx *gorm.DB, where *entity.BusinessPaymentMethod) (*entity.BusinessPaymentMethod, error) {
	var res *entity.BusinessPaymentMethod
	result := tx.Where(where).Take(&res)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return res, nil
}

func (v *businessPaymentMethodDatasource) DeleteBusinessPaymentMethod(tx *gorm.DB, where *entity.BusinessPaymentMethod, ids *[]uuid.UUID) (*[]entity.BusinessPaymentMethod, error) {
	var res *[]entity.BusinessPaymentMethod
	var result *gorm.DB
	if ids != nil {
		result = tx.Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).Where(`id IN ?`, ids).Delete(&res)
	} else {
		result = tx.Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).Where(where).Delete(&res)
	}
	if result.Error != nil {
		return nil, result.Error
	} else if result.RowsAffected == 0 {
		return nil, errors.New("record not found")
	}
	return res, nil
}