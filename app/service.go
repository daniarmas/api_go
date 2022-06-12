package app

import (
	pb "github.com/daniarmas/api_go/pkg"
	"github.com/daniarmas/api_go/usecase"
)

type ItemServer struct {
	pb.UnimplementedItemServiceServer
	itemService usecase.ItemService
}

type BanServer struct {
	pb.UnimplementedBanServiceServer
	banService usecase.BanService
}

type AnalyticsServer struct {
	pb.UnimplementedAnalyticsServiceServer
	analyticsService usecase.AnalyticsService
}

type ObjectStorageServer struct {
	pb.UnimplementedObjectStorageServiceServer
	objectStorageService usecase.ObjectStorageService
}

type OrderServer struct {
	pb.UnimplementedOrderServiceServer
	orderService usecase.OrderService
}

type CartItemServer struct {
	pb.UnimplementedCartItemServiceServer
	cartItemService usecase.CartItemService
}

type UserServer struct {
	pb.UnimplementedUserServiceServer
	userService usecase.UserService
}

type AuthenticationServer struct {
	pb.UnimplementedAuthenticationServiceServer
	authenticationService usecase.AuthenticationService
}

type BusinessServer struct {
	pb.UnimplementedBusinessServiceServer
	businessService usecase.BusinessService
}

func NewOrderServer(
	orderService usecase.OrderService,
) *OrderServer {
	return &OrderServer{
		orderService: orderService,
	}
}

func NewObjectStorageServer(
	objectStorageService usecase.ObjectStorageService,
) *ObjectStorageServer {
	return &ObjectStorageServer{
		objectStorageService: objectStorageService,
	}
}

func NewAnalyticsServer(
	analyticsService usecase.AnalyticsService,
) *AnalyticsServer {
	return &AnalyticsServer{
		analyticsService: analyticsService,
	}
}

func NewBanServer(
	banService usecase.BanService,
) *BanServer {
	return &BanServer{
		banService: banService,
	}
}

func NewItemServer(
	itemService usecase.ItemService,
) *ItemServer {
	return &ItemServer{
		itemService: itemService,
	}
}

func NewCartItemServer(
	cartItemService usecase.CartItemService,
) *CartItemServer {
	return &CartItemServer{
		cartItemService: cartItemService,
	}
}

func NewUserServer(
	userService usecase.UserService,
) *UserServer {
	return &UserServer{
		userService: userService,
	}
}

func NewAuthenticationServer(
	authenticationService usecase.AuthenticationService,
) *AuthenticationServer {
	return &AuthenticationServer{
		authenticationService: authenticationService,
	}
}

func NewBusinessServer(
	businessService usecase.BusinessService,
) *BusinessServer {
	return &BusinessServer{
		businessService: businessService,
	}
}
