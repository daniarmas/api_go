package datasource

import (
	"errors"

	"github.com/daniarmas/api_go/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AuthorizationTokenDatasource interface {
	GetAuthorizationToken(tx *gorm.DB, where *entity.AuthorizationToken) (*entity.AuthorizationToken, error)
	CreateAuthorizationToken(tx *gorm.DB, data *entity.AuthorizationToken) (*entity.AuthorizationToken, error)
	DeleteAuthorizationToken(tx *gorm.DB, where *entity.AuthorizationToken, ids *[]uuid.UUID) (*[]entity.AuthorizationToken, error)
	DeleteAuthorizationTokenByRefreshTokenIds(tx *gorm.DB, ids *[]uuid.UUID) (*[]entity.AuthorizationToken, error)
}

type authorizationTokenDatasource struct{}

func (v *authorizationTokenDatasource) CreateAuthorizationToken(tx *gorm.DB, data *entity.AuthorizationToken) (*entity.AuthorizationToken, error) {
	result := tx.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

func (r *authorizationTokenDatasource) DeleteAuthorizationToken(tx *gorm.DB, where *entity.AuthorizationToken, ids *[]uuid.UUID) (*[]entity.AuthorizationToken, error) {
	var res *[]entity.AuthorizationToken
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

func (r *authorizationTokenDatasource) DeleteAuthorizationTokenByRefreshTokenIds(tx *gorm.DB, ids *[]uuid.UUID) (*[]entity.AuthorizationToken, error) {
	var res *[]entity.AuthorizationToken
	result := tx.Clauses(clause.Returning{}).Where(`refresh_token_id IN ?`, *ids).Delete(&res)
	if result.Error != nil {
		return nil, result.Error
	} else if result.RowsAffected == 0 {
		return nil, errors.New("record not found")
	}
	return res, nil
}

func (v *authorizationTokenDatasource) GetAuthorizationToken(tx *gorm.DB, where *entity.AuthorizationToken) (*entity.AuthorizationToken, error) {
	var res *entity.AuthorizationToken
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
