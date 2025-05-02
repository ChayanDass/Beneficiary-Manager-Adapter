package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

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
