package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func Validate[T any](data T) map[string]string {
	err := validator.New().Struct(data)
	res := map[string]string{}

	if err != nil {
		for _, v := range err.(validator.ValidationErrors) {

			res[v.StructField()] = TranslateTag(v)
		}

	}

	return res

}

func TranslateTag(fd validator.FieldError) string {
	switch fd.ActualTag() {
	case "required":
		return fmt.Sprintf("Field %s is required", fd.StructField())
	case "email":
		return fmt.Sprintf("Field %s must be a valid email address", fd.StructField())
	case "min":
		return fmt.Sprintf("Field %s must be at least %s characters long", fd.StructField(), fd.Param())
	case "max":
		return fmt.Sprintf("Field %s must be at most %s characters long", fd.StructField(), fd.Param())
	case "len":
		return fmt.Sprintf("Field %s must be exactly %s characters long", fd.StructField(), fd.Param())
	case "numeric":
		return fmt.Sprintf("Field %s must be a number", fd.StructField())
	case "alphanum":
		return fmt.Sprintf("Field %s must contain only letters and numbers", fd.StructField())
	case "alpha":
		return fmt.Sprintf("Field %s must contain only letters", fd.StructField())
	case "url":
		return fmt.Sprintf("Field %s must be a valid URL", fd.StructField())
	case "uuid":
		return fmt.Sprintf("Field %s must be a valid UUID", fd.StructField())
	case "oneof":
		return fmt.Sprintf("Field %s must be one of the allowed values: %s", fd.StructField(), fd.Param())
	case "contains":
		return fmt.Sprintf("Field %s must contain '%s'", fd.StructField(), fd.Param())
	case "excludes":
		return fmt.Sprintf("Field %s must not contain '%s'", fd.StructField(), fd.Param())
	case "gt":
		return fmt.Sprintf("Field %s must be greater than %s", fd.StructField(), fd.Param())
	case "gte":
		return fmt.Sprintf("Field %s must be greater than or equal to %s", fd.StructField(), fd.Param())
	case "lt":
		return fmt.Sprintf("Field %s must be less than %s", fd.StructField(), fd.Param())
	case "lte":
		return fmt.Sprintf("Field %s must be less than or equal to %s", fd.StructField(), fd.Param())
	case "ipv4":
		return fmt.Sprintf("Field %s must be a valid IPv4 address", fd.StructField())
	case "ipv6":
		return fmt.Sprintf("Field %s must be a valid IPv6 address", fd.StructField())
	case "boolean":
		return fmt.Sprintf("Field %s must be a boolean value (true/false)", fd.StructField())
	case "datetime":
		return fmt.Sprintf("Field %s must be a valid date-time format", fd.StructField())
	case "ascii":
		return fmt.Sprintf("Field %s must contain only ASCII characters", fd.StructField())
	case "lowercase":
		return fmt.Sprintf("Field %s must be in lowercase", fd.StructField())
	case "uppercase":
		return fmt.Sprintf("Field %s must be in uppercase", fd.StructField())
	case "json":
		return fmt.Sprintf("Field %s must be a valid JSON string", fd.StructField())
	case "base64":
		return fmt.Sprintf("Field %s must be a valid base64-encoded string", fd.StructField())
	case "credit_card":
		return fmt.Sprintf("Field %s must be a valid credit card number", fd.StructField())
	case "uuid4":
		return fmt.Sprintf("Field %s must be a valid UUIDv4", fd.StructField())
	case "e164":
		return fmt.Sprintf("Field %s must be a valid E.164 phone number", fd.StructField())
	case "phone":
		return fmt.Sprintf("Field %s must be a valid phone number", fd.StructField())
	case "dir":
		return fmt.Sprintf("Field %s must be a valid directory path", fd.StructField())
	case "file":
		return fmt.Sprintf("Field %s must be a valid file path", fd.StructField())
	case "iscolor":
		return fmt.Sprintf("Field %s must be a valid color code", fd.StructField())
	default:
		return "Validation failed"
	}
}
