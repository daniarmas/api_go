package datasource

import (
	"errors"
	"fmt"
	"time"

	"github.com/daniarmas/api_go/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	// "gorm.io/gorm/clause"
)

type OrderDatasource interface {
	ListOrder(tx *gorm.DB, where *models.Order) (*[]models.Order, error)
	CreateOrder(tx *gorm.DB, data *models.Order) (*models.Order, error)
	UpdateOrder(tx *gorm.DB, where *models.Order, data *models.Order) (*models.Order, error)
}

type orderDatasource struct{}

func (i *orderDatasource) UpdateOrder(tx *gorm.DB, where *models.Order, data *models.Order) (*models.Order, error) {
	// result := tx.Clauses(clause.Returning{}).Where(where).Updates(&data)
	var response models.Order
	var time = time.Now().UTC()
	result := tx.Raw(`UPDATE "order" SET "status"=?,"update_time"=? WHERE "order"."id" = ? AND "order"."user_fk" = ? AND "order"."delete_time" IS NULL RETURNING "id", "status", "delivery_type", "residence_type", "price", "building_number", "house_number", "business_fk", ST_AsEWKB(coordinates) AS coordinates, "user_fk", "authorization_token_fk", "delivery_date", "create_time", "update_time"`, data.Status, time, where.ID, where.UserFk).Scan(&response)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return &response, nil
}

func (i *orderDatasource) CreateOrder(tx *gorm.DB, data *models.Order) (*models.Order, error) {
	point := fmt.Sprintf("POINT(%v %v)", data.Coordinates.Point.Coords()[1], data.Coordinates.Point.Coords()[0])
	var time = time.Now().UTC()
	var response models.Order
	result := tx.Raw(`INSERT INTO "order" ("id", "authorization_token_fk", "building_number", "business_fk", "coordinates", "delivery_date", "delivery_type", "house_number", "price", "residence_type", "status", "user_fk", "create_time", "update_time") VALUES (?, ?, ?, ?, ST_GeomFromText(?, 4326), ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING "id", "status", "delivery_type", "residence_type", "price", "building_number", "house_number", "business_fk", ST_AsEWKB(coordinates) AS coordinates, "user_fk", "authorization_token_fk", "delivery_date", "create_time", "update_time"`, uuid.New().String(), data.AuthorizationTokenFk, data.BuildingNumber, data.BusinessFk.String(), point, data.DeliveryDate, data.DeliveryType, data.HouseNumber, data.Price, data.ResidenceType, data.Status, data.UserFk, time, time).Scan(&response)
	if result.Error != nil {
		return nil, result.Error
	}
	return &response, nil
}

func (i *orderDatasource) ListOrder(tx *gorm.DB, where *models.Order) (*[]models.Order, error) {
	var order []models.Order
	result := tx.Limit(11).Select("id, status, delivery_type, residence_type, price, building_number, house_number, business_fk, user_fk, device_fk, app_version, delivery_date, create_time, update_time, delete_time, ST_AsEWKB(coordinates) AS coordinates").Where("user_fk = ? AND create_time < ?", where.UserFk, where.CreateTime).Order("create_time desc").Find(&order)
	if result.Error != nil {
		return nil, result.Error
	}
	return &order, nil
}
