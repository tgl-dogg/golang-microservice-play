package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// database is a connected database object
var database *gorm.DB

// GetDB implements Singleton pattern to keep a single connection for all queries.
func GetDB() *gorm.DB {
	return database
}

// Setup database connection based on parameters provided in the receiver.
func (dbConnection DBConnection) Setup() {
	dbInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", dbConnection.Host, dbConnection.Port, dbConnection.User, dbConnection.DBName, dbConnection.Password)
	db, err := gorm.Open(postgres.Open(dbInfo), &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}

	// db.LogMode(false)
	database = db
}

// DBConnection wraps information necessary to connect to a database.
type DBConnection struct {
	Host, Port, User, Password, DBName string
}
