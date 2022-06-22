package repository

import (
	"context"
	"errors"
	"time"

	"github.com/daniarmas/api_go/datasource"
	"github.com/daniarmas/api_go/models"
	"github.com/go-redis/redis/v9"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ApplicationRepository interface {
	CreateApplication(tx *gorm.DB, data *models.Application) (*models.Application, error)
	GetApplication(tx *gorm.DB, where *models.Application, fields *[]string) (*models.Application, error)
	ListApplication(tx *gorm.DB, where *models.Application, cursor *time.Time, fields *[]string) (*[]models.Application, error)
	CheckApplication(tx *gorm.DB, accessToken string) error
}

type applicationRepository struct{}

func (i *applicationRepository) GetApplication(tx *gorm.DB, where *models.Application, fields *[]string) (*models.Application, error) {
	res, err := Datasource.NewApplicationDatasource().GetApplication(tx, where, fields)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *applicationRepository) CreateApplication(tx *gorm.DB, data *models.Application) (*models.Application, error) {
	res, err := Datasource.NewApplicationDatasource().CreateApplication(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *applicationRepository) ListApplication(tx *gorm.DB, where *models.Application, cursor *time.Time, fields *[]string) (*[]models.Application, error) {
	res, err := Datasource.NewApplicationDatasource().ListApplication(tx, where, cursor, fields)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *applicationRepository) CheckApplication(tx *gorm.DB, accessToken string) error {
	jwtAccessToken := datasource.JsonWebTokenMetadata{Token: &accessToken}
	accessTokenParseErr := Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(&jwtAccessToken)
	if accessTokenParseErr != nil {
		switch accessTokenParseErr.Error() {
		case "Token is expired":
			return errors.New("access token expired")
		case "signature is invalid":
			return errors.New("access token signature is invalid")
		case "token contains an invalid number of segments":
			return errors.New("access token contains an invalid number of segments")
		default:
			return accessTokenParseErr
		}
	}
	ctx := context.Background()
	cacheId := "application:" + jwtAccessToken.TokenId.String()
	cacheRes, cacheErr := Rdb.HGetAll(ctx, cacheId).Result()
	// Check if exists in cache
	if len(cacheRes) == 0 || cacheErr == redis.Nil {
		dbRes, dbErr := Datasource.NewApplicationDatasource().GetApplication(tx, &models.Application{ID: jwtAccessToken.TokenId}, nil)
		if dbErr != nil && dbErr.Error() == "record not found" {
			return errors.New("unauthenticated application")
		} else if dbErr != nil {
			return dbErr
		}
		if dbRes != nil && !dbRes.ExpirationTime.IsZero() {
			timeNow := time.Now().UTC()
			if timeNow.After(dbRes.ExpirationTime) {
				return errors.New("access token expired")
			}
		}
		go func() {
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
	} else {
		timeExp, _ := time.Parse(time.RFC3339, cacheRes["expiration_time"])
		if !timeExp.IsZero() {
			timeNow := time.Now().UTC()
			if timeNow.After(timeExp) {
				return errors.New("access token expired")
			}
		}
	}
	return nil
}
