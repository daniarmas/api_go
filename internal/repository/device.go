package repository

import (
	"context"
	"time"

	"github.com/daniarmas/api_go/internal/entity"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type DeviceRepository interface {
	GetDevice(ctx context.Context, tx *gorm.DB, where *entity.Device, fields *[]string) (*entity.Device, error)
	CreateDevice(ctx context.Context, tx *gorm.DB, data *entity.Device) (*entity.Device, error)
	UpdateDevice(ctx context.Context, tx *gorm.DB, where *entity.Device, data *entity.Device) (*entity.Device, error)
}

type deviceRepository struct{}

func (i *deviceRepository) GetDevice(ctx context.Context, tx *gorm.DB, where *entity.Device, fields *[]string) (*entity.Device, error) {
	cacheId := "device:" + where.DeviceIdentifier
	cacheRes, cacheErr := Rdb.HGetAll(ctx, cacheId).Result()
	// Check if exists in cache
	if len(cacheRes) == 0 || cacheErr == redis.Nil {
		dbRes, dbErr := Datasource.NewDeviceDatasource().GetDevice(tx, where, nil)
		if dbErr != nil {
			return nil, dbErr
		}
		// Store in the cache
		go func() {
			cacheErr := Rdb.HSet(context.Background(), cacheId, []string{
				"id", dbRes.ID.String(),
				"platform", dbRes.Platform,
				"system_version", dbRes.SystemVersion,
				"device_identifier", dbRes.DeviceIdentifier,
				"firebase_cloud_messaging_id", dbRes.FirebaseCloudMessagingId,
				"model", dbRes.Model,
				"create_time", dbRes.CreateTime.Format(time.RFC3339),
				"update_time", dbRes.UpdateTime.Format(time.RFC3339),
			}).Err()
			if cacheErr != nil {
				log.Error(cacheErr)
			} else {
				Rdb.Expire(ctx, cacheId, time.Minute*15)
			}
		}()
		return dbRes, nil
	} else {
		id := uuid.MustParse(cacheRes["id"])
		createTime, _ := time.Parse(time.RFC3339, cacheRes["create_time"])
		updateTime, _ := time.Parse(time.RFC3339, cacheRes["update_time"])
		return &entity.Device{
			ID:                       &id,
			Platform:                 cacheRes["platform"],
			SystemVersion:            cacheRes["system_version"],
			DeviceIdentifier:         cacheRes["device_identifier"],
			FirebaseCloudMessagingId: cacheRes["firebase_cloud_messaging_id"],
			Model:                    cacheRes["platform"],
			CreateTime:               createTime,
			UpdateTime:               updateTime,
		}, nil
	}
}

func (v *deviceRepository) CreateDevice(ctx context.Context, tx *gorm.DB, data *entity.Device) (*entity.Device, error) {
	// Store in the database
	dbRes, dbErr := Datasource.NewDeviceDatasource().CreateDevice(tx, data)
	if dbErr != nil {
		return nil, dbErr
	} else {
		// Store in cache
		go func() {
			cacheId := "device:" + dbRes.ID.String()
			cacheErr := Rdb.HSet(context.Background(), cacheId, []string{
				"id", dbRes.ID.String(),
				"platform", dbRes.Platform,
				"system_version", dbRes.SystemVersion,
				"device_identifier", dbRes.DeviceIdentifier,
				"firebase_cloud_messaging_id", dbRes.FirebaseCloudMessagingId,
				"model", dbRes.Model,
				"create_time", dbRes.CreateTime.Format(time.RFC3339),
				"update_time", dbRes.UpdateTime.Format(time.RFC3339),
			}).Err()
			if cacheErr != nil {
				log.Error(cacheErr)
			} else {
				Rdb.Expire(ctx, cacheId, time.Minute*15)
			}
		}()
	}
	return dbRes, nil
}

func (v *deviceRepository) UpdateDevice(ctx context.Context, tx *gorm.DB, where *entity.Device, data *entity.Device) (*entity.Device, error) {
	// Store in the database
	dbRes, dbErr := Datasource.NewDeviceDatasource().UpdateDevice(tx, where, data)
	if dbErr != nil {
		return nil, dbErr
	}
	// Store in cache
	cacheId := "device:" + dbRes.DeviceIdentifier
	cacheErr := Rdb.HSet(ctx, cacheId, []string{
		"id", dbRes.ID.String(),
		"platform", dbRes.Platform,
		"system_version", dbRes.SystemVersion,
		"device_identifier", dbRes.DeviceIdentifier,
		"firebase_cloud_messaging_id", dbRes.FirebaseCloudMessagingId,
		"model", dbRes.Model,
		"create_time", dbRes.CreateTime.Format(time.RFC3339),
		"update_time", dbRes.UpdateTime.Format(time.RFC3339),
	}).Err()
	if cacheErr != nil {
		log.Error(cacheErr)
	} else {
		Rdb.Expire(ctx, cacheId, time.Minute*15)
	}
	return dbRes, nil
}
