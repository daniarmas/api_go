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

type BusinessUserRepository interface {
	GetBusinessUser(ctx context.Context, tx *gorm.DB, where *entity.BusinessUser, fields *[]string) (*entity.BusinessUser, error)
	CreateBusinessUser(ctx context.Context, tx *gorm.DB, data *entity.BusinessUser) (*entity.BusinessUser, error)
	DeleteBusinessUser(ctx context.Context, tx *gorm.DB, where *entity.BusinessUser, ids *[]uuid.UUID) (*[]entity.BusinessUser, error)
}

type businessUserRepository struct{}

func (v *businessUserRepository) CreateBusinessUser(ctx context.Context, tx *gorm.DB, data *entity.BusinessUser) (*entity.BusinessUser, error) {
	// Store in the database
	dbRes, dbErr := Datasource.NewBusinessUserDatasource().CreateBusinessUser(tx, data)
	if dbErr != nil {
		return nil, dbErr
	} else {
		// Store in cache
		go func() {
			ctx := context.Background()
			cacheId := "business_user:" + dbRes.ID.String()
			cacheErr := Rdb.HSet(ctx, cacheId, []string{
				"id", dbRes.ID.String(),
				"is_business_owner", strconv.FormatBool(dbRes.IsBusinessOwner),
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

func (v *businessUserRepository) GetBusinessUser(ctx context.Context, tx *gorm.DB, where *entity.BusinessUser, fields *[]string) (*entity.BusinessUser, error) {
	cacheId := "business_user:" + where.ID.String()
	cacheRes, cacheErr := Rdb.HGetAll(ctx, cacheId).Result()
	// Check if exists in cache
	if len(cacheRes) == 0 || cacheErr == redis.Nil {
		dbRes, dbErr := Datasource.NewBusinessUserDatasource().GetBusinessUser(tx, where, fields)
		if dbErr != nil {
			return nil, dbErr
		}
		// Store in cache
		go func() {
			ctx := context.Background()
			cacheErr := Rdb.HSet(ctx, cacheId, []string{
				"id", dbRes.ID.String(),
				"is_business_owner", strconv.FormatBool(dbRes.IsBusinessOwner),
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
		createTime, _ := time.Parse(time.RFC3339, cacheRes["create_time"])
		updateTime, _ := time.Parse(time.RFC3339, cacheRes["update_time"])
		isBusinessOwner, _ := strconv.ParseBool(cacheRes["is_business_owner"])
		return &entity.BusinessUser{
			ID:              &id,
			UserId:          &userId,
			IsBusinessOwner: isBusinessOwner,
			CreateTime:      createTime,
			UpdateTime:      updateTime,
		}, nil
	}
}

func (v *businessUserRepository) DeleteBusinessUser(ctx context.Context, tx *gorm.DB, where *entity.BusinessUser, ids *[]uuid.UUID) (*[]entity.BusinessUser, error) {
	// Delete in database
	res, err := Datasource.NewBusinessUserDatasource().DeleteBusinessUser(tx, where, ids)
	if err != nil {
		return nil, err
	} else {
		// Delete in cache
		rdbPipe := Rdb.Pipeline()
		for _, item := range *res {
			cacheId := "business_user:" + item.ID.String()
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
	return res, nil
}
