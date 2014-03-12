package udf

import (
	"regexp"
	"strings"
)

func GenerateSlug(title string) string {
	slug := strings.ToLower(title)
	invchars := regexp.MustCompile(`[^a-zA-Z0-9\s-]`)
	underscores := regexp.MustCompile(`\s+`)

	slug = invchars.ReplaceAllString(slug, "")
	slug = underscores.ReplaceAllString(slug, "_")

	return slug
}
