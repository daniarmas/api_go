package repository

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DAO interface {
	NewItemQuery() ItemQuery
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
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
		// log.Fatal("Failed to connect database: ", err)
	}
	return DB, nil
}

// func CloseDB() {
// 	DB.Close()
// }

func (d *dao) NewItemQuery() ItemQuery {
	return &itemQuery{}
}
