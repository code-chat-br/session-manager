package db

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ManagerDB map[string]*gorm.DB

func InitializerDB(db_name string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(db_name+".db"), &gorm.Config{
		Logger: logger.New(log.New(log.Writer(), "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold:             200,
			LogLevel:                  logger.Error,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		}),
	})
}
