package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Messages = map[string]string{
	// Core
	"required":  "%s is required",
	"omitempty": "%s is optional",

	// Length
	"min": "%s must be at least %s characters",
	"max": "%s must be at most %s characters",
	"len": "%s must be exactly %s characters",

	// String content
	"alpha":           "%s must contain only letters",
	"alphanum":        "%s must contain only letters and numbers",
	"alphanumunicode": "%s must contain only letters and numbers",
	"lowercase":       "%s must be lowercase",
	"uppercase":       "%s must be uppercase",
	"contains":        "%s must contain %s",
	"containsany":     "%s must contain one of %s",
	"excludes":        "%s must not contain %s",
	"startswith":      "%s must start with %s",
	"endswith":        "%s must end with %s",

	// Format
	"email":    "%s is not a valid email",
	"url":      "%s is not a valid URL",
	"uri":      "%s is not a valid URI",
	"uuid":     "%s must be a valid UUID",
	"uuid4":    "%s must be a valid UUID",
	"uuid5":    "%s must be a valid UUID",
	"ip":       "%s must be a valid IP address",
	"ipv4":     "%s must be a valid IPv4 address",
	"ipv6":     "%s must be a valid IPv6 address",
	"hostname": "%s must be a valid hostname",

	// Numbers
	"gt":  "%s must be greater than %s",
	"gte": "%s must be greater than or equal to %s",
	"lt":  "%s must be less than %s",
	"lte": "%s must be less than or equal to %s",
	"eq":  "%s must be equal to %s",
	"ne":  "%s must not be equal to %s",

	// Time
	"datetime": "%s must be a valid datetime",

	// Collections
	"unique": "%s must contain unique values",
	"dive":   "%s contains an invalid value",

	// Conditional
	"eqfield":              "%s must be equal to %s",
	"nefield":              "%s must not be equal to %s",
	"required_if":          "%s is required",
	"required_with":        "%s is required",
	"required_without":     "%s is required",
	"required_without_all": "%s is required",

	// Security / enums
	"oneof": "%s must be one of %s",
}

func FormatValidationErrors(err error) []string {
	var errs validator.ValidationErrors
	if !errors.As(err, &errs) {
		return []string{"invalid request"}
	}

	msgs := make([]string, 0, len(errs))
	for _, e := range errs {
		format, ok := Messages[e.Tag()]
		if !ok {
			format = "%s is invalid"
		}
		msgs = append(msgs, fmt.Sprintf(format, strings.ToLower(e.Field()), e.Param()))
	}
	return msgs
}
