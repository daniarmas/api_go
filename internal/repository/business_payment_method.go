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

type BusinessPaymentMethodRepository interface {
	ListBusinessPaymentMethodWithEnabled(ctx context.Context, tx *gorm.DB, where *entity.BusinessPaymentMethod) (*[]entity.BusinessPaymentMethodWithEnabled, error)
	ListBusinessPaymentMethod(ctx context.Context, tx *gorm.DB, where *entity.BusinessPaymentMethod) (*[]entity.BusinessPaymentMethod, error)
	CreateBusinessPaymentMethod(ctx context.Context, tx *gorm.DB, data *entity.BusinessPaymentMethod) (*entity.BusinessPaymentMethod, error)
	UpdateBusinessPaymentMethod(ctx context.Context, tx *gorm.DB, where *entity.BusinessPaymentMethod, data *entity.BusinessPaymentMethod) (*entity.BusinessPaymentMethod, error)
	GetBusinessPaymentMethod(ctx context.Context, tx *gorm.DB, where *entity.BusinessPaymentMethod) (*entity.BusinessPaymentMethod, error)
	DeleteBusinessPaymentMethod(ctx context.Context, tx *gorm.DB, where *entity.BusinessPaymentMethod, ids *[]uuid.UUID) (*[]entity.BusinessPaymentMethod, error)
}

type businessPaymentMethodRepository struct{}

func (i *businessPaymentMethodRepository) ListBusinessPaymentMethodWithEnabled(ctx context.Context, tx *gorm.DB, where *entity.BusinessPaymentMethod) (*[]entity.BusinessPaymentMethodWithEnabled, error) {
	dbRes, err := Datasource.NewBusinessPaymentMethodDatasource().ListBusinessPaymentMethodWithEnabled(tx, where)
	if err != nil {
		return nil, err
	}
	return dbRes, nil
}

func (i *businessPaymentMethodRepository) ListBusinessPaymentMethod(ctx context.Context, tx *gorm.DB, where *entity.BusinessPaymentMethod) (*[]entity.BusinessPaymentMethod, error) {
	result, err := Datasource.NewBusinessPaymentMethodDatasource().ListBusinessPaymentMethod(tx, where)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *businessPaymentMethodRepository) DeleteBusinessPaymentMethod(ctx context.Context, tx *gorm.DB, where *entity.BusinessPaymentMethod, ids *[]uuid.UUID) (*[]entity.BusinessPaymentMethod, error) {
	// Delete in database
	res, err := Datasource.NewBusinessPaymentMethodDatasource().DeleteBusinessPaymentMethod(tx, where, ids)
	if err != nil {
		return nil, err
	}
	// Delete in cache
	rdbPipe := Rdb.Pipeline()
	for _, item := range *res {
		cacheId := "business_payment_method:" + item.ID.String()
		cacheErr := rdbPipe.Del(ctx, cacheId).Err()
		if cacheErr != nil {
			log.Error(cacheErr)
		}
	}
	_, err = rdbPipe.Exec(ctx)
	if err != nil {
		log.Error(err)
	}
	return res, nil
}

func (i *businessPaymentMethodRepository) GetBusinessPaymentMethod(ctx context.Context, tx *gorm.DB, where *entity.BusinessPaymentMethod) (*entity.BusinessPaymentMethod, error) {
	cacheId := "business_payment_method:" + where.ID.String()
	cacheRes, cacheErr := Rdb.HGetAll(ctx, cacheId).Result()
	// Check if exists in cache
	if len(cacheRes) == 0 || cacheErr == redis.Nil {
		dbRes, dbErr := Datasource.NewBusinessPaymentMethodDatasource().GetBusinessPaymentMethod(tx, where)
		if dbErr != nil {
			return nil, dbErr
		}
		// Store in cache
		go func() {
			ctx := context.Background()
			cacheErr := Rdb.HSet(ctx, cacheId, []string{
				"id", dbRes.ID.String(),
				"type", dbRes.Type,
				"address", dbRes.Address,
				"business_id", dbRes.BusinessId.String(),
				"payment_method_id", dbRes.PaymentMethodId.String(),
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
		paymentMethodId := uuid.MustParse(cacheRes["payment_method_id"])
		createTime, _ := time.Parse(time.RFC3339, cacheRes["create_time"])
		updateTime, _ := time.Parse(time.RFC3339, cacheRes["update_time"])
		return &entity.BusinessPaymentMethod{
			ID:              &id,
			Type:            cacheRes["type"],
			Address:         cacheRes["address"],
			PaymentMethodId: &paymentMethodId,
			BusinessId:      &businessId,
			CreateTime:      createTime,
			UpdateTime:      updateTime,
		}, nil
	}
}

func (i *businessPaymentMethodRepository) CreateBusinessPaymentMethod(ctx context.Context, tx *gorm.DB, data *entity.BusinessPaymentMethod) (*entity.BusinessPaymentMethod, error) {
	// Store in the database
	dbRes, dbErr := Datasource.NewBusinessPaymentMethodDatasource().CreateBusinessPaymentMethod(tx, data)
	if dbErr != nil {
		return nil, dbErr
	}
	// Store in cache
	go func() {
		ctx := context.Background()
		cacheId := "business_payment_method:" + dbRes.ID.String()
		cacheErr := Rdb.HSet(ctx, cacheId, []string{
			"id", dbRes.ID.String(),
			"type", dbRes.Type,
			"address", dbRes.Address,
			"business_id", dbRes.BusinessId.String(),
			"payment_method_id", dbRes.PaymentMethodId.String(),
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

func (i *businessPaymentMethodRepository) UpdateBusinessPaymentMethod(ctx context.Context, tx *gorm.DB, where *entity.BusinessPaymentMethod, data *entity.BusinessPaymentMethod) (*entity.BusinessPaymentMethod, error) {
	// Store in the database
	dbRes, dbErr := Datasource.NewBusinessPaymentMethodDatasource().UpdateBusinessPaymentMethod(tx, where, data)
	if dbErr != nil {
		return nil, dbErr
	}
	// Store in cache
	cacheId := "business_payment_method:" + dbRes.ID.String()
	cacheErr := Rdb.HSet(ctx, cacheId, []string{
		"id", dbRes.ID.String(),
		"type", dbRes.Type,
		"address", dbRes.Address,
		"business_id", dbRes.BusinessId.String(),
		"payment_method_id", dbRes.PaymentMethodId.String(),
		"create_time", dbRes.CreateTime.Format(time.RFC3339),
		"update_time", dbRes.UpdateTime.Format(time.RFC3339),
	}).Err()
	if cacheErr != nil {
		log.Error(cacheErr)
	} else {
		Rdb.Expire(ctx, cacheId, time.Minute*15)
	}
	return dbRes, nil
}
