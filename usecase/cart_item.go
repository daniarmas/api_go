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
	ReduceCartItem(request *dto.ReduceCartItem) (*models.CartItem, error)
	DeleteCartItem(request *dto.DeleteCartItemRequest) error
}

type cartItemService struct {
	dao repository.DAO
}

func NewCartItemService(dao repository.DAO) CartItemService {
	return &cartItemService{dao: dao}
}

func (i *cartItemService) ListCartItemAndItem(itemRequest *dto.ListCartItemRequest) (*dto.ListCartItemResponse, error) {
	var items *[]models.CartItemAndItem
	var listItemResponse dto.ListCartItemResponse
	var itemsErr error
	err := datasource.DB.Transaction(func(tx *gorm.DB) error {
		authorizationTokenParseRes, authorizationTokenParseErr := repository.Datasource.NewJwtTokenDatasource().ParseJwtAuthorizationToken(&itemRequest.Metadata.Get("authorization")[0])
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
		user, userErr := i.dao.NewUserQuery().GetUser(tx, &models.User{ID: authorizationTokenRes.UserFk})
		if userErr != nil {
			return userErr
		}
		items, itemsErr = i.dao.NewCartItemRepository().ListCartItemAndItem(tx, &models.CartItem{UserFk: user.ID, Cursor: itemRequest.NextPage})
		if itemsErr != nil {
			return itemsErr
		} else if len(*items) > 10 {
			*items = (*items)[:len(*items)-1]
			listItemResponse.NextPage = (*items)[len(*items)-1].Cursor
		} else if len(*items) == 0 {
			listItemResponse.NextPage = itemRequest.NextPage
		} else {
			listItemResponse.NextPage = (*items)[len(*items)-1].Cursor
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
			CreateTime:           item.CreateTime,
			UpdateTime:           item.UpdateTime,
			Cursor:               item.Cursor,
			Quantity:             item.Quantity,
			ItemFk:               item.ItemFk,
			UserFk:               item.UserFk,
			AuthorizationTokenFk: item.AuthorizationTokenFk,
		})
	}
	listItemResponse.CartItems = *items
	return &listItemResponse, nil
}

func (i *cartItemService) AddCartItem(request *dto.AddCartItem) (*models.CartItem, error) {
	var result *models.CartItem
	var resultErr error
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
		item, itemErr := i.dao.NewItemQuery().GetItem(tx, request.ItemFk, request.Location)
		var itemAvailability int64
		if itemErr != nil {
			return itemErr
		}
		if !item.IsInRange {
			return errors.New("out of range")
		}
		if (item.Availability - int64(request.Quantity)) < 0 {
			return errors.New("no_availability:availability:" + strconv.Itoa(int(item.Availability)))
		} else if item.Availability-int64(request.Quantity) == 0 {
			itemAvailability = -1
		} else {
			itemAvailability = item.Availability - int64(request.Quantity)
		}
		_, updateItemErr := i.dao.NewItemQuery().UpdateItem(tx, &models.Item{ID: item.ID}, &models.Item{Availability: itemAvailability})
		if updateItemErr != nil {
			return updateItemErr
		}
		cartItemRes, cartItemErr := i.dao.NewCartItemRepository().GetCartItem(tx, &models.CartItem{ItemFk: uuid.MustParse(request.ItemFk), UserFk: authorizationTokenRes.UserFk})
		if cartItemErr != nil && cartItemErr.Error() != "record not found" {
			return errors.New("cartitem not found")
		} else if cartItemRes != nil {
			result, resultErr = i.dao.NewCartItemRepository().UpdateCartItem(tx, &models.CartItem{ItemFk: uuid.MustParse(request.ItemFk), UserFk: authorizationTokenRes.UserFk}, &models.CartItem{Quantity: cartItemRes.Quantity + request.Quantity})
			if resultErr != nil {
				return resultErr
			}
		} else if cartItemRes == nil && cartItemErr.Error() == "record not found" {
			result, resultErr = i.dao.NewCartItemRepository().CreateCartItem(tx, &models.CartItem{Name: item.Name, Price: item.Price, Quantity: request.Quantity, ItemFk: item.ID, UserFk: authorizationTokenRes.UserFk, AuthorizationTokenFk: authorizationTokenRes.ID, Cursor: 0})
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
		item, itemErr := i.dao.NewItemQuery().GetItem(tx, request.ItemFk, request.Location)
		if itemErr != nil && itemErr.Error() == "record not found" {
			return errors.New("item not found")
		} else if itemErr != nil {
			return itemErr
		}
		if !item.IsInRange {
			return errors.New("out of range")
		}
		if item.Availability == -1 {
			item.Availability += 1
		}
		result, resultErr = i.dao.NewCartItemRepository().GetCartItem(tx, &models.CartItem{ItemFk: uuid.MustParse(request.ItemFk), UserFk: authorizationTokenRes.UserFk})
		if resultErr != nil && resultErr.Error() != "record not found" {
			return errors.New("cartitem not found")
		}
		_, updateItemErr := i.dao.NewItemQuery().UpdateItem(tx, &models.Item{ID: item.ID}, &models.Item{Availability: item.Availability + 1})
		if updateItemErr != nil {
			return updateItemErr
		}
		if (result.Quantity - 1) == 0 {
			deleteCartItemErr := i.dao.NewCartItemRepository().DeleteCartItem(tx, &models.CartItem{ID: result.ID, UserFk: authorizationTokenRes.UserFk})
			if deleteCartItemErr != nil {
				return deleteCartItemErr
			}
			result = nil
		} else {
			result, resultErr = i.dao.NewCartItemRepository().UpdateCartItem(tx, &models.CartItem{ItemFk: uuid.MustParse(request.ItemFk), UserFk: authorizationTokenRes.UserFk}, &models.CartItem{Quantity: result.Quantity - 1})
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
		cartItemRes, cartItemErr := i.dao.NewCartItemRepository().GetCartItem(tx, &models.CartItem{ID: uuid.MustParse(request.CartItemFk), UserFk: authorizationTokenRes.UserFk})
		if cartItemErr != nil && cartItemErr.Error() != "record not found" {
			return errors.New("cartitem not found")
		}
		item, itemErr := i.dao.NewItemQuery().GetItem(tx, cartItemRes.ItemFk.String(), request.Location)
		if itemErr != nil {
			return itemErr
		}
		if item.Availability == -1 {
			item.Availability += 1
		}
		if !item.IsInRange {
			return errors.New("out of range")
		}
		_, updateItemErr := i.dao.NewItemQuery().UpdateItem(tx, &models.Item{ID: item.ID}, &models.Item{Availability: item.Availability + int64(cartItemRes.Quantity)})
		if updateItemErr != nil {
			return updateItemErr
		}
		deleteCartItemErr := i.dao.NewCartItemRepository().DeleteCartItem(tx, &models.CartItem{ID: cartItemRes.ID, UserFk: authorizationTokenRes.UserFk})
		if deleteCartItemErr != nil {
			return deleteCartItemErr
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
