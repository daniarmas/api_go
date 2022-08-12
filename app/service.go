package app

import (
	"github.com/daniarmas/api_go/internal/usecase"
	pb "github.com/daniarmas/api_go/pkg/grpc"
)

type ItemServer struct {
	pb.UnimplementedItemServiceServer
	itemService usecase.ItemService
}

type PaymentMethodServer struct {
	pb.UnimplementedPaymentMethodServiceServer
	paymentMethodService usecase.PaymentMethodService
}

type ApplicationServer struct {
	pb.UnimplementedApplicationServiceServer
	applicationService usecase.ApplicationService
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

type PermissionServer struct {
	pb.UnimplementedPermissionServiceServer
	permissionService usecase.PermissionService
}

func NewPermissionServer(
	permissionService usecase.PermissionService,
) *PermissionServer {
	return &PermissionServer{
		permissionService: permissionService,
	}
}

func NewOrderServer(
	orderService usecase.OrderService,
) *OrderServer {
	return &OrderServer{
		orderService: orderService,
	}
}

func NewPaymentMethodServer(
	paymentMethodService usecase.PaymentMethodService,
) *PaymentMethodServer {
	return &PaymentMethodServer{
		paymentMethodService: paymentMethodService,
	}
}

func NewApplicationServer(
	applicationService usecase.ApplicationService,
) *ApplicationServer {
	return &ApplicationServer{
		applicationService: applicationService,
	}
}

func NewObjectStorageServer(
	objectStorageService usecase.ObjectStorageService,
) *ObjectStorageServer {
	return &ObjectStorageServer{
		objectStorageService: objectStorageService,
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
