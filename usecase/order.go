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

type OrderService interface {
	ListOrder(request *dto.ListOrderRequest) (*dto.ListOrderResponse, error)
	CreateOrder(request *dto.CreateOrderRequest) (*dto.CreateOrderResponse, error)
	UpdateOrder(request *dto.UpdateOrderRequest) (*dto.UpdateOrderResponse, error)
}

type orderService struct {
	dao repository.DAO
}

func NewOrderService(dao repository.DAO) OrderService {
	return &orderService{dao: dao}
}

func (i *orderService) UpdateOrder(request *dto.UpdateOrderRequest) (*dto.UpdateOrderResponse, error) {
	var response dto.UpdateOrderResponse
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
		updateOrderRes, updateOrderErr := i.dao.NewOrderRepository().UpdateOrder(tx, &models.Order{ID: request.Id, UserFk: authorizationTokenRes.UserFk}, &models.Order{Status: request.OrderStatus})
		if updateOrderErr != nil {
			return updateOrderErr
		}
		response.Order = &models.Order{ID: updateOrderRes.ID, Status: updateOrderRes.Status, DeliveryType: updateOrderRes.DeliveryType, ResidenceType: updateOrderRes.ResidenceType, Price: updateOrderRes.Price, BuildingNumber: updateOrderRes.BuildingNumber, HouseNumber: updateOrderRes.HouseNumber, BusinessFk: updateOrderRes.BusinessFk, UserFk: updateOrderRes.UserFk, Coordinates: updateOrderRes.Coordinates, AuthorizationTokenFk: updateOrderRes.AuthorizationTokenFk, DeliveryDate: updateOrderRes.DeliveryDate, CreateTime: updateOrderRes.CreateTime, UpdateTime: updateOrderRes.UpdateTime}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (i *orderService) CreateOrder(request *dto.CreateOrderRequest) (*dto.CreateOrderResponse, error) {
	var response dto.CreateOrderResponse
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
		listCartItemRes, listCartItemErr := i.dao.NewCartItemRepository().ListCartItemInIds(tx, *request.CartItems)
		if listCartItemErr != nil {
			return listCartItemErr
		}
		var price float64 = 0
		orderedItems := make([]models.OrderedItem, 0, len(*listCartItemRes))
		for _, item := range *listCartItemRes {
			price += item.Price
			orderedItems = append(orderedItems, models.OrderedItem{Price: item.Price, ItemFk: item.ItemFk, UserFk: item.UserFk})
		}
		_, createOrderedItemsErr := i.dao.NewOrderedRepository().BatchCreateOrderedItem(tx, &orderedItems)
		if createOrderedItemsErr != nil {
			return createOrderedItemsErr
		}
		createOrderRes, createOrderErr := i.dao.NewOrderRepository().CreateOrder(tx, &models.Order{Status: request.Status, DeliveryType: request.DeliveryType, ResidenceType: request.ResidenceType, BuildingNumber: request.BuildingNumber, HouseNumber: request.HouseNumber, UserFk: authorizationTokenRes.UserFk, DeliveryDate: request.DeliveryDate, Coordinates: request.Coordinates, AuthorizationTokenFk: authorizationTokenRes.ID, BusinessFk: (*listCartItemRes)[0].BusinessFk, Price: price})
		if createOrderErr != nil {
			return createOrderErr
		}
		unionOrderAndOrderedItems := make([]models.UnionOrderAndOrderedItem, 0, len(orderedItems))
		for _, item := range orderedItems {
			unionOrderAndOrderedItems = append(unionOrderAndOrderedItems, models.UnionOrderAndOrderedItem{OrderFk: createOrderRes.ID, OrderedItemFk: item.ID})
		}
		_, createUnionOrderAndOrderedItemsErr := i.dao.NewUnionOrderAndOrderedItemRepository().BatchCreateUnionOrderAndOrderedItem(tx, &unionOrderAndOrderedItems)
		if createUnionOrderAndOrderedItemsErr != nil {
			return createUnionOrderAndOrderedItemsErr
		}
		deleteCartItemErr := i.dao.NewCartItemRepository().DeleteCartItem(tx, &models.CartItem{UserFk: authorizationTokenRes.UserFk})
		if deleteCartItemErr != nil {
			return deleteCartItemErr
		}
		response.Order = models.Order{Status: createOrderRes.Status, DeliveryType: createOrderRes.DeliveryType, ResidenceType: createOrderRes.ResidenceType, BuildingNumber: createOrderRes.BuildingNumber, HouseNumber: createOrderRes.HouseNumber, BusinessFk: createOrderRes.BusinessFk, AuthorizationTokenFk: createOrderRes.AuthorizationTokenFk, UserFk: createOrderRes.UserFk, DeliveryDate: createOrderRes.DeliveryDate, Coordinates: createOrderRes.Coordinates, Price: price, CreateTime: createOrderRes.CreateTime, UpdateTime: createOrderRes.UpdateTime, ID: createOrderRes.ID}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (i *orderService) ListOrder(request *dto.ListOrderRequest) (*dto.ListOrderResponse, error) {
	var listOrderResponse dto.ListOrderResponse
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
		ordersRes, ordersErr := i.dao.NewOrderRepository().ListOrder(tx, &models.Order{CreateTime: request.NextPage, UserFk: authorizationTokenRes.UserFk})
		if ordersErr != nil {
			return ordersErr
		}
		if len(*ordersRes) > 10 {
			*ordersRes = (*ordersRes)[:len(*ordersRes)-1]
			listOrderResponse.NextPage = (*ordersRes)[len(*ordersRes)-1].CreateTime
		} else if len(*ordersRes) > 0 {
			listOrderResponse.NextPage = (*ordersRes)[len(*ordersRes)-1].CreateTime
		} else {
			listOrderResponse.NextPage = request.NextPage
		}
		listOrderResponse.Orders = ordersRes
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &listOrderResponse, nil
}
