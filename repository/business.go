package repository

import (
	"context"
	"strconv"
	"time"

	"github.com/daniarmas/api_go/models"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"gorm.io/gorm"
)

type BusinessRepository interface {
	Feed(tx *gorm.DB, coordinates ewkb.Point, limit int32, provinceId string, municipalityId string, cursor int32, municipalityNotEqual bool, homeDelivery bool, toPickUp bool) (*[]models.Business, error)
	CreateBusiness(tx *gorm.DB, data *models.Business) (*models.Business, error)
	GetBusiness(tx *gorm.DB, where *models.Business, fields *[]string) (*models.Business, error)
	GetBusinessWithDistance(tx *gorm.DB, where *models.Business) (*models.Business, error)
	UpdateBusiness(tx *gorm.DB, data *models.Business, where *models.Business) (*models.Business, error)
	UpdateBusinessCoordinate(tx *gorm.DB, data *models.Business, where *models.Business) error
	BusinessIsInRange(tx *gorm.DB, coordinates ewkb.Point, businessId *uuid.UUID) (*bool, error)
}

type businessRepository struct{}

func (b *businessRepository) BusinessIsInRange(tx *gorm.DB, coordinates ewkb.Point, businessId *uuid.UUID) (*bool, error) {
	lat := strconv.FormatFloat(coordinates.Point.Coords()[0], 'E', -1, 64)
	long := strconv.FormatFloat(coordinates.Point.Coords()[1], 'E', -1, 64)
	cacheId := "business_is_in_range:" + lat + long + ":" + businessId.String()
	ctx := context.Background()
	cacheRes, cacheErr := Rdb.HGetAll(ctx, cacheId).Result()
	// Check if exists in cache
	if len(cacheRes) == 0 || cacheErr == redis.Nil {
		dbRes, dbErr := Datasource.NewBusinessDatasource().BusinessIsInRange(tx, coordinates, businessId)
		if dbErr != nil {
			return nil, dbErr
		}
		go func() {
			cacheErr := Rdb.HSet(ctx, cacheId, []string{
				"is_in_range", strconv.FormatBool(*dbRes),
			}).Err()
			if cacheErr != nil {
				log.Error(cacheErr)
			} else {
				Rdb.Expire(ctx, cacheId, time.Minute*30)
			}
		}()
		return dbRes, nil
	} else {
		isInRange, _ := strconv.ParseBool(cacheRes["is_in_range"])
		return &isInRange, nil
	}
}

func (b *businessRepository) CreateBusiness(tx *gorm.DB, data *models.Business) (*models.Business, error) {
	res, err := Datasource.NewBusinessDatasource().CreateBusiness(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (b *businessRepository) UpdateBusiness(tx *gorm.DB, data *models.Business, where *models.Business) (*models.Business, error) {
	res, err := Datasource.NewBusinessDatasource().UpdateBusiness(tx, data, where)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (b *businessRepository) UpdateBusinessCoordinate(tx *gorm.DB, data *models.Business, where *models.Business) error {
	err := Datasource.NewBusinessDatasource().UpdateBusinessCoordinate(tx, data, where)
	if err != nil {
		return err
	}
	return nil
}

func (b *businessRepository) Feed(tx *gorm.DB, coordinates ewkb.Point, limit int32, provinceId string, municipalityId string, cursor int32, municipalityNotEqual bool, homeDelivery bool, toPickUp bool) (*[]models.Business, error) {
	result, err := Datasource.NewBusinessDatasource().Feed(tx, coordinates, limit, provinceId, municipalityId, cursor, municipalityNotEqual, homeDelivery, toPickUp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (b *businessRepository) GetBusinessWithDistance(tx *gorm.DB, where *models.Business) (*models.Business, error) {
	result, err := Datasource.NewBusinessDatasource().GetBusinessWithDistance(tx, where)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (b *businessRepository) GetBusiness(tx *gorm.DB, where *models.Business, fields *[]string) (*models.Business, error) {
	result, err := Datasource.NewBusinessDatasource().GetBusiness(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return result, nil
}
