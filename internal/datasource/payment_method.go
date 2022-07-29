package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PaymentMethodDatasource interface {
	ListPaymentMethod(tx *gorm.DB, where *entity.PaymentMethod) (*[]entity.PaymentMethod, error)
	CreatePaymentMethod(tx *gorm.DB, data *entity.PaymentMethod) (*entity.PaymentMethod, error)
	UpdatePaymentMethod(tx *gorm.DB, where *entity.PaymentMethod, data *entity.PaymentMethod) (*entity.PaymentMethod, error)
	GetPaymentMethod(tx *gorm.DB, where *entity.PaymentMethod) (*entity.PaymentMethod, error)
	DeletePaymentMethod(tx *gorm.DB, where *entity.PaymentMethod, ids *[]uuid.UUID) (*[]entity.PaymentMethod, error)
}

type paymentMethodDatasource struct{}

func (i *paymentMethodDatasource) ListPaymentMethod(tx *gorm.DB, where *entity.PaymentMethod) (*[]entity.PaymentMethod, error) {
	var res []entity.PaymentMethod
	result := tx.Where(where).Find(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return &res, nil
}

func (i *paymentMethodDatasource) CreatePaymentMethod(tx *gorm.DB, data *entity.PaymentMethod) (*entity.PaymentMethod, error) {
	result := tx.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

func (i *paymentMethodDatasource) UpdatePaymentMethod(tx *gorm.DB, where *entity.PaymentMethod, data *entity.PaymentMethod) (*entity.PaymentMethod, error) {
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

func (i *paymentMethodDatasource) GetPaymentMethod(tx *gorm.DB, where *entity.PaymentMethod) (*entity.PaymentMethod, error) {
	var res *entity.PaymentMethod
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

func (v *paymentMethodDatasource) DeletePaymentMethod(tx *gorm.DB, where *entity.PaymentMethod, ids *[]uuid.UUID) (*[]entity.PaymentMethod, error) {
	var res *[]entity.PaymentMethod
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
