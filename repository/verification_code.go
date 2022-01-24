package repository

import (
	"errors"

	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type VerificationCodeQuery interface {
	GetVerificationCode(tx *gorm.DB, verificationCode *models.VerificationCode, fields *[]string) (*models.VerificationCode, error)
	CreateVerificationCode(tx *gorm.DB, verificationCode *models.VerificationCode) error
	DeleteVerificationCode(tx *gorm.DB, verificationCode *models.VerificationCode) error
}

type verificationCodeQuery struct{}

func (v *verificationCodeQuery) CreateVerificationCode(tx *gorm.DB, verificationCode *models.VerificationCode) error {
	result := tx.Create(&verificationCode)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (v *verificationCodeQuery) GetVerificationCode(tx *gorm.DB, verificationCode *models.VerificationCode, fields *[]string) (*models.VerificationCode, error) {
	var verificationCodeResult *models.VerificationCode
	result := tx.Where(verificationCode).Select(*fields).Take(&verificationCodeResult)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return verificationCodeResult, nil
}

func (v *verificationCodeQuery) DeleteVerificationCode(tx *gorm.DB, verificationCode *models.VerificationCode) error {
	var verificationCodeResult *[]models.VerificationCode
	result := tx.Where(verificationCode).Delete(&verificationCodeResult)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
