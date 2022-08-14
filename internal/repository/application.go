package repository

import (
	"context"
	"errors"
	"time"

	"github.com/daniarmas/api_go/internal/datasource"
	"github.com/daniarmas/api_go/internal/entity"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ApplicationRepository interface {
	CreateApplication(ctx context.Context, tx *gorm.DB, data *entity.Application) (*entity.Application, error)
	GetApplication(ctx context.Context, tx *gorm.DB, where *entity.Application) (*entity.Application, error)
	ListApplication(ctx context.Context, tx *gorm.DB, where *entity.Application, cursor *time.Time) (*[]entity.Application, error)
	CheckApplication(ctx context.Context, tx *gorm.DB, accessToken string) (*entity.Application, error)
	DeleteApplication(ctx context.Context, tx *gorm.DB, where *entity.Application, ids *[]uuid.UUID) (*[]entity.Application, error)
}

type applicationRepository struct{}

func (i *applicationRepository) GetApplication(ctx context.Context, tx *gorm.DB, where *entity.Application) (*entity.Application, error) {
	res, err := Datasource.NewApplicationDatasource().GetApplication(tx, where)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *applicationRepository) CreateApplication(ctx context.Context, tx *gorm.DB, data *entity.Application) (*entity.Application, error) {
	res, err := Datasource.NewApplicationDatasource().CreateApplication(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *applicationRepository) DeleteApplication(ctx context.Context, tx *gorm.DB, where *entity.Application, ids *[]uuid.UUID) (*[]entity.Application, error) {
	// Delete in database
	dbRes, dbErr := Datasource.NewApplicationDatasource().DeleteApplication(tx, where, ids)
	if dbErr != nil {
		return nil, dbErr
	} else {
		// Delete in cache
		rdbPipe := Rdb.Pipeline()
		for _, item := range *dbRes {
			cacheId := "application:" + item.ID.String()
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

func (i *applicationRepository) ListApplication(ctx context.Context, tx *gorm.DB, where *entity.Application, cursor *time.Time) (*[]entity.Application, error) {
	res, err := Datasource.NewApplicationDatasource().ListApplication(tx, where, cursor)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *applicationRepository) CheckApplication(ctx context.Context, tx *gorm.DB, accessToken string) (*entity.Application, error) {
	jwtAccessToken := datasource.JsonWebTokenMetadata{Token: &accessToken}
	accessTokenParseErr := Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(&jwtAccessToken)
	if accessTokenParseErr != nil {
		switch accessTokenParseErr.Error() {
		case "Token is expired":
			return nil, errors.New("access token expired")
		case "signature is invalid":
			return nil, errors.New("access token signature is invalid")
		case "token contains an invalid number of segments":
			return nil, errors.New("access token contains an invalid number of segments")
		default:
			return nil, accessTokenParseErr
		}
	}
	cacheId := "application:" + jwtAccessToken.TokenId.String()
	cacheRes, cacheErr := Rdb.HGetAll(ctx, cacheId).Result()
	// Check if exists in cache
	if len(cacheRes) == 0 || cacheErr == redis.Nil {
		dbRes, dbErr := Datasource.NewApplicationDatasource().GetApplication(tx, &entity.Application{ID: jwtAccessToken.TokenId})
		if dbErr != nil && dbErr.Error() == "record not found" {
			return nil, errors.New("unauthenticated application")
		} else if dbErr != nil {
			return nil, dbErr
		}
		if dbRes != nil && !dbRes.ExpirationTime.IsZero() {
			timeNow := time.Now().UTC()
			if timeNow.After(dbRes.ExpirationTime) {
				return nil, errors.New("access token expired")
			}
		}
		go func() {
			ctx := context.Background()
			cacheErr := Rdb.HSet(ctx, cacheId, []string{
				"id", dbRes.ID.String(),
				"name", dbRes.Name,
				"version", dbRes.Version,
				"description", dbRes.Description,
				"expiration_time", dbRes.ExpirationTime.Format(time.RFC3339),
				"create_time", dbRes.CreateTime.Format(time.RFC3339),
				"update_time", dbRes.UpdateTime.Format(time.RFC3339),
			}).Err()
			if cacheErr != nil {
				log.Error(cacheErr)
			} else {
				Rdb.Expire(ctx, cacheId, time.Hour*24)
			}
		}()
		return dbRes, nil
	} else {
		timeExp, _ := time.Parse(time.RFC3339, cacheRes["expiration_time"])
		if !timeExp.IsZero() {
			timeNow := time.Now().UTC()
			if timeNow.After(timeExp) {
				return nil, errors.New("access token expired")
			}
		}
	}
	id := uuid.MustParse(cacheRes["id"])
	expirationTime, _ := time.Parse(time.RFC3339, cacheRes["expiration_time"])
	createTime, _ := time.Parse(time.RFC3339, cacheRes["create_time"])
	updateTime, _ := time.Parse(time.RFC3339, cacheRes["update_time"])
	return &entity.Application{
		ID:             &id,
		Name:           cacheRes["name"],
		Version:        cacheRes["version"],
		Description:    cacheRes["version"],
		ExpirationTime: expirationTime,
		CreateTime:     createTime,
		UpdateTime:     updateTime,
	}, nil
}
