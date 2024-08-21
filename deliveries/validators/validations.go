package validators

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/thedevsaddam/govalidator"
)

func ValidatePassword(password string) error {
	var (
		hasMinLen    = len(password) >= 8
		hasMaxLen    = len(password) > 50
		hasUppercase = regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasLowercase = regexp.MustCompile(`[a-z]`).MatchString(password)
		hasNumber    = regexp.MustCompile(`[0-9]`).MatchString(password)
		hasSpecial   = regexp.MustCompile(`[!@#~$%^&*()+|_{}<>?,.;:'"]`).MatchString(password)
	)

	if !hasMinLen {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	if !hasUppercase {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !hasLowercase {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return fmt.Errorf("password must contain at least one number")
	}
	if !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}
	if hasMaxLen {
		return fmt.Errorf("password must be at less than 50 characters long")
	}

	return nil
}

const (
	Numeric string = "^-?[0-9]+$"
	Key     string = "^[-a-zA-Z0-9_-]+$"
)

var (
	regexNumeric = regexp.MustCompile(Numeric)
	regexKey     = regexp.MustCompile(Key)
)

func AddValidatorLibs() {
	govalidator.AddCustomRule("numeric_null_libs", func(field string, rule string, message string, value interface{}) error {
		str := toString(value)
		if str == "" {
			return nil
		}

		err := fmt.Errorf("the %s field must be a valid numeric", field)
		if message != "" {
			err = errors.New(message)
		}

		if !isNumeric(str) {
			return err
		}

		return nil
	})
	govalidator.AddCustomRule("char_libs", func(field string, rule string, message string, value interface{}) error {
		str := toString(value)
		if str == "" {
			return nil
		}

		err := fmt.Errorf("the %s field must be a contains alpha numeric, space, dot, comma, underscore, dash, slash and brackets", field)
		if message != "" {
			err = errors.New(message)
		}

		if !isIdemKey(str) {
			return err
		}

		return nil
	})
}

func toString(v interface{}) string {
	str, ok := v.(string)
	if !ok {
		str = fmt.Sprintf("%v", v)
	}
	return str
}

func isNumeric(str string) bool {
	return regexNumeric.MatchString(str)
}

func isIdemKey(str string) bool {
	return regexKey.MatchString(str)
}
