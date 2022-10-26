package datasource

import (
	"errors"
	"fmt"
	"time"

	"github.com/daniarmas/api_go/internal/entity"
	"github.com/google/uuid"
	"github.com/teris-io/shortid"
	"gorm.io/gorm"
)

type OrderDatasource interface {
	ListOrder(tx *gorm.DB, where *entity.Order) (*[]entity.Order, error)
	ListOrderFilter(tx *gorm.DB, where *entity.OrderBusiness, upcoming bool) (*[]entity.OrderBusiness, error)
	ListOrderWithBusiness(tx *gorm.DB, where *entity.OrderBusiness) (*[]entity.OrderBusiness, error)
	CreateOrder(tx *gorm.DB, data *entity.Order) (*entity.Order, error)
	UpdateOrder(tx *gorm.DB, where *entity.Order, data *entity.Order) (*entity.Order, error)
	GetOrder(tx *gorm.DB, where *entity.Order) (*entity.Order, error)
	ExistsUpcomingOrders(tx *gorm.DB, userId uuid.UUID) (*bool, error)
}

type orderDatasource struct{}

func (i *orderDatasource) ExistsUpcomingOrders(tx *gorm.DB, userId uuid.UUID) (*bool, error) {
	var res []entity.OrderBusiness
	var ret bool
	result := tx.Model(&entity.Order{}).Select(`"order"."id"`).Where(`"order"."user_id" = ? AND (status = 'OrderStatusTypePendingPayment' OR status = 'OrderStatusTypeOrdered' OR status = 'OrderStatusTypeAccepted' OR status = 'OrderStatusTypeReady' OR status = 'OrderStatusTypeAssignedMessenger')`, userId).Limit(1).Scan(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return &ret, nil
		} else {
			return nil, result.Error
		}
	}
	if len(res) != 0 {
		ret = true
		return &ret, nil
	} else {
		return &ret, nil
	}
}

func (i *orderDatasource) GetOrder(tx *gorm.DB, where *entity.Order) (*entity.Order, error) {
	var res entity.Order
	result := tx.Raw(`SELECT "id", "delivery_price_cup", "short_id", "business_name", "business_thumbnail", "items_quantity", "status", "order_type", "price_cup", "number", "address", "business_id", ST_AsEWKB(coordinates) AS coordinates, "user_id", "authorization_token_id", "start_order_time", "end_order_time", "create_time", "update_time", "instructions", "cancel_reasons" FROM "order" WHERE id = ? LIMIT 1`, where.ID).Scan(&res)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return &res, nil
}

func (i *orderDatasource) UpdateOrder(tx *gorm.DB, where *entity.Order, data *entity.Order) (*entity.Order, error) {
	var res entity.Order
	var time = time.Now().UTC()
	result := tx.Raw(`UPDATE "order" SET "status"=?,"update_time"=?,"cancel_reasons"=? WHERE "order"."id" = ? AND "order"."delete_time" IS NULL RETURNING "id", "short_id", "items_quantity", "status", "order_type", "price_cup", "number", "address", "business_id", ST_AsEWKB(coordinates) AS coordinates, "user_id", "authorization_token_id", "start_order_time", "end_order_time", "create_time", "update_time", "instructions", "cancel_reasons", "business_name"`, data.Status, time, data.CancelReasons, where.ID).Scan(&res)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("record not found")
		} else {
			return nil, result.Error
		}
	}
	return &res, nil
}

func (i *orderDatasource) CreateOrder(tx *gorm.DB, data *entity.Order) (*entity.Order, error) {
	point := fmt.Sprintf("POINT(%v %v)", data.Coordinates.Point.Coords()[1], data.Coordinates.Point.Coords()[0])
	var time = time.Now().UTC()
	shortId, err := shortid.Generate()
	if err != nil {
		return nil, err
	}
	var res entity.Order
	result := tx.Raw(`INSERT INTO "order" ("id", "delivery_price_cup", "business_thumbnail", "status", "business_name", "items_quantity", "authorization_token_id", "business_id", "coordinates", "start_order_time", "end_order_time", "order_type", "number", "address", "instructions", "price_cup", "user_id", "create_time", "update_time", "short_id", "payment_method_type", "phone") VALUES (?, ?, ?, ?, ?, ?, ?, ?, ST_GeomFromText(?, 4326), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING "id", "delivery_price_cup", "business_thumbnail", "short_id", "business_name", "items_quantity", "status", "order_type", "price_cup", "number", "address", "business_id", ST_AsEWKB(coordinates) AS coordinates, "user_id", "authorization_token_id", "start_order_time", "end_order_time", "instructions", "cancel_reasons", "create_time", "update_time", "payment_method_type", "phone"`, uuid.New().String(), data.DeliveryPriceCup, data.BusinessThumbnail, data.Status, data.BusinessName, data.ItemsQuantity, data.AuthorizationTokenId, data.BusinessId, point, data.StartOrderTime, data.EndOrderTime, data.OrderType, data.Number, data.Address, data.Instructions, data.PriceCup, data.UserId, time, time, shortId, data.PaymentMethodType, data.Phone).Scan(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return &res, nil
}

func (i *orderDatasource) ListOrder(tx *gorm.DB, where *entity.Order) (*[]entity.Order, error) {
	var res []entity.Order
	selectFields := &[]string{"id", "delivery_price_cup", "business_thumbnail", "status", "business_name", "short_id", "cancel_reasons", "items_quantity", "order_type", "price_cup", "number", "address", "instructions", "business_id", "authorization_token_id", "user_id", "start_order_time", "end_order_time", "create_time", "update_time", "delete_time", "ST_AsEWKB(coordinates) AS coordinates"}
	result := tx.Limit(11).Select(*selectFields).Where("user_id = ? AND create_time < ?", where.UserId, where.CreateTime).Order("create_time desc").Find(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return &res, nil
}

func (i *orderDatasource) ListOrderFilter(tx *gorm.DB, where *entity.OrderBusiness, upcoming bool) (*[]entity.OrderBusiness, error) {
	var res []entity.OrderBusiness
	if upcoming {
		result := tx.Model(&entity.Order{}).Limit(11).Select(`"order"."id", "order"."delivery_price_cup", "order"."business_thumbnail", "order"."cancel_reasons", "order"."short_id", "order"."status", "order"."order_type", "order"."price_cup", "order"."number", "order"."address", "order"."instructions", "order"."business_id", "order"."user_id", "order"."start_order_time", "order"."end_order_time", "order"."create_time", "order"."update_time", "order"."delete_time", ST_AsEWKB("order"."coordinates") AS "coordinates", "order"."business_name", "order"."items_quantity"`).Where(`"order"."user_id" = ? AND "order"."create_time" < ? AND (status = 'OrderStatusTypePendingPayment' OR status = 'OrderStatusTypeOrdered' OR status = 'OrderStatusTypeAccepted' OR status = 'OrderStatusTypeReady' OR status = 'OrderStatusTypeAssignedMessenger')`, where.UserId, where.CreateTime).Order(`"order"."create_time" desc`).Scan(&res)
		if result.Error != nil {
			return nil, result.Error
		}
	} else {
		result := tx.Model(&entity.Order{}).Limit(11).Select(`"order"."id", "order"."delivery_price_cup", "order"."business_thumbnail", "order"."cancel_reasons", "order"."short_id", "order"."status", "order"."order_type", "order"."price_cup", "order"."number", "order"."address", "order"."instructions", "order"."business_id", "order"."user_id", "order"."start_order_time", "order"."end_order_time", "order"."create_time", "order"."update_time", "order"."delete_time", ST_AsEWKB("order"."coordinates") AS "coordinates", "order"."business_name", "order"."items_quantity"`).Where(`"order"."user_id" = ? AND "order"."create_time" < ? AND (status != 'OrderStatusTypePendingPayment' AND status != 'OrderStatusTypeOrdered' AND status != 'OrderStatusTypeAccepted' AND status != 'OrderStatusTypeReady' AND status != 'OrderStatusTypeAssignedMessenger')`, where.UserId, where.CreateTime).Order(`"order"."create_time" desc`).Scan(&res)
		if result.Error != nil {
			return nil, result.Error
		}
	}
	return &res, nil
}

func (i *orderDatasource) ListOrderWithBusiness(tx *gorm.DB, where *entity.OrderBusiness) (*[]entity.OrderBusiness, error) {
	var res []entity.OrderBusiness
	result := tx.Model(&entity.Order{}).Limit(11).Select(`"order"."id", "order"."delivery_price_cup", "order"."business_thumbnail", "order"."cancel_reasons", "order"."short_id", "order"."status", "order"."order_type", "order"."price_cup", "order"."number", "order"."address", "order"."instructions", "order"."business_id", "order"."user_id", "order"."start_order_time", "order"."end_order_time", "order"."create_time", "order"."update_time", "order"."delete_time", ST_AsEWKB("order"."coordinates") AS "coordinates", "order"."business_name", "order"."items_quantity"`).Where(`"order"."user_id" = ? AND "order"."create_time" < ?`, where.UserId, where.CreateTime).Order(`"order"."create_time" desc`).Scan(&res)
	if result.Error != nil {
		return nil, result.Error
	}
	return &res, nil
}
