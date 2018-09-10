package middleware

import (
	"errors"
	"net/http"

	"github.com/Timothylock/inventory-management/responses"
	"github.com/Timothylock/inventory-management/users"
)

func UserRequired(us users.Service, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := ""

		cookie, err := r.Cookie("token")
		if err == nil {
			token = cookie.Value
		}

		valid, err := us.IsValidToken(token)
		if err != nil {
			responses.SendError(w, responses.InternalError(err))
			return
		}

		if !valid {
			responses.SendError(w, responses.Unauthorized(errors.New("user is not authorized to make this request")))
			return
		}

		h.ServeHTTP(w, r)
	})
}
