package repository

import (
	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type VerificationCodeQuery interface {
	GetVerificationCode(tx *gorm.DB, where *models.VerificationCode, fields *[]string) (*models.VerificationCode, error)
	CreateVerificationCode(tx *gorm.DB, data *models.VerificationCode) error
	DeleteVerificationCode(tx *gorm.DB, where *models.VerificationCode) error
}

type verificationCodeQuery struct{}

func (v *verificationCodeQuery) CreateVerificationCode(tx *gorm.DB, data *models.VerificationCode) error {
	err := Datasource.NewVerificationCodeDatasource().CreateVerificationCode(tx, data)
	if err != nil {
		return err
	}
	return nil
}

func (v *verificationCodeQuery) GetVerificationCode(tx *gorm.DB, where *models.VerificationCode, fields *[]string) (*models.VerificationCode, error) {
	result, err := Datasource.NewVerificationCodeDatasource().GetVerificationCode(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *verificationCodeQuery) DeleteVerificationCode(tx *gorm.DB, where *models.VerificationCode) error {
	err := Datasource.NewVerificationCodeDatasource().DeleteVerificationCode(tx, where)
	if err != nil {
		return err
	}
	return nil
}
