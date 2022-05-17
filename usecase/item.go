package usecase

import (
	// "context"
	"context"
	"errors"

	"github.com/daniarmas/api_go/datasource"
	"github.com/daniarmas/api_go/dto"
	"github.com/daniarmas/api_go/models"
	pb "github.com/daniarmas/api_go/pkg"
	"github.com/daniarmas/api_go/repository"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

type ItemService interface {
	GetItem(request *dto.GetItemRequest) (*models.ItemBusiness, error)
	ListItem(itemRequest *dto.ListItemRequest) (*dto.ListItemResponse, error)
	SearchItem(name string, provinceId string, municipalityId string, cursor int64, searchMunicipalityType string) (*dto.SearchItemResponse, error)
	CreateItem(request *dto.CreateItemRequest) (*models.Item, error)
	UpdateItem(request *dto.UpdateItemRequest) (*models.Item, error)
	DeleteItem(request *dto.DeleteItemRequest) error
}

type itemService struct {
	dao repository.DAO
}

func NewItemService(dao repository.DAO) ItemService {
	return &itemService{dao: dao}
}

func (i *itemService) UpdateItem(request *dto.UpdateItemRequest) (*models.Item, error) {
	var updateItemRes *models.Item
	var updateItemErr error
	err := datasource.DB.Transaction(func(tx *gorm.DB) error {
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
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId})
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		} else if authorizationTokenRes == nil {
			return errors.New("unauthenticated")
		}
		permissionExistsErr := i.dao.NewUserPermissionRepository().UserPermissionExists(tx, &models.UserPermission{Name: "update_item", UserId: *&authorizationTokenRes.RefreshToken.UserId, BusinessId: request.BusinessId})
		if permissionExistsErr != nil && permissionExistsErr.Error() == "record not found" {
			return errors.New("permission denied")
		} else if permissionExistsErr != nil {
			return permissionExistsErr
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
		getCartItemRes, getCartItemErr := i.dao.NewCartItemRepository().GetCartItem(tx, &models.CartItem{ItemId: request.ItemId})
		if getCartItemErr != nil && getCartItemErr.Error() != "record not found" {
			return getCartItemErr
		} else if getCartItemRes != nil {
			return errors.New("item in the cart")
		}
		getItemRes, getItemErr := i.dao.NewItemQuery().GetItem(tx, &models.Item{ID: request.ItemId})
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
	err := datasource.DB.Transaction(func(tx *gorm.DB) error {
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
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId})
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		} else if authorizationTokenRes == nil {
			return errors.New("unauthenticated")
		}
		getItemRes, getItemErr := i.dao.NewItemQuery().GetItem(tx, &models.Item{ID: request.ItemId})
		if getItemErr != nil {
			return getItemErr
		}
		permissionExistsErr := i.dao.NewUserPermissionRepository().UserPermissionExists(tx, &models.UserPermission{Name: "delete_item", UserId: *authorizationTokenRes.UserId, BusinessId: getItemRes.BusinessId})
		if permissionExistsErr != nil && permissionExistsErr.Error() == "record not found" {
			return errors.New("permission denied")
		} else if permissionExistsErr != nil {
			return permissionExistsErr
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
		getCartItemRes, getCartItemErr := i.dao.NewCartItemRepository().GetCartItem(tx, &models.CartItem{ItemId: request.ItemId})
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
	err := datasource.DB.Transaction(func(tx *gorm.DB) error {
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
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId})
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		} else if authorizationTokenRes == nil {
			return errors.New("unauthenticated")
		}
		permissionExistsErr := i.dao.NewUserPermissionRepository().UserPermissionExists(tx, &models.UserPermission{Name: "create_item", UserId: *authorizationTokenRes.UserId, BusinessId: request.BusinessId})
		if permissionExistsErr != nil && permissionExistsErr.Error() == "record not found" {
			return errors.New("permission denied")
		} else if permissionExistsErr != nil {
			return permissionExistsErr
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
		businessRes, businessErr := i.dao.NewBusinessQuery().GetBusinessProvinceAndMunicipality(tx, request.BusinessId)
		if businessErr != nil {
			return businessErr
		}
		itemRes, itemErr = i.dao.NewItemQuery().CreateItem(tx, &models.Item{Name: request.Name, Description: request.Description, Price: request.Price, Availability: -1, BusinessId: request.BusinessId, BusinessCollectionId: uuid.MustParse(request.BusinessCollectionId), HighQualityPhoto: request.HighQualityPhoto, LowQualityPhoto: request.LowQualityPhoto, Thumbnail: request.Thumbnail, HighQualityPhotoBlurHash: request.HighQualityPhotoBlurHash, LowQualityPhotoBlurHash: request.LowQualityPhotoBlurHash, ThumbnailBlurHash: request.ThumbnailBlurHash, Status: "Unavailable", ProvinceId: businessRes.ProvinceId, MunicipalityId: businessRes.MunicipalityId})
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

func (i *itemService) ListItem(itemRequest *dto.ListItemRequest) (*dto.ListItemResponse, error) {
	var items *[]models.Item
	var listItemResponse dto.ListItemResponse
	var itemsErr error
	err := datasource.DB.Transaction(func(tx *gorm.DB) error {
		if itemRequest.BusinessId != "" && itemRequest.BusinessCollectionId == "" {
			itemCategoryRes, itemCategoryErr := i.dao.NewBusinessCollectionQuery().GetBusinessCollection(tx, &models.BusinessCollection{Index: 0, BusinessId: uuid.MustParse(itemRequest.BusinessId)})
			if itemCategoryErr != nil {
				return itemCategoryErr
			}
			items, itemsErr = i.dao.NewItemQuery().ListItem(tx, &models.Item{BusinessId: uuid.MustParse(itemRequest.BusinessId), BusinessCollectionId: itemCategoryRes.ID}, itemRequest.NextPage)
			if itemsErr != nil {
				return itemsErr
			} else if len(*items) > 10 {
				*items = (*items)[:len(*items)-1]
				listItemResponse.NextPage = (*items)[len(*items)-1].CreateTime
			} else if len(*items) == 0 {
				listItemResponse.NextPage = itemRequest.NextPage
			} else {
				listItemResponse.NextPage = (*items)[len(*items)-1].CreateTime
			}
		} else if itemRequest.BusinessId != "" && itemRequest.BusinessCollectionId != "" {
			items, itemsErr = i.dao.NewItemQuery().ListItem(tx, &models.Item{BusinessId: uuid.MustParse(itemRequest.BusinessId), BusinessCollectionId: uuid.MustParse(itemRequest.BusinessCollectionId)}, itemRequest.NextPage)
			if itemsErr != nil {
				return itemsErr
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	listItemResponse.Items = *items
	return &listItemResponse, nil
}

func (i *itemService) GetItem(request *dto.GetItemRequest) (*models.ItemBusiness, error) {
	var item *models.ItemBusiness
	var itemErr error
	err := datasource.DB.Transaction(func(tx *gorm.DB) error {
		item, itemErr = i.dao.NewItemQuery().GetItemWithLocation(tx, request.Id, request.Location)
		if itemErr != nil {
			return itemErr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (i *itemService) SearchItem(name string, provinceId string, municipalityId string, cursor int64, searchMunicipalityType string) (*dto.SearchItemResponse, error) {
	var response *[]models.Item
	var searchItemResponse dto.SearchItemResponse
	var responseErr error
	err := datasource.DB.Transaction(func(tx *gorm.DB) error {
		if searchMunicipalityType == "More" {
			response, responseErr = i.dao.NewItemQuery().SearchItem(tx, name, provinceId, municipalityId, cursor, false, 10)
			if responseErr != nil {
				return responseErr
			}
			if len(*response) > 10 {
				*response = (*response)[:len(*response)-1]
				searchItemResponse.NextPage = int32((*response)[len(*response)-1].Cursor)
				searchItemResponse.SearchMunicipalityType = pb.SearchMunicipalityType_More.String()
			} else if len(*response) <= 10 && len(*response) != 0 {
				length := 10 - len(*response)
				responseAdd, responseErr := i.dao.NewItemQuery().SearchItem(tx, name, provinceId, municipalityId, cursor, true, int64(length))
				if responseErr != nil {
					return responseErr
				}
				if len(*responseAdd) > length {
					*responseAdd = (*responseAdd)[:len(*responseAdd)-1]
				}
				*response = append(*response, *responseAdd...)
				searchItemResponse.NextPage = int32((*response)[len(*response)-1].Cursor)
				searchItemResponse.SearchMunicipalityType = pb.SearchMunicipalityType_NoMore.String()
			} else if len(*response) == 0 {
				response, responseErr = i.dao.NewItemQuery().SearchItem(tx, name, provinceId, municipalityId, cursor, true, 10)
				if responseErr != nil {
					return responseErr
				}
				if len(*response) > 10 {
					*response = (*response)[:len(*response)-1]
					searchItemResponse.NextPage = int32((*response)[len(*response)-1].Cursor)
				} else if len(*response) <= 10 && len(*response) != 0 {
					searchItemResponse.NextPage = int32((*response)[len(*response)-1].Cursor)
				}
				searchItemResponse.SearchMunicipalityType = pb.SearchMunicipalityType_NoMore.String()
			}
		} else {
			response, responseErr = i.dao.NewItemQuery().SearchItem(tx, name, provinceId, municipalityId, cursor, true, 10)
			if responseErr != nil {
				return responseErr
			}
			if len(*response) > 10 {
				*response = (*response)[:len(*response)-1]
				searchItemResponse.NextPage = int32((*response)[len(*response)-1].Cursor)
			} else if len(*response) <= 10 && len(*response) != 0 {
				searchItemResponse.NextPage = int32((*response)[len(*response)-1].Cursor)
			}
			searchItemResponse.SearchMunicipalityType = pb.SearchMunicipalityType_NoMore.String()
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	searchItemResponse.Items = response
	return &searchItemResponse, nil
}
