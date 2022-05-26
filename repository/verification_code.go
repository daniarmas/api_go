package repository

import (
	"github.com/daniarmas/api_go/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VerificationCodeQuery interface {
	GetVerificationCode(tx *gorm.DB, where *models.VerificationCode, fields *[]string) (*models.VerificationCode, error)
	CreateVerificationCode(tx *gorm.DB, data *models.VerificationCode) (*models.VerificationCode, error)
	DeleteVerificationCode(tx *gorm.DB, where *models.VerificationCode, ids *[]uuid.UUID) (*[]models.VerificationCode, error)
}

type verificationCodeQuery struct{}

func (v *municipalityRepository) CreateVerificationCode(tx *gorm.DB, data *models.VerificationCode) (*models.VerificationCode, error) {
	res, err := Datasource.NewVerificationCodeDatasource().CreateVerificationCode(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (v *municipalityRepository) GetVerificationCode(tx *gorm.DB, where *models.VerificationCode, fields *[]string) (*models.VerificationCode, error) {
	res, err := Datasource.NewVerificationCodeDatasource().GetVerificationCode(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (v *municipalityRepository) DeleteVerificationCode(tx *gorm.DB, where *models.VerificationCode, ids *[]uuid.UUID) (*[]models.VerificationCode, error) {
	res, err := Datasource.NewVerificationCodeDatasource().DeleteVerificationCode(tx, where, ids)
	if err != nil {
		return nil, err
	}
	return res, nil
}
