package infra

import (
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const MYSQL_ERROR_FOREIGN_DOESNT_EXISTS uint16 = 1452

type Database struct {
	DB  *gorm.DB
	SQL *sql.DB
}

func ConfigDB(username, password, host, port, dbname string,
	connMaxLifetime time.Duration,
	maxIdleConns, maxOpenConns int) (*Database, error) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)

	sql, err := db.DB()
	if err != nil {
		return nil, err
	}

	return &Database{
		DB:  db,
		SQL: sql,
	}, nil
}

// Refers: https://gorm.io/docs/connecting_to_the_database.html#MySQL
// 		   https://gorm.io/docs/generic_interface.html#Connection-Pool
