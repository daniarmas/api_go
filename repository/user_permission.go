package repository

import (
	"context"
	"time"

	"github.com/daniarmas/api_go/models"
	"github.com/daniarmas/api_go/utils"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserPermissionRepository interface {
	GetUserPermission(tx *gorm.DB, where *models.UserPermission, fields *[]string) (*models.UserPermission, error)
}

type userPermissionRepository struct{}

func (v *userPermissionRepository) GetUserPermission(tx *gorm.DB, where *models.UserPermission, fields *[]string) (*models.UserPermission, error) {
	var cacheId string
	if where.BusinessId != nil {
		cacheId = "user_permission:" + where.Name + ":" + where.BusinessId.String() + ":" + where.UserId.String()
	} else {
		cacheId = "user_permission:" + where.Name + ":" + where.UserId.String()
	}
	ctx := context.Background()
	cacheRes, cacheErr := Rdb.HGetAll(ctx, cacheId).Result()
	// Check if exists in cache
	if len(cacheRes) == 0 || cacheErr == redis.Nil {
		dbRes, dbErr := Datasource.NewUserPermissionDatasource().GetUserPermission(tx, where, nil)
		if dbErr != nil {
			return nil, dbErr
		}
		ctx := context.Background()
		rdbPipe := Rdb.Pipeline()
		var businessId string
		if dbRes.BusinessId != nil {
			businessId = dbRes.BusinessId.String()
		}
		cacheErr := rdbPipe.HSet(ctx, cacheId, []string{
			"id", dbRes.ID.String(),
			"name", dbRes.Name,
			"user_id", dbRes.UserId.String(),
			"business_id", businessId,
			"permission_id", dbRes.PermissionId.String(),
			"create_time", dbRes.CreateTime.Format(time.RFC3339),
			"update_time", dbRes.UpdateTime.Format(time.RFC3339),
		}).Err()
		if cacheErr != nil {
			return nil, cacheErr
		} else {
			rdbPipe.Expire(ctx, cacheId, time.Minute*5)
		}
		_, err := rdbPipe.Exec(ctx)
		if err != nil {
			return nil, err
		}
		return dbRes, nil
	} else {
		id := uuid.MustParse(cacheRes["id"])
		userId := uuid.MustParse(cacheRes["user_id"])
		permissionId := uuid.MustParse(cacheRes["permission_id"])
		createTime, _ := time.Parse(time.RFC3339, cacheRes["create_time"])
		updateTime, _ := time.Parse(time.RFC3339, cacheRes["update_time"])
		return &models.UserPermission{
			ID:           &id,
			Name:         cacheRes["name"],
			UserId:       &userId,
			BusinessId:   utils.UuidParse(cacheRes["business_id"]),
			PermissionId: &permissionId,
			CreateTime:   createTime,
			UpdateTime:   updateTime,
		}, nil
	}
}
