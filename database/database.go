// This package connects the application to the database and works with the ORM to setup an instance
package database

import (
	"log"
	"os"

	"github.com/adramelech-123/fiber-api/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// The database instance will point to gorm.DB struct and methods
type DbInstance struct {
	Db *gorm.DB
}

var Database DbInstance

func ConnectDb() {
	// Variable db will open our sqlite database and provide the initial config
	db, err := gorm.Open(sqlite.Open("fiberapi.db"), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to the database! \n", err.Error())
		os.Exit(2)
	}

	log.Println("Database connection successful! 😁")

	// A logger to log into our database
	db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("Running migrations...")
	
	// Add migrations
	db.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{})

	Database = DbInstance{Db: db}
}