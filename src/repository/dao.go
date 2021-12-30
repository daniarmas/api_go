package repository

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DAO interface {
	NewItemQuery() ItemQuery
	NewVerificationCodeQuery() VerificationCodeQuery
}

type dao struct{}

var DB *gorm.DB

func NewDAO(db *gorm.DB) DAO {
	DB = db
	return &dao{}
}

func NewDB() (*gorm.DB, error) {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	host := viper.Get("DB_HOST").(string)
	port := viper.Get("DB_PORT").(string)
	user := viper.Get("DB_USER").(string)
	dbName := viper.Get("DB_DATABASE").(string)
	password := viper.Get("DB_PASSWORD").(string)

	// Starting a database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbName, port)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Millisecond * 200,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		},
	)
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
		// log.Fatal("Failed to connect database: ", err)
	}
	return DB, nil
}

func (d *dao) NewItemQuery() ItemQuery {
	return &itemQuery{}
}

func (d *dao) NewVerificationCodeQuery() VerificationCodeQuery {
	return &verificationCodeQuery{}
}
