package requests

import "github.com/thedevsaddam/govalidator"

type CreateUserRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
}

func (h *CreateUserRequest) ValiadateCreateUser() interface{} {

	validator := govalidator.New(govalidator.Options{
		Data: h,
		Rules: govalidator.MapData{
			"username":     []string{"required", "char_libs", "min:3", "max:50"},
			"password":     []string{"required"},
			"phone_number": []string{"numeric_null_libs", "max:15"},
		},
		RequiredDefault: true,
	}).ValidateStruct()

	if len(validator) > 0 {
		return validator
	}

	return nil
}

func (h *CreateUserRequest) ValiadateCheckUsername() interface{} {

	validator := govalidator.New(govalidator.Options{
		Data: h,
		Rules: govalidator.MapData{
			"username": []string{"required", "char_libs", "max:50"},
		},
		RequiredDefault: true,
	}).ValidateStruct()

	if len(validator) > 0 {
		return validator
	}

	return nil
}

type UpdateUserRequest struct {
	Username    string `json:"username"`
	NewPassword string `json:"new_password"`
	OldPassword string `json:"old_password"`
	PhoneNumber string `json:"phone_number"`
}

func (h *UpdateUserRequest) ValiadateUpdateUser() interface{} {

	validator := govalidator.New(govalidator.Options{
		Data: h,
		Rules: govalidator.MapData{
			"username":     []string{"required", "char_libs", "max:50"},
			"email":        []string{"required", "email", "max:50"},
			"password":     []string{"required"},
			"phone_number": []string{"numeric_null_libs", "max:15"},
		},
		RequiredDefault: true,
	}).ValidateStruct()

	if len(validator) > 0 {
		return validator
	}

	return nil
}

type VerifyUser struct {
	VANumber string  `json:"va_number"`
	Amount   float64 `json:"amount"`
}

func (h *VerifyUser) VerifyUser() interface{} {

	validator := govalidator.New(govalidator.Options{
		Data: h,
		Rules: govalidator.MapData{
			"va_number": []string{"required", "number", "max:16"},
			"amount":    []string{"required", "numeric"},
		},
		RequiredDefault: true,
	}).ValidateStruct()

	if len(validator) > 0 {
		return validator
	}

	return nil
}

type CountRequest struct {
	Count int `json:"count"`
}
