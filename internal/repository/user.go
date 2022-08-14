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

type UserRepository interface {
	GetUser(ctx context.Context, tx *gorm.DB, where *entity.User) (*entity.User, error)
	GetUserWithAddress(ctx context.Context, tx *gorm.DB, where *entity.User) (*entity.User, error)
	CreateUser(ctx context.Context, tx *gorm.DB, data *entity.User) (*entity.User, error)
	UpdateUser(ctx context.Context, tx *gorm.DB, where *entity.User, data *entity.User) (*entity.User, error)
}

type userRepository struct{}

func (u *userRepository) GetUserWithAddress(ctx context.Context, tx *gorm.DB, where *entity.User) (*entity.User, error) {
	res, err := Datasource.NewUserDatasource().GetUserWithAddress(tx, where)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (u *userRepository) GetUser(ctx context.Context, tx *gorm.DB, where *entity.User) (*entity.User, error) {
	var cacheId string
	if where.ID != nil {
		cacheId = "user:" + where.ID.String()
	} else {
		cacheId = "user:" + where.Email
	}
	cacheRes, cacheErr := Rdb.HGetAll(ctx, cacheId).Result()
	// Check if exists in cache
	if len(cacheRes) == 0 || cacheErr == redis.Nil {
		dbRes, dbErr := Datasource.NewUserDatasource().GetUser(tx, where)
		if dbErr != nil {
			return nil, dbErr
		}
		// Store in cache
		go func() {
			ctx := context.Background()
			cacheErr := Rdb.HSet(ctx, cacheId, []string{
				"id", dbRes.ID.String(),
				"email", dbRes.Email,
				"fullname", dbRes.FullName,
				"is_legal_age", strconv.FormatBool(dbRes.IsLegalAge),
				"thumbnail", dbRes.Thumbnail,
				"high_quality_photo", dbRes.HighQualityPhoto,
				"low_quality_photo", dbRes.LowQualityPhoto,
				"blurhash", dbRes.BlurHash,
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
		isLegalAge, _ := strconv.ParseBool(cacheRes["is_legal_age"])
		createTime, _ := time.Parse(time.RFC3339, cacheRes["create_time"])
		updateTime, _ := time.Parse(time.RFC3339, cacheRes["update_time"])
		return &entity.User{
			ID:               &id,
			FullName:         cacheRes["fullname"],
			Email:            cacheRes["email"],
			IsLegalAge:       isLegalAge,
			HighQualityPhoto: cacheRes["high_quality_photo"],
			LowQualityPhoto:  cacheRes["low_quality_photo"],
			Thumbnail:        cacheRes["thumbnail"],
			BlurHash:         cacheRes["blurhash"],
			CreateTime:       createTime,
			UpdateTime:       updateTime,
		}, nil
	}
}

func (u *userRepository) CreateUser(ctx context.Context, tx *gorm.DB, data *entity.User) (*entity.User, error) {
	// Store in the database
	dbRes, dbErr := Datasource.NewUserDatasource().CreateUser(tx, data)
	if dbErr != nil {
		return nil, dbErr
	} else {
		// Store in cache
		go func() {
			ctx := context.Background()
			cacheId := "user:" + dbRes.ID.String()
			cacheErr := Rdb.HSet(ctx, cacheId, []string{
				"id", dbRes.ID.String(),
				"email", dbRes.Email,
				"fullname", dbRes.FullName,
				"is_legal_age", strconv.FormatBool(dbRes.IsLegalAge),
				"thumbnail", dbRes.Thumbnail,
				"high_quality_photo", dbRes.HighQualityPhoto,
				"low_quality_photo", dbRes.LowQualityPhoto,
				"blurhash", dbRes.BlurHash,
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

func (u *userRepository) UpdateUser(ctx context.Context, tx *gorm.DB, where *entity.User, data *entity.User) (*entity.User, error) {
	// Store in the database
	dbRes, dbErr := Datasource.NewUserDatasource().UpdateUser(tx, where, data)
	if dbErr != nil {
		return nil, dbErr
	} else {
		// Store in cache
		cacheId := "user:" + dbRes.ID.String()
		cacheErr := Rdb.HSet(ctx, cacheId, []string{
			"id", dbRes.ID.String(),
			"email", dbRes.Email,
			"fullname", dbRes.FullName,
			"is_legal_age", strconv.FormatBool(dbRes.IsLegalAge),
			"thumbnail", dbRes.Thumbnail,
			"high_quality_photo", dbRes.HighQualityPhoto,
			"low_quality_photo", dbRes.LowQualityPhoto,
			"blurhash", dbRes.BlurHash,
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
