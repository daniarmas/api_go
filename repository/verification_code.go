package repository

import (
	"github.com/daniarmas/api_go/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VerificationCodeRepository interface {
	GetVerificationCode(tx *gorm.DB, where *models.VerificationCode, fields *[]string) (*models.VerificationCode, error)
	CreateVerificationCode(tx *gorm.DB, data *models.VerificationCode) (*models.VerificationCode, error)
	DeleteVerificationCode(tx *gorm.DB, where *models.VerificationCode, ids *[]uuid.UUID) (*[]models.VerificationCode, error)
}

type verificationCodeRepository struct{}

func (v *verificationCodeRepository) CreateVerificationCode(tx *gorm.DB, data *models.VerificationCode) (*models.VerificationCode, error) {
	res, err := Datasource.NewVerificationCodeDatasource().CreateVerificationCode(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (v *verificationCodeRepository) GetVerificationCode(tx *gorm.DB, where *models.VerificationCode, fields *[]string) (*models.VerificationCode, error) {
	res, err := Datasource.NewVerificationCodeDatasource().GetVerificationCode(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (v *verificationCodeRepository) DeleteVerificationCode(tx *gorm.DB, where *models.VerificationCode, ids *[]uuid.UUID) (*[]models.VerificationCode, error) {
	res, err := Datasource.NewVerificationCodeDatasource().DeleteVerificationCode(tx, where, ids)
	if err != nil {
		return nil, err
	}
	return res, nil
}
