package util

import (
	"strings"
)

// IsApplicationJSONMediaType will return true if mediaType is in the format of
// application/*json (application/json, application/hal+json...)
func IsApplicationJSONMediaType(mediaType string) bool {
	return strings.HasPrefix(mediaType, "application/") &&
		strings.HasSuffix(mediaType, "json")
}
