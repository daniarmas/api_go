package repository

import (
	"github.com/daniarmas/api_go/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentMethodRepository interface {
	ListPaymentMethod(tx *gorm.DB, where *entity.PaymentMethod) (*[]entity.PaymentMethod, error)
	CreatePaymentMethod(tx *gorm.DB, data *entity.PaymentMethod) (*entity.PaymentMethod, error)
	UpdatePaymentMethod(tx *gorm.DB, where *entity.PaymentMethod, data *entity.PaymentMethod) (*entity.PaymentMethod, error)
	GetPaymentMethod(tx *gorm.DB, where *entity.PaymentMethod) (*entity.PaymentMethod, error)
	DeletePaymentMethod(tx *gorm.DB, where *entity.PaymentMethod, ids *[]uuid.UUID) (*[]entity.PaymentMethod, error)
}

type paymentMethodRepository struct{}

func (i *paymentMethodRepository) ListPaymentMethod(tx *gorm.DB, where *entity.PaymentMethod) (*[]entity.PaymentMethod, error) {
	result, err := Datasource.NewPaymentMethodDatasource().ListPaymentMethod(tx, where)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *paymentMethodRepository) DeletePaymentMethod(tx *gorm.DB, where *entity.PaymentMethod, ids *[]uuid.UUID) (*[]entity.PaymentMethod, error) {
	res, err := Datasource.NewPaymentMethodDatasource().DeletePaymentMethod(tx, where, ids)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *paymentMethodRepository) GetPaymentMethod(tx *gorm.DB, where *entity.PaymentMethod) (*entity.PaymentMethod, error) {
	res, err := Datasource.NewPaymentMethodDatasource().GetPaymentMethod(tx, where)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *paymentMethodRepository) CreatePaymentMethod(tx *gorm.DB, data *entity.PaymentMethod) (*entity.PaymentMethod, error) {
	res, err := Datasource.NewPaymentMethodDatasource().CreatePaymentMethod(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *paymentMethodRepository) UpdatePaymentMethod(tx *gorm.DB, where *entity.PaymentMethod, data *entity.PaymentMethod) (*entity.PaymentMethod, error) {
	result, err := Datasource.NewPaymentMethodDatasource().UpdatePaymentMethod(tx, where, data)
	if err != nil {
		return nil, err
	}
	return result, nil
}
