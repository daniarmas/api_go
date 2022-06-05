package datasource

import (
	"errors"
	"fmt"
	"time"

	"github.com/daniarmas/api_go/models"
	"github.com/google/uuid"
	"github.com/teris-io/shortid"
	"gorm.io/gorm"
)

type OrderDatasource interface {
	ListOrder(tx *gorm.DB, where *models.Order) (*[]models.Order, error)
	ListOrderWithBusiness(tx *gorm.DB, where *models.OrderBusiness) (*[]models.OrderBusiness, error)
	CreateOrder(tx *gorm.DB, data *models.Order) (*models.Order, error)
	UpdateOrder(tx *gorm.DB, where *models.Order, data *models.Order) (*models.Order, error)
	GetOrder(tx *gorm.DB, where *models.Order) (*models.Order, error)
}

type orderDatasource struct{}

func (i *orderDatasource) GetOrder(tx *gorm.DB, where *models.Order) (*models.Order, error) {
	var res models.Order
	result := tx.Raw(`SELECT id, items_quantity, status, order_type, price, building_number, house_number, business_id, ST_AsEWKB(coordinates) AS coordinates, user_id, authorization_token_id, order_time, create_time, update_time FROM "order" WHERE id = ? LIMIT 1`, where.ID).Scan(&res)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return &res, nil
}

func (i *orderDatasource) UpdateOrder(tx *gorm.DB, where *models.Order, data *models.Order) (*models.Order, error) {
	// result := tx.Clauses(clause.Returning{}).Where(where).Updates(&data)
	var res models.Order
	var time = time.Now().UTC()
	result := tx.Raw(`UPDATE "order" SET "status"=?,"update_time"=?,"cancel_reasons"=? WHERE "order"."id" = ? AND "order"."user_id" = ? AND "order"."delete_time" IS NULL RETURNING "id", "status", "order_type", "price", "building_number", "cancel_reasons", "house_number", "business_id", ST_AsEWKB(coordinates) AS coordinates, "user_id", "authorization_token_id", "order_time", "create_time", "update_time"`, data.Status, time, where.CancelReasons, where.ID, where.UserId).Scan(&res)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return &res, nil
}

func (i *orderDatasource) CreateOrder(tx *gorm.DB, data *models.Order) (*models.Order, error) {
	point := fmt.Sprintf("POINT(%v %v)", data.Coordinates.Point.Coords()[1], data.Coordinates.Point.Coords()[0])
	var time = time.Now().UTC()
	shortId, err := shortid.Generate()
	if err != nil {
		return nil, err
	}
	var res models.Order
	result := tx.Raw(`INSERT INTO "order" ("id", "business_name", "items_quantity", "authorization_token_id", "business_id", "coordinates", "order_time", "order_type", "number", "address", "instructions", "price", "user_id", "create_time", "update_time", "short_id") VALUES (?, ?, ?, ?, ?, ST_GeomFromText(?, 4326), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING "id", "short_id", "items_quantity", "status", "order_type", "price", "number", "address", "business_id", ST_AsEWKB(coordinates) AS coordinates, "user_id", "authorization_token_id", "order_time", "create_time", "update_time"`, uuid.New().String(), data.BusinessName, data.ItemsQuantity, data.AuthorizationTokenId, data.BusinessId, point, data.OrderTime, data.OrderType, data.Number, data.Address, data.Instructions, data.Price, data.UserId, time, time, shortId).Scan(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return &res, nil
}

func (i *orderDatasource) ListOrder(tx *gorm.DB, where *models.Order) (*[]models.Order, error) {
	var res []models.Order
	result := tx.Limit(11).Select("id, status, cancel_reasons, items_quantity, order_type, price, number, address, instructions, business_id, authorization_token_id, user_id, order_time, create_time, update_time, delete_time, ST_AsEWKB(coordinates) AS coordinates").Where("user_id = ? AND create_time < ?", where.UserId, where.CreateTime).Order("create_time desc").Find(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return &res, nil
}

func (i *orderDatasource) ListOrderWithBusiness(tx *gorm.DB, where *models.OrderBusiness) (*[]models.OrderBusiness, error) {
	var res []models.OrderBusiness
	result := tx.Model(&models.Order{}).Limit(11).Select(`"order"."id", "order"."cancel_reasons", "order"."short_id", "order"."status", "order"."order_type", "order"."price", "order"."number", "order"."address", "order"."instructions", "order"."business_id", "order"."user_id", "order"."order_time", "order"."create_time", "order"."update_time", "order"."delete_time", ST_AsEWKB("order"."coordinates") AS "coordinates", "order"."business_name", "order"."items_quantity"`).Where(`"order"."user_id" = ? AND "order"."create_time" < ?`, where.UserId, where.CreateTime).Order(`"order"."create_time" desc`).Scan(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return &res, nil
}
