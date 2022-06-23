package usecase

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/daniarmas/api_go/datasource"
	"github.com/daniarmas/api_go/models"
	pb "github.com/daniarmas/api_go/pkg"
	"github.com/daniarmas/api_go/repository"
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
	dao    repository.DAO
	config *utils.Config
	stDb   *sql.DB
}

func NewItemService(dao repository.DAO, config *utils.Config, stDb *sql.DB) ItemService {
	return &itemService{dao: dao, config: config, stDb: stDb}
}

func (i *itemService) UpdateItem(ctx context.Context, req *pb.UpdateItemRequest, md *utils.ClientMetadata) (*pb.Item, error) {
	var updateItemRes *models.Item
	var updateItemErr error
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
		appErr := i.dao.NewApplicationRepository().CheckApplication(tx, *md.AccessToken)
		if appErr != nil {
			return appErr
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
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id", "refresh_token_id", "device_id", "user_id", "app", "app_version", "create_time", "update_time"})
		if authorizationTokenErr != nil && authorizationTokenErr.Error() == "record not found" {
			return errors.New("unauthenticated")
		} else if authorizationTokenErr != nil {
			return authorizationTokenErr
		}
		id := uuid.MustParse(req.Item.Id)
		getItemRes, getItemErr := i.dao.NewItemRepository().GetItem(tx, &models.Item{ID: &id}, nil)
		if getItemErr != nil {
			return getItemErr
		}
		_, permissionErr := i.dao.NewUserPermissionRepository().GetUserPermission(tx, &models.UserPermission{UserId: authorizationTokenRes.UserId, Name: "update_item", BusinessId: getItemRes.BusinessId}, &[]string{"id"})
		if permissionErr != nil && permissionErr.Error() == "record not found" {
			return errors.New("permission denied")
		}
		getCartItemRes, getCartItemErr := i.dao.NewCartItemRepository().GetCartItem(tx, &models.CartItem{ItemId: &id}, nil)
		if getCartItemErr != nil && getCartItemErr.Error() != "record not found" {
			return getCartItemErr
		} else if getCartItemRes != nil {
			return errors.New("item in the cart")
		}
		if req.Item.HighQualityPhoto != "" || req.Item.LowQualityPhoto != "" || req.Item.Thumbnail != "" {
			_, hqErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), datasource.Config.ItemsBulkName, req.Item.HighQualityPhoto)
			if hqErr != nil && hqErr.Error() == "ObjectMissing" {
				return errors.New("HighQualityPhotoObject missing")
			} else if hqErr != nil {
				return hqErr
			}
			_, lqErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), datasource.Config.ItemsBulkName, req.Item.LowQualityPhoto)
			if lqErr != nil && lqErr.Error() == "ObjectMissing" {
				return errors.New("LowQualityPhotoObject missing")
			} else if lqErr != nil {
				return lqErr
			}
			_, tnErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), datasource.Config.ItemsBulkName, req.Item.Thumbnail)
			if tnErr != nil && tnErr.Error() == "ObjectMissing" {
				return errors.New("ThumbnailObject missing")
			} else if tnErr != nil {
				return tnErr
			}
			_, copyHqErr := repository.Datasource.NewObjectStorageDatasource().CopyObject(context.Background(), minio.CopyDestOptions{Bucket: repository.Config.ItemsDeletedBulkName, Object: getItemRes.HighQualityPhoto}, minio.CopySrcOptions{Bucket: repository.Config.ItemsBulkName, Object: getItemRes.HighQualityPhoto})
			if copyHqErr != nil {
				return copyHqErr
			}
			_, copyLqErr := repository.Datasource.NewObjectStorageDatasource().CopyObject(context.Background(), minio.CopyDestOptions{Bucket: repository.Config.ItemsDeletedBulkName, Object: getItemRes.LowQualityPhoto}, minio.CopySrcOptions{Bucket: repository.Config.ItemsBulkName, Object: getItemRes.LowQualityPhoto})
			if copyLqErr != nil {
				return copyLqErr
			}
			_, copyThErr := repository.Datasource.NewObjectStorageDatasource().CopyObject(context.Background(), minio.CopyDestOptions{Bucket: repository.Config.ItemsDeletedBulkName, Object: getItemRes.Thumbnail}, minio.CopySrcOptions{Bucket: repository.Config.ItemsBulkName, Object: getItemRes.Thumbnail})
			if copyThErr != nil {
				return copyThErr
			}
			rmHqErr := repository.Datasource.NewObjectStorageDatasource().RemoveObject(context.Background(), repository.Config.ItemsBulkName, getItemRes.HighQualityPhoto, minio.RemoveObjectOptions{})
			if rmHqErr != nil {
				return rmHqErr
			}
			rmLqErr := repository.Datasource.NewObjectStorageDatasource().RemoveObject(context.Background(), repository.Config.ItemsBulkName, getItemRes.LowQualityPhoto, minio.RemoveObjectOptions{})
			if rmLqErr != nil {
				return rmLqErr
			}
			rmThErr := repository.Datasource.NewObjectStorageDatasource().RemoveObject(context.Background(), repository.Config.ItemsBulkName, getItemRes.Thumbnail, minio.RemoveObjectOptions{})
			if rmThErr != nil {
				return rmThErr
			}
		}
		updateItemRes, updateItemErr = i.dao.NewItemRepository().UpdateItem(tx, &models.Item{ID: &id}, &models.Item{Name: req.Item.Name, Description: req.Item.Description, PriceCup: req.Item.PriceCup, Availability: int64(req.Item.Availability), HighQualityPhoto: req.Item.HighQualityPhoto, LowQualityPhoto: req.Item.LowQualityPhoto, Thumbnail: req.Item.Thumbnail, BlurHash: req.Item.BlurHash})
		if updateItemErr != nil {
			return updateItemErr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &pb.Item{
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
		CreateTime:           timestamppb.New(updateItemRes.CreateTime),
		UpdateTime:           timestamppb.New(updateItemRes.UpdateTime),
	}, nil
}

func (i *itemService) DeleteItem(ctx context.Context, req *pb.DeleteItemRequest, md *utils.ClientMetadata) error {
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
		appErr := i.dao.NewApplicationRepository().CheckApplication(tx, *md.AccessToken)
		if appErr != nil {
			return appErr
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
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id", "refresh_token_id", "device_id", "user_id", "app", "app_version", "create_time", "update_time"})
		if authorizationTokenErr != nil && authorizationTokenErr.Error() == "record not found" {
			return errors.New("unauthenticated")
		} else if authorizationTokenErr != nil {
			return authorizationTokenErr
		}
		id := uuid.MustParse(req.Id)
		getItemRes, getItemErr := i.dao.NewItemRepository().GetItem(tx, &models.Item{ID: &id}, nil)
		if getItemErr != nil && getItemErr.Error() == "record not found" {
			return errors.New("item not found")
		} else if getItemErr != nil {
			return getItemErr
		}
		_, err := i.dao.NewUserPermissionRepository().GetUserPermission(tx, &models.UserPermission{Name: "delete_item", UserId: authorizationTokenRes.UserId, BusinessId: getItemRes.BusinessId}, &[]string{"id"})
		if err != nil && err.Error() == "record not found" {
			return errors.New("permission denied")
		} else if err != nil {
			return err
		}
		businessRes, businessErr := i.dao.NewBusinessScheduleRepository().BusinessIsOpen(tx, &models.BusinessSchedule{BusinessId: getItemRes.BusinessId}, "OrderTypePickUp")
		if businessErr != nil && businessErr.Error() != "business closed" {
			return businessErr
		} else if businessRes {
			return errors.New("business is open")
		}
		businessHomeDeliveryRes, businessHomeDeliveryErr := i.dao.NewBusinessScheduleRepository().BusinessIsOpen(tx, &models.BusinessSchedule{BusinessId: getItemRes.BusinessId}, "OrderTypeHomeDelivery")
		if businessHomeDeliveryErr != nil && businessHomeDeliveryErr.Error() != "business closed" {
			return businessHomeDeliveryErr
		} else if businessHomeDeliveryRes {
			return errors.New("business is open")
		}
		getCartItemRes, getCartItemErr := i.dao.NewCartItemRepository().GetCartItem(tx, &models.CartItem{ItemId: &id}, nil)
		if getCartItemErr != nil && getCartItemErr.Error() != "record not found" {
			return getCartItemErr
		} else if getCartItemRes != nil {
			return errors.New("item in the cart")
		}
		deleteItemErr := i.dao.NewItemRepository().DeleteItem(tx, &models.Item{ID: &id})
		if deleteItemErr != nil {
			return deleteItemErr
		}
		_, copyHqErr := repository.Datasource.NewObjectStorageDatasource().CopyObject(context.Background(), minio.CopyDestOptions{Bucket: repository.Config.ItemsDeletedBulkName, Object: getItemRes.HighQualityPhoto}, minio.CopySrcOptions{Bucket: repository.Config.ItemsBulkName, Object: getItemRes.HighQualityPhoto})
		if copyHqErr != nil {
			return copyHqErr
		}
		_, copyLqErr := repository.Datasource.NewObjectStorageDatasource().CopyObject(context.Background(), minio.CopyDestOptions{Bucket: repository.Config.ItemsDeletedBulkName, Object: getItemRes.LowQualityPhoto}, minio.CopySrcOptions{Bucket: repository.Config.ItemsBulkName, Object: getItemRes.LowQualityPhoto})
		if copyLqErr != nil {
			return copyLqErr
		}
		_, copyThErr := repository.Datasource.NewObjectStorageDatasource().CopyObject(context.Background(), minio.CopyDestOptions{Bucket: repository.Config.ItemsDeletedBulkName, Object: getItemRes.Thumbnail}, minio.CopySrcOptions{Bucket: repository.Config.ItemsBulkName, Object: getItemRes.Thumbnail})
		if copyThErr != nil {
			return copyThErr
		}
		rmHqErr := repository.Datasource.NewObjectStorageDatasource().RemoveObject(context.Background(), repository.Config.ItemsBulkName, getItemRes.HighQualityPhoto, minio.RemoveObjectOptions{})
		if rmHqErr != nil {
			return rmHqErr
		}
		rmLqErr := repository.Datasource.NewObjectStorageDatasource().RemoveObject(context.Background(), repository.Config.ItemsBulkName, getItemRes.LowQualityPhoto, minio.RemoveObjectOptions{})
		if rmLqErr != nil {
			return rmLqErr
		}
		rmThErr := repository.Datasource.NewObjectStorageDatasource().RemoveObject(context.Background(), repository.Config.ItemsBulkName, getItemRes.Thumbnail, minio.RemoveObjectOptions{})
		if rmThErr != nil {
			return rmThErr
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (i *itemService) CreateItem(ctx context.Context, req *pb.CreateItemRequest, md *utils.ClientMetadata) (*pb.Item, error) {
	var itemRes *models.Item
	var itemErr error
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
		appErr := i.dao.NewApplicationRepository().CheckApplication(tx, *md.AccessToken)
		if appErr != nil {
			return appErr
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
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenRepository().GetAuthorizationToken(ctx, tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, &[]string{"id", "refresh_token_id", "device_id", "user_id", "app", "app_version", "create_time", "update_time"})
		if authorizationTokenErr != nil && authorizationTokenErr.Error() == "record not found" {
			return errors.New("unauthenticated")
		} else if authorizationTokenErr != nil {
			return authorizationTokenErr
		}
		businessId := uuid.MustParse(req.Item.BusinessId)
		businessRes, businessErr := i.dao.NewBusinessRepository().GetBusiness(tx, &models.Business{ID: &businessId}, &[]string{"province_id", "municipality_id"})
		if businessErr != nil {
			return businessErr
		}
		_, err := i.dao.NewUserPermissionRepository().GetUserPermission(tx, &models.UserPermission{Name: "create_item", UserId: authorizationTokenRes.UserId, BusinessId: &businessId}, &[]string{"id"})
		if err != nil && err.Error() == "record not found" {
			return errors.New("permission denied")
		} else if err != nil {
			return err
		}
		_, hqErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), datasource.Config.ItemsBulkName, req.Item.HighQualityPhoto)
		if hqErr != nil && hqErr.Error() == "ObjectMissing" {
			return errors.New("highQualityPhotoObject missing")
		} else if hqErr != nil {
			return hqErr
		}
		_, lqErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), datasource.Config.ItemsBulkName, req.Item.LowQualityPhoto)
		if lqErr != nil && lqErr.Error() == "ObjectMissing" {
			return errors.New("lowQualityPhotoObject missing")
		} else if lqErr != nil {
			return lqErr
		}
		_, tnErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), datasource.Config.ItemsBulkName, req.Item.Thumbnail)
		if tnErr != nil && tnErr.Error() == "ObjectMissing" {
			return errors.New("thumbnailObject missing")
		} else if tnErr != nil {
			return tnErr
		}
		businessCollectionId := uuid.MustParse(req.Item.BusinessCollectionId)
		itemRes, itemErr = i.dao.NewItemRepository().CreateItem(tx, &models.Item{Name: req.Item.Name, Description: req.Item.Description, PriceCup: req.Item.PriceCup, CostCup: req.Item.CostCup, ProfitCup: req.Item.ProfitCup, PriceUsd: req.Item.PriceUsd, ProfitUsd: req.Item.ProfitUsd, CostUsd: req.Item.CostUsd, Availability: int64(req.Item.Availability), BusinessId: &businessId, BusinessCollectionId: &businessCollectionId, HighQualityPhoto: req.Item.HighQualityPhoto, LowQualityPhoto: req.Item.LowQualityPhoto, Thumbnail: req.Item.Thumbnail, BlurHash: req.Item.BlurHash, ProvinceId: businessRes.ProvinceId, MunicipalityId: businessRes.MunicipalityId, AvailableFlag: req.Item.AvailableFlag, EnabledFlag: req.Item.EnabledFlag})
		if itemErr != nil {
			return itemErr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &pb.Item{Id: itemRes.ID.String(), Name: itemRes.Name, Description: itemRes.Description, PriceCup: itemRes.PriceCup, CostCup: itemRes.CostCup, ProfitCup: itemRes.ProfitCup, PriceUsd: itemRes.PriceUsd, CostUsd: itemRes.CostUsd, ProfitUsd: itemRes.ProfitUsd, EnabledFlag: itemRes.EnabledFlag, AvailableFlag: itemRes.AvailableFlag, Availability: int32(itemRes.Availability), BusinessId: itemRes.BusinessId.String(), BusinessCollectionId: itemRes.BusinessCollectionId.String(), HighQualityPhoto: itemRes.HighQualityPhoto, LowQualityPhoto: itemRes.LowQualityPhoto, Thumbnail: itemRes.Thumbnail, BlurHash: itemRes.BlurHash, ProvinceId: itemRes.ProvinceId.String(), MunicipalityId: itemRes.MunicipalityId.String(), ThumbnailUrl: i.config.ItemsBulkName + "/" + itemRes.Thumbnail, HighQualityPhotoUrl: i.config.ItemsBulkName + "/" + itemRes.HighQualityPhoto, LowQualityPhotoUrl: i.config.ItemsBulkName + "/" + itemRes.LowQualityPhoto, CreateTime: timestamppb.New(itemRes.CreateTime), UpdateTime: timestamppb.New(itemRes.UpdateTime)}, nil
}

func (i *itemService) ListItem(ctx context.Context, req *pb.ListItemRequest, md *utils.ClientMetadata) (*pb.ListItemResponse, error) {
	where := models.Item{}
	var nextPage time.Time
	if req.NextPage == nil || (req.NextPage.Nanos == 0 && req.NextPage.Seconds == 0) {
		nextPage = time.Now()
	} else {
		nextPage = req.NextPage.AsTime()
	}
	var items *[]models.Item
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
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
		appErr := i.dao.NewApplicationRepository().CheckApplication(tx, *md.AccessToken)
		if appErr != nil {
			return appErr
		}
		items, itemsErr = i.dao.NewItemRepository().ListItem(tx, &where, nextPage, nil)
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
	var item *models.Item
	var itemErr error
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
		appErr := i.dao.NewApplicationRepository().CheckApplication(tx, *md.AccessToken)
		if appErr != nil {
			return appErr
		}
		id := uuid.MustParse(req.Id)
		item, itemErr = i.dao.NewItemRepository().GetItem(tx, &models.Item{ID: &id}, nil)
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
	var response *[]models.Item
	var searchItemResponse pb.SearchItemResponse
	var responseErr error
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
		appErr := i.dao.NewApplicationRepository().CheckApplication(tx, *md.AccessToken)
		if appErr != nil {
			return appErr
		}
		if req.SearchMunicipalityType.String() == "More" {
			response, responseErr = i.dao.NewItemRepository().SearchItem(tx, req.Name, req.ProvinceId, req.MunicipalityId, int64(req.NextPage), false, 10, &[]string{"id", "name", "price_cup", "thumbnail", "blurhash", "cursor"})
			if responseErr != nil {
				return responseErr
			}
			if len(*response) > 10 {
				*response = (*response)[:len(*response)-1]
				searchItemResponse.NextPage = int32((*response)[len(*response)-1].Cursor)
				searchItemResponse.SearchMunicipalityType = pb.SearchMunicipalityType_More
			} else if len(*response) <= 10 && len(*response) != 0 {
				length := 10 - len(*response)
				responseAdd, responseErr := i.dao.NewItemRepository().SearchItem(tx, req.Name, req.ProvinceId, req.MunicipalityId, int64(req.NextPage), true, int64(length), &[]string{"id", "name", "price_cup", "thumbnail", "blurhash", "cursor"})
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
				response, responseErr = i.dao.NewItemRepository().SearchItem(tx, req.Name, req.ProvinceId, req.MunicipalityId, int64(req.NextPage), true, 10, &[]string{"id", "name", "price_cup", "thumbnail", "blurhash", "cursor"})
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
			response, responseErr = i.dao.NewItemRepository().SearchItem(tx, req.Name, req.ProvinceId, req.MunicipalityId, int64(req.NextPage), true, 10, &[]string{"id", "name", "priceCup", "thumbnail", "thumbnail_blurhash", "cursor"})
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
				Cursor:       int32(e.Cursor),
			})
		}
		searchItemResponse.Items = itemsResponse
	}
	return &searchItemResponse, nil
}

func (i *itemService) SearchItemByBusiness(ctx context.Context, req *pb.SearchItemByBusinessRequest, md *utils.ClientMetadata) (*pb.SearchItemByBusinessResponse, error) {
	var response *[]models.Item
	var searchItemResponse pb.SearchItemByBusinessResponse
	var responseErr error
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
		appErr := i.dao.NewApplicationRepository().CheckApplication(tx, *md.AccessToken)
		if appErr != nil {
			return appErr
		}
		response, responseErr = i.dao.NewItemRepository().SearchItemByBusiness(tx, req.Name, int64(req.NextPage), req.BusinessId, &[]string{"id", "name", "price", "thumbnail", "thumbnail_blurhash", "cursor"})
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
			})
		}
		searchItemResponse.Items = itemsResponse
	}
	return &searchItemResponse, nil
}
