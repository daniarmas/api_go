package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/daniarmas/api_go/datasource"
	"github.com/daniarmas/api_go/dto"
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
	GetItem(ctx context.Context, req *pb.GetItemRequest, md *utils.ClientMetadata) (*pb.GetItemResponse, error)
	ListItem(ctx context.Context, req *pb.ListItemRequest, md *utils.ClientMetadata) (*pb.ListItemResponse, error)
	SearchItem(ctx context.Context, req *pb.SearchItemRequest, md *utils.ClientMetadata) (*pb.SearchItemResponse, error)
	CreateItem(request *dto.CreateItemRequest) (*models.Item, error)
	UpdateItem(request *dto.UpdateItemRequest) (*models.Item, error)
	DeleteItem(request *dto.DeleteItemRequest) error
}

type itemService struct {
	dao    repository.DAO
	config *utils.Config
}

func NewItemService(dao repository.DAO, config *utils.Config) ItemService {
	return &itemService{dao: dao, config: config}
}

func (i *itemService) UpdateItem(request *dto.UpdateItemRequest) (*models.Item, error) {
	var updateItemRes *models.Item
	var updateItemErr error
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: &request.Metadata.Get("authorization")[0]}
		authorizationTokenParseErr := repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(jwtAuthorizationToken)
		if authorizationTokenParseErr != nil {
			switch authorizationTokenParseErr.Error() {
			case "Token is expired":
				return errors.New("authorizationtoken expired")
			case "signature is invalid":
				return errors.New("signature is invalid")
			case "token contains an invalid number of segments":
				return errors.New("token contains an invalid number of segments")
			default:
				return authorizationTokenParseErr
			}
		}
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, nil)
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		} else if authorizationTokenRes == nil {
			return errors.New("unauthenticated")
		}
		_, err := i.dao.NewUserPermissionRepository().GetUserPermission(tx, &models.UserPermission{Name: "update_item", UserId: authorizationTokenRes.RefreshToken.UserId, BusinessId: request.BusinessId}, &[]string{"id"})
		if err != nil && err.Error() == "record not found" {
			return errors.New("permission denied")
		} else if err != nil {
			return err
		}
		businessRes, businessErr := i.dao.NewBusinessScheduleRepository().BusinessIsOpen(tx, &models.BusinessSchedule{BusinessId: request.BusinessId}, "OrderTypePickUp")
		if businessErr != nil {
			return businessErr
		} else if businessRes {
			return errors.New("business is open")
		}
		businessHomeDeliveryRes, businessHomeDeliveryErr := i.dao.NewBusinessScheduleRepository().BusinessIsOpen(tx, &models.BusinessSchedule{BusinessId: request.BusinessId}, "OrderTypeHomeDelivery")
		if businessHomeDeliveryErr != nil {
			return businessHomeDeliveryErr
		} else if businessHomeDeliveryRes {
			return errors.New("business is open")
		}
		getCartItemRes, getCartItemErr := i.dao.NewCartItemRepository().GetCartItem(tx, &models.CartItem{ItemId: request.ItemId}, nil)
		if getCartItemErr != nil && getCartItemErr.Error() != "record not found" {
			return getCartItemErr
		} else if getCartItemRes != nil {
			return errors.New("item in the cart")
		}
		getItemRes, getItemErr := i.dao.NewItemQuery().GetItem(tx, &models.Item{ID: request.ItemId}, &[]string{})
		if getItemErr != nil {
			return getItemErr
		}
		if request.HighQualityPhoto != "" || request.LowQualityPhoto != "" || request.Thumbnail != "" {
			_, hqErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), datasource.Config.ItemsBulkName, request.HighQualityPhoto)
			if hqErr != nil && hqErr.Error() == "ObjectMissing" {
				return errors.New("HighQualityPhotoObject missing")
			} else if hqErr != nil {
				return hqErr
			}
			_, lqErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), datasource.Config.ItemsBulkName, request.LowQualityPhoto)
			if lqErr != nil && lqErr.Error() == "ObjectMissing" {
				return errors.New("LowQualityPhotoObject missing")
			} else if lqErr != nil {
				return lqErr
			}
			_, tnErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), datasource.Config.ItemsBulkName, request.Thumbnail)
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
		updateItemRes, updateItemErr = i.dao.NewItemQuery().UpdateItem(tx, &models.Item{ID: request.ItemId}, &models.Item{Name: request.Name, Description: request.Description, Price: request.Price, Availability: request.Availability, BusinessCollectionId: request.BusinessColletionId, HighQualityPhoto: request.HighQualityPhoto, HighQualityPhotoBlurHash: request.HighQualityPhotoBlurHash, LowQualityPhoto: request.LowQualityPhoto, LowQualityPhotoBlurHash: request.LowQualityPhotoBlurHash, Thumbnail: request.Thumbnail, ThumbnailBlurHash: request.ThumbnailBlurHash})
		if updateItemErr != nil {
			return updateItemErr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return updateItemRes, nil
}

func (i *itemService) DeleteItem(request *dto.DeleteItemRequest) error {
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: &request.Metadata.Get("authorization")[0]}
		authorizationTokenParseErr := repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(jwtAuthorizationToken)
		if authorizationTokenParseErr != nil {
			switch authorizationTokenParseErr.Error() {
			case "Token is expired":
				return errors.New("authorizationtoken expired")
			case "signature is invalid":
				return errors.New("signature is invalid")
			case "token contains an invalid number of segments":
				return errors.New("token contains an invalid number of segments")
			default:
				return authorizationTokenParseErr
			}
		}
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, nil)
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		} else if authorizationTokenRes == nil {
			return errors.New("unauthenticated")
		}
		getItemRes, getItemErr := i.dao.NewItemQuery().GetItem(tx, &models.Item{ID: request.ItemId}, &[]string{})
		if getItemErr != nil {
			return getItemErr
		}
		_, err := i.dao.NewUserPermissionRepository().GetUserPermission(tx, &models.UserPermission{Name: "delete_item", UserId: authorizationTokenRes.UserId, BusinessId: getItemRes.BusinessId}, &[]string{"id"})
		if err != nil && err.Error() == "record not found" {
			return errors.New("permission denied")
		} else if err != nil {
			return err
		}
		businessRes, businessErr := i.dao.NewBusinessScheduleRepository().BusinessIsOpen(tx, &models.BusinessSchedule{BusinessId: getItemRes.BusinessId}, "OrderTypePickUp")
		if businessErr != nil {
			return businessErr
		} else if businessRes {
			return errors.New("business is open")
		}
		businessHomeDeliveryRes, businessHomeDeliveryErr := i.dao.NewBusinessScheduleRepository().BusinessIsOpen(tx, &models.BusinessSchedule{BusinessId: getItemRes.BusinessId}, "OrderTypeHomeDelivery")
		if businessHomeDeliveryErr != nil {
			return businessHomeDeliveryErr
		} else if businessHomeDeliveryRes {
			return errors.New("business is open")
		}
		getCartItemRes, getCartItemErr := i.dao.NewCartItemRepository().GetCartItem(tx, &models.CartItem{ItemId: request.ItemId}, nil)
		if getCartItemErr != nil && getCartItemErr.Error() != "record not found" {
			return getCartItemErr
		} else if getCartItemRes != nil {
			return errors.New("item in the cart")
		}
		deleteItemErr := i.dao.NewItemQuery().DeleteItem(tx, &models.Item{ID: request.ItemId})
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

func (i *itemService) CreateItem(request *dto.CreateItemRequest) (*models.Item, error) {
	var itemRes *models.Item
	var itemErr error
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: &request.Metadata.Get("authorization")[0]}
		authorizationTokenParseErr := repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(jwtAuthorizationToken)
		if authorizationTokenParseErr != nil {
			switch authorizationTokenParseErr.Error() {
			case "Token is expired":
				return errors.New("authorizationtoken expired")
			case "signature is invalid":
				return errors.New("signature is invalid")
			case "token contains an invalid number of segments":
				return errors.New("token contains an invalid number of segments")
			default:
				return authorizationTokenParseErr
			}
		}
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, nil)
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		} else if authorizationTokenRes == nil {
			return errors.New("unauthenticated")
		}
		_, err := i.dao.NewUserPermissionRepository().GetUserPermission(tx, &models.UserPermission{Name: "create_item", UserId: authorizationTokenRes.UserId, BusinessId: request.BusinessId}, &[]string{"id"})
		if err != nil && err.Error() == "record not found" {
			return errors.New("permission denied")
		} else if err != nil {
			return err
		}
		_, hqErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), datasource.Config.ItemsBulkName, request.HighQualityPhoto)
		if hqErr != nil && hqErr.Error() == "ObjectMissing" {
			return errors.New("HighQualityPhotoObject missing")
		} else if hqErr != nil {
			return hqErr
		}
		_, lqErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), datasource.Config.ItemsBulkName, request.LowQualityPhoto)
		if lqErr != nil && lqErr.Error() == "ObjectMissing" {
			return errors.New("LowQualityPhotoObject missing")
		} else if lqErr != nil {
			return lqErr
		}
		_, tnErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), datasource.Config.ItemsBulkName, request.Thumbnail)
		if tnErr != nil && tnErr.Error() == "ObjectMissing" {
			return errors.New("ThumbnailObject missing")
		} else if tnErr != nil {
			return tnErr
		}
		businessRes, businessErr := i.dao.NewBusinessQuery().GetBusinessProvinceAndMunicipality(tx, *request.BusinessId)
		if businessErr != nil {
			return businessErr
		}
		businessCollectionId := uuid.MustParse(request.BusinessCollectionId)
		itemRes, itemErr = i.dao.NewItemQuery().CreateItem(tx, &models.Item{Name: request.Name, Description: request.Description, Price: request.Price, Availability: -1, BusinessId: request.BusinessId, BusinessCollectionId: &businessCollectionId, HighQualityPhoto: request.HighQualityPhoto, LowQualityPhoto: request.LowQualityPhoto, Thumbnail: request.Thumbnail, HighQualityPhotoBlurHash: request.HighQualityPhotoBlurHash, LowQualityPhotoBlurHash: request.LowQualityPhotoBlurHash, ThumbnailBlurHash: request.ThumbnailBlurHash, Status: "Unavailable", ProvinceId: businessRes.ProvinceId, MunicipalityId: businessRes.MunicipalityId})
		if itemErr != nil {
			return itemErr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return itemRes, nil
}

func (i *itemService) ListItem(ctx context.Context, req *pb.ListItemRequest, md *utils.ClientMetadata) (*pb.ListItemResponse, error) {
	var where models.Item
	var nextPage time.Time
	if req.NextPage == nil {
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
	}
	if req.BusinessId != "" {
		businessId = uuid.MustParse(req.BusinessId)
	}
	if req.BusinessId != "" || req.BusinessCollectionId != "" {
		where = models.Item{BusinessId: &businessId, BusinessCollectionId: &businessCollectionId}
	}
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
		items, itemsErr = i.dao.NewItemQuery().ListItem(tx, &where, nextPage)
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
				Id:                       item.ID.String(),
				Name:                     item.Name,
				Description:              item.Description,
				Price:                    item.Price,
				Availability:             int32(item.Availability),
				BusinessId:               item.BusinessId.String(),
				BusinessCollectionId:     item.BusinessCollectionId.String(),
				HighQualityPhoto:         item.HighQualityPhoto,
				HighQualityPhotoUrl:      i.config.ItemsBulkName + "/" + item.HighQualityPhoto,
				HighQualityPhotoBlurHash: item.HighQualityPhotoBlurHash,
				LowQualityPhoto:          item.LowQualityPhoto,
				LowQualityPhotoUrl:       i.config.ItemsBulkName + "/" + item.LowQualityPhoto,
				LowQualityPhotoBlurHash:  item.LowQualityPhotoBlurHash,
				Thumbnail:                item.Thumbnail,
				ThumbnailUrl:             i.config.ItemsBulkName + "/" + item.Thumbnail,
				ThumbnailBlurHash:        item.ThumbnailBlurHash,
				Cursor:                   int32(item.Cursor),
				CreateTime:               timestamppb.New(item.CreateTime),
				UpdateTime:               timestamppb.New(item.UpdateTime),
			})
		}
		res.Items = itemsResponse
	}
	return &res, nil
}

func (i *itemService) GetItem(ctx context.Context, req *pb.GetItemRequest, md *utils.ClientMetadata) (*pb.GetItemResponse, error) {
	var res pb.GetItemResponse
	var item *models.ItemBusiness
	var itemErr error
	err := datasource.Connection.Transaction(func(tx *gorm.DB) error {
		item, itemErr = i.dao.NewItemQuery().GetItemWithLocation(tx, req.Id, ewkb.Point{Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{req.Location.Latitude, req.Location.Longitude}).SetSRID(4326)})
		if itemErr != nil {
			return itemErr
		}
		res.Item = &pb.Item{
			Id:                       item.ID.String(),
			Name:                     item.Name,
			Description:              item.Description,
			Price:                    item.Price,
			Availability:             int32(item.Availability),
			BusinessId:               item.BusinessId.String(),
			BusinessCollectionId:     item.BusinessCollectionId.String(),
			HighQualityPhoto:         item.HighQualityPhoto,
			HighQualityPhotoUrl:      i.config.ItemsBulkName + "/" + item.HighQualityPhoto,
			HighQualityPhotoBlurHash: item.HighQualityPhotoBlurHash,
			LowQualityPhoto:          item.LowQualityPhoto,
			LowQualityPhotoUrl:       i.config.ItemsBulkName + "/" + item.LowQualityPhoto,
			LowQualityPhotoBlurHash:  item.LowQualityPhotoBlurHash,
			Thumbnail:                item.Thumbnail,
			ThumbnailUrl:             i.config.ItemsBulkName + "/" + item.Thumbnail,
			ThumbnailBlurHash:        item.ThumbnailBlurHash,
			Cursor:                   item.Cursor,
			CreateTime:               timestamppb.New(item.CreateTime),
			UpdateTime:               timestamppb.New(item.UpdateTime),
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
		if req.SearchMunicipalityType.String() == "More" {
			response, responseErr = i.dao.NewItemQuery().SearchItem(tx, req.Name, req.ProvinceId, req.MunicipalityId, int64(req.NextPage), false, 10, &[]string{})
			if responseErr != nil {
				return responseErr
			}
			if len(*response) > 10 {
				*response = (*response)[:len(*response)-1]
				searchItemResponse.NextPage = int32((*response)[len(*response)-1].Cursor)
				searchItemResponse.SearchMunicipalityType = pb.SearchMunicipalityType_More
			} else if len(*response) <= 10 && len(*response) != 0 {
				length := 10 - len(*response)
				responseAdd, responseErr := i.dao.NewItemQuery().SearchItem(tx, req.Name, req.ProvinceId, req.MunicipalityId, int64(req.NextPage), true, int64(length), &[]string{})
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
				response, responseErr = i.dao.NewItemQuery().SearchItem(tx, req.Name, req.ProvinceId, req.MunicipalityId, int64(req.NextPage), true, 10, &[]string{})
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
			response, responseErr = i.dao.NewItemQuery().SearchItem(tx, req.Name, req.ProvinceId, req.MunicipalityId, int64(req.NextPage), true, 10, &[]string{})
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
				Id:                e.ID.String(),
				Name:              e.Name,
				Thumbnail:         e.Thumbnail,
				ThumbnailUrl:      i.config.ItemsBulkName + "/" + e.Thumbnail,
				ThumbnailBlurHash: e.ThumbnailBlurHash,
				Price:             e.Price,
				Cursor:            int32(e.Cursor),
				Status:            *utils.ParseItemStatusType(&e.Status),
			})
		}
		searchItemResponse.Items = itemsResponse
	}
	return &searchItemResponse, nil
}
