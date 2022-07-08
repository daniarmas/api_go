package repository

import (
	"context"
	"time"

	"github.com/daniarmas/api_go/internal/entity"
	"github.com/daniarmas/api_go/utils"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserPermissionRepository interface {
	CreateUserPermission(tx *gorm.DB, data *[]entity.UserPermission) (*[]entity.UserPermission, error)
	GetUserPermission(tx *gorm.DB, where *entity.UserPermission, fields *[]string) (*entity.UserPermission, error)
	DeleteUserPermission(tx *gorm.DB, where *entity.UserPermission, ids *[]uuid.UUID) (*[]entity.UserPermission, error)
	DeleteUserPermissionByBusinessRoleId(tx *gorm.DB, where *entity.UserPermission) (*[]entity.UserPermission, error)
	DeleteUserPermissionByPermissionId(tx *gorm.DB, permissionIds *[]uuid.UUID) (*[]entity.UserPermission, error)
}

type userPermissionRepository struct{}

func (v *userPermissionRepository) DeleteUserPermissionByBusinessRoleId(tx *gorm.DB, where *entity.UserPermission) (*[]entity.UserPermission, error) {
	res, err := Datasource.NewUserPermissionDatasource().DeleteUserPermissionByBusinessRoleId(tx, where)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (v *userPermissionRepository) CreateUserPermission(tx *gorm.DB, data *[]entity.UserPermission) (*[]entity.UserPermission, error) {
	res, err := Datasource.NewUserPermissionDatasource().CreateUserPermission(tx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (v *userPermissionRepository) DeleteUserPermissionByPermissionId(tx *gorm.DB, permission_ids *[]uuid.UUID) (*[]entity.UserPermission, error) {
	res, err := Datasource.NewUserPermissionDatasource().DeleteUserPermissionByPermissionId(tx, permission_ids)
	if err != nil {
		return nil, err
	} else {
		ctx := context.Background()
		rdbPipe := Rdb.Pipeline()
		for _, i := range *res {
			cacheId := "user_permission:" + i.ID.String()
			cacheErr := rdbPipe.Del(ctx, cacheId).Err()
			if cacheErr != nil {
				log.Error(cacheErr)
			}
		}
		_, err := rdbPipe.Exec(ctx)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (v *userPermissionRepository) GetUserPermission(tx *gorm.DB, where *entity.UserPermission, fields *[]string) (*entity.UserPermission, error) {
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
		return &entity.UserPermission{
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

func (v *userPermissionRepository) DeleteUserPermission(tx *gorm.DB, where *entity.UserPermission, ids *[]uuid.UUID) (*[]entity.UserPermission, error) {
	ctx := context.Background()
	cacheId := "user_permission:" + where.ID.String()
	cacheErr := Rdb.Del(ctx, cacheId).Err()
	if cacheErr != nil {
		log.Error(cacheErr)
	}
	res, err := Datasource.NewUserPermissionDatasource().DeleteUserPermission(tx, where, ids)
	if err != nil {
		return nil, err
	}
	return res, nil
}
