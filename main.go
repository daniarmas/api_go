package main

import (
	"fmt"
	// "log"

	"github.com/daniarmas/api_go/app"
	"github.com/daniarmas/api_go/datasource"
	pb "github.com/daniarmas/api_go/pkg"
	log "github.com/sirupsen/logrus"
	"github.com/daniarmas/api_go/repository"
	"github.com/daniarmas/api_go/tlscert"
	"github.com/daniarmas/api_go/usecase"
	"github.com/daniarmas/api_go/utils"
	"google.golang.org/grpc"
)

func main() {
	// Configurations
	config, err := datasource.NewConfig()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	// Starting gRPC server
	//if we crash the go code, we get the file name and line number
	// log.SetFlags(log.LstdFlags | log.Lshortfile)
	builder := utils.GrpcServerBuilder{}
	addInterceptors(&builder)
	if config.Environment == "development" {
		builder.EnableReflection(true)
	}
	if config.Tls == "true" {
		builder.SetTlsCert(&tlscert.Cert)
	}
	s := builder.Build()
	s.RegisterService(serviceRegister)
	grpcServerAddress := fmt.Sprintf("0.0.0.0:%d", config.ApiPort)
	startErr := s.Start(grpcServerAddress)
	if startErr != nil {
		log.Fatalf("%v", startErr)
	}
	s.AwaitTermination(func() {
		log.Print("Shutting down the server")
	})
}

func serviceRegister(sv *grpc.Server) {
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
	// ObjectStorageServer5
	objectStorage, objectStorageErr := datasource.NewMinioClient(config)
	if objectStorageErr != nil {
		log.Fatalf("Error connecting to minio: %v", objectStorageErr)
	}
	// Datasource
	datasourceDao := datasource.NewDAO(db, config, objectStorage)
	// Repository
	repositoryDao := repository.NewDAO(db, config, datasourceDao)
	itemService := usecase.NewItemService(repositoryDao, config)
	authenticationService := usecase.NewAuthenticationService(repositoryDao, config)
	businessService := usecase.NewBusinessService(repositoryDao, config)
	userService := usecase.NewUserService(repositoryDao)
	cartItemService := usecase.NewCartItemService(repositoryDao, config)
	orderService := usecase.NewOrderService(repositoryDao)
	banService := usecase.NewBanService(repositoryDao)
	objectStorageService := usecase.NewObjectStorageService(repositoryDao)
	pb.RegisterItemServiceServer(sv, app.NewItemServer(
		itemService,
	))
	pb.RegisterAuthenticationServiceServer(sv, app.NewAuthenticationServer(
		authenticationService,
	))
	pb.RegisterBusinessServiceServer(sv, app.NewBusinessServer(
		businessService,
	))
	pb.RegisterUserServiceServer(sv, app.NewUserServer(
		userService,
	))
	pb.RegisterCartItemServiceServer(sv, app.NewCartItemServer(
		cartItemService,
	))
	pb.RegisterOrderServiceServer(sv, app.NewOrderServer(
		orderService,
	))
	pb.RegisterBanServiceServer(sv, app.NewBanServer(
		banService,
	))
	pb.RegisterObjectStorageServiceServer(sv, app.NewObjectStorageServer(
		objectStorageService,
	))
}

func addInterceptors(s *utils.GrpcServerBuilder) {
	s.SetUnaryInterceptors(utils.GetDefaultUnaryServerInterceptors())
	s.SetStreamInterceptors(utils.GetDefaultStreamServerInterceptors())
}
