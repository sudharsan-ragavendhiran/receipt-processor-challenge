package models

import (
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

// Validator instance
var Validate = validator.New()

// Receipt struct with validation rules
type Receipt struct {
	Retailer     string `json:"retailer" binding:"required" validate:"retailerformat"`
	PurchaseDate string `json:"purchaseDate" binding:"required" validate:"dateformat"`
	PurchaseTime string `json:"purchaseTime" binding:"required" validate:"timeformat"`
	Items        []Item `json:"items" binding:"required,min=1"`
	Total        string `json:"total" binding:"required" validate:"priceformat"`
}

// Item struct with validation rules
type Item struct {
	ShortDescription string `json:"shortDescription" binding:"required" validate:"descriptionformat"`
	Price            string `json:"price" binding:"required" validate:"priceformat"`
}

// Custom validation for `YYYY-MM-DD` date format
func ValidateDateFormat(fl validator.FieldLevel) bool {
	_, err := time.Parse("2006-01-02", fl.Field().String()) // Strict parsing
	return err == nil
}

// Custom validation for `HH:MM` (24-hour format)
func ValidateTimeFormat(fl validator.FieldLevel) bool {
	_, err := time.Parse("15:04", fl.Field().String()) // 24-hour format
	return err == nil
}

// Custom validation for `XX.XX` format (price and total)
func ValidatePriceFormat(fl validator.FieldLevel) bool {
	matched, _ := regexp.MatchString(`^\d+\.\d{2}$`, fl.Field().String()) // Enforces exactly two decimal places
	return matched
}

// Custom validation for `retailer` (Only letters, numbers, spaces, &, and -)
func ValidateRetailerFormat(fl validator.FieldLevel) bool {
	matched, _ := regexp.MatchString(`^[\w\s\-&]+$`, fl.Field().String())
	return matched
}

// Custom validation for `shortDescription` (Only letters, numbers, spaces, and -)
func ValidateDescriptionFormat(fl validator.FieldLevel) bool {
	matched, _ := regexp.MatchString(`^[\w\s\-]+$`, fl.Field().String())
	return matched
}

// Register custom validation rules
func init() {
	Validate.RegisterValidation("dateformat", ValidateDateFormat)
	Validate.RegisterValidation("timeformat", ValidateTimeFormat)
	Validate.RegisterValidation("priceformat", ValidatePriceFormat)
	Validate.RegisterValidation("retailerformat", ValidateRetailerFormat)
	Validate.RegisterValidation("descriptionformat", ValidateDescriptionFormat)
}
