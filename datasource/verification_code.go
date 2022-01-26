package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/models"
	"gorm.io/gorm"
)

type VerificationCodeDatasource interface {
	GetVerificationCode(tx *gorm.DB, verificationCode *models.VerificationCode, fields *[]string) (*models.VerificationCode, error)
	CreateVerificationCode(tx *gorm.DB, verificationCode *models.VerificationCode) error
	DeleteVerificationCode(tx *gorm.DB, verificationCode *models.VerificationCode) error
}

type verificationCodeDatasource struct{}

func (v *verificationCodeDatasource) CreateVerificationCode(tx *gorm.DB, verificationCode *models.VerificationCode) error {
	result := tx.Create(&verificationCode)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (v *verificationCodeDatasource) GetVerificationCode(tx *gorm.DB, verificationCode *models.VerificationCode, fields *[]string) (*models.VerificationCode, error) {
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

func (v *verificationCodeDatasource) DeleteVerificationCode(tx *gorm.DB, verificationCode *models.VerificationCode) error {
	var verificationCodeResult *[]models.VerificationCode
	result := tx.Where(verificationCode).Delete(&verificationCodeResult)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
