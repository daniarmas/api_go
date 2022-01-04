package main

import (
	// "fmt"
	"fmt"
	"log"
	"net"

	// "runtime"
	pb "github.com/daniarmas/api_go/pkg"
	"github.com/daniarmas/api_go/app"
	"github.com/daniarmas/api_go/repository"
	"github.com/daniarmas/api_go/service"
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
	authenticationService := service.NewAuthenticationService(dao)

	// Starting gRPC server
	listener, err := net.Listen("tcp", "localhost:8282")
	if err != nil {
		log.Fatalln(err)
	}

	grpcServer := grpc.NewServer()
	// Registring the services
	pb.RegisterItemServiceServer(grpcServer, app.NewItemServer(
		itemService,
	))
	pb.RegisterAuthenticationServiceServer(grpcServer, app.NewAuthenticationServer(
		authenticationService,
	))
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Server running at localhost:8081")
}
