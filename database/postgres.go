package database

import (
	"fmt"
	"onemildata/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return &postgresDb{Db: db}
}

func (p *postgresDb) GetDb() *gorm.DB {
	return p.Db
}
