package utils

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func StartApplication(db *gorm.DB, userID uint, schemeID uint) (*Application, error) {
	var app Application
	err := db.Where("user_id = ? AND scheme_id = ? AND is_draft = true", userID, schemeID).First(&app).Error

	if err == nil {
		// Draft already exists, return it
		return &app, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Check or create StudentProfile
	var profile StudentProfile
	if err := db.Where("user_id = ?", userID).First(&profile).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		// Create a new empty StudentProfile
		profile = StudentProfile{
			UserID:          userID,
			FullName:        "", // can be filled later
			Email:           "",
			PhoneNumber:     "",
			Gender:          "",
			Qualification:   "",
			Nationality:     "",
			Income:          0,
			IsInternational: false,
		}
		if err := db.Create(&profile).Error; err != nil {
			return nil, fmt.Errorf("failed to create student profile: %v", err)
		}
	} else if err != nil {
		return nil, err
	}

	// Create a new Application in draft state
	app = Application{
		UserID:           userID,
		SchemeID:         schemeID,
		StudentProfileID: profile.ID,
		IsDraft:          true,
	}
	if err := db.Create(&app).Error; err != nil {
		return nil, err
	}

	return &app, nil
}
