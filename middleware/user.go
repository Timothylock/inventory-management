package middleware

import (
	"errors"
	"net/http"

	"github.com/Timothylock/inventory-management/responses"
)

func UserRequired(h http.Handler) http.Handler  {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allow := true

		if !allow {
			responses.SendError(w, responses.NotLoggedIn(errors.New("must be logged in to access this endpoint")))
			return
		}

		h.ServeHTTP(w, r)
	})
}
