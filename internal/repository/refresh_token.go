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

type RefreshTokenRepository interface {
	GetRefreshToken(ctx context.Context, tx *gorm.DB, where *entity.RefreshToken) (*entity.RefreshToken, error)
	CreateRefreshToken(ctx context.Context, tx *gorm.DB, data *entity.RefreshToken) (*entity.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, tx *gorm.DB, where *entity.RefreshToken, ids *[]uuid.UUID) (*[]entity.RefreshToken, error)
	DeleteRefreshTokenDeviceIdNotEqual(ctx context.Context, tx *gorm.DB, where *entity.RefreshToken, ids *[]uuid.UUID) (*[]entity.RefreshToken, error)
}

type refreshTokenRepository struct{}

func (v *refreshTokenRepository) CreateRefreshToken(ctx context.Context, tx *gorm.DB, data *entity.RefreshToken) (*entity.RefreshToken, error) {
	// Store in the database
	dbRes, dbErr := Datasource.NewRefreshTokenDatasource().CreateRefreshToken(tx, data)
	if dbErr != nil {
		return nil, dbErr
	} else {
		// Store in cache
		go func() {
			ctx := context.Background()
			cacheId := "refresh_token:" + dbRes.ID.String()
			cacheErr := Rdb.HSet(ctx, cacheId, []string{
				"id", dbRes.ID.String(),
				"user_id", dbRes.UserId.String(),
				"device_id", dbRes.DeviceId.String(),
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

func (r *refreshTokenRepository) DeleteRefreshToken(ctx context.Context, tx *gorm.DB, where *entity.RefreshToken, ids *[]uuid.UUID) (*[]entity.RefreshToken, error) {
	// Delete in database
	dbRes, dbErr := Datasource.NewRefreshTokenDatasource().DeleteRefreshToken(tx, where, ids)
	if dbErr != nil {
		return nil, dbErr
	} else {
		// Delete in cache
		rdbPipe := Rdb.Pipeline()
		for _, item := range *dbRes {
			cacheId := "refresh_token:" + item.ID.String()
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

func (r *refreshTokenRepository) DeleteRefreshTokenDeviceIdNotEqual(ctx context.Context, tx *gorm.DB, where *entity.RefreshToken, ids *[]uuid.UUID) (*[]entity.RefreshToken, error) {
	// Delete in database
	dbRes, dbErr := Datasource.NewRefreshTokenDatasource().DeleteRefreshTokenDeviceIdNotEqual(tx, where, ids)
	if dbErr != nil {
		return nil, dbErr
	} else {
		// Delete in cache
		rdbPipe := Rdb.Pipeline()
		for _, item := range *dbRes {
			cacheId := "refresh_token:" + item.ID.String()
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

func (r *refreshTokenRepository) GetRefreshToken(ctx context.Context, tx *gorm.DB, where *entity.RefreshToken) (*entity.RefreshToken, error) {
	cacheId := "refresh_token:" + where.ID.String()
	cacheRes, cacheErr := Rdb.HGetAll(ctx, cacheId).Result()
	// Check if exists in cache
	if len(cacheRes) == 0 || cacheErr == redis.Nil {
		dbRes, dbErr := Datasource.NewRefreshTokenDatasource().GetRefreshToken(tx, where)
		if dbErr != nil {
			return nil, dbErr
		}
		// Store in cache
		go func() {
			ctx := context.Background()
			cacheErr := Rdb.HSet(ctx, cacheId, []string{
				"id", dbRes.ID.String(),
				"user_id", dbRes.UserId.String(),
				"device_id", dbRes.DeviceId.String(),
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
		deviceId := uuid.MustParse(cacheRes["device_id"])
		createTime, _ := time.Parse(time.RFC3339, cacheRes["create_time"])
		updateTime, _ := time.Parse(time.RFC3339, cacheRes["update_time"])
		return &entity.RefreshToken{
			ID:         &id,
			UserId:     &userId,
			DeviceId:   &deviceId,
			CreateTime: createTime,
			UpdateTime: updateTime,
		}, nil
	}
}
