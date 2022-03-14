package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/daniarmas/api_go/app"
	"github.com/daniarmas/api_go/datasource"
	pb "github.com/daniarmas/api_go/pkg"
	"github.com/daniarmas/api_go/repository"
	"github.com/daniarmas/api_go/usecase"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	// Configurations
	now := time.Now()
	now = now.Add(time.Duration(20) * time.Minute)
	now = now.Add(time.Duration(3) * time.Hour * 24)
	fmt.Println(timestamppb.New(now.UTC()))
	fmt.Println(now.UTC())
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
	// Register all services
	repositoryDao := repository.NewDAO(db, config, datasourceDao)
	itemService := usecase.NewItemService(repositoryDao)
	authenticationService := usecase.NewAuthenticationService(repositoryDao)
	businessService := usecase.NewBusinessService(repositoryDao)
	userService := usecase.NewUserService(repositoryDao)
	cartItemService := usecase.NewCartItemService(repositoryDao)
	orderService := usecase.NewOrderService(repositoryDao)
	banService := usecase.NewBanService(repositoryDao)
	objectStorageService := usecase.NewObjectStorageService(repositoryDao)
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
	pb.RegisterCartItemServiceServer(grpcServer, app.NewCartItemServer(
		cartItemService,
	))
	pb.RegisterOrderServiceServer(grpcServer, app.NewOrderServer(
		orderService,
	))
	pb.RegisterBanServiceServer(grpcServer, app.NewBanServer(
		banService,
	))
	pb.RegisterObjectStorageServiceServer(grpcServer, app.NewObjectStorageServer(
		objectStorageService,
	))
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Server running at localhost:8081")
}
