package app

import (
	pb "github.com/daniarmas/api_go/pkg"
	"github.com/daniarmas/api_go/usecase"
)

type ItemServer struct {
	pb.UnimplementedItemServiceServer
	itemService usecase.ItemService
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
