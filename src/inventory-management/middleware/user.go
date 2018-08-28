package middleware

import (
	"errors"
	"net/http"

	"inventory-management/responses"
)

func UserRequired(h http.Handler) http.Handler  {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allow := false

		if !allow {
			responses.WithError(w, responses.NotLoggedIn(errors.New("must be logged in to access this endpoint")))
			return
		}

		h.ServeHTTP(w, r)
	})
}
