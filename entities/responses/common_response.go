package responses

import (
	"dating-app-api/entities/requests"
	"net/http"
)

var (
	statusCreated       = http.StatusCreated
	statusOk            = http.StatusOK
	statusBadRequest    = http.StatusBadRequest
	statusNotFound      = http.StatusNotFound
	statusUnAuthorize   = http.StatusUnauthorized
	internalServerError = http.StatusInternalServerError
)

type Response struct {
	StatusCode int                             `json:"status_code"`
	Message    string                          `json:"message"`
	Meta       *requests.MetaPaginationRequest `json:"meta,omitempty"`
	Data       interface{}                     `json:"data,omitempty"`
	Validation interface{}                     `json:"validation,omitempty"`
}

type CommondResponse struct{}

func NewResponseAPI() *CommondResponse {
	return &CommondResponse{}
}

func (cmd CommondResponse) StatusOk(data interface{}, meta *requests.MetaPaginationRequest, message string) Response {
	jsonResp := Response{
		StatusCode: statusOk,
		Message:    message,
		Meta:       meta,
		Data:       data,
	}
	return jsonResp
}

func (cmd CommondResponse) StatusCreated(data interface{}, message string) Response {
	jsonResp := Response{
		StatusCode: statusCreated,
		Message:    message,
		Data:       data,
	}
	return jsonResp
}

func (cmd CommondResponse) StatusBadRequest(validation interface{}, message string) Response {
	jsonResp := Response{
		StatusCode: statusBadRequest,
		Message:    message,
		Validation: validation,
	}
	return jsonResp
}

func (cmd CommondResponse) StatusNotFound(message string) Response {
	jsonResp := Response{
		StatusCode: statusNotFound,
		Message:    message,
	}
	return jsonResp
}

func (cmd CommondResponse) StatusUnAuthorize(err string) Response {
	jsonResp := Response{
		StatusCode: statusUnAuthorize,
		Message:    err,
	}
	return jsonResp
}

func (cmd CommondResponse) StatusServerError(err string) Response {
	jsonResp := Response{
		StatusCode: internalServerError,
		Message:    err,
	}
	return jsonResp
}
