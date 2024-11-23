package apispec

import (
	"time"

	"github.com/xeipuuv/gojsonschema"
)

var formats = []string{
	"date",
	"time",
	"date-time",
	"email",
	"ipv4",
	"ipv6",
	"uuid",
	"json-pointer",
	// "relative-json-pointer", // matched with "1.147.1"
	// "hostname",
	// "regex",
	// "uri",           // can be also iri
	// "uri-reference", // can be also iri-reference
	// "uri-template",
}

func getStringFormat(value any) string {
	str, ok := value.(string)
	if !ok || str == "" {
		return ""
	}

	for _, format := range formats {
		if gojsonschema.FormatCheckers.IsFormat(format, value) {
			return format
		}
	}

	return ""
}

// isDateFormat checks if input is a correctly formatted date with spaces (excluding RFC3339 = "2006-01-02T15:04:05Z07:00")
// This is useful to identify date string instead of an array.
func isDateFormat(input any) bool {
	asString, ok := input.(string)
	if !ok {
		return false
	}
	if _, err := time.Parse(time.ANSIC, asString); err == nil {
		return true
	}
	if _, err := time.Parse(time.UnixDate, asString); err == nil {
		return true
	}
	if _, err := time.Parse(time.RubyDate, asString); err == nil {
		return true
	}
	if _, err := time.Parse(time.RFC822, asString); err == nil {
		return true
	}
	if _, err := time.Parse(time.RFC822Z, asString); err == nil {
		return true
	}
	if _, err := time.Parse(time.RFC850, asString); err == nil {
		return true
	}
	if _, err := time.Parse(time.RFC1123, asString); err == nil {
		return true
	}
	if _, err := time.Parse(time.RFC1123Z, asString); err == nil {
		return true
	}
	if _, err := time.Parse(time.Stamp, asString); err == nil {
		return true
	}
	if _, err := time.Parse(time.StampMilli, asString); err == nil {
		return true
	}
	if _, err := time.Parse(time.StampMicro, asString); err == nil {
		return true
	}
	if _, err := time.Parse(time.StampNano, asString); err == nil {
		return true
	}
	return false
}
