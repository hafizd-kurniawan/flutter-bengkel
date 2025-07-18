package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// GenerateRandomString generates a random string of specified length
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// ValidatePhoneNumber validates Indonesian phone number format
func ValidatePhoneNumber(phone string) bool {
	// Remove all non-digit characters
	phone = regexp.MustCompile(`\D`).ReplaceAllString(phone, "")
	
	// Indonesian phone number patterns
	patterns := []string{
		`^08\d{8,11}$`,     // Mobile: 08xxxxxxxxxx
		`^62\d{9,12}$`,     // International: 62xxxxxxxxx
		`^021\d{7,8}$`,     // Jakarta landline
		`^0\d{2,3}\d{6,8}$`, // Other landlines
	}
	
	for _, pattern := range patterns {
		if matched, _ := regexp.MatchString(pattern, phone); matched {
			return true
		}
	}
	
	return false
}

// FormatCurrency formats a float64 value to Indonesian Rupiah format
func FormatCurrency(amount float64) string {
	// Convert to string with 2 decimal places
	str := fmt.Sprintf("%.2f", amount)
	
	// Split integer and decimal parts
	parts := strings.Split(str, ".")
	intPart := parts[0]
	decPart := parts[1]
	
	// Add thousand separators to integer part
	var result []string
	for i, digit := range intPart {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			result = append(result, ".")
		}
		result = append(result, string(digit))
	}
	
	// Return formatted currency
	if decPart == "00" {
		return "Rp " + strings.Join(result, "")
	}
	return "Rp " + strings.Join(result, "") + "," + decPart
}

// ParseCurrency parses Indonesian Rupiah format to float64
func ParseCurrency(currency string) (float64, error) {
	// Remove "Rp" prefix and whitespace
	currency = strings.TrimSpace(currency)
	currency = strings.TrimPrefix(currency, "Rp")
	currency = strings.TrimSpace(currency)
	
	// Replace Indonesian separators with standard format
	currency = strings.ReplaceAll(currency, ".", "")
	currency = strings.ReplaceAll(currency, ",", ".")
	
	return strconv.ParseFloat(currency, 64)
}

// ValidateEmail validates email format
func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// SanitizeString removes potentially dangerous characters from string
func SanitizeString(input string) string {
	// Remove HTML tags
	input = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(input, "")
	
	// Remove script tags and content
	input = regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`).ReplaceAllString(input, "")
	
	// Trim whitespace
	input = strings.TrimSpace(input)
	
	return input
}

// GenerateSlug generates a URL-friendly slug from string
func GenerateSlug(text string) string {
	// Convert to lowercase
	text = strings.ToLower(text)
	
	// Replace spaces and special characters with hyphens
	text = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(text, "-")
	
	// Remove leading/trailing hyphens
	text = strings.Trim(text, "-")
	
	return text
}

// Pagination helper
type PaginationParams struct {
	Page  int
	Limit int
}

// GetPaginationParams extracts pagination parameters with defaults
func GetPaginationParams(page, limit int) PaginationParams {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	
	return PaginationParams{
		Page:  page,
		Limit: limit,
	}
}

// CalculateOffset calculates database offset from page and limit
func CalculateOffset(page, limit int) int {
	return (page - 1) * limit
}

// CalculateTotalPages calculates total pages from total records and limit
func CalculateTotalPages(total int64, limit int) int {
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}
	return totalPages
}