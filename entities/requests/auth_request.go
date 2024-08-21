package requests

import "github.com/thedevsaddam/govalidator"

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthRequest) ValiadateAuthLogin() interface{} {

	validator := govalidator.New(govalidator.Options{
		Data: h,
		Rules: govalidator.MapData{
			"username": []string{"required", "max:50"},
			"password": []string{"required"},
		},
		RequiredDefault: true,
	}).ValidateStruct()

	if len(validator) > 0 {
		return validator
	}

	return nil
}

func (h *AuthRequest) ValiadateChangePassword() interface{} {

	validator := govalidator.New(govalidator.Options{
		Data: h,
		Rules: govalidator.MapData{
			"password": []string{"required"},
		},
		RequiredDefault: true,
	}).ValidateStruct()

	if len(validator) > 0 {
		return validator
	}

	return nil
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthRequest) ValiadateRefreshToken() interface{} {

	validator := govalidator.New(govalidator.Options{
		Data: h,
		Rules: govalidator.MapData{
			"refresh_token": []string{"required"},
		},
		RequiredDefault: true,
	}).ValidateStruct()

	if len(validator) > 0 {
		return validator
	}

	return nil
}

type ChangePasswordRequest struct {
	Password string `json:"password"`
	UserId   string
}

func (h *ChangePasswordRequest) ValiadateNewPassword() interface{} {

	validator := govalidator.New(govalidator.Options{
		Data: h,
		Rules: govalidator.MapData{
			"password": []string{"required"},
		},
		RequiredDefault: true,
	}).ValidateStruct()

	if len(validator) > 0 {
		return validator
	}

	return nil
}
