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

type ItemRepository interface {
	GetItem(ctx context.Context, tx *gorm.DB, where *entity.ItemBusiness) (*entity.ItemBusiness, error)
	ListItem(ctx context.Context, tx *gorm.DB, where *entity.ItemBusiness, cursor time.Time) (*[]entity.ItemBusiness, error)
	ListItemInIds(ctx context.Context, tx *gorm.DB, ids []uuid.UUID) (*[]entity.ItemBusiness, error)
	CreateItem(ctx context.Context, tx *gorm.DB, data *entity.ItemBusiness) (*entity.ItemBusiness, error)
	SearchItem(ctx context.Context, tx *gorm.DB, name string, provinceId string, municipalityId string, cursor int64, municipalityNotEqual bool, limit int64) (*[]entity.ItemBusiness, error)
	SearchItemByBusiness(ctx context.Context, tx *gorm.DB, name string, cursor int64, businessId string) (*[]entity.ItemBusiness, error)
	UpdateItem(ctx context.Context, tx *gorm.DB, where *entity.ItemBusiness, data *entity.ItemBusiness) (*entity.ItemBusiness, error)
	UpdateItems(ctx context.Context, tx *gorm.DB, data *[]entity.ItemBusiness) (*[]entity.ItemBusiness, error)
	DeleteItem(ctx context.Context, tx *gorm.DB, where *entity.ItemBusiness, ids *[]uuid.UUID) (*[]entity.ItemBusiness, error)
}

type itemRepository struct{}

func (v *itemRepository) DeleteItem(ctx context.Context, tx *gorm.DB, where *entity.ItemBusiness, ids *[]uuid.UUID) (*[]entity.ItemBusiness, error) {
	// Delete in database
	dbRes, dbErr := Datasource.NewItemDatasource().DeleteItem(tx, where, ids)
	if dbErr != nil {
		return nil, dbErr
	} else {
		// Delete in cache
		rdbPipe := Rdb.Pipeline()
		for _, item := range *dbRes {
			cacheId := "item:" + item.ID.String()
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

func (v *itemRepository) CreateItem(ctx context.Context, tx *gorm.DB, data *entity.ItemBusiness) (*entity.ItemBusiness, error) {
	// Store in the database
	dbRes, dbErr := Datasource.NewItemDatasource().CreateItem(tx, data)
	if dbErr != nil {
		return nil, dbErr
	} else {
		// Store in cache
		go func() {
			ctx := context.Background()
			cacheId := "item:" + dbRes.ID.String()
			cacheErr := Rdb.HSet(ctx, cacheId, []string{
				"id", dbRes.ID.String(),
				"name", dbRes.Name,
				"business_name", dbRes.BusinessName,
				"description", dbRes.Description,
				"thumbnail", dbRes.Thumbnail,
				"high_quality_photo", dbRes.HighQualityPhoto,
				"low_quality_photo", dbRes.LowQualityPhoto,
				"blurhash", dbRes.BlurHash,
				"cost_cup", dbRes.CostCup,
				"price_cup", dbRes.PriceCup,
				"profit_cup", dbRes.ProfitCup,
				"cost_usd", dbRes.CostUsd,
				"price_usd", dbRes.PriceCup,
				"profit_usd", dbRes.ProfitUsd,
				"cursor", strconv.FormatInt(int64(dbRes.Cursor), 10),
				"enabled_flag", strconv.FormatBool(dbRes.EnabledFlag),
				"available_flag", strconv.FormatBool(dbRes.AvailableFlag),
				"availability", strconv.FormatInt(dbRes.Availability, 10),
				"province_id", dbRes.ProvinceId.String(),
				"municipality_id", dbRes.MunicipalityId.String(),
				"business_id", dbRes.BusinessId.String(),
				"business_collection_id", dbRes.BusinessCollectionId.String(),
				"create_time", dbRes.CreateTime.Format(time.RFC3339),
				"update_time", dbRes.UpdateTime.Format(time.RFC3339),
			}).Err()
			if cacheErr != nil {
				log.Error(cacheErr)
			} else {
				Rdb.Expire(ctx, cacheId, time.Second*15)
			}
		}()
	}
	return dbRes, nil
}

func (i *itemRepository) ListItem(ctx context.Context, tx *gorm.DB, where *entity.ItemBusiness, cursor time.Time) (*[]entity.ItemBusiness, error) {
	// Get from database
	dbRes, dbErr := Datasource.NewItemDatasource().ListItem(tx, where, cursor)
	if dbErr != nil {
		return nil, dbErr
	} else {
		// Delete in cache
		go func() {
			ctx := context.Background()
			rdbPipe := Rdb.Pipeline()
			for _, item := range *dbRes {
				cacheId := "item:" + item.ID.String()
				cacheErr := rdbPipe.HSet(ctx, cacheId, []string{
					"id", item.ID.String(),
					"name", item.Name,
					"business_name", item.BusinessName,
					"description", item.Description,
					"thumbnail", item.Thumbnail,
					"high_quality_photo", item.HighQualityPhoto,
					"low_quality_photo", item.LowQualityPhoto,
					"blurhash", item.BlurHash,
					"cost_cup", item.CostCup,
					"price_cup", item.PriceCup,
					"profit_cup", item.ProfitCup,
					"cost_usd", item.CostUsd,
					"price_usd", item.PriceCup,
					"profit_usd", item.ProfitUsd,
					"cursor", strconv.FormatInt(int64(item.Cursor), 10),
					"enabled_flag", strconv.FormatBool(item.EnabledFlag),
					"available_flag", strconv.FormatBool(item.AvailableFlag),
					"availability", strconv.FormatInt(item.Availability, 10),
					"province_id", item.ProvinceId.String(),
					"municipality_id", item.MunicipalityId.String(),
					"business_id", item.BusinessId.String(),
					"business_collection_id", item.BusinessCollectionId.String(),
					"create_time", item.CreateTime.Format(time.RFC3339),
					"update_time", item.UpdateTime.Format(time.RFC3339),
				}).Err()
				if cacheErr != nil {
					log.Error(cacheErr)
				} else {
					rdbPipe.Expire(ctx, cacheId, time.Second*15)
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

func (i *itemRepository) ListItemInIds(ctx context.Context, tx *gorm.DB, ids []uuid.UUID) (*[]entity.ItemBusiness, error) {
	// Delete in database
	dbRes, dbErr := Datasource.NewItemDatasource().ListItemInIds(tx, ids)
	if dbErr != nil {
		return nil, dbErr
	} else {
		// Delete in cache
		go func() {
			ctx := context.Background()
			rdbPipe := Rdb.Pipeline()
			for _, item := range *dbRes {
				cacheId := "item:" + item.ID.String()
				cacheErr := rdbPipe.HSet(ctx, cacheId, []string{
					"id", item.ID.String(),
					"name", item.Name,
					"business_name", item.BusinessName,
					"description", item.Description,
					"thumbnail", item.Thumbnail,
					"high_quality_photo", item.HighQualityPhoto,
					"low_quality_photo", item.LowQualityPhoto,
					"blurhash", item.BlurHash,
					"cost_cup", item.CostCup,
					"price_cup", item.PriceCup,
					"profit_cup", item.ProfitCup,
					"cost_usd", item.CostUsd,
					"price_usd", item.PriceCup,
					"profit_usd", item.ProfitUsd,
					"cursor", strconv.FormatInt(int64(item.Cursor), 10),
					"business_open_flag", strconv.FormatBool(item.BusinessOpenFlag),
					"enabled_flag", strconv.FormatBool(item.EnabledFlag),
					"available_flag", strconv.FormatBool(item.AvailableFlag),
					"availability", strconv.FormatInt(item.Availability, 10),
					"province_id", item.ProvinceId.String(),
					"municipality_id", item.MunicipalityId.String(),
					"business_id", item.BusinessId.String(),
					"business_collection_id", item.BusinessCollectionId.String(),
					"create_time", item.CreateTime.Format(time.RFC3339),
					"update_time", item.UpdateTime.Format(time.RFC3339),
				}).Err()
				if cacheErr != nil {
					log.Error(cacheErr)
				} else {
					rdbPipe.Expire(ctx, cacheId, time.Second*15)
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

func (i *itemRepository) GetItem(ctx context.Context, tx *gorm.DB, where *entity.ItemBusiness) (*entity.ItemBusiness, error) {
	cacheId := "item:" + where.ID.String()
	cacheRes, cacheErr := Rdb.HGetAll(ctx, cacheId).Result()
	// Check if exists in cache
	if len(cacheRes) == 0 || cacheErr == redis.Nil {
		dbRes, dbErr := Datasource.NewItemDatasource().GetItem(tx, where)
		if dbErr != nil {
			return nil, dbErr
		}
		// Store in cache
		go func() {
			ctx := context.Background()
			cacheErr := Rdb.HSet(ctx, cacheId, []string{
				"id", dbRes.ID.String(),
				"name", dbRes.Name,
				"business_name", dbRes.BusinessName,
				"description", dbRes.Description,
				"thumbnail", dbRes.Thumbnail,
				"high_quality_photo", dbRes.HighQualityPhoto,
				"low_quality_photo", dbRes.LowQualityPhoto,
				"blurhash", dbRes.BlurHash,
				"cost_cup", dbRes.CostCup,
				"price_cup", dbRes.PriceCup,
				"profit_cup", dbRes.ProfitCup,
				"cost_usd", dbRes.CostUsd,
				"price_usd", dbRes.PriceCup,
				"profit_usd", dbRes.ProfitUsd,
				"cursor", strconv.FormatInt(int64(dbRes.Cursor), 10),
				"business_open_flag", strconv.FormatBool(dbRes.BusinessOpenFlag),
				"enabled_flag", strconv.FormatBool(dbRes.EnabledFlag),
				"available_flag", strconv.FormatBool(dbRes.AvailableFlag),
				"availability", strconv.FormatInt(dbRes.Availability, 10),
				"province_id", dbRes.ProvinceId.String(),
				"municipality_id", dbRes.MunicipalityId.String(),
				"business_id", dbRes.BusinessId.String(),
				"business_collection_id", dbRes.BusinessCollectionId.String(),
				"create_time", dbRes.CreateTime.Format(time.RFC3339),
				"update_time", dbRes.UpdateTime.Format(time.RFC3339),
			}).Err()
			if cacheErr != nil {
				log.Error(cacheErr)
			} else {
				Rdb.Expire(ctx, cacheId, time.Second*15)
			}
		}()
		return dbRes, nil
	} else {
		id := uuid.MustParse(cacheRes["id"])
		businessId := uuid.MustParse(cacheRes["business_id"])
		businessCollectionId := uuid.MustParse(cacheRes["business_collection_id"])
		provinceId := uuid.MustParse(cacheRes["province_id"])
		municipalityId := uuid.MustParse(cacheRes["municipality_id"])
		availableFlag, _ := strconv.ParseBool(cacheRes["available_flag"])
		enabledFlag, _ := strconv.ParseBool(cacheRes["enabled_flag"])
		openFlag, _ := strconv.ParseBool(cacheRes["business_open_flag"])
		createTime, _ := time.Parse(time.RFC3339, cacheRes["create_time"])
		updateTime, _ := time.Parse(time.RFC3339, cacheRes["update_time"])
		availability, _ := strconv.ParseInt(cacheRes["availability"], 10, 0)
		cursor, _ := strconv.ParseInt(cacheRes["cursor"], 10, 0)
		return &entity.ItemBusiness{
			ID:                   &id,
			Name:                 cacheRes["name"],
			BusinessName:         cacheRes["business_name"],
			Description:          cacheRes["description"],
			PriceCup:             cacheRes["price_cup"],
			CostCup:              cacheRes["cost_cup"],
			ProfitCup:            cacheRes["profit_cup"],
			PriceUsd:             cacheRes["price_usd"],
			CostUsd:              cacheRes["cost_usd"],
			ProfitUsd:            cacheRes["profit_usd"],
			BusinessOpenFlag:     openFlag,
			AvailableFlag:        availableFlag,
			EnabledFlag:          enabledFlag,
			Availability:         availability,
			BusinessId:           &businessId,
			BusinessCollectionId: &businessCollectionId,
			ProvinceId:           &provinceId,
			MunicipalityId:       &municipalityId,
			HighQualityPhoto:     cacheRes["high_quality_photo"],
			LowQualityPhoto:      cacheRes["low_quality_photo"],
			Thumbnail:            cacheRes["thumbnail"],
			BlurHash:             cacheRes["blurhash"],
			Cursor:               int32(cursor),
			CreateTime:           createTime,
			UpdateTime:           updateTime,
		}, nil
	}
}

func (i *itemRepository) UpdateItem(ctx context.Context, tx *gorm.DB, where *entity.ItemBusiness, data *entity.ItemBusiness) (*entity.ItemBusiness, error) {
	// Store in the database
	dbRes, dbErr := Datasource.NewItemDatasource().UpdateItem(tx, where, data)
	if dbErr != nil {
		return nil, dbErr
	} else {
		// Store in cache
		cacheId := "item:" + dbRes.ID.String()
		cacheErr := Rdb.HSet(ctx, cacheId, []string{
			"id", dbRes.ID.String(),
			"name", dbRes.Name,
			"business_name", dbRes.BusinessName,
			"description", dbRes.Description,
			"thumbnail", dbRes.Thumbnail,
			"high_quality_photo", dbRes.HighQualityPhoto,
			"low_quality_photo", dbRes.LowQualityPhoto,
			"blurhash", dbRes.BlurHash,
			"cost_cup", dbRes.CostCup,
			"price_cup", dbRes.PriceCup,
			"profit_cup", dbRes.ProfitCup,
			"cost_usd", dbRes.CostUsd,
			"price_usd", dbRes.PriceCup,
			"profit_usd", dbRes.ProfitUsd,
			"cursor", strconv.FormatInt(int64(dbRes.Cursor), 10),
			"enabled_flag", strconv.FormatBool(dbRes.EnabledFlag),
			"available_flag", strconv.FormatBool(dbRes.AvailableFlag),
			"availability", strconv.FormatInt(dbRes.Availability, 10),
			"province_id", dbRes.ProvinceId.String(),
			"municipality_id", dbRes.MunicipalityId.String(),
			"business_id", dbRes.BusinessId.String(),
			"business_collection_id", dbRes.BusinessCollectionId.String(),
			"create_time", dbRes.CreateTime.Format(time.RFC3339),
			"update_time", dbRes.UpdateTime.Format(time.RFC3339),
		}).Err()
		if cacheErr != nil {
			log.Error(cacheErr)
		} else {
			Rdb.Expire(ctx, cacheId, time.Second*15)
		}
	}
	return dbRes, nil
}

func (i *itemRepository) UpdateItems(ctx context.Context, tx *gorm.DB, data *[]entity.ItemBusiness) (*[]entity.ItemBusiness, error) {
	// Update in database
	dbRes, dbErr := Datasource.NewItemDatasource().UpdateItems(tx, data)
	if dbErr != nil {
		return nil, dbErr
	} else {
		// Update in cache
		rdbPipe := Rdb.Pipeline()
		for _, item := range *dbRes {
			cacheId := "item:" + item.ID.String()
			cacheErr := rdbPipe.HSet(ctx, cacheId, []string{
				"id", item.ID.String(),
				"name", item.Name,
				"business_name", item.BusinessName,
				"description", item.Description,
				"thumbnail", item.Thumbnail,
				"high_quality_photo", item.HighQualityPhoto,
				"low_quality_photo", item.LowQualityPhoto,
				"blurhash", item.BlurHash,
				"cost_cup", item.CostCup,
				"price_cup", item.PriceCup,
				"profit_cup", item.ProfitCup,
				"cost_usd", item.CostUsd,
				"price_usd", item.PriceCup,
				"profit_usd", item.ProfitUsd,
				"cursor", strconv.FormatInt(int64(item.Cursor), 10),
				"enabled_flag", strconv.FormatBool(item.EnabledFlag),
				"available_flag", strconv.FormatBool(item.AvailableFlag),
				"availability", strconv.FormatInt(item.Availability, 10),
				"province_id", item.ProvinceId.String(),
				"municipality_id", item.MunicipalityId.String(),
				"business_id", item.BusinessId.String(),
				"business_collection_id", item.BusinessCollectionId.String(),
				"create_time", item.CreateTime.Format(time.RFC3339),
				"update_time", item.UpdateTime.Format(time.RFC3339),
			}).Err()
			if cacheErr != nil {
				log.Error(cacheErr)
			} else {
				rdbPipe.Expire(ctx, cacheId, time.Second*15)
			}
		}
		_, err := rdbPipe.Exec(ctx)
		if err != nil {
			log.Error(err)
		}
	}
	return dbRes, nil
}

func (i *itemRepository) SearchItem(ctx context.Context, tx *gorm.DB, name string, provinceId string, municipalityId string, cursor int64, municipalityNotEqual bool, limit int64) (*[]entity.ItemBusiness, error) {
	res, err := Datasource.NewItemDatasource().SearchItem(tx, name, provinceId, municipalityId, cursor, municipalityNotEqual, limit)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *itemRepository) SearchItemByBusiness(ctx context.Context, tx *gorm.DB, name string, cursor int64, businessId string) (*[]entity.ItemBusiness, error) {
	res, err := Datasource.NewItemDatasource().SearchItemByBusiness(tx, name, cursor, businessId)
	if err != nil {
		return nil, err
	}
	return res, nil
}
