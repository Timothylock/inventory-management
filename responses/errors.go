package responses

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type httpError struct {
	StatusCode int    `json:"-"`
	ErrorCode  int    `json:"code"`
	Message    string `json:"details"`
}

func SendError(w http.ResponseWriter, err httpError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.StatusCode)

	b, e := json.Marshal(err)
	if e != nil {
		fmt.Sprintf("failed marshalling error - %s", e)
	}

	w.Write(b)
}

func InternalError(err error) httpError {
	return httpError{
		StatusCode: http.StatusInternalServerError,
		ErrorCode:  1000,
		Message:    err.Error(),
	}
}

func MissingParamError(param string) httpError {
	return httpError{
		StatusCode: http.StatusBadRequest,
		ErrorCode:  1002,
		Message:    fmt.Sprintf("Missing param - %s", param),
	}
}

func Unauthorized(err error) httpError {
	return httpError{
		StatusCode: http.StatusUnauthorized,
		ErrorCode:  1001,
		Message:    err.Error(),
	}
}

func ItemNotFound(err error) httpError {
	return httpError{
		StatusCode: http.StatusNotFound,
		ErrorCode:  1100,
		Message:    err.Error(),
	}
}

func ItemAlreadyExists(err error) httpError {
	return httpError{
		StatusCode: http.StatusBadRequest,
		ErrorCode:  1101,
		Message:    err.Error(),
	}
}
