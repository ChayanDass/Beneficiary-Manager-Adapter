package models

import (
	"time"
)

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"unique;not null" json:"username"`
}
type UploadDocument struct {
	ID             uint           `gorm:"primaryKey" json:"-"`
	StudentID      uint           `gorm:"not null" json:"-"` // Foreign key to StudentProfile
	Name           string         `json:"name"`
	URL            string         `json:"url"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	StudentProfile StudentProfile `gorm:"foreignKey:StudentID" json:"-"`
}

type Address struct {
	ID        uint      `gorm:"primaryKey" json:"-"`
	StudentID uint      `gorm:"not null" json:"-"`
	Type      string    `gorm:"type:varchar(20)" json:"type"` // permanent, current
	Street    string    `json:"street"`
	City      string    `json:"city"`
	State     string    `json:"state"`
	Pincode   string    `gorm:"type:varchar(10)" json:"pincode"`
	Country   string    `json:"country"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type StudentAcademicQualification struct {
	ID            uint      `gorm:"primaryKey" json:"-"`
	StudentID     uint      `gorm:"not null" json:"-"`
	Degree        string    `json:"degree"`
	University    string    `json:"university"`
	YearOfPassing int       `json:"year_of_passing"`
	Grade         string    `json:"grade"`
	Course        string    `json:"course"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type StudentProfile struct {
	ID               uint                           `gorm:"primaryKey" json:"-"`
	UserID           uint                           `gorm:"not null;unique" json:"-"`
	FullName         string                         `gorm:"not null" json:"full_name"`
	DateOfBirth      time.Time                      `json:"date_of_birth"`
	Gender           string                         `gorm:"type:varchar(10)" json:"gender"`
	PhoneNumber      string                         `gorm:"type:varchar(15)" json:"phone_number"`
	Qualification    string                         `gorm:"type:varchar(50)" json:"qualification"`
	Email            string                         `gorm:"type:varchar(100)" json:"email"`
	AadhaarNumber    string                         `gorm:"type:varchar(12)" json:"aadhaar_number"`
	Nationality      string                         `json:"nationality"` // Added Nationality
	Category         string                         `gorm:"type:varchar(20)" json:"category"`
	Income           float64                        `json:"income"`
	PassportNumber   string                         `json:"passport_number"`  // Passport number for international students
	IsInternational  bool                           `json:"is_international"` // Flag to mark international students
	CreatedAt        time.Time                      `json:"created_at"`
	UpdatedAt        time.Time                      `json:"updated_at"`
	Documents        []UploadDocument               `gorm:"foreignKey:StudentID" json:"documents"`
	EducationHistory []StudentAcademicQualification `gorm:"foreignKey:StudentID" json:"education_history"` // Academic qualifications
	Addresses        []Address                      `gorm:"foreignKey:StudentID" json:"addresses"`         // List of addresses
}

// Application represents a scholarship application
type Application struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	UserID           uint           `gorm:"not null" json:"user_id"`
	SchemeID         uint           `gorm:"not null" json:"scheme_id"`
	StudentProfileID uint           `gorm:"not null" json:"student_profile_id"`
	IsDraft          bool           `gorm:"default:true" json:"is_draft"`
	Verified         bool           `gorm:"default:false" json:"verified"`
	SubmittedAt      *time.Time     `json:"submitted_at,omitempty"`
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	User             User           `gorm:"foreignKey:UserID" json:"user"`
	Scheme           Scheme         `gorm:"foreignKey:SchemeID" json:"-"`
	StudentProfile   StudentProfile `gorm:"foreignKey:StudentProfileID" json:"student_profile"`
}
