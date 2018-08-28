package responses

import (
	"encoding/json"
	"net/http"
	"fmt"
)

type httpError struct {
	StatusCode   int    `json:"-"`
	InternalCode int    `json:"code"`
	Message      string `json:"message"`
}

func WithError(w http.ResponseWriter, err httpError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.StatusCode)

	b, merr := json.Marshal(err)
	if merr != nil {
		fmt.Sprintf("failed marshalling error - %s", merr)
	}

	w.Write(b)
}

func NotLoggedIn(err error) httpError {
	return httpError{
		StatusCode:   http.StatusUnauthorized,
		InternalCode: 1001,
		Message:      err.Error(),
	}
}


