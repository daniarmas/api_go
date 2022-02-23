package datasource

import (
	"fmt"
	"time"

	"github.com/daniarmas/api_go/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderDatasource interface {
	ListOrder(tx *gorm.DB, where *models.Order) (*[]models.Order, error)
	CreateOrder(tx *gorm.DB, data *models.Order) (*models.Order, error)
}

type orderDatasource struct{}

func (i *orderDatasource) CreateOrder(tx *gorm.DB, data *models.Order) (*models.Order, error) {
	point := fmt.Sprintf("POINT(%v %v)", data.Coordinates.Point.Coords()[1], data.Coordinates.Point.Coords()[0])
	// result := tx.Model(&models.Order{}).Create(map[string]interface{}{
	// 	"Status":               data.Status,
	// 	"DeliveryType":         data.DeliveryDate,
	// 	"ResidenceType":        data.ResidenceType,
	// 	"Price":                data.Price,
	// 	"BuildingNumber":       data.BuildingNumber,
	// 	"HouseNumber":          data.HouseNumber,
	// 	"BusinessFk":           data.BusinessFk.String(),
	// 	"Coordinates":          point[1 : len(point)-2],
	// 	"UserFk":               data.UserFk.String(),
	// 	"AuthorizationTokenFk": data.AuthorizationTokenFk.String(),
	// 	"DeliveryDate":         data.DeliveryDate,
	// })
	var time = time.Now().UTC()
	var response models.Order
	// query := fmt.Sprintf(`INSERT INTO "order" ("id", "authorization_token_fk", "building_number", "business_fk", "coordinates", "delivery_date", "delivery_type", "house_number", "price", "residence_type", "status", "user_fk", "create_time", "update_time") VALUES (%v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v) RETURNING "id", "status", "delivery_type", "residence_type", "price", "building_number", "house_number", "business_fk", "coordinates", "user_fk", "authorization_token_fk", "delivery_date", "create_time", "update_time"`, uuid.New().String(), data.AuthorizationTokenFk, data.BuildingNumber, data.BusinessFk.String(), point, data.DeliveryDate, data.DeliveryType, data.HouseNumber, data.Price, data.ResidenceType, data.Status, data.UserFk, time.String(), time.String())
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
