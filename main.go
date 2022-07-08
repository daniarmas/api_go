package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/daniarmas/api_go/app"
	"github.com/daniarmas/api_go/cli"
	"github.com/daniarmas/api_go/datasource"
	pb "github.com/daniarmas/api_go/pkg"
	"github.com/daniarmas/api_go/pkg/sqldb"
	"github.com/daniarmas/api_go/repository"
	"github.com/daniarmas/api_go/tlscert"
	"github.com/daniarmas/api_go/usecase"
	"github.com/daniarmas/api_go/utils"
	"github.com/go-redis/redis/v9"
	_ "github.com/jackc/pgx/v4/stdlib"
	log "github.com/sirupsen/logrus"
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
	cfg, err := datasource.NewConfig()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	// Redis
	rdbAddres := fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort)
	rdb := redis.NewClient(&redis.Options{
		Addr: rdbAddres,
		// Password: config.RedisPassword, // no password set
		// DB:       config.RedisDb,       // use default DB
	})
	// Standard Library Database Connection
	stDb, err := sql.Open("pgx",
		cfg.DBDsn)
	if err != nil {
		log.Fatal(err)
	}
	// defer stDb.Close()
	// Database GORM
	db, err := sqldb.New(cfg)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
		return
	}
	// ObjectStorageServer
	objectStorage, objectStorageErr := datasource.NewMinioClient(cfg)
	if objectStorageErr != nil {
		log.Fatalf("Error connecting to minio: %v", objectStorageErr)
	}
	// Datasource
	datasourceDao := datasource.NewDAO(db.Gorm, cfg, objectStorage)
	// Repository
	repositoryDao := repository.NewDAO(db.Gorm, cfg, datasourceDao, rdb)
	// Handle the cli
	cli.HandleCli(os.Args, db.Gorm, cfg, repositoryDao)
	itemService := usecase.NewItemService(repositoryDao, cfg, stDb)
	authenticationService := usecase.NewAuthenticationService(repositoryDao, cfg)
	businessService := usecase.NewBusinessService(repositoryDao, cfg, stDb)
	userService := usecase.NewUserService(repositoryDao, cfg, rdb)
	cartItemService := usecase.NewCartItemService(repositoryDao, cfg)
	orderService := usecase.NewOrderService(repositoryDao)
	banService := usecase.NewBanService(repositoryDao)
	objectStorageService := usecase.NewObjectStorageService(repositoryDao)
	analyicsService := usecase.NewAnalyticsService(repositoryDao, stDb)
	applicationService := usecase.NewApplicationService(repositoryDao)
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
