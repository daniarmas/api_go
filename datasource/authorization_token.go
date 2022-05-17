package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AuthorizationTokenDatasource interface {
	GetAuthorizationToken(tx *gorm.DB, where *models.AuthorizationToken, fields *[]string) (*models.AuthorizationToken, error)
	CreateAuthorizationToken(tx *gorm.DB, data *models.AuthorizationToken) (*models.AuthorizationToken, error)
	DeleteAuthorizationToken(tx *gorm.DB, where *models.AuthorizationToken, ids *[]uuid.UUID) (*[]models.AuthorizationToken, error)
	DeleteAuthorizationTokenByRefreshTokenIds(tx *gorm.DB, ids *[]uuid.UUID) (*[]models.AuthorizationToken, error)
}

type authorizationTokenDatasource struct{}

func (v *authorizationTokenDatasource) CreateAuthorizationToken(tx *gorm.DB, data *models.AuthorizationToken) (*models.AuthorizationToken, error) {
	result := tx.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

func (r *authorizationTokenDatasource) DeleteAuthorizationToken(tx *gorm.DB, where *models.AuthorizationToken, ids *[]uuid.UUID) (*[]models.AuthorizationToken, error) {
	var res *[]models.AuthorizationToken
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

func (r *authorizationTokenDatasource) DeleteAuthorizationTokenByRefreshTokenIds(tx *gorm.DB, ids *[]uuid.UUID) (*[]models.AuthorizationToken, error) {
	var res *[]models.AuthorizationToken
	result := tx.Clauses(clause.Returning{}).Where(`refresh_token_id IN ?`, ids).Delete(&res)
	if result.Error != nil {
		return nil, result.Error
	} else if result.RowsAffected == 0 {
		return nil, errors.New("record not found")
	}
	return res, nil
}

func (v *authorizationTokenDatasource) GetAuthorizationToken(tx *gorm.DB, where *models.AuthorizationToken, fields *[]string) (*models.AuthorizationToken, error) {
	var res *models.AuthorizationToken
	result := tx.Where(where).Select(fields).Take(&res)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return res, nil
}
