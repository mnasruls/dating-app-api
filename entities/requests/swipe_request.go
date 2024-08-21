package requests

import "github.com/thedevsaddam/govalidator"

type SwipeRequest struct {
	UserId string `json:"user_id"`
	Type   string `json:"type"`
}

func (h *SwipeRequest) ValiadateSwipe() interface{} {

	validator := govalidator.New(govalidator.Options{
		Data: h,
		Rules: govalidator.MapData{
			"user_id": []string{"required", "uuid"},
			"type":    []string{"required", "in:left,right"},
		},
		RequiredDefault: true,
	}).ValidateStruct()

	if len(validator) > 0 {
		return validator
	}

	return nil
}
