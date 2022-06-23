package repository

import (
	"context"
	"strconv"
	"time"

	"github.com/daniarmas/api_go/models"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ItemRepository interface {
	GetItem(tx *gorm.DB, where *models.Item, fields *[]string) (*models.Item, error)
	ListItem(tx *gorm.DB, where *models.Item, cursor time.Time, fields *[]string) (*[]models.Item, error)
	ListItemInIds(tx *gorm.DB, ids []uuid.UUID, fields *[]string) (*[]models.Item, error)
	CreateItem(tx *gorm.DB, data *models.Item) (*models.Item, error)
	SearchItem(tx *gorm.DB, name string, provinceId string, municipalityId string, cursor int64, municipalityNotEqual bool, limit int64, fields *[]string) (*[]models.Item, error)
	SearchItemByBusiness(tx *gorm.DB, name string, cursor int64, businessId string, fields *[]string) (*[]models.Item, error)
	UpdateItem(tx *gorm.DB, where *models.Item, data *models.Item) (*models.Item, error)
	UpdateItems(tx *gorm.DB, data *[]models.Item) (*[]models.Item, error)
	DeleteItem(tx *gorm.DB, where *models.Item) error
}

type itemRepository struct{}

func (v *itemRepository) DeleteItem(tx *gorm.DB, where *models.Item) error {
	err := Datasource.NewItemDatasource().DeleteItem(tx, where)
	if err != nil {
		return err
	} else {
		cacheId := "item:" + where.ID.String()
		ctx := context.Background()
		cacheRes, cacheErr := Rdb.HDel(ctx, cacheId).Result()
		if cacheRes == 0 || cacheErr == redis.Nil {
			log.Error(cacheErr)
		}
	}
	return nil
}

func (v *itemRepository) CreateItem(tx *gorm.DB, data *models.Item) (*models.Item, error) {
	res, err := Datasource.NewItemDatasource().CreateItem(tx, data)
	if err != nil {
		return nil, err
	}
	cacheId := "item:" + res.ID.String()
	ctx := context.Background()
	rdbPipe := Rdb.Pipeline()
	cacheErr := rdbPipe.HSet(ctx, cacheId, []string{
		"id", res.ID.String(),
		"name", res.Name,
		"description", res.Description,
		"thumbnail", res.Thumbnail,
		"high_quality_photo", res.HighQualityPhoto,
		"low_quality_photo", res.LowQualityPhoto,
		"blurhash", res.BlurHash,
		"price_cup", res.PriceCup,
		"cost_cup", res.CostCup,
		"profit_cup", res.ProfitCup,
		"price_usd", res.PriceUsd,
		"cost_usd", res.CostUsd,
		"profit_usd", res.ProfitUsd,
		"cursor", strconv.Itoa(int(res.Cursor)),
		"province_id", res.ProvinceId.String(),
		"municipality_id", res.MunicipalityId.String(),
		"enabled_flag", strconv.FormatBool(res.EnabledFlag),
		"available_flag", strconv.FormatBool(res.AvailableFlag),
		"availability", strconv.Itoa(int(res.Availability)),
		"business_id", res.BusinessId.String(),
		"business_collection_id", res.BusinessCollectionId.String(),
		"create_time", res.CreateTime.Format(time.RFC3339),
		"update_time", res.UpdateTime.Format(time.RFC3339),
	}).Err()
	if cacheErr != nil {
		log.Error(cacheErr)
	} else {
		rdbPipe.Expire(ctx, cacheId, time.Minute*5)
	}
	_, err = rdbPipe.Exec(ctx)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *itemRepository) ListItem(tx *gorm.DB, where *models.Item, cursor time.Time, fields *[]string) (*[]models.Item, error) {
	result, err := Datasource.NewItemDatasource().ListItem(tx, where, cursor, fields)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *itemRepository) ListItemInIds(tx *gorm.DB, ids []uuid.UUID, fields *[]string) (*[]models.Item, error) {
	result, err := Datasource.NewItemDatasource().ListItemInIds(tx, ids, fields)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *itemRepository) GetItem(tx *gorm.DB, where *models.Item, fields *[]string) (*models.Item, error) {
	cacheId := "item:" + where.ID.String()
	ctx := context.Background()
	cacheRes, cacheErr := Rdb.HGetAll(ctx, cacheId).Result()
	// Check if exists in cache
	if len(cacheRes) == 0 || cacheErr == redis.Nil {
		dbRes, dbErr := Datasource.NewItemDatasource().GetItem(tx, where, fields)
		if dbErr != nil {
			return nil, dbErr
		}
		ctx := context.Background()
		rdbPipe := Rdb.Pipeline()
		cacheErr := rdbPipe.HSet(ctx, cacheId, []string{
			"id", dbRes.ID.String(),
			"name", dbRes.Name,
			"description", dbRes.Description,
			"thumbnail", dbRes.Thumbnail,
			"high_quality_photo", dbRes.HighQualityPhoto,
			"low_quality_photo", dbRes.LowQualityPhoto,
			"blurhash", dbRes.BlurHash,
			"price_cup", dbRes.PriceCup,
			"cost_cup", dbRes.CostCup,
			"profit_cup", dbRes.ProfitCup,
			"price_usd", dbRes.PriceUsd,
			"cost_usd", dbRes.CostUsd,
			"profit_usd", dbRes.ProfitUsd,
			"cursor", strconv.Itoa(int(dbRes.Cursor)),
			"province_id", dbRes.ProvinceId.String(),
			"municipality_id", dbRes.MunicipalityId.String(),
			"enabled_flag", strconv.FormatBool(dbRes.EnabledFlag),
			"available_flag", strconv.FormatBool(dbRes.AvailableFlag),
			"availability", strconv.Itoa(int(dbRes.Availability)),
			"business_id", dbRes.BusinessId.String(),
			"business_collection_id", dbRes.BusinessCollectionId.String(),
			"create_time", dbRes.CreateTime.Format(time.RFC3339),
			"update_time", dbRes.UpdateTime.Format(time.RFC3339),
		}).Err()
		if cacheErr != nil {
			log.Error(cacheErr)
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
		createTime, _ := time.Parse(time.RFC3339, cacheRes["create_time"])
		updateTime, _ := time.Parse(time.RFC3339, cacheRes["update_time"])
		availableFlag, _ := strconv.ParseBool(cacheRes["available_flag"])
		enabledFlag, _ := strconv.ParseBool(cacheRes["enabled_flag"])
		availability, _ := strconv.ParseInt(cacheRes["availability"], 0, 64)
		cursor, _ := strconv.ParseInt(cacheRes["cursor"], 0, 32)
		businessId := uuid.MustParse(cacheRes["business_id"])
		businessCollectionId := uuid.MustParse(cacheRes["business_collection_id"])
		provinceId := uuid.MustParse(cacheRes["province_id"])
		municipalityId := uuid.MustParse(cacheRes["municipality_id"])
		return &models.Item{
			ID:                   &id,
			Name:                 cacheRes["name"],
			Description:          cacheRes["description"],
			PriceCup:             cacheRes["price_cup"],
			CostCup:              cacheRes["cost_cup"],
			ProfitCup:            cacheRes["profit_cup"],
			PriceUsd:             cacheRes["price_usd"],
			CostUsd:              cacheRes["cost_usd"],
			ProfitUsd:            cacheRes["profit_usd"],
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

func (i *itemRepository) UpdateItem(tx *gorm.DB, where *models.Item, data *models.Item) (*models.Item, error) {
	result, err := Datasource.NewItemDatasource().UpdateItem(tx, where, data)
	if err != nil {
		return nil, err
	}
	cacheId := "item:" + where.ID.String()
	ctx := context.Background()
	rdbPipe := Rdb.Pipeline()
	cacheErr := rdbPipe.HSet(ctx, cacheId, []string{
		"id", result.ID.String(),
		"name", result.Name,
		"description", result.Description,
		"thumbnail", result.Thumbnail,
		"high_quality_photo", result.HighQualityPhoto,
		"low_quality_photo", result.LowQualityPhoto,
		"blurhash", result.BlurHash,
		"price_cup", result.PriceCup,
		"cost_cup", result.CostCup,
		"profit_cup", result.ProfitCup,
		"price_usd", result.PriceUsd,
		"cost_usd", result.CostUsd,
		"profit_usd", result.ProfitUsd,
		"cursor", strconv.Itoa(int(result.Cursor)),
		"province_id", result.ProvinceId.String(),
		"municipality_id", result.MunicipalityId.String(),
		"enabled_flag", strconv.FormatBool(result.EnabledFlag),
		"available_flag", strconv.FormatBool(result.AvailableFlag),
		"availability", strconv.Itoa(int(result.Availability)),
		"business_id", result.BusinessId.String(),
		"business_collection_id", result.BusinessCollectionId.String(),
		"create_time", result.CreateTime.Format(time.RFC3339),
		"update_time", result.UpdateTime.Format(time.RFC3339),
	}).Err()
	if cacheErr != nil {
		log.Error(cacheErr)
	} else {
		rdbPipe.Expire(ctx, cacheId, time.Minute*5)
	}
	_, err = rdbPipe.Exec(ctx)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *itemRepository) UpdateItems(tx *gorm.DB, data *[]models.Item) (*[]models.Item, error) {
	result, err := Datasource.NewItemDatasource().UpdateItems(tx, data)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	rdbPipe := Rdb.Pipeline()
	for _, i := range *result {
		cacheId := "item:" + i.ID.String()
		cacheErr := rdbPipe.HSet(ctx, cacheId, []string{
			"id", i.ID.String(),
			"name", i.Name,
			"description", i.Description,
			"thumbnail", i.Thumbnail,
			"high_quality_photo", i.HighQualityPhoto,
			"low_quality_photo", i.LowQualityPhoto,
			"blurhash", i.BlurHash,
			"price_cup", i.PriceCup,
			"cost_cup", i.CostCup,
			"profit_cup", i.ProfitCup,
			"price_usd", i.PriceUsd,
			"cost_usd", i.CostUsd,
			"profit_usd", i.ProfitUsd,
			"cursor", strconv.Itoa(int(i.Cursor)),
			"province_id", i.ProvinceId.String(),
			"municipality_id", i.MunicipalityId.String(),
			"enabled_flag", strconv.FormatBool(i.EnabledFlag),
			"available_flag", strconv.FormatBool(i.AvailableFlag),
			"availability", strconv.Itoa(int(i.Availability)),
			"business_id", i.BusinessId.String(),
			"business_collection_id", i.BusinessCollectionId.String(),
			"create_time", i.CreateTime.Format(time.RFC3339),
			"update_time", i.UpdateTime.Format(time.RFC3339),
		}).Err()
		if cacheErr != nil {
			log.Error(cacheErr)
		} else {
			rdbPipe.Expire(ctx, cacheId, time.Minute*5)
		}
	}
	_, err = rdbPipe.Exec(ctx)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *itemRepository) SearchItem(tx *gorm.DB, name string, provinceId string, municipalityId string, cursor int64, municipalityNotEqual bool, limit int64, fields *[]string) (*[]models.Item, error) {
	result, err := Datasource.NewItemDatasource().SearchItem(tx, name, provinceId, municipalityId, cursor, municipalityNotEqual, limit, fields)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *itemRepository) SearchItemByBusiness(tx *gorm.DB, name string, cursor int64, businessId string, fields *[]string) (*[]models.Item, error) {
	result, err := Datasource.NewItemDatasource().SearchItemByBusiness(tx, name, cursor, businessId, fields)
	if err != nil {
		return nil, err
	}
	return result, nil
}
