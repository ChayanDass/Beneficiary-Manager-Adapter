package api

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/ChayanDass/beneficiary-manager/pkg/db"
	"github.com/ChayanDass/beneficiary-manager/pkg/models"
	"github.com/ChayanDass/beneficiary-manager/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// func GetSchemes(c *gin.Context) {
// 	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
// 	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
// 	pagination := models.PaginationInput{
// 		Page:  page,
// 		Limit: limit,
// 	}
// 	var schemes []models.Scheme
// 	offset := pagination.GetOffset()
// 	var filter models.SchemeFilter

// 	// Bind query parameters into the filter struct
// 	if err := c.ShouldBindQuery(&filter); err != nil {
// 		c.JSON(http.StatusBadRequest, models.ErrorResponse{
// 			Code:    http.StatusBadRequest,
// 			Message: "Invalid query parameters",
// 			Error:   err.Error(),
// 		})
// 		return
// 	}

// 	// Start building the query with the joins directly
// 	query := db.DB.
// 		Preload("Eligibility").
// 		Preload("Eligibility.DocumentMappings").
// 		Preload("Eligibility.DocumentMappings.Document").
// 		Model(&models.Scheme{}).
// 		Joins("JOIN eligibilities ON eligibilities.id = schemes.eligibility_id")
// 	if err := query.Find(&schemes).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
// 			Code:    http.StatusInternalServerError,
// 			Message: "Failed to fetch schemes",
// 			Error:   err.Error(),
// 		})
// 		return
// 	}

// 	// Print the result (schemes) to see the raw data before filtering
// 	fmt.Println("Schemes with Eligibility before applying filters:")
// 	for _, scheme := range schemes {
// 		fmt.Println(scheme.Eligibility.Gender)
// 	}

// 	fmt.Println(filter)

// 	// Apply Scheme Filters
// 	if filter.Name != nil {
// 		query = query.Where("schemes.name ILIKE ?", "%"+*filter.Name+"%")
// 	}
// 	if filter.Status != nil {
// 		query = query.Where("schemes.status = ?", *filter.Status)
// 	}
// 	if filter.MinAmount != nil {
// 		query = query.Where("schemes.amount >= ?", *filter.MinAmount)
// 	}
// 	if filter.MaxAmount != nil {
// 		query = query.Where("schemes.amount <= ?", *filter.MaxAmount)
// 	}
// 	if filter.StartAfter != nil {
// 		query = query.Where("schemes.start_date >= ?", *filter.StartAfter)
// 	}
// 	if filter.EndBefore != nil {
// 		query = query.Where("schemes.end_date <= ?", *filter.EndBefore)
// 	}

// 	// Apply Eligibility Filters
// 	if filter.Gender != nil {
// 		query = query.Where("eligibilities.gender = ?", *filter.Gender)
// 	}
// 	if filter.AcademicQualification != nil {
// 		query = query.Where("eligibilities.academic_qualification = ?", *filter.AcademicQualification)
// 	}
// 	if filter.IncomeLimit != nil {
// 		query = query.Where("eligibilities.income_limit >= ?", *filter.IncomeLimit)
// 	}
// 	if filter.Category != nil {
// 		query = query.Where("eligibilities.category = ?", *filter.Category)
// 	}

// 	// Count total schemes first (after the filters)
// 	var totalCount int64
// 	if err := query.Count(&totalCount).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
// 			Code:    http.StatusInternalServerError,
// 			Message: "Failed to fetch total count",
// 			Error:   err.Error(),
// 		})
// 		return
// 	}

// 	// Calculate total pages
// 	totalPages := int64(math.Ceil(float64(totalCount) / float64(pagination.GetLimit())))

// 	// Pagination
// 	// Build pagination links
// 	params := c.Request.URL.Query()
// 	basePath := c.Request.URL.Path
// 	var previous, next string
// 	if pagination.Page > 1 {
// 		params.Set("page", strconv.FormatInt(pagination.Page-1, 10))
// 		previous = basePath + "?" + params.Encode()
// 	}
// 	if pagination.Page < totalPages {
// 		params.Set("page", strconv.FormatInt(pagination.Page+1, 10))
// 		next = basePath + "?" + params.Encode()
// 	}

// 	meta := &models.PaginationMeta{
// 		ResourceCount: int(totalCount),
// 		TotalPages:    totalPages,
// 		Page:          pagination.Page,
// 		Limit:         pagination.Limit,
// 		Previous:      previous,
// 		Next:          next,
// 	}

// 	// Fetch the filtered schemes
// 	if err := query.Offset(int(offset)).Limit(int(pagination.GetLimit())).Find(&schemes).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
// 			Code:    http.StatusInternalServerError,
// 			Message: "Failed to fetch schemes",
// 			Error:   err.Error(),
// 		})
// 		return
// 	}

// 	// Respond with the filtered schemes and pagination info
// 	res := models.SchemeResponse{
// 		Data:    schemes,
// 		Code:    http.StatusOK,
// 		Message: "Schemes fetched successfully",
// 		Meta:    meta,
// 	}

// 	c.JSON(http.StatusOK, res)
// }

func GetSchemes(c *gin.Context) {
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	pagination := models.PaginationInput{
		Page:  page,
		Limit: limit,
	}
	offset := pagination.GetOffset()

	var filter models.SchemeFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid query parameters",
			Error:   err.Error(),
		})
		return
	}

	var schemes []models.Scheme
	query := db.DB.
		Preload("Eligibility").
		Preload("Eligibility.DocumentMappings").
		Preload("Eligibility.DocumentMappings.Document").
		Model(&models.Scheme{}).
		Joins("JOIN eligibilities ON eligibilities.id = schemes.eligibility_id")

	// Apply filters using utility
	query = utils.ApplySchemeFilters(query, filter)

	var totalCount int64
	if err := query.Count(&totalCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to fetch total count",
			Error:   err.Error(),
		})
		return
	}

	totalPages := int64(math.Ceil(float64(totalCount) / float64(pagination.Limit)))

	params := c.Request.URL.Query()
	basePath := c.Request.URL.Path
	var previous, next string
	if pagination.Page > 1 {
		params.Set("page", strconv.FormatInt(pagination.Page-1, 10))
		previous = basePath + "?" + params.Encode()
	}
	if pagination.Page < totalPages {
		params.Set("page", strconv.FormatInt(pagination.Page+1, 10))
		next = basePath + "?" + params.Encode()
	}

	meta := &models.PaginationMeta{
		ResourceCount: int(totalCount),
		TotalPages:    totalPages,
		Page:          pagination.Page,
		Limit:         pagination.Limit,
		Previous:      previous,
		Next:          next,
	}

	if err := query.Offset(int(offset)).Limit(int(pagination.Limit)).Find(&schemes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to fetch schemes",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SchemeResponse{
		Data:    schemes,
		Code:    http.StatusOK,
		Message: "Schemes fetched successfully",
		Meta:    meta,
	})
}

func GetSchemeByID(c *gin.Context) {
	id := c.Param("id")
	var scheme models.Scheme

	if err := db.DB.
		Preload("Eligibility").
		Preload("Eligibility.DocumentMappings").
		Preload("Eligibility.DocumentMappings.Document").
		First(&scheme, id).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Scheme not found",
				Error:   err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to fetch scheme",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Scheme retrieved successfully",
		Data:    scheme,
	})
}

// @Summary Get Scheme Status
// @Description Get the current status of a scheme by its ID
// @Tags scheme
// @Accept json
// @Produce json
// @Param id path string true "Scheme ID"
// @Success 200 {object} models.Scheme
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/schemes/status/{id} [get]
func GetSchemeStatus(c *gin.Context) {
	// Extract parameters from URL
	id := c.Param("id")
	var scheme models.Scheme

	err := db.DB.Where("id = ?", id).First(&scheme).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Scheme not found",
				Error:   err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to retrieve scheme",
			Error:   err.Error(),
		})
		return
	}

	res := models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Application status fetched successfully",
		Data:    fmt.Sprintf("application status is %s", scheme.Status),
	}
	c.JSON(http.StatusOK, res)
}
