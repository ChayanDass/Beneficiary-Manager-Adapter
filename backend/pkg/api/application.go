package api

import (
	"net/http"
	"time"

	"github.com/ChayanDass/beneficiary-manager/pkg/db"
	"github.com/ChayanDass/beneficiary-manager/pkg/models"
	"github.com/ChayanDass/beneficiary-manager/pkg/utils"
	"github.com/gin-gonic/gin"
)

func GetApplications(c *gin.Context) {
	var applications []models.Application

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "User ID not found in context",
		})
		return
	}

	// Ensure userID is of the correct type (uint)
	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "Invalid user ID type",
		})
		return
	}

	// Fetch all applications for the user
	if err := db.DB.
		Preload("User").
		Preload("StudentProfile").
		Preload("StudentProfile.Addresses").
		Preload("StudentProfile.EducationHistory").
		Preload("StudentProfile.Documents").
		Where("user_id = ?", userIDUint).Find(&applications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to fetch applications",
			Error:   err.Error(),
		})
		return
	}

	// If no applications are found
	if len(applications) == 0 {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "No applications found for this user",
		})
		return
	}

	// Return the list of applications
	c.JSON(http.StatusOK, applications)
}

type SubmitExistingApplicationRequest struct {
	ApplicationID uint `json:"application_id" binding:"required"`
}

func SubmitApplication(c *gin.Context) {
	var req SubmitExistingApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request",
			Error:   err.Error(),
		})
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDVal.(uint)

	var application models.Application
	err := db.DB.
		Preload("User").
		Preload("Scheme").
		Preload("StudentProfile").
		Preload("StudentProfile.Documents").
		Preload("StudentProfile.EducationHistory").
		Preload("StudentProfile.Addresses").
		Where("id = ? AND user_id = ?", req.ApplicationID, userID).
		First(&application).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	if !application.IsDraft {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusFound,
			Message: "Application is already submitted",
			Error:   "Application is already submitted",
		})
		return
	}

	// ✅ Check completeness before submission
	if err := utils.CheckApplicationCompleteness(&application); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Application is incomplete",
			Error:   err.Error(),
		})
		return
	}

	now := time.Now()
	application.IsDraft = false
	application.SubmittedAt = &now

	if err := db.DB.
		Save(&application).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to submit application",
			Error:   err.Error(),
		})
		return
	}

	res := models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Application submitted successfully",
		Data:    &application,
	}
	// Return the success response
	// c.JSON(http.StatusOK, res)
	// Return the success response

	c.JSON(http.StatusOK, res)
}

func WithdrawApplication(c *gin.Context) {
	var req SubmitExistingApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request",
			Error:   err.Error(),
		})
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Code: http.StatusUnauthorized, Message: "Unauthorized"})
		return
	}
	userID := userIDVal.(uint)

	var application models.Application
	err := db.DB.
		Where("id = ? AND user_id = ?", req.ApplicationID, userID).
		First(&application).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	if application.IsDraft {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Draft applications cannot be withdrawn"})
		return
	}

	if application.SubmittedAt == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Application not submitted"})
		return
	}

	application.IsDraft = true
	application.SubmittedAt = nil

	if err := db.DB.Save(&application).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to withdraw application"})
		return
	}

	res := models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Application withdrawn successfully",
	}

	c.JSON(http.StatusOK, res)
}
