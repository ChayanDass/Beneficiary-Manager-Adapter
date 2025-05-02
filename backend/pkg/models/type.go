package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ------------------ Enums ------------------

type AcademicQualification string

const (
	AcademicQualificationNone         AcademicQualification = "None"
	AcademicQualificationClassX       AcademicQualification = "Class-X"
	AcademicQualificationClassXII     AcademicQualification = "Class-XII"
	AcademicQualificationDiploma      AcademicQualification = "Diploma"
	AcademicQualificationGraduate     AcademicQualification = "Graduate"
	AcademicQualificationPostGraduate AcademicQualification = "Post-Graduate"
)

type Gender string

const (
	GenderMale   Gender = "Male"
	GenderFemale Gender = "Female"
	GenderOther  Gender = "Other"
)

type Category string

const (
	CategoryGeneral Category = "General"
	CategorySC      Category = "SC"
	CategoryST      Category = "ST"
	CategoryOBC     Category = "OBC"
	CategoryOther   Category = "Other"
)

type Document string

const (
	DocumentAadharCard            Document = "aadhar_card"
	DocumentPanCard               Document = "pan_card"
	DocumentDrivingLic            Document = "driving_license"
	DocumentClassXCert            Document = "class_x_certificate"
	DocumentClassXIICertificate   Document = "class_xii_certificate"
	DocumentDiplomaCertificate    Document = "diploma_certificate"
	DocumentGraduationCertificate Document = "graduation_certificate"
	DocumentPostGradCertificate   Document = "post_grad_certificate"
	DocumentPassport              Document = "passport"
	DocumentOther                 Document = "other"
)

// ------------------ Core Models ------------------

// User represents a user
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Scheme represents a scholarship scheme
type Scheme struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	Name            string         `json:"name" gorm:"not null"`
	Description     string         `json:"description"`
	EligibilityID   uint           `json:"eligibility_id"` // foreign key to Eligibility
	Eligibility     Eligibility    `gorm:"foreignKey:EligibilityID" json:"eligibility"`
	Amount          float64        `json:"amount"`
	ApplicationLink string         `json:"application_link"`
	StartDate       time.Time      `json:"start_date"`
	EndDate         time.Time      `json:"end_date"`
	Status          string         `json:"status" gorm:"type:varchar(20);default:'upcoming'"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

// Eligibility represents eligibility criteria
type Eligibility struct {
	ID                    uint                  `gorm:"primaryKey" json:"id"`
	Gender                Gender                `gorm:"type:varchar(10)" json:"gender"`
	AgeMin                int                   `json:"age_min"`
	AgeMax                int                   `json:"age_max"`
	IncomeLimit           float64               `json:"income_limit"`
	AcademicQualification AcademicQualification `gorm:"type:varchar(20)" json:"academic_qualification"`
	Category              Category              `gorm:"type:varchar(20)" json:"category"`
	DocumentsRequired     datatypes.JSON        `gorm:"type:jsonb" json:"documents_required"`
	CreatedAt             time.Time             `json:"created_at"`
	UpdatedAt             time.Time             `json:"updated_at"`
}

// Application represents a scholarship application
type Application struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	SchemeID  uint      `gorm:"not null" json:"scheme_id"`
	Status    string    `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, approved, rejected
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      User      `gorm:"foreignKey:UserID;references:ID" json:"user"`
	Scheme    Scheme    `gorm:"foreignKey:SchemeID;references:ID" json:"scheme"`
}

// ------------------ Filtering ------------------

// SchemeFilter represents filter criteria

type SchemeFilter struct {
	Name                  *string    `form:"name" example:"Scholar Scheme"`
	Status                *string    `form:"status" example:"upcoming"`
	MinAmount             *float64   `form:"min_amount" example:"1000"`
	MaxAmount             *float64   `form:"max_amount" example:"5000"`
	StartAfter            *time.Time `form:"start_after" example:"2023-01-01T00:00:00Z"`
	EndBefore             *time.Time `form:"end_before" example:"2023-12-31T23:59:59Z"`
	Gender                *string    `form:"gender"  `
	AcademicQualification *string    `form:"academic_qualification"`
	IncomeLimit           *float64   `form:"income_limit"`
	Category              *string    `form:"category"`
}

// ------------------ Pagination ------------------

// PaginationMeta contains metadata for paginated responses
type PaginationMeta struct {
	ResourceCount int    `json:"resource_count" example:"200"`
	TotalPages    int64  `json:"total_pages,omitempty" example:"20"`
	Page          int64  `json:"page,omitempty" example:"10"`
	Limit         int64  `json:"limit,omitempty" example:"10"`
	Next          string `json:"next,omitempty" example:"/api/v1/schemes?limit=10&page=11"`
	Previous      string `json:"previous,omitempty" example:"/api/v1/schemes?limit=10&page=9"`
}

// PaginationInput is the input model for pagination
type PaginationInput struct {
	Page  int64 `json:"page" example:"10"`
	Limit int64 `json:"limit" example:"10"`
}

// PaginationParse defines behavior for pagination inputs
type PaginationParse interface {
	GetOffset() int64
	GetLimit() int64
}

// GetOffset returns offset value
func (p PaginationInput) GetOffset() int64 {
	return (p.Page - 1) * p.Limit
}

// GetLimit returns limit value
func (p PaginationInput) GetLimit() int64 {
	return p.Limit
}

// ------------------ API Responses ------------------

// ErrorResponse for API error output
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// SuccessResponse for success output
type SuccessResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SchemeResponse for paginated schemes
type SchemeResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    interface{}     `json:"data,omitempty"`
	Meta    *PaginationMeta `json:"meta,omitempty"`
}
