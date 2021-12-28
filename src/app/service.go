package app

import (
	pb "github.com/daniarmas/api_go/pkg"
	"github.com/daniarmas/api_go/src/service"
)

type ItemServer struct {
	pb.UnimplementedItemServiceServer
	itemService service.ItemService
}

func NewItemServer(
	itemService service.ItemService,
) *ItemServer {
	return &ItemServer{
		itemService: itemService,
	}
}
