package datasource

import (
	"errors"
	"fmt"
	"time"

	"github.com/daniarmas/api_go/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserAddressDatasource interface {
	ListUserAddress(tx *gorm.DB, where *entity.UserAddress) (*[]entity.UserAddress, error)
	CreateUserAddress(tx *gorm.DB, data *entity.UserAddress) (*entity.UserAddress, error)
	UpdateUserAddress(tx *gorm.DB, where *entity.UserAddress, data *entity.UserAddress) (*entity.UserAddress, error)
	UpdateUserAddressByUserId(tx *gorm.DB, where *entity.UserAddress, data *entity.UserAddress) (*entity.UserAddress, error)
	UpdateUserAddressSelected(tx *gorm.DB, where *entity.UserAddress, data *entity.UserAddress) (*entity.UserAddress, error)
	GetUserAddress(tx *gorm.DB, where *entity.UserAddress) (*entity.UserAddress, error)
	DeleteUserAddress(tx *gorm.DB, where *entity.UserAddress, ids *[]uuid.UUID) (*[]entity.UserAddress, error)
}

type userAddressDatasource struct{}

func (i *userAddressDatasource) ListUserAddress(tx *gorm.DB, where *entity.UserAddress) (*[]entity.UserAddress, error) {
	var res []entity.UserAddress
	selectFields := &[]string{"id", "selected", "name", "address", "number", "ST_AsEWKB(coordinates) AS coordinates", "instructions", "user_id", "province_id", "municipality_id", "create_time", "update_time"}
	result := tx.Where(where).Select(*selectFields).Find(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return &res, nil
}

func (i *userAddressDatasource) CreateUserAddress(tx *gorm.DB, data *entity.UserAddress) (*entity.UserAddress, error) {
	point := fmt.Sprintf("POINT(%v %v)", data.Coordinates.Point.Coords()[1], data.Coordinates.Point.Coords()[0])
	var time = time.Now().UTC()
	var res entity.UserAddress
	result := tx.Raw(`INSERT INTO "user_address" ("name", "selected", "address", "number", "coordinates", "instructions", "user_id", "province_id", "municipality_id", "create_time", "update_time") VALUES (?, ?, ?, ?, ST_GeomFromText(?, 4326), ?, ?, ?, ?, ?, ?) RETURNING "id", "selected", "name", "address", "number", ST_AsEWKB(coordinates) AS coordinates, "instructions", "user_id", "province_id", "municipality_id", "create_time", "update_time"`, data.Name, data.Selected, data.Address, data.Number, point, data.Instructions, data.UserId, data.ProvinceId, data.MunicipalityId, time, time).Scan(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return &res, nil
}

func (i *userAddressDatasource) UpdateUserAddressSelected(tx *gorm.DB, where *entity.UserAddress, data *entity.UserAddress) (*entity.UserAddress, error) {
	var res entity.UserAddress
	var time = time.Now().UTC()
	result := tx.Raw(`UPDATE "user_address" SET "selected"='true',"update_time"=? WHERE "user_address"."id" = ? AND "user_address"."delete_time" IS NULL RETURNING "id", "name", "selected", "address", "number", ST_AsEWKB(coordinates) AS coordinates, "instructions", "user_id", "province_id", "municipality_id", "create_time", "update_time"`, time, where.ID).Scan(&res)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	} else if result.RowsAffected == 0 {
		return nil, errors.New("record not found")
	}
	return &res, nil
}

func (i *userAddressDatasource) UpdateUserAddress(tx *gorm.DB, where *entity.UserAddress, data *entity.UserAddress) (*entity.UserAddress, error) {
	var point *string
	if data.Coordinates.Point != nil {
		value := fmt.Sprintf("POINT(%v %v)", data.Coordinates.Point.Coords()[1], data.Coordinates.Point.Coords()[0])
		point = &value
	}
	var res entity.UserAddress
	var time = time.Now().UTC()
	result := tx.Raw(`UPDATE "user_address" SET "selected"=?,"name"=?,"address"=?,"number"=?,"coordinates"=ST_GeomFromText(?, 4326),"instructions"=?,"user_id"=?,"province_id"=?,"municipality_id"=?,"create_time"=?,"update_time"=? WHERE "user_address"."id" = ? AND "user_address"."delete_time" IS NULL RETURNING "id", "name", "selected", "address", "number", ST_AsEWKB(coordinates) AS coordinates, "instructions", "user_id", "province_id", "municipality_id", "create_time", "update_time"`, data.Selected, data.Name, data.Address, data.Number, point, data.Instructions, data.UserId, data.ProvinceId, data.MunicipalityId, time, time, where.ID).Scan(&res)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	} else if result.RowsAffected == 0 {
		return nil, errors.New("record not found")
	}
	return &res, nil
}

func (i *userAddressDatasource) UpdateUserAddressByUserId(tx *gorm.DB, where *entity.UserAddress, data *entity.UserAddress) (*entity.UserAddress, error) {
	var res entity.UserAddress
	var time = time.Now().UTC()
	result := tx.Raw(`UPDATE "user_address" SET "selected"='false',"update_time"=? WHERE "user_address"."user_id" = ? AND "user_address"."delete_time" IS NULL RETURNING "id", "name", "selected", "address", "number", ST_AsEWKB(coordinates) AS coordinates, "instructions", "user_id", "province_id", "municipality_id", "create_time", "update_time"`, time, where.UserId).Scan(&res)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	} else if result.RowsAffected == 0 {
		return nil, errors.New("record not found")
	}
	return &res, nil
}

func (i *userAddressDatasource) GetUserAddress(tx *gorm.DB, where *entity.UserAddress) (*entity.UserAddress, error) {
	var res entity.UserAddress
	result := tx.Raw(`SELECT "id", "name", "selected", "address", "number", ST_AsEWKB(coordinates) AS coordinates, "instructions", "user_id", "province_id", "municipality_id", "create_time", "update_time" FROM "user_address" WHERE id = ? LIMIT 1`, where.ID).Scan(&res)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return &res, nil
}

func (v *userAddressDatasource) DeleteUserAddress(tx *gorm.DB, where *entity.UserAddress, ids *[]uuid.UUID) (*[]entity.UserAddress, error) {
	var res *[]entity.UserAddress
	var result *gorm.DB
	if ids != nil {
		result = tx.Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).Where(`id IN ?`, ids).Delete(&res)
	} else {
		result = tx.Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).Where(where).Delete(&res)
	}
	if result.Error != nil {
		return nil, result.Error
	} else if result.RowsAffected == 0 {
		return nil, errors.New("record not found")
	}
	return res, nil
}
