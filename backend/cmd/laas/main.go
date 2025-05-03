// @title Beneficiary Manager API
// @version 1.0
// @description This is a backend service for managing schemes, users, and applications.
// @contact.name Chayan
// @host localhost:8080
// @BasePath /api
package main

import (
	"flag"
	"log"
	"os"

	"github.com/ChayanDass/beneficiary-manager/pkg/api"
	"github.com/ChayanDass/beneficiary-manager/pkg/db"
	"github.com/ChayanDass/beneficiary-manager/pkg/models"
	"github.com/joho/godotenv"
	"gorm.io/gorm/clause"
)

// declare flags to input the basic requirement of database connection and the path of the data file
var (
	dbhost   = flag.String("host", getEnv("DB_HOST", "localhost"), "host name")
	port     = flag.String("port", getEnv("DB_PORT", "5432"), "port number")
	user     = flag.String("user", getEnv("DB_USER", "postgres"), "user name")
	dbname   = flag.String("dbname", getEnv("DB_NAME", "onset_adaptar"), "database name")
	password = flag.String("password", getEnv("DB_PASSWORD", "postgres"), "password")
)

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	flag.Parse()
	db.Connect(dbhost, port, user, dbname, password)
	r := api.Router()

	if err := db.DB.AutoMigrate(&models.Application{}); err != nil {
		log.Fatalf("Failed to automigrate database: %v", err)
	}
	if err := db.DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Failed to automigrate database: %v", err)

	}
	if err := db.DB.AutoMigrate(&models.StudentAcademicQualification{}); err != nil {
		log.Fatalf("Failed to automigrate database: %v", err)

	}
	if err := db.DB.AutoMigrate(&models.Scheme{}); err != nil {
		log.Fatalf("Failed to automigrate database: %v", err)
	}
	if err := db.DB.AutoMigrate(&models.Eligibility{}); err != nil {
		log.Fatalf("Failed to automigrate database: %v", err)
	}
	if err := db.DB.AutoMigrate(&models.Address{}); err != nil {
		log.Fatalf("Failed to automigrate database: %v", err)
	}
	if err := db.DB.AutoMigrate(&models.StudentProfile{}); err != nil {
		log.Fatalf("Failed to automigrate database: %v", err)

	}

	if err := db.DB.AutoMigrate(&models.UploadDocument{}); err != nil {
		log.Fatalf("Failed to automigrate database: %v", err)
	}

	if err := db.DB.AutoMigrate(&models.DocumentsRequired{}); err != nil {
		log.Fatalf("Failed to automigrate database: %v", err)
	}

	if err := db.DB.AutoMigrate(&models.EligibilityDocumentMap{}); err != nil {
		log.Fatalf("Failed to automigrate database: %v", err)
	}
	if err := db.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&models.DefaultDocumentsRequired).Error; err != nil {
		log.Fatalf("Failed to seed database with default documents types: %s", err.Error())
	}

	if err := r.Run(); err != nil {
		log.Fatalf("Error while running the server: %v", err)
	}

}
