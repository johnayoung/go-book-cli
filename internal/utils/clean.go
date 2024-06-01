package utils

import (
	"regexp"
	"strings"
)

// CleanName generates a clean, repeatable folder name from the given input string.
func CleanName(input string) string {
	// Convert to lowercase
	name := strings.ToLower(input)

	// Replace spaces with hyphens
	name = strings.ReplaceAll(name, " ", "-")

	// Remove any special characters
	re := regexp.MustCompile("[^a-z0-9-]+")
	name = re.ReplaceAllString(name, "")

	return name
}
