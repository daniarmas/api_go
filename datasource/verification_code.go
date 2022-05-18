package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VerificationCodeDatasource interface {
	GetVerificationCode(tx *gorm.DB, where *models.VerificationCode, fields *[]string) (*models.VerificationCode, error)
	CreateVerificationCode(tx *gorm.DB, data *models.VerificationCode) (*models.VerificationCode, error)
	DeleteVerificationCode(tx *gorm.DB, where *models.VerificationCode, ids *[]uuid.UUID) (*[]models.VerificationCode, error)
}

type verificationCodeDatasource struct{}

func (v *verificationCodeDatasource) CreateVerificationCode(tx *gorm.DB, data *models.VerificationCode) (*models.VerificationCode, error) {
	result := tx.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

func (v *verificationCodeDatasource) GetVerificationCode(tx *gorm.DB, verificationCode *models.VerificationCode, fields *[]string) (*models.VerificationCode, error) {
	var res *models.VerificationCode
	result := tx.Where(verificationCode).Select(*fields).Take(&res)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return res, nil
}

func (v *verificationCodeDatasource) DeleteVerificationCode(tx *gorm.DB, verificationCode *models.VerificationCode, ids *[]uuid.UUID) (*[]models.VerificationCode, error) {
	var res *[]models.VerificationCode
	result := tx.Where(verificationCode).Delete(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return res, nil
}
