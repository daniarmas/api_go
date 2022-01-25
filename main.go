package main

import (
	"fmt"
	"log"
	"net"

	"github.com/daniarmas/api_go/app"
	"github.com/daniarmas/api_go/datasource"
	pb "github.com/daniarmas/api_go/pkg"
	"github.com/daniarmas/api_go/repository"
	"github.com/daniarmas/api_go/usecase"
	"google.golang.org/grpc"
)

func main() {
	// Configurations
	config, err := datasource.NewConfig()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	// Database
	db, err := datasource.NewDB(config)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
		return
	}
	// ObjectStorageServer
	objectStorage, objectStorageErr := datasource.NewMinioClient(config)
	if objectStorageErr != nil {
		log.Fatalf("Error connecting to minio: %v", objectStorageErr)
	}
	// Datasource
	datasourceDao := datasource.NewDAO(db, config, objectStorage)
	fmt.Println(datasourceDao)
	// Register all services
	dao := repository.NewDAO(db, config)
	itemService := usecase.NewItemService(dao)
	authenticationService := usecase.NewAuthenticationService(dao)
	businessService := usecase.NewBusinessService(dao)
	userService := usecase.NewUserService(dao)

	// Starting gRPC server
	address := fmt.Sprintf("0.0.0.0:%s", config.ApiPort)
	listener, err := net.Listen("tcp", address)
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
	pb.RegisterBusinessServiceServer(grpcServer, app.NewBusinessServer(
		businessService,
	))
	pb.RegisterUserServiceServer(grpcServer, app.NewUserServer(
		userService,
	))
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Server running at localhost:8081")
}
