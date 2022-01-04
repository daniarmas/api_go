package app

import (
	pb "github.com/daniarmas/api_go/pkg"
	"github.com/daniarmas/api_go/service"
)

type ItemServer struct {
	pb.UnimplementedItemServiceServer
	itemService service.ItemService
}

type AuthenticationServer struct {
	pb.UnimplementedAuthenticationServiceServer
	authenticationService service.AuthenticationService
}

func NewItemServer(
	itemService service.ItemService,
) *ItemServer {
	return &ItemServer{
		itemService: itemService,
	}
}

func NewAuthenticationServer(
	authenticationService service.AuthenticationService,
) *AuthenticationServer {
	return &AuthenticationServer{
		authenticationService: authenticationService,
	}
}
