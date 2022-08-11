package repository

import (
	"context"
	"strconv"
	"time"

	"github.com/daniarmas/api_go/internal/entity"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserConfigurationRepository interface {
	ListUserConfiguration(ctx context.Context, tx *gorm.DB, where *entity.UserConfiguration) (*[]entity.UserConfiguration, error)
	CreateUserConfiguration(ctx context.Context, tx *gorm.DB, data *entity.UserConfiguration) (*entity.UserConfiguration, error)
	UpdateUserConfiguration(ctx context.Context, tx *gorm.DB, where *entity.UserConfiguration, data *entity.UserConfiguration) (*entity.UserConfiguration, error)
	GetUserConfiguration(ctx context.Context, tx *gorm.DB, where *entity.UserConfiguration) (*entity.UserConfiguration, error)
	DeleteUserConfiguration(ctx context.Context, tx *gorm.DB, where *entity.UserConfiguration, ids *[]uuid.UUID) (*[]entity.UserConfiguration, error)
}

type userConfigurationRepository struct{}

func (i *userConfigurationRepository) ListUserConfiguration(ctx context.Context, tx *gorm.DB, where *entity.UserConfiguration) (*[]entity.UserConfiguration, error) {
	// Get from database
	dbRes, dbErr := Datasource.NewUserConfigurationDatasource().ListUserConfiguration(tx, where)
	if dbErr != nil {
		return nil, dbErr
	} else {
		// Delete in cache
		go func() {
			ctx := context.Background()
			rdbPipe := Rdb.Pipeline()
			for _, item := range *dbRes {
				cacheId := "user_configuration:" + item.ID.String()
				cacheErr := rdbPipe.HSet(ctx, cacheId, []string{
					"id", item.ID.String(),
					"data_saving", strconv.FormatBool(*item.DataSaving),
					"high_quality_images_wifi", strconv.FormatBool(*item.HighQualityImagesWifi),
					"high_quality_images_data", strconv.FormatBool(*item.HighQualityImagesData),
					"payment_method", item.PaymentMethod,
					"user_id", item.UserId.String(),
					"create_time", item.CreateTime.Format(time.RFC3339),
					"update_time", item.UpdateTime.Format(time.RFC3339),
				}).Err()
				if cacheErr != nil {
					log.Error(cacheErr)
				} else {
					rdbPipe.Expire(ctx, cacheId, time.Minute*15)
				}
			}
			_, err := rdbPipe.Exec(ctx)
			if err != nil {
				log.Error(err)
			}
		}()
	}
	return dbRes, nil
}

func (i *userConfigurationRepository) DeleteUserConfiguration(ctx context.Context, tx *gorm.DB, where *entity.UserConfiguration, ids *[]uuid.UUID) (*[]entity.UserConfiguration, error) {
	// Delete in database
	dbRes, dbErr := Datasource.NewUserConfigurationDatasource().DeleteUserConfiguration(tx, where, ids)
	if dbErr != nil {
		return nil, dbErr
	} else {
		// Delete in cache
		rdbPipe := Rdb.Pipeline()
		for _, item := range *dbRes {
			cacheId := "user_configuration:" + item.ID.String()
			cacheErr := rdbPipe.Del(ctx, cacheId).Err()
			if cacheErr != nil {
				log.Error(cacheErr)
			}
		}
		_, err := rdbPipe.Exec(ctx)
		if err != nil {
			log.Error(err)
		}
	}
	return dbRes, nil
}

func (i *userConfigurationRepository) GetUserConfiguration(ctx context.Context, tx *gorm.DB, where *entity.UserConfiguration) (*entity.UserConfiguration, error) {
	var cacheId string
	if where.ID != nil {
		cacheId = "user_configuration:" + where.ID.String()
	} else {
		cacheId = "user_configuration:" + where.UserId.String()
	}
	cacheRes, cacheErr := Rdb.HGetAll(ctx, cacheId).Result()
	// Check if exists in cache
	if len(cacheRes) == 0 || cacheErr == redis.Nil {
		dbRes, dbErr := Datasource.NewUserConfigurationDatasource().GetUserConfiguration(tx, where)
		if dbErr != nil {
			return nil, dbErr
		}
		// Store in cache
		go func() {
			ctx := context.Background()
			cacheErr := Rdb.HSet(ctx, cacheId, []string{
				"id", dbRes.ID.String(),
				"data_saving", strconv.FormatBool(*dbRes.DataSaving),
				"high_quality_images_wifi", strconv.FormatBool(*dbRes.HighQualityImagesWifi),
				"high_quality_images_data", strconv.FormatBool(*dbRes.HighQualityImagesData),
				"payment_method", dbRes.PaymentMethod,
				"user_id", dbRes.UserId.String(),
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
		userId := uuid.MustParse(cacheRes["user_id"])
		dataSaving, _ := strconv.ParseBool(cacheRes["data_saving"])
		highQualityImagesWifi, _ := strconv.ParseBool(cacheRes["high_quality_images_wifi"])
		highQualityImagesData, _ := strconv.ParseBool(cacheRes["high_quality_images_data"])
		createTime, _ := time.Parse(time.RFC3339, cacheRes["create_time"])
		updateTime, _ := time.Parse(time.RFC3339, cacheRes["update_time"])
		return &entity.UserConfiguration{
			ID:                    &id,
			DataSaving:            &dataSaving,
			HighQualityImagesWifi: &highQualityImagesWifi,
			HighQualityImagesData: &highQualityImagesData,
			UserId:                &userId,
			PaymentMethod:         cacheRes["payment_method"],
			CreateTime:            createTime,
			UpdateTime:            updateTime,
		}, nil
	}
}

func (i *userConfigurationRepository) CreateUserConfiguration(ctx context.Context, tx *gorm.DB, data *entity.UserConfiguration) (*entity.UserConfiguration, error) {
	// Store in the database
	dbRes, dbErr := Datasource.NewUserConfigurationDatasource().CreateUserConfiguration(tx, data)
	if dbErr != nil {
		return nil, dbErr
	} else {
		// Store in cache
		go func() {
			ctx := context.Background()
			cacheId := "user_configuration:" + dbRes.ID.String()
			cacheErr := Rdb.HSet(ctx, cacheId, []string{
				"id", dbRes.ID.String(),
				"data_saving", strconv.FormatBool(*dbRes.DataSaving),
				"high_quality_images_wifi", strconv.FormatBool(*dbRes.HighQualityImagesWifi),
				"high_quality_images_data", strconv.FormatBool(*dbRes.HighQualityImagesData),
				"payment_method", dbRes.PaymentMethod,
				"user_id", dbRes.UserId.String(),
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

func (i *userConfigurationRepository) UpdateUserConfiguration(ctx context.Context, tx *gorm.DB, where *entity.UserConfiguration, data *entity.UserConfiguration) (*entity.UserConfiguration, error) {
	// Store in the database
	dbRes, dbErr := Datasource.NewUserConfigurationDatasource().UpdateUserConfiguration(tx, where, data)
	if dbErr != nil {
		return nil, dbErr
	} else {
		// Store in cache
		cacheId := "user_configuration:" + dbRes.ID.String()
		cacheErr := Rdb.HSet(ctx, cacheId, []string{
			"id", dbRes.ID.String(),
			"data_saving", strconv.FormatBool(*dbRes.DataSaving),
			"high_quality_images_wifi", strconv.FormatBool(*dbRes.HighQualityImagesWifi),
			"high_quality_images_data", strconv.FormatBool(*dbRes.HighQualityImagesData),
			"payment_method", dbRes.PaymentMethod,
			"user_id", dbRes.UserId.String(),
			"create_time", dbRes.CreateTime.Format(time.RFC3339),
			"update_time", dbRes.UpdateTime.Format(time.RFC3339),
		}).Err()
		if cacheErr != nil {
			log.Error(cacheErr)
		} else {
			Rdb.Expire(ctx, cacheId, time.Minute*15)
		}
	}
	return dbRes, nil
}
