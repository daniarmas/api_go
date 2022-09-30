package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	firebase "firebase.google.com/go"
	"github.com/daniarmas/api_go/app"
	"github.com/daniarmas/api_go/cli"
	"github.com/daniarmas/api_go/config"
	"github.com/daniarmas/api_go/internal/datasource"
	"github.com/daniarmas/api_go/internal/repository"
	"github.com/daniarmas/api_go/internal/usecase"
	pb "github.com/daniarmas/api_go/pkg/grpc"
	"github.com/daniarmas/api_go/pkg/rdb"
	"github.com/daniarmas/api_go/pkg/s3"
	"github.com/daniarmas/api_go/pkg/sqldb"
	"github.com/daniarmas/api_go/tlscert"
	"github.com/daniarmas/api_go/utils"
	"github.com/getsentry/sentry-go"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"google.golang.org/api/option"
)

func CustomMatcher(key string) (string, bool) {
	switch key {
	case "x-user-id":
		return key, true
	case "Access-Token":
		return key, true
	case "Device-Id":
		return key, true
	case "Platform":
		return key, true
	case "Firebase-Cloud-Messaging-Id":
		return key, true
	case "Model":
		return key, true
	case "Authorization":
		return key, true
	case "System-Version":
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}

func main() {
	// Sentry
	err := sentry.Init(sentry.ClientOptions{
		Dsn: "http://b60953b9e22c48a9bce354aeb2f65854@sentry.mool.cu:30957/2",
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)

	sentry.CaptureMessage("It works!")
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	// Starting API REST
	go func() {
		// mux
		mux := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(CustomMatcher))
		if _, err := os.Stat("mool-for-shopping-firebase-adminsdk-4vkol-f4cc371851.json"); err == nil {
			fmt.Printf("File exists\n")
		} else {
			fmt.Printf("File does not exist\n")
		}
		// Configurations
		cfg, err := config.New()
		if err != nil {
			log.Fatal("cannot load config:", err)
		}
		// Mool for shopping - Firebase Client
		// opt := []option.ClientOption{option.WithCredentialsJSON([]byte(cfg.MoolShoppingFirebase))}
		opt := option.WithCredentialsFile("mool-for-shopping-firebase-adminsdk-4vkol-f4cc371851.json")
		moolShoppingApp, err := firebase.NewApp(context.Background(), nil, opt)
		if err != nil && cfg.Environment != "development" {
			log.Fatalf("error initializing app: %v", err)
		}
		// Obtain a messaging.Client from the App.
		ctx := context.Background()
		moolShoppingClient, err := moolShoppingApp.Messaging(ctx)
		if err != nil && cfg.Environment != "development" {
			log.Fatalf("error getting Mool Shopping Firebase Messaging Client: %v\n", err)
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
		itemService := usecase.NewItemService(repository, cfg, stDb, sqlDb)
		authenticationService := usecase.NewAuthenticationService(repository, cfg, sqlDb, moolShoppingClient)
		businessService := usecase.NewBusinessService(repository, cfg, stDb, sqlDb)
		userService := usecase.NewUserService(repository, cfg, rdb, sqlDb)
		cartItemService := usecase.NewCartItemService(repository, cfg, sqlDb)
		orderService := usecase.NewOrderService(repository, sqlDb, cfg)
		objectStorageService := usecase.NewObjectStorageService(repository, sqlDb, cfg)
		applicationService := usecase.NewApplicationService(repository, sqlDb)
		paymentMethodService := usecase.NewPaymentMethodService(repository, cfg, rdb, sqlDb)
		permissionService := usecase.NewPermissionService(repository, sqlDb)
		pb.RegisterItemServiceHandlerServer(context.Background(), mux, app.NewItemServer(itemService))
		pb.RegisterAuthenticationServiceHandlerServer(context.Background(), mux, app.NewAuthenticationServer(authenticationService))
		pb.RegisterBusinessServiceHandlerServer(context.Background(), mux, app.NewBusinessServer(
			businessService,
		))
		pb.RegisterUserServiceHandlerServer(context.Background(), mux, app.NewUserServer(
			userService,
		))
		pb.RegisterCartItemServiceHandlerServer(context.Background(), mux, app.NewCartItemServer(
			cartItemService,
		))
		pb.RegisterOrderServiceHandlerServer(context.Background(), mux, app.NewOrderServer(
			orderService,
		))
		pb.RegisterObjectStorageServiceHandlerServer(context.Background(), mux, app.NewObjectStorageServer(
			objectStorageService,
		))
		pb.RegisterApplicationServiceHandlerServer(context.Background(), mux, app.NewApplicationServer(
			applicationService,
		))
		pb.RegisterPaymentMethodServiceHandlerServer(context.Background(), mux, app.NewPaymentMethodServer(
			paymentMethodService,
		))
		pb.RegisterPermissionServiceHandlerServer(context.Background(), mux, app.NewPermissionServer(
			permissionService,
		))
		pb.RegisterApplicationServiceHandlerServer(context.Background(), mux, app.NewApplicationServer(applicationService))
		// http server
		// cors.Default() setup the middleware with default options being
		// all origins accepted with simple methods (GET, POST). See
		// documentation below for more options.
		handler := cors.AllowAll().Handler(mux)
		port := fmt.Sprintf(":%s", cfg.ApiRestPort)
		log.Infof("Rest Server started on %s ", cfg.ApiRestPort)
		http.ListenAndServe(port, handler)
	}()
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
	grpcServerAddress := fmt.Sprintf("0.0.0.0:%s", cfg.ApiPort)
	startErr := s.Start(grpcServerAddress)
	if startErr != nil {
		log.Fatalf("%v", startErr)
	}
	s.AwaitTermination(func() {
		log.Print("Shutting down the server")
	})
}

func serviceRegister(sv *grpc.Server) {
	if _, err := os.Stat("mool-for-shopping-firebase-adminsdk-4vkol-f4cc371851.json"); err == nil {
		fmt.Printf("File exists\n")
	} else {
		fmt.Printf("File does not exist\n")
	}
	// Configurations
	cfg, err := config.New()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	// Mool for shopping - Firebase Client
	// opt := []option.ClientOption{option.WithCredentialsJSON([]byte(cfg.MoolShoppingFirebase))}
	opt := option.WithCredentialsFile("mool-for-shopping-firebase-adminsdk-4vkol-f4cc371851.json")
	moolShoppingApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil && cfg.Environment != "development" {
		log.Fatalf("error initializing app: %v", err)
	}
	// Obtain a messaging.Client from the App.
	ctx := context.Background()
	moolShoppingClient, err := moolShoppingApp.Messaging(ctx)
	if err != nil && cfg.Environment != "development" {
		log.Fatalf("error getting Mool Shopping Firebase Messaging Client: %v\n", err)
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
	authenticationService := usecase.NewAuthenticationService(repository, cfg, sqlDb, moolShoppingClient)
	businessService := usecase.NewBusinessService(repository, cfg, stDb, sqlDb)
	userService := usecase.NewUserService(repository, cfg, rdb, sqlDb)
	cartItemService := usecase.NewCartItemService(repository, cfg, sqlDb)
	orderService := usecase.NewOrderService(repository, sqlDb, cfg)
	objectStorageService := usecase.NewObjectStorageService(repository, sqlDb, cfg)
	applicationService := usecase.NewApplicationService(repository, sqlDb)
	paymentMethodService := usecase.NewPaymentMethodService(repository, cfg, rdb, sqlDb)
	permissionService := usecase.NewPermissionService(repository, sqlDb)
	pb.RegisterPermissionServiceServer(sv, app.NewPermissionServer(permissionService))
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
	pb.RegisterObjectStorageServiceServer(sv, app.NewObjectStorageServer(
		objectStorageService,
	))
	pb.RegisterApplicationServiceServer(sv, app.NewApplicationServer(applicationService))
	pb.RegisterPaymentMethodServiceServer(sv, app.NewPaymentMethodServer(
		paymentMethodService,
	))
}

func addInterceptors(s *utils.GrpcServerBuilder) {
	s.SetUnaryInterceptors(utils.GetDefaultUnaryServerInterceptors())
	s.SetStreamInterceptors(utils.GetDefaultStreamServerInterceptors())
}
