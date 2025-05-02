// @title Beneficiary Manager API
// @version 1.0
// @description This is a sample API for beneficiary management.
// @host localhost:8080
// @BasePath /api/v1

package api

import (
	"net/http"
	"os"

	"github.com/ChayanDass/beneficiary-manager/pkg/middleware"
	"github.com/ChayanDass/beneficiary-manager/pkg/models"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	DEFAULT_PORT = "8080"
)

func Router() *gin.Engine {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = DEFAULT_PORT
	}

	// Initialize Gin router
	r := gin.Default()

	// Apply global middlewares
	r.Use(middleware.CORSMiddleware())

	// Handle invalid routes
	r.NoRoute(HandleInvalidUrl)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 group
	api := r.Group("/api/v1")
	{
		// Scheme Routes
		scheme := api.Group("/schemes")
		{
			scheme.GET("", GetSchemes)                 // Fetch available schemes
			scheme.GET("/:id", GetSchemeByID)          // Get a specific scheme
			scheme.GET("/status/:id", GetSchemeStatus) // Fetch scheme status
		}

		// Application Routes
		application := api.Group("/applications")
		application.Use(middleware.BasicAuth())
		{
			application.POST("/", SubmitApplication)                       // Submit application
			application.GET("/", GetApplications)                          // Get application status
			application.POST("/withdraw-application", WithdrawApplication) // Submit application without user ID
			application.POST("/init-application", InitApplication)         // Initialize application
			application.PUT("/:id", ModifyApplication)                     // Update application
			application.GET("/status/:id", GetApplicationStatus)           // Get application by ID

		}

	}

	return r
}

func HandleInvalidUrl(c *gin.Context) {
	er := models.ErrorResponse{
		Code:    http.StatusNotFound,
		Message: "No such path exists, please check the URL",
		Error:   "invalid path",
	}
	c.JSON(http.StatusNotFound, er)
}
