package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ChayanDass/beneficiary-manager/pkg/db"
	"github.com/ChayanDass/beneficiary-manager/pkg/models"
	"github.com/ChayanDass/beneficiary-manager/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetApplications retrieves all applications for the authenticated user.
// @Summary Get user applications
// @Description Fetches all applications associated with the authenticated user.
// @Tags Applications
// @Accept json
// @Produce json
// @Success 200 {array} models.Application "List of applications"
// @Failure 401 {object} models.ErrorResponse "Unauthorized, user ID not found in context"
// @Failure 404 {object} models.ErrorResponse "No applications found for this user"
// @Failure 500 {object} models.ErrorResponse "Failed to fetch applications"
// @Router /applications [get]
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

// SubmitApplication submits an existing draft application for the authenticated user.
// @Summary Submit application
// @Description Submits a draft application after validating its completeness.
// @Tags Applications
// @Accept json
// @Produce json
// @Param request body models.SubmitExistingApplicationRequest true "Submit application request"
// @Success 200 {object} models.SuccessResponse "Application submitted successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request or application is incomplete"
// @Failure 401 {object} models.ErrorResponse "Unauthorized, user ID not found in context"
// @Failure 404 {object} models.ErrorResponse "Application not found"
// @Failure 500 {object} models.ErrorResponse "Failed to submit application"
// @Router /applications/submit [post]
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
	now := time.Now()
	application.IsDraft = false
	application.Status = "submitted"
	application.SubmittedAt = &now

	// ✅ Check completeness before submission
	if err := utils.CheckApplicationCompleteness(&application); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Application is incomplete",
			Error:   err.Error(),
		})
		return
	}

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

// WithdrawApplication withdraws a submitted application for the authenticated user.
//
// @Summary Withdraw application
// @Description Withdraws a submitted application and marks it as a draft.
// @Tags Applications
// @Accept json
// @Produce json
// @Param request body models.SubmitExistingApplicationRequest true "Withdraw application request"
// @Success 200 {object} models.SuccessResponse "Application withdrawn successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request or application cannot be withdrawn"
// @Failure 401 {object} models.ErrorResponse "Unauthorized, user ID not found in context"
// @Failure 404 {object} models.ErrorResponse "Application not found"
// @Failure 500 {object} models.ErrorResponse "Failed to withdraw application"
// @Router /applications/withdraw [post]
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
	application.Status = "draft"
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

// InitApplication initializes a new draft application for the authenticated user.
//
// @Summary Initialize application
// @Description Creates a new draft application for a specific scheme.
// @Tags Applications
// @Accept json
// @Produce json
// @Param request body models.InitApplicationRequest true "Initialize application request"
// @Success 201 {object} models.SuccessResponse "Application initialized successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized, user ID not found in context"
// @Failure 404 {object} models.ErrorResponse "Scheme not found"
// @Failure 409 {object} models.ErrorResponse "Application already exists"
// @Failure 500 {object} models.ErrorResponse "Failed to initialize application"
// @Router /applications/init [post]
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
		Status:           "draft",
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
		"message": "Application initialized successfully",
	})
}

// ModifyApplication modifies an existing draft application for the authenticated user.
//
// @Summary Modify application
// @Description Updates the details of a draft application, including the student profile and related data.
// @Tags Applications
// @Accept json
// @Produce json
// @Param id path string true "Application ID"
// @Param request body models.StudentProfileInput true "Modify application request"
// @Success 200 {object} models.SuccessResponse "Application modified successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid input"
// @Failure 401 {object} models.ErrorResponse "Unauthorized, user ID not found in context"
// @Failure 403 {object} models.ErrorResponse "Cannot modify application, it is already submitted"
// @Failure 404 {object} models.ErrorResponse "Application not found"
// @Failure 500 {object} models.ErrorResponse "Failed to update application"
// @Router /applications/{id}/modify [put]
func ModifyApplication(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "User ID not found in context",
		})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "Invalid user ID type",
		})
		return
	}

	schemeID := c.Param("id")

	var application models.Application
	if err := db.DB.
		Preload("User").
		Preload("Scheme").
		Preload("StudentProfile").
		Preload("StudentProfile.Documents").
		Preload("StudentProfile.EducationHistory").
		Preload("StudentProfile.Addresses").
		Where("user_id = ? AND id = ?", userIDUint, schemeID).
		First(&application).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	studentProfile := application.StudentProfile

	var input models.StudentProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if !application.IsDraft {
		c.JSON(http.StatusForbidden,
			models.ErrorResponse{
				Code:    http.StatusForbidden,
				Message: "Cannot modify application, it is already submitted.",
				Error:   "Application is already submitted",
			})
		return
	}

	// --- Update student profile fields ---
	if input.FullName != "" {
		studentProfile.FullName = input.FullName
	}
	if input.Email != "" {
		studentProfile.Email = input.Email
	}
	if input.PhoneNumber != "" {
		studentProfile.PhoneNumber = input.PhoneNumber
	}
	if input.DateOfBirth != nil {
		studentProfile.DateOfBirth = *input.DateOfBirth
	}
	if input.Qualification != "" {
		studentProfile.Qualification = input.Qualification
	}
	if input.Category != "" {
		studentProfile.Category = input.Category
	}
	if input.Income != nil {
		studentProfile.Income = *input.Income
	}
	if input.Nationality != "" {
		studentProfile.Nationality = input.Nationality
	}
	if input.Gender != "" {
		studentProfile.Gender = input.Gender
	}
	if input.AadhaarNumber != "" {
		studentProfile.AadhaarNumber = input.AadhaarNumber
	}
	if err := db.DB.Omit("Documents", "Addresses", "EducationHistory").Save(&studentProfile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update student profile"})
		return
	}

	if err := utils.UpsertStudentDocuments(db.DB, studentProfile.ID, input.Documents); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := utils.UpsertStudentAddresses(db.DB, studentProfile.ID, input.Addresses); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to upsert address",
			"details": err.Error(),
		})
		return
	}
	if err := utils.UpsertEducationHistory(db.DB, studentProfile.ID, input.EducationHistory); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Application modified successfully",
		Data:    &application,
	}
	// Return the success response

	c.JSON(http.StatusOK, res)
}

// GetApplicationStatus retrieves the status of a specific application for the authenticated user.
//
// @Summary Get application status
// @Description Fetches the status of an application based on the application ID and the authenticated user's ID.
// @Tags Applications
// @Accept json
// @Produce json
// @Param id path string true "Application ID"
// @Success 200 {object} models.SuccessResponse "Application status fetched successfully"
// @Failure 401 {object} models.ErrorResponse "Unauthorized, user ID not found in context"
// @Failure 404 {object} models.ErrorResponse "Application not found"
// @Router /applications/{id}/status [get]
func GetApplicationStatus(c *gin.Context) {
	applicationID := c.Param("id")
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Code: http.StatusUnauthorized, Message: "Unauthorized"})
		return
	}
	userID := userIDVal.(uint)
	var application models.Application
	if err := db.DB.
		Where("id = ? AND user_id=?", applicationID, userID).
		First(&application).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}
	res := models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Application status fetched successfully",
		Data:    fmt.Sprintf("application status is %s", application.Status),
	}
	c.JSON(http.StatusOK, res)
}
