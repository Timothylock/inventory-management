package middleware

import (
	"errors"
	"net/http"

	"github.com/Timothylock/inventory-management/responses"
)

func UserRequired(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			responses.SendError(w, responses.InternalError(err))
		}

		allow := cookie.Value == "atoken"

		if !allow {
			responses.SendError(w, responses.Unauthorized(errors.New("user is not authorized to make this request")))
			return
		}

		h.ServeHTTP(w, r)
	})
}
