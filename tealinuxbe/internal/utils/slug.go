package utils

import (
	"regexp"
	"strings"
)

func MakeSlug(text string) string {
	// Convert to lower case
	slug := strings.ToLower(text)
	// Remove non-alphanumeric characters
	reg, _ := regexp.Compile("[^a-z0-9]+")
	slug = reg.ReplaceAllString(slug, "-")
	// Remove leading and trailing hyphens
	slug = strings.Trim(slug, "-")
	return slug
}
