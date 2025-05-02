package utils

import (
	"errors"
	"fmt"

	"github.com/ChayanDass/beneficiary-manager/pkg/models"
)

func CheckApplicationCompleteness(app *models.Application) error {
	profile := app.StudentProfile

	// Basic Profile Checks
	if profile.FullName == "" {
		return errors.New("full name is missing")
	}
	if profile.Email == "" {
		return errors.New("email is missing")
	}
	if profile.AadhaarNumber == "" {
		return errors.New("aadhaar number is missing")
	}
	if profile.PhoneNumber == "" {
		return errors.New("phone number is missing")
	}
	if profile.DateOfBirth.IsZero() {
		return errors.New("date of birth is missing")
	}
	if profile.Qualification == "" {
		return errors.New("qualification is missing")
	}
	if profile.Nationality == "" {
		return errors.New("nationality is missing")
	}
	if profile.Category == "" {
		return errors.New("category is missing")
	}
	if profile.Income <= 0 {
		return errors.New("income must be greater than 0")
	}

	// Documents Check
	if len(profile.Documents) == 0 {
		return errors.New("no documents uploaded")
	}
	for i, doc := range profile.Documents {
		if doc.Name == "" || doc.URL == "" {
			return fmt.Errorf("document %d is incomplete", i+1)
		}
	}

	// Education Check
	if len(profile.EducationHistory) == 0 {
		return errors.New("no education history found")
	}
	for i, edu := range profile.EducationHistory {
		if edu.Degree == "" || edu.University == "" || edu.YearOfPassing == 0 {
			return fmt.Errorf("education history %d is incomplete", i+1)
		}
	}

	// Address Check
	if len(profile.Addresses) == 0 {
		return errors.New("no addresses provided")
	}
	hasPermanent := false
	for _, addr := range profile.Addresses {
		if addr.Type == "permanent" {
			hasPermanent = true
		}
		if addr.Street == "" || addr.City == "" || addr.State == "" || addr.Pincode == "" || addr.Country == "" {
			return errors.New("an address entry is incomplete")
		}
	}
	if !hasPermanent {
		return errors.New("permanent address is required")
	}

	return nil
}
