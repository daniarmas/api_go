package sqldb

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/daniarmas/api_go/config"
	logg "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Sql struct {
	Gorm  *gorm.DB
	SqlDb *sql.DB
}

func New(cfg config.Config) (*Sql, error) {
	// Starting a database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBDatabase, cfg.DBPort)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Millisecond * 200,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		},
	)
	gorm, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 newLogger,
	})
	if err != nil {
		logg.Error(err)
	}
	connect, err := gorm.DB()
	if err != nil {
		return nil, err
	}
	var dbError error
	maxAttempts := 20
	for attempts := 1; attempts <= maxAttempts; attempts++ {
		dbError = connect.Ping()
		if dbError == nil {
			break
		}
		logg.Error(dbError)
		time.Sleep(time.Duration(attempts) * time.Second)
	}
	if dbError != nil {
		logg.Error(dbError)
	}
	return &Sql{
		SqlDb: connect,
		Gorm:  gorm,
	}, nil
}

// Close -.
func (g *Sql) Close() {
	if g.SqlDb != nil {
		g.SqlDb.Close()
		logg.Info("sql server connection closed")
	}
}
