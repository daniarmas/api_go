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
	"gorm.io/gorm"
)

type ItemService interface {
	GetItem(request *dto.GetItemRequest) (*models.ItemBusiness, error)
	ListItem(itemRequest *dto.ListItemRequest) (*dto.ListItemResponse, error)
	SearchItem(name string, provinceFk string, municipalityFk string, cursor int64, searchMunicipalityType string) (*dto.SearchItemResponse, error)
	CreateItem(request *dto.CreateItemRequest) (*models.Item, error)
}

type itemService struct {
	dao repository.DAO
}

func NewItemService(dao repository.DAO) ItemService {
	return &itemService{dao: dao}
}

func (i *itemService) CreateItem(request *dto.CreateItemRequest) (*models.Item, error) {
	var itemRes *models.Item
	var itemErr error
	err := datasource.DB.Transaction(func(tx *gorm.DB) error {
		authorizationTokenParseRes, authorizationTokenParseErr := repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(&request.Metadata.Get("authorization")[0])
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
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: uuid.MustParse(*authorizationTokenParseRes)}, &[]string{"id", "user_fk"})
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		} else if authorizationTokenRes == nil {
			return errors.New("unauthenticated")
		}
		permissionExistsErr := i.dao.NewPermissionRepository().PermissionExists(tx, &models.Permission{Name: "create_item", UserFk: authorizationTokenRes.UserFk, BusinessFk: request.BusinessFk})
		if permissionExistsErr != nil && permissionExistsErr.Error() == "record not found" {
			return errors.New("permission denied")
		} else if permissionExistsErr != nil {
			return permissionExistsErr
		}
		_, hqErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), datasource.Config.ItemsBulkName, request.HighQualityPhotoObject)
		if hqErr != nil && hqErr.Error() == "ObjectMissing" {
			return errors.New("HighQualityPhotoObject missing")
		} else if hqErr != nil {
			return hqErr
		}
		_, lqErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), datasource.Config.ItemsBulkName, request.LowQualityPhotoObject)
		if lqErr != nil && lqErr.Error() == "ObjectMissing" {
			return errors.New("LowQualityPhotoObject missing")
		} else if lqErr != nil {
			return lqErr
		}
		_, tnErr := i.dao.NewObjectStorageRepository().ObjectExists(context.Background(), datasource.Config.ItemsBulkName, request.ThumbnailObject)
		if tnErr != nil && tnErr.Error() == "ObjectMissing" {
			return errors.New("ThumbnailObject missing")
		} else if tnErr != nil {
			return tnErr
		}
		businessRes, businessErr := i.dao.NewBusinessQuery().GetBusinessProvinceAndMunicipality(tx, request.BusinessFk)
		if businessErr != nil {
			return businessErr
		}
		itemRes, itemErr = i.dao.NewItemQuery().CreateItem(tx, &models.Item{Name: request.Name, Description: request.Description, Price: float64(request.Price), Availability: -1, BusinessFk: request.BusinessFk, BusinessItemCategoryFk: uuid.MustParse(request.BusinessItemCategoryFk), HighQualityPhoto: datasource.Config.ItemsBulkName + "/" + request.HighQualityPhotoObject, LowQualityPhoto: datasource.Config.ItemsBulkName + "/" + request.LowQualityPhotoObject, Thumbnail: datasource.Config.ItemsBulkName + "/" + request.ThumbnailObject, HighQualityPhotoObject: request.HighQualityPhotoObject, LowQualityPhotoObject: datasource.Config.ItemsBulkName + "/" + request.LowQualityPhotoObject, ThumbnailObject: datasource.Config.ItemsBulkName + "/" + request.ThumbnailObject, HighQualityPhotoBlurHash: request.HighQualityPhotoBlurHash, LowQualityPhotoBlurHash: request.LowQualityPhotoBlurHash, ThumbnailBlurHash: request.ThumbnailBlurHash, Status: "Unavailable", ProvinceFk: businessRes.ProvinceFk, MunicipalityFk: businessRes.MunicipalityFk})
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
		if itemRequest.BusinessFk != "" && itemRequest.BusinessItemCategoryFk == "" {
			itemCategoryRes, itemCategoryErr := i.dao.NewItemCategoryQuery().GetItemCategory(tx, &models.BusinessItemCategory{Index: 0, BusinessFk: uuid.MustParse(itemRequest.BusinessFk)})
			if itemCategoryErr != nil {
				return itemCategoryErr
			}
			items, itemsErr = i.dao.NewItemQuery().ListItem(tx, &models.Item{BusinessFk: uuid.MustParse(itemRequest.BusinessFk), BusinessItemCategoryFk: itemCategoryRes.ID}, itemRequest.NextPage)
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
		} else if itemRequest.BusinessFk != "" && itemRequest.BusinessItemCategoryFk != "" {
			items, itemsErr = i.dao.NewItemQuery().ListItem(tx, &models.Item{BusinessFk: uuid.MustParse(itemRequest.BusinessFk), BusinessItemCategoryFk: uuid.MustParse(itemRequest.BusinessItemCategoryFk)}, itemRequest.NextPage)
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
		item, itemErr = i.dao.NewItemQuery().GetItem(tx, request.Id, request.Location)
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

func (i *itemService) SearchItem(name string, provinceFk string, municipalityFk string, cursor int64, searchMunicipalityType string) (*dto.SearchItemResponse, error) {
	var response *[]models.Item
	var searchItemResponse dto.SearchItemResponse
	var responseErr error
	err := datasource.DB.Transaction(func(tx *gorm.DB) error {
		if searchMunicipalityType == "More" {
			response, responseErr = i.dao.NewItemQuery().SearchItem(tx, name, provinceFk, municipalityFk, cursor, false, 10)
			if responseErr != nil {
				return responseErr
			}
			if len(*response) > 10 {
				*response = (*response)[:len(*response)-1]
				searchItemResponse.NextPage = int32((*response)[len(*response)-1].Cursor)
				searchItemResponse.SearchMunicipalityType = pb.SearchMunicipalityType_More.String()
			} else if len(*response) <= 10 && len(*response) != 0 {
				length := 10 - len(*response)
				responseAdd, responseErr := i.dao.NewItemQuery().SearchItem(tx, name, provinceFk, municipalityFk, cursor, true, int64(length))
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
				response, responseErr = i.dao.NewItemQuery().SearchItem(tx, name, provinceFk, municipalityFk, cursor, true, 10)
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
			response, responseErr = i.dao.NewItemQuery().SearchItem(tx, name, provinceFk, municipalityFk, cursor, true, 10)
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
