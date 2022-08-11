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

type PermissionRepository interface {
	GetPermission(ctx context.Context, tx *gorm.DB, where *entity.Permission) (*entity.Permission, error)
	ListPermissionAll(ctx context.Context, tx *gorm.DB, where *entity.Permission) (*[]entity.Permission, error)
	ListPermissionByIdAll(ctx context.Context, tx *gorm.DB, where *entity.Permission, ids *[]uuid.UUID) (*[]entity.Permission, error)
}

type permissionRepository struct{}

func (i *permissionRepository) GetPermission(ctx context.Context, tx *gorm.DB, where *entity.Permission) (*entity.Permission, error) {
	cacheId := "permission:" + where.ID.String()
	cacheRes, cacheErr := Rdb.HGetAll(ctx, cacheId).Result()
	// Check if exists in cache
	if len(cacheRes) == 0 || cacheErr == redis.Nil {
		dbRes, dbErr := Datasource.NewPermissionDatasource().GetPermission(tx, where)
		if dbErr != nil {
			return nil, dbErr
		}
		// Store in cache
		go func() {
			ctx := context.Background()
			cacheErr := Rdb.HSet(ctx, cacheId, []string{
				"id", dbRes.ID.String(),
				"name", dbRes.Name,
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
		return &entity.Permission{
			ID:         &id,
			Name:       cacheRes["name"],
			CreateTime: createTime,
			UpdateTime: updateTime,
		}, nil
	}
}

func (i *permissionRepository) ListPermissionAll(ctx context.Context, tx *gorm.DB, where *entity.Permission) (*[]entity.Permission, error) {
	// Get from database
	dbRes, dbErr := Datasource.NewPermissionDatasource().ListPermissionAll(tx, where)
	if dbErr != nil {
		return nil, dbErr
	} else {
		// Delete in cache
		go func() {
			ctx := context.Background()
			rdbPipe := Rdb.Pipeline()
			for _, item := range *dbRes {
				cacheId := "permission:" + item.ID.String()
				cacheErr := rdbPipe.HSet(ctx, cacheId, []string{
					"id", item.ID.String(),
					"name", item.Name,
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

func (i *permissionRepository) ListPermissionByIdAll(ctx context.Context, tx *gorm.DB, where *entity.Permission, ids *[]uuid.UUID) (*[]entity.Permission, error) {
	// Get from database
	dbRes, dbErr := Datasource.NewPermissionDatasource().ListPermissionByIdAll(tx, where, ids)
	if dbErr != nil {
		return nil, dbErr
	} else {
		// Delete in cache
		go func() {
			ctx := context.Background()
			rdbPipe := Rdb.Pipeline()
			for _, item := range *dbRes {
				cacheId := "permission:" + item.ID.String()
				cacheErr := rdbPipe.HSet(ctx, cacheId, []string{
					"id", item.ID.String(),
					"name", item.Name,
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
