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

type BusinessRoleRepository interface {
	CreateBusinessRole(ctx context.Context, tx *gorm.DB, data *entity.BusinessRole) (*entity.BusinessRole, error)
	UpdateBusinessRole(ctx context.Context, tx *gorm.DB, where *entity.BusinessRole, data *entity.BusinessRole) (*entity.BusinessRole, error)
	GetBusinessRole(ctx context.Context, tx *gorm.DB, where *entity.BusinessRole) (*entity.BusinessRole, error)
	ListBusinessRole(ctx context.Context, tx *gorm.DB, where *entity.BusinessRole, cursor *time.Time) (*[]entity.BusinessRole, error)
	DeleteBusinessRole(ctx context.Context, tx *gorm.DB, where *entity.BusinessRole, ids *[]uuid.UUID) (*[]entity.BusinessRole, error)
}

type businessRoleRepository struct{}

func (v *businessRoleRepository) UpdateBusinessRole(ctx context.Context, tx *gorm.DB, where *entity.BusinessRole, data *entity.BusinessRole) (*entity.BusinessRole, error) {
	// Store in the database
	dbRes, dbErr := Datasource.NewBusinessRoleDatasource().UpdateBusinessRole(tx, where, data)
	if dbErr != nil {
		return nil, dbErr
	} else {
		// Store in cache
		cacheId := "business_role:" + dbRes.ID.String()
		cacheErr := Rdb.HSet(ctx, cacheId, []string{
			"id", dbRes.ID.String(),
			"name", dbRes.Name,
			"business_id", dbRes.BusinessId.String(),
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

func (v *businessRoleRepository) CreateBusinessRole(ctx context.Context, tx *gorm.DB, data *entity.BusinessRole) (*entity.BusinessRole, error) {
	// Store in the database
	dbRes, dbErr := Datasource.NewBusinessRoleDatasource().CreateBusinessRole(tx, data)
	if dbErr != nil {
		return nil, dbErr
	} else {
		// Store in cache
		go func() {
			ctx := context.Background()
			cacheId := "business_role:" + dbRes.ID.String()
			cacheErr := Rdb.HSet(ctx, cacheId, []string{
				"id", dbRes.ID.String(),
				"name", dbRes.Name,
				"business_id", dbRes.BusinessId.String(),
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

func (v *businessRoleRepository) DeleteBusinessRole(ctx context.Context, tx *gorm.DB, where *entity.BusinessRole, ids *[]uuid.UUID) (*[]entity.BusinessRole, error) {
	// Delete in database
	dbRes, dbErr := Datasource.NewBusinessRoleDatasource().DeleteBusinessRole(tx, where, ids)
	if dbErr != nil {
		return nil, dbErr
	} else {
		// Delete in cache
		rdbPipe := Rdb.Pipeline()
		for _, item := range *dbRes {
			cacheId := "business_role:" + item.ID.String()
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

func (i *businessRoleRepository) ListBusinessRole(ctx context.Context, tx *gorm.DB, where *entity.BusinessRole, cursor *time.Time) (*[]entity.BusinessRole, error) {
	// Get from database
	dbRes, dbErr := Datasource.NewBusinessRoleDatasource().ListBusinessRole(tx, where, cursor)
	if dbErr != nil {
		return nil, dbErr
	} else {
		// Delete in cache
		go func() {
			ctx := context.Background()
			rdbPipe := Rdb.Pipeline()
			for _, item := range *dbRes {
				cacheId := "business_role:" + item.ID.String()
				cacheErr := rdbPipe.HSet(ctx, cacheId, []string{
					"id", item.ID.String(),
					"name", item.Name,
					"business_id", item.BusinessId.String(),
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

func (i *businessRoleRepository) GetBusinessRole(ctx context.Context, tx *gorm.DB, where *entity.BusinessRole) (*entity.BusinessRole, error) {
	cacheId := "business_role:" + where.ID.String()
	cacheRes, cacheErr := Rdb.HGetAll(ctx, cacheId).Result()
	// Check if exists in cache
	if len(cacheRes) == 0 || cacheErr == redis.Nil {
		dbRes, dbErr := Datasource.NewBusinessRoleDatasource().GetBusinessRole(tx, where)
		if dbErr != nil {
			return nil, dbErr
		}
		// Store in cache
		go func() {
			ctx := context.Background()
			cacheErr := Rdb.HSet(ctx, cacheId, []string{
				"id", dbRes.ID.String(),
				"name", dbRes.Name,
				"business_id", dbRes.BusinessId.String(),
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
		businessId := uuid.MustParse(cacheRes["business_id"])
		createTime, _ := time.Parse(time.RFC3339, cacheRes["create_time"])
		updateTime, _ := time.Parse(time.RFC3339, cacheRes["update_time"])
		return &entity.BusinessRole{
			ID:         &id,
			Name:       cacheRes["name"],
			BusinessId: &businessId,
			CreateTime: createTime,
			UpdateTime: updateTime,
		}, nil
	}
}
