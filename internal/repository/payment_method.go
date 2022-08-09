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

type PaymentMethodRepository interface {
	ListPaymentMethod(ctx context.Context, tx *gorm.DB, where *entity.PaymentMethod) (*[]entity.PaymentMethod, error)
	CreatePaymentMethod(ctx context.Context, tx *gorm.DB, data *entity.PaymentMethod) (*entity.PaymentMethod, error)
	UpdatePaymentMethod(ctx context.Context, tx *gorm.DB, where *entity.PaymentMethod, data *entity.PaymentMethod) (*entity.PaymentMethod, error)
	GetPaymentMethod(ctx context.Context, tx *gorm.DB, where *entity.PaymentMethod) (*entity.PaymentMethod, error)
	DeletePaymentMethod(ctx context.Context, tx *gorm.DB, where *entity.PaymentMethod, ids *[]uuid.UUID) (*[]entity.PaymentMethod, error)
}

type paymentMethodRepository struct{}

func (i *paymentMethodRepository) ListPaymentMethod(ctx context.Context, tx *gorm.DB, where *entity.PaymentMethod) (*[]entity.PaymentMethod, error) {
	// Get from database
	dbRes, dbErr := Datasource.NewPaymentMethodDatasource().ListPaymentMethod(tx, where)
	if dbErr != nil {
		return nil, dbErr
	} else {
		// Store in cache
		go func() {
			ctx := context.Background()
			rdbPipe := Rdb.Pipeline()
			for _, item := range *dbRes {
				cacheId := "payment_method:" + item.ID.String()
				cacheErr := rdbPipe.HSet(ctx, cacheId, []string{
					"id", item.ID.String(),
					"type", item.Type,
					"address", item.Address,
					"enabled", strconv.FormatBool(item.Enabled),
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

func (i *paymentMethodRepository) DeletePaymentMethod(ctx context.Context, tx *gorm.DB, where *entity.PaymentMethod, ids *[]uuid.UUID) (*[]entity.PaymentMethod, error) {
	// Delete in database
	dbRes, dbErr := Datasource.NewPaymentMethodDatasource().DeletePaymentMethod(tx, where, ids)
	if dbErr != nil {
		return nil, dbErr
	} else {
		// Delete in cache
		rdbPipe := Rdb.Pipeline()
		for _, item := range *dbRes {
			cacheId := "payment_method:" + item.ID.String()
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

func (i *paymentMethodRepository) GetPaymentMethod(ctx context.Context, tx *gorm.DB, where *entity.PaymentMethod) (*entity.PaymentMethod, error) {
	cacheId := "payment_method:" + where.ID.String()
	cacheRes, cacheErr := Rdb.HGetAll(ctx, cacheId).Result()
	// Check if exists in cache
	if len(cacheRes) == 0 || cacheErr == redis.Nil {
		dbRes, dbErr := Datasource.NewPaymentMethodDatasource().GetPaymentMethod(tx, where)
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
				"enabled", strconv.FormatBool(dbRes.Enabled),
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
		enabled, _ := strconv.ParseBool(cacheRes["type"])
		return &entity.PaymentMethod{
			ID:         &id,
			Type:       cacheRes["type"],
			Address:    cacheRes["address"],
			Enabled:    enabled,
			CreateTime: createTime,
			UpdateTime: updateTime,
		}, nil
	}
}

func (i *paymentMethodRepository) CreatePaymentMethod(ctx context.Context, tx *gorm.DB, data *entity.PaymentMethod) (*entity.PaymentMethod, error) {
	// Store in the database
	dbRes, dbErr := Datasource.NewPaymentMethodDatasource().CreatePaymentMethod(tx, data)
	if dbErr != nil {
		return nil, dbErr
	} else {
		// Store in cache
		go func() {
			ctx := context.Background()
			cacheId := "payment_method:" + dbRes.ID.String()
			cacheErr := Rdb.HSet(ctx, cacheId, []string{
				"id", dbRes.ID.String(),
				"type", dbRes.Type,
				"address", dbRes.Address,
				"enabled", strconv.FormatBool(dbRes.Enabled),
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

func (i *paymentMethodRepository) UpdatePaymentMethod(ctx context.Context, tx *gorm.DB, where *entity.PaymentMethod, data *entity.PaymentMethod) (*entity.PaymentMethod, error) {
	// Store in the database
	dbRes, dbErr := Datasource.NewPaymentMethodDatasource().UpdatePaymentMethod(tx, where, data)
	if dbErr != nil {
		return nil, dbErr
	} else {
		// Store in cache
		cacheId := "payment_method:" + dbRes.ID.String()
		cacheErr := Rdb.HSet(ctx, cacheId, []string{
			"id", dbRes.ID.String(),
			"type", dbRes.Type,
			"address", dbRes.Address,
			"enabled", strconv.FormatBool(dbRes.Enabled),
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
