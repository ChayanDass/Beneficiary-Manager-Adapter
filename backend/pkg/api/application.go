package api

import (
	"net/http"

	"github.com/ChayanDass/beneficiary-manager/pkg/db"
	"github.com/ChayanDass/beneficiary-manager/pkg/models"
	"github.com/gin-gonic/gin"
)

// SubmitApplication godoc
// @Summary Submit an application
// @Description Submits an application for a scheme
// @Tags Applications
// @Accept json
// @Produce json
// @Param application body models.Application true "Application Data"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/applications [post]

func SubmitApplication(c *gin.Context) {
	var application models.Application

	if err := c.ShouldBindJSON(&application); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid input data",
			Error:   err.Error(),
		})
		return
	}

	var scheme models.Scheme
	if err := db.DB.First(&scheme, application.SchemeID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "Scheme not found",
			Error:   err.Error(),
		})
		return
	}

	if err := db.DB.Create(&application).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to submit application",
			Error:   err.Error(),
		})
		return
	}

	db.DB.Preload("User").First(&application, application.ID)

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Application submitted successfully",
		Data:    application,
	})
}

// GetApplication godoc
// @Summary Get application details
// @Description Retrieves details of a submitted application
// @Tags Applications
// @Accept json
// @Produce json
// @Param id path string true "Application ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/v1/applications/{id} [get]
func GetApplication(c *gin.Context) {
	id := c.Param("id")

	var application models.Application
	if err := db.DB.Preload("User").First(&application, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "Application not found",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Application fetched successfully",
		Data:    application,
	})
}

// GetApplicationStatus godoc
// @Summary Get application status
// @Description Retrieves the status of an application
// @Tags Applications
// @Accept json
// @Produce json
// @Param id path string true "Application ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/v1/applications/{id}/status [get]
func GetApplicationStatus(c *gin.Context) {
	id := c.Param("id")

	var application models.Application
	if err := db.DB.First(&application, id).Preload("User").Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "Application not found",
			Error:   err.Error(),
		})
		return
	}

}
