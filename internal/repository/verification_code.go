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

type VerificationCodeRepository interface {
	GetVerificationCode(ctx context.Context, tx *gorm.DB, where *entity.VerificationCode, fields *[]string) (*entity.VerificationCode, error)
	CreateVerificationCode(ctx context.Context, tx *gorm.DB, data *entity.VerificationCode) (*entity.VerificationCode, error)
	DeleteVerificationCode(ctx context.Context, tx *gorm.DB, where *entity.VerificationCode, ids *[]uuid.UUID) (*[]entity.VerificationCode, error)
}

type verificationCodeRepository struct{}

func (v *verificationCodeRepository) CreateVerificationCode(ctx context.Context, tx *gorm.DB, data *entity.VerificationCode) (*entity.VerificationCode, error) {
	// Store in the database
	dbRes, dbErr := Datasource.NewVerificationCodeDatasource().CreateVerificationCode(tx, data)
	if dbErr != nil {
		return nil, dbErr
	}
	// Store in cache
	go func() {
		cacheId := "verification_code:" + dbRes.Code
		cacheErr := Rdb.HSet(context.Background(), cacheId, []string{
			"id", dbRes.ID.String(),
			"code", dbRes.Code,
			"email", dbRes.Email,
			"type", dbRes.Type,
			"device_identifier", dbRes.DeviceIdentifier,
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

func (v *verificationCodeRepository) GetVerificationCode(ctx context.Context, tx *gorm.DB, where *entity.VerificationCode, fields *[]string) (*entity.VerificationCode, error) {
	cacheId := "verification_code:" + where.Code
	cacheRes, cacheErr := Rdb.HGetAll(ctx, cacheId).Result()
	// Check if exists in cache
	if len(cacheRes) == 0 || cacheErr == redis.Nil {
		dbRes, dbErr := Datasource.NewVerificationCodeDatasource().GetVerificationCode(tx, where, nil)
		if dbErr != nil {
			return nil, dbErr
		}
		// Store in cache
		go func() {
			cacheId := "verification_code:" + dbRes.ID.String()
			cacheErr := Rdb.HSet(context.Background(), cacheId, []string{
				"id", dbRes.ID.String(),
				"code", dbRes.Code,
				"email", dbRes.Email,
				"type", dbRes.Type,
				"device_identifier", dbRes.DeviceIdentifier,
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
		return &entity.VerificationCode{
			ID:               &id,
			Code:             cacheRes["code"],
			Email:            cacheRes["email"],
			Type:             cacheRes["type"],
			DeviceIdentifier: cacheRes["device_identifier"],
			CreateTime:       createTime,
			UpdateTime:       updateTime,
		}, nil
	}
}

func (v *verificationCodeRepository) DeleteVerificationCode(ctx context.Context, tx *gorm.DB, where *entity.VerificationCode, ids *[]uuid.UUID) (*[]entity.VerificationCode, error) {
	res, err := Datasource.NewVerificationCodeDatasource().DeleteVerificationCode(tx, where, ids)
	if err != nil {
		return nil, err
	}
	return res, nil
}
