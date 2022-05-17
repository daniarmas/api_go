package usecase

import (
	"errors"
	"strconv"

	"github.com/daniarmas/api_go/datasource"
	"github.com/daniarmas/api_go/dto"
	"github.com/daniarmas/api_go/models"
	"github.com/daniarmas/api_go/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartItemService interface {
	ListCartItemAndItem(itemRequest *dto.ListCartItemRequest) (*dto.ListCartItemResponse, error)
	AddCartItem(request *dto.AddCartItem) (*models.CartItem, error)
	CartItemQuantity(request *dto.CartItemQuantity) (*bool, error)
	ReduceCartItem(request *dto.ReduceCartItem) (*models.CartItem, error)
	DeleteCartItem(request *dto.DeleteCartItemRequest) error
	EmptyCartItem(request *dto.EmptyCartItemRequest) error
}

type cartItemService struct {
	dao repository.DAO
}

func NewCartItemService(dao repository.DAO) CartItemService {
	return &cartItemService{dao: dao}
}

func (i *cartItemService) EmptyCartItem(request *dto.EmptyCartItemRequest) error {
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
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, nil)
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		} else if authorizationTokenRes == nil {
			return errors.New("unauthenticated")
		}
		listCartItemsRes, listCartItemsErr := i.dao.NewCartItemRepository().ListCartItemAll(tx, &models.CartItem{UserId: *authorizationTokenRes.UserId})
		if listCartItemsErr != nil {
			return listCartItemsErr
		}
		itemFks := make([]uuid.UUID, 0, len(*listCartItemsRes))
		for _, item := range *listCartItemsRes {
			itemFks = append(itemFks, item.ItemId)
		}
		itemsRes, itemsErr := i.dao.NewItemQuery().ListItemInIds(tx, itemFks)
		if itemsErr != nil {
			return itemsErr
		}
		for _, item := range *listCartItemsRes {
			var index = -1
			for i, n := range *itemsRes {
				if n.ID == item.ItemId {
					index = i
				}
			}
			(*itemsRes)[index].Availability += int64(item.Quantity)
		}
		for _, item := range *itemsRes {
			_, updateItemsErr := i.dao.NewItemQuery().UpdateItem(tx, &models.Item{ID: item.ID}, &item)
			if updateItemsErr != nil {
				return updateItemsErr
			}
		}
		_, err := i.dao.NewCartItemRepository().DeleteCartItem(tx, &models.CartItem{UserId: *authorizationTokenRes.UserId}, nil)
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

func (i *cartItemService) CartItemQuantity(request *dto.CartItemQuantity) (*bool, error) {
	var cartItemQuantityRes *bool
	var cartItemQuantityErr error
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
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, nil)
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		} else if authorizationTokenRes == nil {
			return errors.New("unauthenticated")
		}
		cartItemQuantityRes, cartItemQuantityErr = i.dao.NewCartItemRepository().CartItemQuantity(tx, &models.CartItem{UserId: *authorizationTokenRes.UserId})
		if cartItemQuantityErr != nil {
			return cartItemQuantityErr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return cartItemQuantityRes, nil
}

func (i *cartItemService) ListCartItemAndItem(itemRequest *dto.ListCartItemRequest) (*dto.ListCartItemResponse, error) {
	var items *[]models.CartItem
	var listItemResponse dto.ListCartItemResponse
	var itemsErr error
	err := datasource.DB.Transaction(func(tx *gorm.DB) error {
		jwtAuthorizationToken := &datasource.JsonWebTokenMetadata{Token: &itemRequest.Metadata.Get("authorization")[0]}
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
		user, userErr := i.dao.NewUserQuery().GetUser(tx, &models.User{ID: *authorizationTokenRes.UserId}, &[]string{})
		if userErr != nil {
			return userErr
		}
		items, itemsErr = i.dao.NewCartItemRepository().ListCartItemAndItem(tx, &models.CartItem{UserId: user.ID}, &itemRequest.NextPage)
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
		return nil
	})
	if err != nil {
		return nil, err
	}
	itemsResponse := make([]dto.CartItem, 0, len(*items))
	for _, item := range *items {
		itemsResponse = append(itemsResponse, dto.CartItem{
			ID:                   item.ID,
			Name:                 item.Name,
			Price:                item.Price,
			Thumbnail:            item.Thumbnail,
			ThumbnailBlurHash:    item.ThumbnailBlurHash,
			Quantity:             item.Quantity,
			ItemId:               item.ItemId,
			UserId:               item.UserId,
			AuthorizationTokenId: item.AuthorizationTokenId,
			CreateTime:           item.CreateTime,
			UpdateTime:           item.UpdateTime,
		})
	}
	listItemResponse.CartItems = *items
	return &listItemResponse, nil
}

func (i *cartItemService) AddCartItem(request *dto.AddCartItem) (*models.CartItem, error) {
	var result *models.CartItem
	var resultErr error
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
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, nil)
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		} else if authorizationTokenRes == nil {
			return errors.New("unauthenticated")
		}
		item, itemErr := i.dao.NewItemQuery().GetItemWithLocation(tx, request.ItemId, request.Location)
		var itemAvailability int64
		if itemErr != nil {
			return itemErr
		}
		if (item.Availability - int64(request.Quantity)) < 0 {
			return errors.New("no_availability:availability:" + strconv.Itoa(int(item.Availability)))
		} else if item.Availability-int64(request.Quantity) == 0 {
			itemAvailability = -1
		} else {
			itemAvailability = item.Availability - int64(request.Quantity)
		}
		_, err := i.dao.NewUnionBusinessAndMunicipalityRepository().GetUnionBusinessAndMunicipality(tx, &models.UnionBusinessAndMunicipality{BusinessId: item.BusinessId, MunicipalityId: request.MunicipalityId}, &[]string{"id"})
		if err != nil && err.Error() == "record not found" {
			return errors.New("out of range")
		} else if err != nil {
			return err
		}
		_, updateItemErr := i.dao.NewItemQuery().UpdateItem(tx, &models.Item{ID: item.ID}, &models.Item{Availability: itemAvailability})
		if updateItemErr != nil {
			return updateItemErr
		}
		cartItemRes, cartItemErr := i.dao.NewCartItemRepository().GetCartItem(tx, &models.CartItem{ItemId: uuid.MustParse(request.ItemId), UserId: *authorizationTokenRes.UserId}, nil)
		if cartItemErr != nil && cartItemErr.Error() != "record not found" {
			return errors.New("cartitem not found")
		} else if cartItemRes != nil {
			result, resultErr = i.dao.NewCartItemRepository().UpdateCartItem(tx, &models.CartItem{ItemId: uuid.MustParse(request.ItemId), UserId: *authorizationTokenRes.UserId}, &models.CartItem{Quantity: cartItemRes.Quantity + request.Quantity})
			if resultErr != nil {
				return resultErr
			}
		} else if cartItemRes == nil && cartItemErr.Error() == "record not found" {
			result, resultErr = i.dao.NewCartItemRepository().CreateCartItem(tx, &models.CartItem{Name: item.Name, Price: item.Price, Quantity: request.Quantity, ItemId: item.ID, UserId: *authorizationTokenRes.UserId, AuthorizationTokenId: *authorizationTokenRes.ID, BusinessId: item.BusinessId})
			if resultErr != nil {
				return resultErr
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *cartItemService) ReduceCartItem(request *dto.ReduceCartItem) (*models.CartItem, error) {
	var result *models.CartItem
	var resultErr error
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
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, nil)
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		} else if authorizationTokenRes == nil {
			return errors.New("unauthenticated")
		}
		item, itemErr := i.dao.NewItemQuery().GetItemWithLocation(tx, request.ItemId, request.Location)
		if itemErr != nil && itemErr.Error() == "record not found" {
			return errors.New("item not found")
		} else if itemErr != nil {
			return itemErr
		}
		if item.Availability == -1 {
			item.Availability += 1
		}
		_, err := i.dao.NewUnionBusinessAndMunicipalityRepository().GetUnionBusinessAndMunicipality(tx, &models.UnionBusinessAndMunicipality{BusinessId: item.BusinessId, MunicipalityId: request.MunicipalityId}, &[]string{"id"})
		if err != nil && err.Error() == "record not found" {
			return errors.New("out of range")
		} else if err != nil {
			return err
		}
		result, resultErr = i.dao.NewCartItemRepository().GetCartItem(tx, &models.CartItem{ItemId: uuid.MustParse(request.ItemId), UserId: *authorizationTokenRes.UserId}, nil)
		if resultErr != nil && resultErr.Error() != "record not found" {
			return errors.New("cartitem not found")
		}
		_, updateItemErr := i.dao.NewItemQuery().UpdateItem(tx, &models.Item{ID: item.ID}, &models.Item{Availability: item.Availability + 1})
		if updateItemErr != nil {
			return updateItemErr
		}
		if (result.Quantity - 1) == 0 {
			_, err := i.dao.NewCartItemRepository().DeleteCartItem(tx, &models.CartItem{ID: result.ID, UserId: *authorizationTokenRes.UserId}, nil)
			if err != nil {
				return err
			}
			result = nil
		} else {
			result, resultErr = i.dao.NewCartItemRepository().UpdateCartItem(tx, &models.CartItem{ItemId: uuid.MustParse(request.ItemId), UserId: *authorizationTokenRes.UserId}, &models.CartItem{Quantity: result.Quantity - 1})
			if resultErr != nil {
				return resultErr
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *cartItemService) DeleteCartItem(request *dto.DeleteCartItemRequest) error {
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
		authorizationTokenRes, authorizationTokenErr := i.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: jwtAuthorizationToken.TokenId}, nil)
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		} else if authorizationTokenRes == nil {
			return errors.New("unauthenticated")
		}
		cartItemRes, cartItemErr := i.dao.NewCartItemRepository().GetCartItem(tx, &models.CartItem{ID: uuid.MustParse(request.CartItemId), UserId: *authorizationTokenRes.UserId}, nil)
		if cartItemErr != nil && cartItemErr.Error() != "record not found" {
			return errors.New("cartitem not found")
		}
		item, itemErr := i.dao.NewItemQuery().GetItemWithLocation(tx, cartItemRes.ItemId.String(), request.Location)
		if itemErr != nil {
			return itemErr
		}
		if item.Availability == -1 {
			item.Availability += 1
		}
		_, municipalityErr := i.dao.NewUnionBusinessAndMunicipalityRepository().GetUnionBusinessAndMunicipality(tx, &models.UnionBusinessAndMunicipality{BusinessId: item.BusinessId, MunicipalityId: request.MunicipalityId}, &[]string{"id"})
		if municipalityErr != nil && municipalityErr.Error() == "record not found" {
			return errors.New("out of range")
		} else if municipalityErr != nil {
			return municipalityErr
		}
		_, updateItemErr := i.dao.NewItemQuery().UpdateItem(tx, &models.Item{ID: item.ID}, &models.Item{Availability: item.Availability + int64(cartItemRes.Quantity)})
		if updateItemErr != nil {
			return updateItemErr
		}
		_, err := i.dao.NewCartItemRepository().DeleteCartItem(tx, &models.CartItem{ID: cartItemRes.ID, UserId: *authorizationTokenRes.UserId}, nil)
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
