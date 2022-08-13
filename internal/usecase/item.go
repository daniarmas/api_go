package usecase

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/daniarmas/api_go/config"
	"github.com/daniarmas/api_go/internal/datasource"
	"github.com/daniarmas/api_go/internal/entity"
	"github.com/daniarmas/api_go/internal/repository"
	pb "github.com/daniarmas/api_go/pkg/grpc"
	"github.com/daniarmas/api_go/pkg/sqldb"
	"github.com/daniarmas/api_go/utils"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type ItemService interface {
	GetItem(ctx context.Context, req *pb.GetItemRequest, md *utils.ClientMetadata) (*pb.Item, error)
	ListItem(ctx context.Context, req *pb.ListItemRequest, md *utils.ClientMetadata) (*pb.ListItemResponse, error)
	SearchItem(ctx context.Context, req *pb.SearchItemRequest, md *utils.ClientMetadata) (*pb.SearchItemResponse, error)
	SearchItemByBusiness(ctx context.Context, req *pb.SearchItemByBusinessRequest, md *utils.ClientMetadata) (*pb.SearchItemByBusinessResponse, error)
	CreateItem(ctx context.Context, req *pb.CreateItemRequest, md *utils.ClientMetadata) (*pb.Item, error)
	UpdateItem(ctx context.Context, req *pb.UpdateItemRequest, md *utils.ClientMetadata) (*pb.Item, error)
	DeleteItem(ctx context.Context, req *pb.DeleteItemRequest, md *utils.ClientMetadata) error
}

type itemService struct {
	dao    repository.Repository
	config *config.Config
	stDb   *sql.DB
	sqldb  *sqldb.Sql
}

func NewItemService(dao repository.Repository, config *config.Config, stDb *sql.DB, sqldb *sqldb.Sql) ItemService {
	return &itemService{dao: dao, config: config, stDb: stDb, sqldb: sqldb}
}

func (i *itemService) UpdateItem(ctx context.Context, req *pb.UpdateItemRequest, md *utils.ClientMetadata) (*pb.Item, error) {
	var res pb.Item
	err := i.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		_, err := i.dao.NewApplicationRepository().CheckApplication(ctx, tx, *md.AccessToken)
		if err != nil {
			return err
		}
		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: md.Authorization}
		err = repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(jwtAuthorizationToken)
		if err != nil {
			switch err.Error() {
			case "Token is expired":
				return errors.New("authorization token expired")
			case "signature is invalid":
				return errors.New("authorization token signature is invalid")
			case "token contains an invalid number of segments":
				return errors.New("authorization token contains an invalid number of segments")
			default:
				return err
			}
		}
		authorizationTokenRes, err := i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: jwtAuthorizationToken.TokenId})
		if err != nil && err.Error() == "record not found" {
			return errors.New("unauthenticated user")
		} else if err != nil {
			return err
		}
		id := uuid.MustParse(req.Item.Id)
		getItemRes, err := i.dao.NewItemRepository().GetItem(ctx, tx, &entity.Item{ID: &id})
		if err != nil && err.Error() == "record not found" {
			return errors.New("item not found")
		} else if err != nil {
			return err
		}
		_, err = i.dao.NewUserPermissionRepository().GetUserPermission(ctx, tx, &entity.UserPermission{UserId: authorizationTokenRes.UserId, Name: "update_item", BusinessId: getItemRes.BusinessId})
		if err != nil && err.Error() == "record not found" {
			return errors.New("permission denied")
		}
		getCartItemRes, err := i.dao.NewCartItemRepository().GetCartItem(tx, &entity.CartItem{ItemId: &id})
		if err != nil && err.Error() != "record not found" {
			return err
		} else if getCartItemRes != nil {
			return errors.New("item in the cart")
		}
		if req.Item.HighQualityPhoto != "" || req.Item.LowQualityPhoto != "" || req.Item.Thumbnail != "" {
			_, err := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), i.config.ItemsBulkName, req.Item.HighQualityPhoto)
			if err != nil && err.Error() == "ObjectMissing" {
				return errors.New("HighQualityPhotoObject missing")
			} else if err != nil {
				return err
			}
			_, err = i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), i.config.ItemsBulkName, req.Item.LowQualityPhoto)
			if err != nil && err.Error() == "ObjectMissing" {
				return errors.New("LowQualityPhotoObject missing")
			} else if err != nil {
				return err
			}
			_, err = i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), i.config.ItemsBulkName, req.Item.Thumbnail)
			if err != nil && err.Error() == "ObjectMissing" {
				return errors.New("ThumbnailObject missing")
			} else if err != nil {
				return err
			}
			_, err = repository.Datasource.NewObjectStorageDatasource().CopyObject(context.Background(), minio.CopyDestOptions{Bucket: i.config.ItemsDeletedBulkName, Object: getItemRes.HighQualityPhoto}, minio.CopySrcOptions{Bucket: i.config.ItemsBulkName, Object: getItemRes.HighQualityPhoto})
			if err != nil {
				return err
			}
			_, err = repository.Datasource.NewObjectStorageDatasource().CopyObject(context.Background(), minio.CopyDestOptions{Bucket: i.config.ItemsDeletedBulkName, Object: getItemRes.LowQualityPhoto}, minio.CopySrcOptions{Bucket: i.config.ItemsBulkName, Object: getItemRes.LowQualityPhoto})
			if err != nil {
				return err
			}
			_, err = repository.Datasource.NewObjectStorageDatasource().CopyObject(context.Background(), minio.CopyDestOptions{Bucket: i.config.ItemsDeletedBulkName, Object: getItemRes.Thumbnail}, minio.CopySrcOptions{Bucket: i.config.ItemsBulkName, Object: getItemRes.Thumbnail})
			if err != nil {
				return err
			}
			err = repository.Datasource.NewObjectStorageDatasource().RemoveObject(context.Background(), i.config.ItemsBulkName, getItemRes.HighQualityPhoto, minio.RemoveObjectOptions{})
			if err != nil {
				return err
			}
			err = repository.Datasource.NewObjectStorageDatasource().RemoveObject(context.Background(), i.config.ItemsBulkName, getItemRes.LowQualityPhoto, minio.RemoveObjectOptions{})
			if err != nil {
				return err
			}
			err = repository.Datasource.NewObjectStorageDatasource().RemoveObject(context.Background(), i.config.ItemsBulkName, getItemRes.Thumbnail, minio.RemoveObjectOptions{})
			if err != nil {
				return err
			}
		}
		updateItemRes, err := i.dao.NewItemRepository().UpdateItem(ctx, tx, &entity.Item{ID: &id}, &entity.Item{Name: req.Item.Name, Description: req.Item.Description, PriceCup: req.Item.PriceCup, Availability: int64(req.Item.Availability), HighQualityPhoto: req.Item.HighQualityPhoto, LowQualityPhoto: req.Item.LowQualityPhoto, Thumbnail: req.Item.Thumbnail, BlurHash: req.Item.BlurHash})
		if err != nil {
			return err
		}
		res = pb.Item{
			Id:                   updateItemRes.ID.String(),
			Name:                 updateItemRes.Name,
			Description:          updateItemRes.Description,
			PriceCup:             updateItemRes.PriceCup,
			Availability:         int32(updateItemRes.Availability),
			BusinessId:           updateItemRes.BusinessId.String(),
			BusinessCollectionId: updateItemRes.BusinessCollectionId.String(),
			HighQualityPhoto:     updateItemRes.HighQualityPhoto,
			LowQualityPhoto:      updateItemRes.LowQualityPhoto,
			Thumbnail:            updateItemRes.Thumbnail,
			BlurHash:             updateItemRes.BlurHash,
			Cursor:               updateItemRes.Cursor,
			BusinessName:         updateItemRes.BusinessName,
			CreateTime:           timestamppb.New(updateItemRes.CreateTime),
			UpdateTime:           timestamppb.New(updateItemRes.UpdateTime),
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (i *itemService) DeleteItem(ctx context.Context, req *pb.DeleteItemRequest, md *utils.ClientMetadata) error {
	err := i.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		_, err := i.dao.NewApplicationRepository().CheckApplication(ctx, tx, *md.AccessToken)
		if err != nil {
			return err
		}
		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: md.Authorization}
		err = repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(jwtAuthorizationToken)
		if err != nil {
			switch err.Error() {
			case "Token is expired":
				return errors.New("authorization token expired")
			case "signature is invalid":
				return errors.New("authorization token signature is invalid")
			case "token contains an invalid number of segments":
				return errors.New("authorization token contains an invalid number of segments")
			default:
				return err
			}
		}
		authorizationTokenRes, err := i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: jwtAuthorizationToken.TokenId})
		if err != nil && err.Error() == "record not found" {
			return errors.New("unauthenticated")
		} else if err != nil {
			return err
		}
		id := uuid.MustParse(req.Id)
		getItemRes, err := i.dao.NewItemRepository().GetItem(ctx, tx, &entity.Item{ID: &id})
		if err != nil && err.Error() == "record not found" {
			return errors.New("item not found")
		} else if err != nil {
			return err
		}
		_, err = i.dao.NewUserPermissionRepository().GetUserPermission(ctx, tx, &entity.UserPermission{Name: "delete_item", UserId: authorizationTokenRes.UserId, BusinessId: getItemRes.BusinessId})
		if err != nil && err.Error() == "record not found" {
			return errors.New("permission denied")
		} else if err != nil {
			return err
		}
		businessIsOpenRes, err := i.dao.NewBusinessScheduleRepository().BusinessIsOpen(tx, &entity.BusinessSchedule{BusinessId: getItemRes.BusinessId})
		if err != nil && err.Error() != "business closed" {
			return err
		} else if businessIsOpenRes {
			return errors.New("business is open")
		}
		getCartItemRes, err := i.dao.NewCartItemRepository().GetCartItem(tx, &entity.CartItem{ItemId: &id})
		if err != nil && err.Error() != "record not found" {
			return err
		} else if getCartItemRes != nil {
			return errors.New("item in the cart")
		}
		_, err = i.dao.NewItemRepository().DeleteItem(ctx, tx, &entity.Item{ID: &id}, nil)
		if err != nil {
			return err
		}
		_, err = repository.Datasource.NewObjectStorageDatasource().CopyObject(context.Background(), minio.CopyDestOptions{Bucket: i.config.ItemsDeletedBulkName, Object: getItemRes.HighQualityPhoto}, minio.CopySrcOptions{Bucket: i.config.ItemsBulkName, Object: getItemRes.HighQualityPhoto})
		if err != nil {
			return err
		}
		_, err = repository.Datasource.NewObjectStorageDatasource().CopyObject(context.Background(), minio.CopyDestOptions{Bucket: i.config.ItemsDeletedBulkName, Object: getItemRes.LowQualityPhoto}, minio.CopySrcOptions{Bucket: i.config.ItemsBulkName, Object: getItemRes.LowQualityPhoto})
		if err != nil {
			return err
		}
		_, err = repository.Datasource.NewObjectStorageDatasource().CopyObject(context.Background(), minio.CopyDestOptions{Bucket: i.config.ItemsDeletedBulkName, Object: getItemRes.Thumbnail}, minio.CopySrcOptions{Bucket: i.config.ItemsBulkName, Object: getItemRes.Thumbnail})
		if err != nil {
			return err
		}
		err = repository.Datasource.NewObjectStorageDatasource().RemoveObject(context.Background(), i.config.ItemsBulkName, getItemRes.HighQualityPhoto, minio.RemoveObjectOptions{})
		if err != nil {
			return err
		}
		err = repository.Datasource.NewObjectStorageDatasource().RemoveObject(context.Background(), i.config.ItemsBulkName, getItemRes.LowQualityPhoto, minio.RemoveObjectOptions{})
		if err != nil {
			return err
		}
		err = repository.Datasource.NewObjectStorageDatasource().RemoveObject(context.Background(), i.config.ItemsBulkName, getItemRes.Thumbnail, minio.RemoveObjectOptions{})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (i *itemService) CreateItem(ctx context.Context, req *pb.CreateItemRequest, md *utils.ClientMetadata) (*pb.Item, error) {
	var itemRes *entity.Item
	var itemErr error
	err := i.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		_, err := i.dao.NewApplicationRepository().CheckApplication(ctx, tx, *md.AccessToken)
		if err != nil {
			return err
		}
		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: md.Authorization}
		authorizationTokenParseErr := repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(jwtAuthorizationToken)
		if authorizationTokenParseErr != nil {
			switch authorizationTokenParseErr.Error() {
			case "Token is expired":
				return errors.New("authorization token expired")
			case "signature is invalid":
				return errors.New("authorization token signature is invalid")
			case "token contains an invalid number of segments":
				return errors.New("authorization token contains an invalid number of segments")
			default:
				return authorizationTokenParseErr
			}
		}
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &entity.AuthorizationToken{ID: jwtAuthorizationToken.TokenId})
		if authorizationTokenErr != nil && authorizationTokenErr.Error() == "record not found" {
			return errors.New("unauthenticated")
		} else if authorizationTokenErr != nil {
			return authorizationTokenErr
		}
		businessId := uuid.MustParse(req.Item.BusinessId)
		businessRes, businessErr := i.dao.NewBusinessRepository().GetBusiness(tx, &entity.Business{ID: &businessId})
		if businessErr != nil {
			return businessErr
		}
		_, err = i.dao.NewUserPermissionRepository().GetUserPermission(ctx, tx, &entity.UserPermission{Name: "create_item", UserId: authorizationTokenRes.UserId, BusinessId: &businessId})
		if err != nil && err.Error() == "record not found" {
			return errors.New("permission denied")
		} else if err != nil {
			return err
		}
		_, hqErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), i.config.ItemsBulkName, req.Item.HighQualityPhoto)
		if hqErr != nil && hqErr.Error() == "ObjectMissing" {
			return errors.New("highQualityPhotoObject missing")
		} else if hqErr != nil {
			return hqErr
		}
		_, lqErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), i.config.ItemsBulkName, req.Item.LowQualityPhoto)
		if lqErr != nil && lqErr.Error() == "ObjectMissing" {
			return errors.New("lowQualityPhotoObject missing")
		} else if lqErr != nil {
			return lqErr
		}
		_, tnErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), i.config.ItemsBulkName, req.Item.Thumbnail)
		if tnErr != nil && tnErr.Error() == "ObjectMissing" {
			return errors.New("thumbnailObject missing")
		} else if tnErr != nil {
			return tnErr
		}
		businessCollectionId := uuid.MustParse(req.Item.BusinessCollectionId)
		itemRes, itemErr = i.dao.NewItemRepository().CreateItem(ctx, tx, &entity.Item{Name: req.Item.Name, Description: req.Item.Description, PriceCup: req.Item.PriceCup, CostCup: req.Item.CostCup, ProfitCup: req.Item.ProfitCup, PriceUsd: req.Item.PriceUsd, ProfitUsd: req.Item.ProfitUsd, CostUsd: req.Item.CostUsd, Availability: int64(req.Item.Availability), BusinessId: &businessId, BusinessCollectionId: &businessCollectionId, HighQualityPhoto: req.Item.HighQualityPhoto, LowQualityPhoto: req.Item.LowQualityPhoto, Thumbnail: req.Item.Thumbnail, BlurHash: req.Item.BlurHash, ProvinceId: businessRes.ProvinceId, MunicipalityId: businessRes.MunicipalityId, AvailableFlag: req.Item.AvailableFlag, EnabledFlag: req.Item.EnabledFlag})
		if itemErr != nil {
			return itemErr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &pb.Item{Id: itemRes.ID.String(), BusinessName: itemRes.BusinessName, Name: itemRes.Name, Description: itemRes.Description, PriceCup: itemRes.PriceCup, CostCup: itemRes.CostCup, ProfitCup: itemRes.ProfitCup, PriceUsd: itemRes.PriceUsd, CostUsd: itemRes.CostUsd, ProfitUsd: itemRes.ProfitUsd, EnabledFlag: itemRes.EnabledFlag, AvailableFlag: itemRes.AvailableFlag, Availability: int32(itemRes.Availability), BusinessId: itemRes.BusinessId.String(), BusinessCollectionId: itemRes.BusinessCollectionId.String(), HighQualityPhoto: itemRes.HighQualityPhoto, LowQualityPhoto: itemRes.LowQualityPhoto, Thumbnail: itemRes.Thumbnail, BlurHash: itemRes.BlurHash, ProvinceId: itemRes.ProvinceId.String(), MunicipalityId: itemRes.MunicipalityId.String(), ThumbnailUrl: i.config.ItemsBulkName + "/" + itemRes.Thumbnail, HighQualityPhotoUrl: i.config.ItemsBulkName + "/" + itemRes.HighQualityPhoto, LowQualityPhotoUrl: i.config.ItemsBulkName + "/" + itemRes.LowQualityPhoto, CreateTime: timestamppb.New(itemRes.CreateTime), UpdateTime: timestamppb.New(itemRes.UpdateTime)}, nil
}

func (i *itemService) ListItem(ctx context.Context, req *pb.ListItemRequest, md *utils.ClientMetadata) (*pb.ListItemResponse, error) {
	where := entity.Item{}
	var nextPage time.Time
	if req.NextPage == nil || (req.NextPage.Nanos == 0 && req.NextPage.Seconds == 0) {
		nextPage = time.Now()
	} else {
		nextPage = req.NextPage.AsTime()
	}
	var items *[]entity.Item
	var res pb.ListItemResponse
	var itemsErr error
	var businessCollectionId, businessId uuid.UUID
	if req.BusinessCollectionId != "" {
		businessCollectionId = uuid.MustParse(req.BusinessCollectionId)
		where.BusinessCollectionId = &businessCollectionId
	}
	if req.BusinessId != "" {
		businessId = uuid.MustParse(req.BusinessId)
		where.BusinessId = &businessId
	}
	err := i.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		_, err := i.dao.NewApplicationRepository().CheckApplication(ctx, tx, *md.AccessToken)
		if err != nil {
			return err
		}
		items, itemsErr = i.dao.NewItemRepository().ListItem(ctx, tx, &where, nextPage)
		if itemsErr != nil {
			return itemsErr
		} else if len(*items) > 10 {
			*items = (*items)[:len(*items)-1]
			res.NextPage = timestamppb.New((*items)[len(*items)-1].CreateTime)
		} else if len(*items) == 0 {
			res.NextPage = timestamppb.New(nextPage)
		} else {
			res.NextPage = timestamppb.New((*items)[len(*items)-1].CreateTime)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if items != nil {
		itemsResponse := make([]*pb.Item, 0, len(*items))
		for _, item := range *items {
			itemsResponse = append(itemsResponse, &pb.Item{
				Id:                   item.ID.String(),
				Name:                 item.Name,
				BusinessName:         item.BusinessName,
				Description:          item.Description,
				PriceCup:             item.PriceCup,
				Availability:         int32(item.Availability),
				BusinessId:           item.BusinessId.String(),
				BusinessCollectionId: item.BusinessCollectionId.String(),
				HighQualityPhoto:     item.HighQualityPhoto,
				HighQualityPhotoUrl:  i.config.ItemsBulkName + "/" + item.HighQualityPhoto,
				LowQualityPhoto:      item.LowQualityPhoto,
				LowQualityPhotoUrl:   i.config.ItemsBulkName + "/" + item.LowQualityPhoto,
				Thumbnail:            item.Thumbnail,
				ThumbnailUrl:         i.config.ItemsBulkName + "/" + item.Thumbnail,
				BlurHash:             item.BlurHash,
				Cursor:               int32(item.Cursor),
				CreateTime:           timestamppb.New(item.CreateTime),
				UpdateTime:           timestamppb.New(item.UpdateTime),
			})
		}
		res.Items = itemsResponse
	}
	return &res, nil
}

func (i *itemService) GetItem(ctx context.Context, req *pb.GetItemRequest, md *utils.ClientMetadata) (*pb.Item, error) {
	var res pb.Item
	var item *entity.Item
	var itemErr error
	err := i.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		_, err := i.dao.NewApplicationRepository().CheckApplication(ctx, tx, *md.AccessToken)
		if err != nil {
			return err
		}
		id := uuid.MustParse(req.Id)
		item, itemErr = i.dao.NewItemRepository().GetItem(ctx, tx, &entity.Item{ID: &id})
		if itemErr != nil {
			return itemErr
		}
		isInRangeRes, isInRangeErr := i.dao.NewBusinessRepository().BusinessIsInRange(tx, ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.Location.Latitude, req.Location.Longitude}).SetSRID(4326)}, item.BusinessId)
		if isInRangeErr != nil {
			return isInRangeErr
		}
		res = pb.Item{
			Id:                   item.ID.String(),
			Name:                 item.Name,
			BusinessName:         item.BusinessName,
			Description:          item.Description,
			PriceCup:             item.PriceCup,
			Availability:         int32(item.Availability),
			BusinessId:           item.BusinessId.String(),
			BusinessCollectionId: item.BusinessCollectionId.String(),
			HighQualityPhoto:     item.HighQualityPhoto,
			HighQualityPhotoUrl:  i.config.ItemsBulkName + "/" + item.HighQualityPhoto,
			LowQualityPhoto:      item.LowQualityPhoto,
			LowQualityPhotoUrl:   i.config.ItemsBulkName + "/" + item.LowQualityPhoto,
			Thumbnail:            item.Thumbnail,
			ThumbnailUrl:         i.config.ItemsBulkName + "/" + item.Thumbnail,
			BlurHash:             item.BlurHash,
			Cursor:               item.Cursor,
			IsInRange:            *isInRangeRes,
			CreateTime:           timestamppb.New(item.CreateTime),
			UpdateTime:           timestamppb.New(item.UpdateTime),
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (i *itemService) SearchItem(ctx context.Context, req *pb.SearchItemRequest, md *utils.ClientMetadata) (*pb.SearchItemResponse, error) {
	var response *[]entity.Item
	var searchItemResponse pb.SearchItemResponse
	var responseErr error
	err := i.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		_, err := i.dao.NewApplicationRepository().CheckApplication(ctx, tx, *md.AccessToken)
		if err != nil {
			return err
		}
		if req.SearchMunicipalityType.String() == "More" {
			response, responseErr = i.dao.NewItemRepository().SearchItem(ctx, tx, req.Name, req.ProvinceId, req.MunicipalityId, int64(req.NextPage), false, 10)
			if responseErr != nil {
				return responseErr
			}
			if len(*response) > 10 {
				*response = (*response)[:len(*response)-1]
				searchItemResponse.NextPage = int32((*response)[len(*response)-1].Cursor)
				searchItemResponse.SearchMunicipalityType = pb.SearchMunicipalityType_More
			} else if len(*response) <= 10 && len(*response) != 0 {
				length := 10 - len(*response)
				responseAdd, responseErr := i.dao.NewItemRepository().SearchItem(ctx, tx, req.Name, req.ProvinceId, req.MunicipalityId, int64(req.NextPage), true, int64(length))
				if responseErr != nil {
					return responseErr
				}
				if len(*responseAdd) > length {
					*responseAdd = (*responseAdd)[:len(*responseAdd)-1]
				}
				*response = append(*response, *responseAdd...)
				searchItemResponse.NextPage = int32((*response)[len(*response)-1].Cursor)
				searchItemResponse.SearchMunicipalityType = pb.SearchMunicipalityType_NoMore
			} else if len(*response) == 0 {
				response, responseErr = i.dao.NewItemRepository().SearchItem(ctx, tx, req.Name, req.ProvinceId, req.MunicipalityId, int64(req.NextPage), true, 10)
				if responseErr != nil {
					return responseErr
				}
				if len(*response) > 10 {
					*response = (*response)[:len(*response)-1]
					searchItemResponse.NextPage = int32((*response)[len(*response)-1].Cursor)
				} else if len(*response) <= 10 && len(*response) != 0 {
					searchItemResponse.NextPage = int32((*response)[len(*response)-1].Cursor)
				}
				searchItemResponse.SearchMunicipalityType = pb.SearchMunicipalityType_NoMore
			}
		} else {
			response, responseErr = i.dao.NewItemRepository().SearchItem(ctx, tx, req.Name, req.ProvinceId, req.MunicipalityId, int64(req.NextPage), true, 10)
			if responseErr != nil {
				return responseErr
			}
			if len(*response) > 10 {
				*response = (*response)[:len(*response)-1]
				searchItemResponse.NextPage = int32((*response)[len(*response)-1].Cursor)
			} else if len(*response) <= 10 && len(*response) != 0 {
				searchItemResponse.NextPage = int32((*response)[len(*response)-1].Cursor)
			}
			searchItemResponse.SearchMunicipalityType = pb.SearchMunicipalityType_NoMore
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if response != nil {
		itemsResponse := make([]*pb.SearchItem, 0, len(*response))
		for _, e := range *response {
			itemsResponse = append(itemsResponse, &pb.SearchItem{
				Id:           e.ID.String(),
				Name:         e.Name,
				Thumbnail:    e.Thumbnail,
				ThumbnailUrl: i.config.ItemsBulkName + "/" + e.Thumbnail,
				BlurHash:     e.BlurHash,
				PriceCup:     e.PriceCup,
				BusinessId:   e.BusinessId.String(),
				Cursor:       int32(e.Cursor),
			})
		}
		searchItemResponse.Items = itemsResponse
	}
	return &searchItemResponse, nil
}

func (i *itemService) SearchItemByBusiness(ctx context.Context, req *pb.SearchItemByBusinessRequest, md *utils.ClientMetadata) (*pb.SearchItemByBusinessResponse, error) {
	var response *[]entity.Item
	var searchItemResponse pb.SearchItemByBusinessResponse
	var responseErr error
	err := i.sqldb.Gorm.Transaction(func(tx *gorm.DB) error {
		_, err := i.dao.NewApplicationRepository().CheckApplication(ctx, tx, *md.AccessToken)
		if err != nil {
			return err
		}
		response, responseErr = i.dao.NewItemRepository().SearchItemByBusiness(ctx, tx, req.Name, int64(req.NextPage), req.BusinessId)
		if responseErr != nil {
			return responseErr
		}
		if len(*response) <= 10 && len(*response) > 1 {
			*response = (*response)[:len(*response)]
			searchItemResponse.NextPage = int32((*response)[len(*response)-1].Cursor)
		} else if len(*response) == 1 {
			*response = (*response)[:len(*response)]
			searchItemResponse.NextPage = int32((*response)[len(*response)-1].Cursor)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(*response) != 0 {
		itemsResponse := make([]*pb.SearchItem, 0, len(*response))
		for _, e := range *response {
			itemsResponse = append(itemsResponse, &pb.SearchItem{
				Id:           e.ID.String(),
				Name:         e.Name,
				Thumbnail:    e.Thumbnail,
				ThumbnailUrl: i.config.ItemsBulkName + "/" + e.Thumbnail,
				BlurHash:     e.BlurHash,
				PriceCup:     e.PriceCup,
				Cursor:       int32(e.Cursor),
				BusinessId:   e.BusinessId.String(),
			})
		}
		searchItemResponse.Items = itemsResponse
	}
	return &searchItemResponse, nil
}
