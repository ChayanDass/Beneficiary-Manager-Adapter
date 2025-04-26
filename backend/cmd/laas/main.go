// @title Beneficiary Manager API
// @version 1.0
// @description This is a sample API for beneficiary management.
// @host localhost:8080
// @BasePath /api/v1
// @schemes http https
// @accept json
// @produce json
// @openapi 3.0.0
package main

import (
	"flag"
	"log"
	"os"

	"github.com/ChayanDass/beneficiary-manager/pkg/api"
	"github.com/ChayanDass/beneficiary-manager/pkg/db"
	"github.com/ChayanDass/beneficiary-manager/pkg/models"
	"github.com/joho/godotenv"
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

	if err := db.DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Failed to automigrate database: %v", err)
	}

	if err := db.DB.AutoMigrate(&models.Application{}); err != nil {
		log.Fatalf("Failed to automigrate database: %v", err)
	}

	if err := db.DB.AutoMigrate(&models.Scheme{}); err != nil {
		log.Fatalf("Failed to automigrate database: %v", err)
	}

	if err := r.Run(); err != nil {
		log.Fatalf("Error while running the server: %v", err)
	}

}
