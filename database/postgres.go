package database

import (
	"fmt"
	"log"
	"onemildata/config"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type (
	Database interface {
		GetDb() *gorm.DB
	}

	postgresDb struct {
		Db *gorm.DB
	}
)

func NewPostgresDb(cfg *config.Config) Database {

	fmt.Printf(`%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True`,
		cfg.Db.User,
		cfg.Db.Password,
		cfg.Db.Host,
		cfg.Db.Port,
		cfg.Db.DbName)

	dsn := fmt.Sprintf(
		`%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True`,
		cfg.Db.User,
		cfg.Db.Password,
		cfg.Db.Host,
		cfg.Db.Port,
		cfg.Db.DbName,
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,         // Don't include params in the SQL log
			Colorful:                  false,         // Disable color
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("failed to connect database")
	}

	return &postgresDb{Db: db}
}

func (p *postgresDb) GetDb() *gorm.DB {
	return p.Db
}
