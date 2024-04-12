package utils

import (
	"reflect"
	"regexp"
	"time"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate
var Translator ut.Translator

// InitReqDataValidationTranslation translate request data validation error message from default to custom for each defined field.
func InitReqDataValidationTranslation() {
	Validate = validator.New()
	en := en.New()
	uni := ut.New(en, en)
	Translator, _ = uni.GetTranslator("en")

	Validate.RegisterValidation("alpha_with_spaces", CustomAlphaWithSpaceValidator)
	Validate.RegisterValidation("alphanum_with_spaces", CustomAlphaNumWithSpaceValidator)
	Validate.RegisterValidation("time", CustomTimeValidator)
	Validate.RegisterValidation("slice_of_numbers", CustomSliceOfNumberValidator)

	Validate.RegisterTranslation("required", Translator, func(ut ut.Translator) error {
		return ut.Add("required", "{0} field is required.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field(), fe.Param())
		return t
	})

	Validate.RegisterTranslation("alpha", Translator, func(ut ut.Translator) error {
		return ut.Add("alpha", "{0} field must contain alphabets only.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("alpha", fe.Field(), fe.Param())
		return t
	})

	Validate.RegisterTranslation("alpha_with_spaces", Translator, func(ut ut.Translator) error {
		return ut.Add("alpha_with_spaces", "{0} field must contain alphabets and space only.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("alpha_with_spaces", fe.Field(), fe.Param())
		return t
	})

	Validate.RegisterTranslation("alphanum_with_spaces", Translator, func(ut ut.Translator) error {
		return ut.Add("alphanum_with_spaces", "{0} field must contain alphabets, numbers and space only.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("alphanum_with_spaces", fe.Field(), fe.Param())
		return t
	})

	Validate.RegisterTranslation("number", Translator, func(ut ut.Translator) error {
		return ut.Add("number", "{0} field must contain numbers only.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("number", fe.Field(), fe.Param())
		return t
	})

	Validate.RegisterTranslation("boolean", Translator, func(ut ut.Translator) error {
		return ut.Add("boolean", "{0} field must contain either true or false.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("boolean", fe.Field(), fe.Param())
		return t
	})

	Validate.RegisterTranslation("slice_of_numbers", Translator, func(ut ut.Translator) error {
		return ut.Add("slice_of_numbers", "{0} field must be slice of numbers only and must contain at least 1 value.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("slice_of_numbers", fe.Field(), fe.Param())
		return t
	})

	Validate.RegisterTranslation("time", Translator, func(ut ut.Translator) error {
		return ut.Add("time", "{0} field must be time only and must be greater than now.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("time", fe.Field(), fe.Param())
		return t
	})

	Validate.RegisterTranslation("min", Translator, func(ut ut.Translator) error {
		return ut.Add("min", "{0} field violates minimum length/value constraint. length/value must be at least {1} long.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("min", fe.Field(), fe.Param())
		return t
	})

	Validate.RegisterTranslation("max", Translator, func(ut ut.Translator) error {
		return ut.Add("max", "{0} field violates maximum length/value constraint. length/value must be at most {1} long.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("max", fe.Field(), fe.Param())
		return t
	})

	Validate.RegisterTranslation("gte", Translator, func(ut ut.Translator) error {
		return ut.Add("max", "{0} field must be greater than {1}.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("max", fe.Field(), fe.Param())
		return t
	})

	Validate.RegisterTranslation("email", Translator, func(ut ut.Translator) error {
		return ut.Add("email", "{0} must be a valid email address", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("email", fe.Field(), fe.Param())
		return t
	})

	Validate.RegisterTranslation("oneof", Translator, func(ut ut.Translator) error {
		return ut.Add("oneof", "{0} must be one of {1}", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("oneof", fe.Field(), fe.Param())
		return t
	})

}

func CustomAlphaWithSpaceValidator(fl validator.FieldLevel) bool {
	str := fl.Field().String()
	pattern := "^[a-zA-Z ]+$"
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(str)
}

func CustomAlphaNumWithSpaceValidator(fl validator.FieldLevel) bool {
	str := fl.Field().String()
	pattern := "^[a-zA-Z0-9 -.]+$"
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(str)
}

func CustomSliceOfNumberValidator(fl validator.FieldLevel) bool {
	slice := fl.Field()
	if slice.Kind() != reflect.Slice {
		return false
	}
	if slice.Len() == 0 {
		return false
	}
	for i := 0; i < slice.Len(); i++ {
		if slice.Index(i).Kind() != reflect.Int64 {
			return false
		}
	}
	return true
}

func CustomTimeValidator(fl validator.FieldLevel) bool {
	datetime := fl.Field().Interface().(time.Time)
	return datetime.After(time.Now())
}
