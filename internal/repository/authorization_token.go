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

type AuthorizationTokenRepository interface {
	GetAuthorizationToken(ctx context.Context, tx *gorm.DB, where *entity.AuthorizationToken, fields *[]string) (*entity.AuthorizationToken, error)
	CreateAuthorizationToken(ctx context.Context, tx *gorm.DB, data *entity.AuthorizationToken) (*entity.AuthorizationToken, error)
	DeleteAuthorizationToken(ctx context.Context, tx *gorm.DB, where *entity.AuthorizationToken, ids *[]uuid.UUID) (*[]entity.AuthorizationToken, error)
	DeleteAuthorizationTokenByRefreshTokenIds(ctx context.Context, tx *gorm.DB, ids *[]uuid.UUID) (*[]entity.AuthorizationToken, error)
}

type authorizationTokenRepository struct{}

func (v *authorizationTokenRepository) CreateAuthorizationToken(ctx context.Context, tx *gorm.DB, data *entity.AuthorizationToken) (*entity.AuthorizationToken, error) {
	// Store in the database
	dbRes, dbErr := Datasource.NewAuthorizationTokenDatasource().CreateAuthorizationToken(tx, data)
	if dbErr != nil {
		return nil, dbErr
	}
	// Store in cache
	go func() {
		cacheId := "authorization_token:" + dbRes.ID.String()
		cacheErr := Rdb.HSet(context.Background(), cacheId, []string{
			"id", dbRes.ID.String(),
			"refresh_token_id", dbRes.RefreshTokenId.String(),
			"user_id", dbRes.UserId.String(),
			"device_id", dbRes.DeviceId.String(),
			"app", *dbRes.App,
			"app_version", *dbRes.AppVersion,
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
}

func (r *authorizationTokenRepository) DeleteAuthorizationToken(ctx context.Context, tx *gorm.DB, where *entity.AuthorizationToken, ids *[]uuid.UUID) (*[]entity.AuthorizationToken, error) {
	// Delete in cache
	if where.ID != nil {
		go func() {
			cacheId := "authorization_token:" + where.ID.String()
			cacheErr := Rdb.Del(ctx, cacheId).Err()
			if cacheErr != nil {
				log.Error(cacheErr)
			}
		}()
	}
	// Delete in database
	dbRes, dbErr := Datasource.NewAuthorizationTokenDatasource().DeleteAuthorizationToken(tx, where, ids)
	if dbErr != nil {
		return nil, dbErr
	}
	return dbRes, nil
}

func (r *authorizationTokenRepository) DeleteAuthorizationTokenByRefreshTokenIds(ctx context.Context, tx *gorm.DB, ids *[]uuid.UUID) (*[]entity.AuthorizationToken, error) {
	// Delete in cache
	go func() {
		for _, i := range *ids {
			cacheId := "authorization_token:" + i.String()
			cacheErr := Rdb.Del(ctx, cacheId).Err()
			if cacheErr != nil {
				log.Error(cacheErr)
			}
		}
	}()
	// Delete in database
	res, err := Datasource.NewAuthorizationTokenDatasource().DeleteAuthorizationTokenByRefreshTokenIds(tx, ids)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (v *authorizationTokenRepository) GetAuthorizationToken(ctx context.Context, tx *gorm.DB, where *entity.AuthorizationToken, fields *[]string) (*entity.AuthorizationToken, error) {
	cacheId := "authorization_token:" + where.ID.String()
	cacheRes, cacheErr := Rdb.HGetAll(ctx, cacheId).Result()
	// Check if exists in cache
	if len(cacheRes) == 0 || cacheErr == redis.Nil {
		dbRes, dbErr := Datasource.NewAuthorizationTokenDatasource().GetAuthorizationToken(tx, where, fields)
		if dbErr != nil {
			return nil, dbErr
		}
		// Store in cache
		go func() {
			cacheErr := Rdb.HSet(ctx, cacheId, []string{
				"id", dbRes.ID.String(),
				"refresh_token_id", dbRes.RefreshTokenId.String(),
				"user_id", dbRes.UserId.String(),
				"device_id", dbRes.DeviceId.String(),
				"app", *dbRes.App,
				"app_version", *dbRes.AppVersion,
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
		refreshTokenId := uuid.MustParse(cacheRes["refresh_token_id"])
		userId := uuid.MustParse(cacheRes["user_id"])
		deviceId := uuid.MustParse(cacheRes["device_id"])
		app := cacheRes["app"]
		appVersion := cacheRes["app_version"]
		createTime, _ := time.Parse(time.RFC3339, cacheRes["create_time"])
		updateTime, _ := time.Parse(time.RFC3339, cacheRes["update_time"])
		return &entity.AuthorizationToken{
			ID:             &id,
			RefreshTokenId: &refreshTokenId,
			UserId:         &userId,
			DeviceId:       &deviceId,
			App:            &app,
			AppVersion:     &appVersion,
			CreateTime:     createTime,
			UpdateTime:     updateTime,
		}, nil
	}
}
