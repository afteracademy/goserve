package network

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func CustomTagNameFunc() validator.TagNameFunc {
	return func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if len(name) > 1 {
			return name
		}
		if name == "" {
			name = strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]
		}
		if len(name) > 1 {
			return name
		}
		if name == "" {
			name = strings.SplitN(fld.Tag.Get("uri"), ",", 2)[0]
		}
		if name == "-" {
			return ""
		}
		return name
	}
}
