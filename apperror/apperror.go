package apperror

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	validator "github.com/go-playground/validator/v10"
)

var (
	customErrors = map[string]error{
		// user related errors
		"userID.required":           errors.New("is required"),
		"userID.uuid":               errors.New("should be uuid"),
		"email.required":            errors.New("is required"),
		"email.email":               errors.New("should be email"),
		"profileImageUrl.required":  errors.New("is required"),
		"role.required":             errors.New("is required"),
		"role.oneof":                errors.New("should be one of librarian or patrons"),
		"password.required":         errors.New("is required"),
		"password.validatepassword": errors.New("password should contain uppercase and lowercase letter and no special character"),
		// book related errors
		"ISBN.required":          errors.New("is required"),
		"title.required":         errors.New("is required"),
		"author.required":        errors.New("is required"),
		"publishedDate.required": errors.New("is required"),
		"coverImage.required":    errors.New("is required"),
		"shelfNumber.required":   errors.New("is required"),
		"inLibrary.required":     errors.New("is required"),
		"views.required":         errors.New("is required"),
		"booksLeft.required":     errors.New("is required"),
		"wishlistCount.required": errors.New("is required"),
		"rating.required":        errors.New("is required"),
		"reviewCount.required":   errors.New("is required"),
		// checkout related errors
		"ID.required":     errors.New("is required"),
		"bookID.required": errors.New("is required"),
		// review related errors
	}
)

// CustomValidationError converts validation and json marshal error into custom error type
func CustomValidationError(err error) []map[string]string {
	errs := make([]map[string]string, 0)
	switch errTypes := err.(type) {
	case validator.ValidationErrors:
		for _, e := range errTypes {
			errorMap := make(map[string]string)

			key := e.Field() + "." + e.Tag()

			if v, ok := customErrors[key]; ok {
				errorMap[e.Field()] = v.Error()
			} else {
				errorMap[e.Field()] = fmt.Sprintf("custom message is not available: %v", err)
			}
			errs = append(errs, errorMap)
		}
		return errs
	case *json.UnmarshalTypeError:
		errs = append(errs, map[string]string{errTypes.Field: fmt.Sprintf("%v cannot be a %v", errTypes.Field, errTypes.Value)})
		return errs
	}
	errs = append(errs, map[string]string{"unknown": fmt.Sprintf("unsupported custom error for: %v", err)})
	return errs
}

func RegisterTags(v *validator.Validate) {
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		tags := []string{"json", "uri", "form"}
		for _, key := range tags {
			tag := fld.Tag.Get(key)
			name := strings.SplitN(tag, ",", 2)[0]
			if name == "-" {
				return ""
			} else if len(name) != 0 {
				return name
			}
		}
		return ""
	})
}
