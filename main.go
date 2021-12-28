package main

import (
	// "fmt"
	"fmt"
	"log"
	"net"

	// "runtime"
	pb "github.com/daniarmas/api_go/pkg"
	"github.com/daniarmas/api_go/src/app"
	"github.com/daniarmas/api_go/src/repository"
	"github.com/daniarmas/api_go/src/service"
	"google.golang.org/grpc"
)

func main() {
	// fmt.Println(runtime.NumCPU())

	// DB
	db, err := repository.NewDB()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
		return
	}

	// Register all services
	dao := repository.NewDAO(db)
	itemService := service.NewItemService(dao)

	// Starting gRPC server
	listener, err := net.Listen("tcp", "localhost:8282")
	if err != nil {
		log.Fatalln(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterItemServiceServer(grpcServer, app.NewItemServer(
		itemService,
	))

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Server running at localhost:8081")
}
