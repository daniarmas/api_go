package datasource

import (
	"errors"
	"fmt"
	"time"

	"github.com/daniarmas/api_go/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserAddressDatasource interface {
	ListUserAddress(tx *gorm.DB, where *models.UserAddress, fields *[]string) (*[]models.UserAddress, error)
}

type userAddressDatasource struct{}

func (i *userAddressDatasource) ListUserAddress(tx *gorm.DB, where *models.UserAddress, fields *[]string) (*[]models.UserAddress, error) {
	var res []models.UserAddress
	selectFields := &[]string{"id", "tag", "address", "number", "ST_AsEWKB(coordinates) AS coordinates", "instructions", "user_id", "province_id", "municipality_id", "create_time", "update_time"}
	if fields != nil {
		selectFields = fields
	}
	result := tx.Where(where).Select(*selectFields).Find(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return &res, nil
}

func (i *userAddressDatasource) CreateUserAddress(tx *gorm.DB, data *models.UserAddress) (*models.UserAddress, error) {
	point := fmt.Sprintf("POINT(%v %v)", data.Coordinates.Point.Coords()[1], data.Coordinates.Point.Coords()[0])
	var time = time.Now().UTC()
	var res models.UserAddress
	result := tx.Raw(`INSERT INTO "user_address" ("tag", "address", "number", "coordinates", "instructions", "user_id", "province_id", "municipality_id", "create_time", "update_time") VALUES (?, ?, ?, ST_GeomFromText(?, 4326), ?, ?, ?, ?, ?, ?) RETURNING "id", "tag", "address", "number", ST_AsEWKB(coordinates) AS coordinates, "instructions", "user_id", "province_id", "municipality_id", "create_time", "update_time"`, data.Tag, data.Address, data.Number, point, data.Instructions, data.UserId, data.ProvinceId, data.MunicipalityId, time, time).Scan(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return &res, nil
}

func (i *userAddressDatasource) GetUserAddress(tx *gorm.DB, where *models.UserAddress) (*models.UserAddress, error) {
	var res models.UserAddress
	result := tx.Raw(`SELECT "id", "tag", "address", "number", ST_AsEWKB(coordinates) AS coordinates, "instructions", "user_id", "province_id", "municipality_id", "create_time", "update_time" FROM "user_address" WHERE id = ? LIMIT 1`, where.ID).Scan(&res)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return &res, nil
}

func (v *userAddressDatasource) DeleteUserAddress(tx *gorm.DB, where *models.UserAddress, ids *[]uuid.UUID) (*[]models.UserAddress, error) {
	var res *[]models.UserAddress
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