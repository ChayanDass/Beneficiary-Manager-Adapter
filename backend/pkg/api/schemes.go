package api

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/ChayanDass/beneficiary-manager/pkg/db"
	"github.com/ChayanDass/beneficiary-manager/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetSchemes retrieves all schemes with optional filters
// GetSchemes retrieves all schemes with optional filters
// GetSchemes retrieves all schemes with optional filters
func GetSchemes(c *gin.Context) {
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	pagination := models.PaginationInput{
		Page:  page,
		Limit: limit,
	}
	var schemes []models.Scheme
	offset := pagination.GetOffset()
	var filter models.SchemeFilter

	// Bind query parameters into the filter struct
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid query parameters",
			Error:   err.Error(),
		})
		return
	}

	// Start building the query with the joins directly
	query := db.DB.Model(&models.Scheme{}).
		Joins("JOIN eligibilities ON eligibilities.id = schemes.eligibility_id").
		Preload("Eligibility") // Preload Eligibility if you want to load it as part of the scheme data
	// Print the query output before filtering
	if err := query.Find(&schemes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to fetch schemes",
			Error:   err.Error(),
		})
		return
	}

	// Print the result (schemes) to see the raw data before filtering
	fmt.Println("Schemes with Eligibility before applying filters:")
	for _, scheme := range schemes {
		fmt.Println(scheme.Eligibility.Gender)
	}

	fmt.Println(filter)

	// Apply Scheme Filters
	if filter.Name != nil {
		query = query.Where("schemes.name ILIKE ?", "%"+*filter.Name+"%")
	}
	if filter.Status != nil {
		query = query.Where("schemes.status = ?", *filter.Status)
	}
	if filter.MinAmount != nil {
		query = query.Where("schemes.amount >= ?", *filter.MinAmount)
	}
	if filter.MaxAmount != nil {
		query = query.Where("schemes.amount <= ?", *filter.MaxAmount)
	}
	if filter.StartAfter != nil {
		query = query.Where("schemes.start_date >= ?", *filter.StartAfter)
	}
	if filter.EndBefore != nil {
		query = query.Where("schemes.end_date <= ?", *filter.EndBefore)
	}

	// Apply Eligibility Filters
	if filter.Gender != nil {
		query = query.Where("eligibilities.gender = ?", *filter.Gender)
	}
	if filter.AcademicQualification != nil {
		query = query.Where("eligibilities.academic_qualification = ?", *filter.AcademicQualification)
	}
	if filter.IncomeLimit != nil {
		query = query.Where("eligibilities.income_limit <= ?", *filter.IncomeLimit)
	}
	if filter.Category != nil {
		query = query.Where("eligibilities.category = ?", *filter.Category)
	}

	// Count total schemes first (after the filters)
	var totalCount int64
	if err := query.Count(&totalCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to fetch total count",
			Error:   err.Error(),
		})
		return
	}

	// Calculate total pages
	totalPages := int64(math.Ceil(float64(totalCount) / float64(pagination.GetLimit())))

	// Pagination
	// Build pagination links
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

	// Fetch the filtered schemes
	if err := query.Offset(int(offset)).Limit(int(pagination.GetLimit())).Find(&schemes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to fetch schemes",
			Error:   err.Error(),
		})
		return
	}

	// Respond with the filtered schemes and pagination info
	res := models.SchemeResponse{
		Data:    schemes,
		Code:    http.StatusOK,
		Message: "Schemes fetched successfully",
		Meta:    meta,
	}

	c.JSON(http.StatusOK, res)
}

// GetSchemeByID retrieves a specific scheme by its ID
func GetSchemeByID(c *gin.Context) {
	id := c.Param("id")
	var scheme models.Scheme

	if err := db.DB.First(&scheme, id).Error; err != nil {
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

	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Scheme retrieved successfully",
		Data:    scheme,
	})
}

// GetSchemeStatus retrieves the status of a specific scheme by its ID

func GetSchemeStatus(c *gin.Context) {
	// Extract parameters from URL
	id := c.Param("id")
	name := c.DefaultQuery("name", "") // Get "name" query parameter, default to empty string if not provided
	var scheme models.Scheme

	// Build the query
	query := db.DB.Model(&models.Scheme{})

	// Apply filter for ID or Name
	if id != "" {
		query = query.Where("id = ?", id)
	} else if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	} else {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Either 'id' or 'name' must be provided",
			Error:   "Missing 'id' or 'name' parameter",
		})
		return
	}

	// Fetch the scheme status based on the provided filters
	if err := query.Select("status").First(&scheme).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Scheme not found",
				Error:   err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to fetch scheme status",
				Error:   err.Error(),
			})
		}
		return
	}

	// Return the scheme status
	c.JSON(http.StatusOK, models.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Scheme status retrieved successfully",
		Data:    scheme.Status,
	})
}
