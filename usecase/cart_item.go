package usecase

import (
	"errors"

	"github.com/daniarmas/api_go/datasource"
	"github.com/daniarmas/api_go/dto"
	"github.com/daniarmas/api_go/models"
	"github.com/daniarmas/api_go/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartItemService interface {
	ListCartItemAndItem(itemRequest *dto.ListCartItemRequest) (*dto.ListCartItemResponse, error)
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
