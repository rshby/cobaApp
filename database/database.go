package database

import (
	"cobaApp/config"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"time"
)

func ConnectDatabase(cfg config.IConfig, log *logrus.Logger) *sql.DB {
	config := cfg.GetConfig().Database

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		config.User, config.Password, config.Host, config.Port, config.Name)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("cant connect database : %v", err)
	}

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(30)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(20 * time.Minute)

	log.Infof("success connect database")

	return db
}
