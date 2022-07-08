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

type PartnerApplicationDatasource interface {
	ListPartnerApplication(tx *gorm.DB, where *entity.PartnerApplication, fields *[]string, cursor *time.Time) (*[]entity.PartnerApplication, error)
	CreatePartnerApplication(tx *gorm.DB, data *entity.PartnerApplication) (*entity.PartnerApplication, error)
	UpdatePartnerApplication(tx *gorm.DB, where *entity.PartnerApplication, data *entity.PartnerApplication) (*entity.PartnerApplication, error)
	GetPartnerApplication(tx *gorm.DB, where *entity.PartnerApplication) (*entity.PartnerApplication, error)
	DeletePartnerApplication(tx *gorm.DB, where *entity.PartnerApplication, ids *[]uuid.UUID) (*[]entity.PartnerApplication, error)
}

type partnerApplicationDatasource struct{}

func (r *partnerApplicationDatasource) DeletePartnerApplication(tx *gorm.DB, where *entity.PartnerApplication, ids *[]uuid.UUID) (*[]entity.PartnerApplication, error) {
	var res *[]entity.PartnerApplication
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

func (i *partnerApplicationDatasource) GetPartnerApplication(tx *gorm.DB, where *entity.PartnerApplication) (*entity.PartnerApplication, error) {
	var res entity.PartnerApplication
	result := tx.Raw(`SELECT "id", "business_name", "description", "status", ST_AsEWKB(coordinates) AS coordinates, "user_id", "province_id", "municipality_id", "create_time", "update_time" FROM "partner_application" WHERE id = ? LIMIT 1`, where.ID).Scan(&res)
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

func (i *partnerApplicationDatasource) UpdatePartnerApplication(tx *gorm.DB, where *entity.PartnerApplication, data *entity.PartnerApplication) (*entity.PartnerApplication, error) {
	var res entity.PartnerApplication
	var time = time.Now().UTC()
	result := tx.Raw(`UPDATE "partner_application" SET "status"=?,"update_time"=? WHERE "partner_application"."id" = ? AND "partner_application"."delete_time" IS NULL RETURNING "id", "business_name", "description", "status", ST_AsEWKB(coordinates) AS coordinates, "province_id", "municipality_id", "user_id", "create_time", "update_time"`, data.Status, time, where.ID).Scan(&res)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return &res, nil
}

func (i *partnerApplicationDatasource) CreatePartnerApplication(tx *gorm.DB, data *entity.PartnerApplication) (*entity.PartnerApplication, error) {
	point := fmt.Sprintf("POINT(%v %v)", data.Coordinates.Point.Coords()[1], data.Coordinates.Point.Coords()[0])
	var time = time.Now().UTC()
	var res entity.PartnerApplication
	result := tx.Raw(`INSERT INTO "partner_application" ("user_id", "business_name", "description", "coordinates", "province_id", "municipality_id", "create_time", "update_time") VALUES (?, ?, ?, ST_GeomFromText(?, 4326), ?, ?, ?, ?) RETURNING "id", "business_name", "description", "status", ST_AsEWKB(coordinates) AS coordinates, "province_id", "municipality_id", "user_id", "create_time", "update_time"`, data.UserId, data.BusinessName, data.Description, point, data.ProvinceId, data.MunicipalityId, time, time).Scan(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return &res, nil
}

func (i *partnerApplicationDatasource) ListPartnerApplication(tx *gorm.DB, where *entity.PartnerApplication, fields *[]string, cursor *time.Time) (*[]entity.PartnerApplication, error) {
	var res []entity.PartnerApplication
	selectFields := &[]string{"id", "business_name", "description", "user_id", "status", "ST_AsEWKB(coordinates) AS coordinates", "province_id", "municipality_id", "create_time", "update_time"}
	if fields != nil {
		selectFields = fields
	}
	result := tx.Limit(11).Select(*selectFields).Where(where).Where(`partner_application.create_time < ?`, cursor).Order("create_time desc").Find(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return &res, nil
}
