package service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Timothylock/inventory-management/responses"
)

type LoginBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a *API) Login() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ad := LoginBody{}
		err := parseBody(r, &ad)
		if err != nil {
			responses.SendError(w, responses.InternalError(err))
			return
		}

		expiration := time.Now().Add(31 * 24 * time.Hour)
		cookie := http.Cookie{Name: "token", Value: "atoken", Expires: expiration}

		http.SetCookie(w, &cookie)

		fmt.Fprint(w, "OK")
	})
}

func (a *API) LoginCheck() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "OK")
	})
}
