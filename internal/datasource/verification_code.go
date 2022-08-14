package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type VerificationCodeDatasource interface {
	GetVerificationCode(tx *gorm.DB, where *entity.VerificationCode) (*entity.VerificationCode, error)
	CreateVerificationCode(tx *gorm.DB, data *entity.VerificationCode) (*entity.VerificationCode, error)
	DeleteVerificationCode(tx *gorm.DB, where *entity.VerificationCode, ids *[]uuid.UUID) (*[]entity.VerificationCode, error)
}

type verificationCodeDatasource struct{}

func (v *verificationCodeDatasource) CreateVerificationCode(tx *gorm.DB, data *entity.VerificationCode) (*entity.VerificationCode, error) {
	result := tx.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

func (v *verificationCodeDatasource) GetVerificationCode(tx *gorm.DB, verificationCode *entity.VerificationCode) (*entity.VerificationCode, error) {
	var res *entity.VerificationCode
	result := tx.Where(verificationCode).Take(&res)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return res, nil
}

func (v *verificationCodeDatasource) DeleteVerificationCode(tx *gorm.DB, where *entity.VerificationCode, ids *[]uuid.UUID) (*[]entity.VerificationCode, error) {
	var res *[]entity.VerificationCode
	var result *gorm.DB
	if ids != nil {
		result = tx.Clauses(clause.Returning{}).Where(`id IN ?`, ids).Delete(&res)
	} else {
		result = tx.Clauses(clause.Returning{}).Where(where).Delete(&res)
	}
	if result.Error != nil {
		return nil, result.Error
	} else if result.RowsAffected == 0 {
		return nil, errors.New("record not found")
	}
	return res, nil
}
