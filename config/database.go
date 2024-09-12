package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"my-fiber-project/model"
)

var Database *gorm.DB

func Connect() error {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	// Log the values (be careful not to log sensitive information in production)
	log.Printf("DB_HOST: %s, DB_USER: %s, DB_NAME: %s, DB_PORT: %s", dbHost, dbUser, dbName, dbPort)

	if dbHost == "" || dbUser == "" || dbPassword == "" || dbName == "" || dbPort == "" {
		return fmt.Errorf("missing required environment variables")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	Database, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	err = Database.AutoMigrate(&model.Post{}, &model.Category{})
	if err != nil {
		return fmt.Errorf("failed to auto migrate: %v", err)
	}

	log.Println("Database connected successfully")
	return nil
}
