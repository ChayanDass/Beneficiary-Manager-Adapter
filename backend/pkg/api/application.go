package api

import (
	"net/http"
	"time"

	"github.com/ChayanDass/beneficiary-manager/pkg/db"
	"github.com/ChayanDass/beneficiary-manager/pkg/models"
	"github.com/ChayanDass/beneficiary-manager/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

func SubmitApplication(c *gin.Context) {
	var req models.SubmitExistingApplicationRequest
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
	var req models.SubmitExistingApplicationRequest
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

// InitApplication initializes the application form for a user
// @Summary Initialize application for a user
// @Description Initialize application with all required data like student profile, scheme, etc.
// @Tags application
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param scheme_id path int true "Scheme ID"
// @Security BasicAuth
// @Success 200 {object} models.Application
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/v1/applications/init-application/{user_id}/{scheme_id} [post]
func InitApplication(c *gin.Context) {
	var req models.InitApplicationRequest
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

	// Optional: Check if there's already a draft application for this user & scheme
	var existing models.Application
	if err := db.DB.
		Where("user_id = ? AND scheme_id = ? ", userID, req.SchemeID).
		First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": " application already exists"})
		return
	}
	var schema models.Scheme
	// Check if the scheme exists
	if err := db.DB.First(&schema, req.SchemeID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Scheme not found",
				Error:   err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to fetch scheme",
				Error:   err.Error(),
			})
		}
		return
	}

	// Continue with the logic after checking for the scheme

	student := models.StudentProfile{
		UserID:   userID,
		FullName: "Unknown", // Use default or empty values
		// ... other defaults
	}
	if err := db.DB.Create(&student).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create student profile"})
		return
	}

	application := models.Application{
		UserID:           userID,
		SchemeID:         req.SchemeID,
		IsDraft:          true,
		StudentProfileID: student.ID,
	}

	if err := db.DB.Create(&application).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to initialize application",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Application initialized successfully",
		"application": application,
	})
}
