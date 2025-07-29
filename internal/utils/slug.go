package utils

import (
	"regexp"
	"strings"
)

// GenerateSlug generates a URL-friendly slug from a string
func GenerateSlug(text string) string {
	// Convert to lowercase
	slug := strings.ToLower(text)
	
	// Replace spaces and special characters with hyphens
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	slug = reg.ReplaceAllString(slug, "-")
	
	// Remove leading and trailing hyphens
	slug = strings.Trim(slug, "-")
	
	// Limit length to 100 characters
	if len(slug) > 100 {
		slug = slug[:100]
		// Remove trailing hyphen if present
		slug = strings.TrimSuffix(slug, "-")
	}
	
	return slug
}