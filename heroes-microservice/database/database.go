package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is a connected database object
var DB *gorm.DB

func (dbConnection DBConnection) Setup() {
	dbInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", dbConnection.Host, dbConnection.Port, dbConnection.User, dbConnection.DBName, dbConnection.Password)
	db, err := gorm.Open(postgres.Open(dbInfo), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// db.LogMode(false)
	// db.AutoMigrate([]heroes.Race{})
	DB = db
}

// GetDB implements Singleton pattern to keep a single connection for all queries.
func GetDB() *gorm.DB {
	return DB
}

type DBConnection struct {
	Host, Port, User, Password, DBName string
}
