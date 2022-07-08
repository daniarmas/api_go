package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/daniarmas/api_go/app"
	"github.com/daniarmas/api_go/cli"
	"github.com/daniarmas/api_go/config"
	"github.com/daniarmas/api_go/datasource"
	pb "github.com/daniarmas/api_go/pkg"
	"github.com/daniarmas/api_go/pkg/rdb"
	"github.com/daniarmas/api_go/pkg/s3"
	"github.com/daniarmas/api_go/pkg/sqldb"
	"github.com/daniarmas/api_go/repository"
	"github.com/daniarmas/api_go/tlscert"
	"github.com/daniarmas/api_go/usecase"
	"github.com/daniarmas/api_go/utils"
	_ "github.com/jackc/pgx/v4/stdlib"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	// Configurations
	cfg, err := config.New()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	// Starting gRPC server
	//if we crash the go code, we get the file name and line number
	// log.SetFlags(log.LstdFlags | log.Lshortfile)
	builder := utils.GrpcServerBuilder{}
	addInterceptors(&builder)
	if cfg.Environment == "development" {
		builder.EnableReflection(true)
	}
	if cfg.Tls == "true" {
		builder.SetTlsCert(&tlscert.Cert)
	}
	s := builder.Build()
	s.RegisterService(serviceRegister)
	grpcServerAddress := fmt.Sprintf("0.0.0.0:%d", cfg.ApiPort)
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
	cfg, err := config.New()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	// Redis
	rdb := rdb.New(cfg)
	// Standard Library Database Connection
	stDb, err := sql.Open("pgx",
		cfg.DBDsn)
	if err != nil {
		log.Fatal(err)
	}
	// defer stDb.Close()
	// Database GORM
	sqlDb, err := sqldb.New(cfg)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
		return
	}
	// ObjectStorageServer
	s3, err := s3.New(cfg)
	if err != nil {
		log.Fatalf("Error connecting to minio: %v", err)
	}
	// Datasource
	datasource := datasource.New(sqlDb.Gorm, cfg, s3)
	// Repository
	repository := repository.New(sqlDb.Gorm, cfg, datasource, rdb)
	// Handle the cli
	cli.HandleCli(os.Args, sqlDb.Gorm, cfg, repository)
	itemService := usecase.NewItemService(repository, cfg, stDb, sqlDb)
	authenticationService := usecase.NewAuthenticationService(repository, cfg, sqlDb)
	businessService := usecase.NewBusinessService(repository, cfg, stDb, sqlDb)
	userService := usecase.NewUserService(repository, cfg, rdb, sqlDb)
	cartItemService := usecase.NewCartItemService(repository, cfg, sqlDb)
	orderService := usecase.NewOrderService(repository, sqlDb)
	banService := usecase.NewBanService(repository, sqlDb)
	objectStorageService := usecase.NewObjectStorageService(repository, sqlDb, cfg)
	analyicsService := usecase.NewAnalyticsService(repository, stDb)
	applicationService := usecase.NewApplicationService(repository, sqlDb)
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
	pb.RegisterAnalyticsServiceServer(sv, app.NewAnalyticsServer(analyicsService))
	pb.RegisterApplicationServiceServer(sv, app.NewApplicationServer(applicationService))
}

func addInterceptors(s *utils.GrpcServerBuilder) {
	s.SetUnaryInterceptors(utils.GetDefaultUnaryServerInterceptors())
	s.SetStreamInterceptors(utils.GetDefaultStreamServerInterceptors())
}
