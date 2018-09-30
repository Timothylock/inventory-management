package service

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Timothylock/inventory-management/responses"
	"github.com/Timothylock/inventory-management/users"
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

		u, err := a.userService.CheckUser(ad.Username, ad.Password)
		if err != nil {
			responses.SendError(w, responses.InternalError(err))
			return
		}
		if !u.Valid {
			responses.SendError(w, responses.Unauthorized(errors.New("incorrect username or password")))
			return
		}

		expiration := time.Now().Add(31 * 24 * time.Hour)
		cookie := http.Cookie{Name: "token", Value: u.Token, Expires: expiration, Path: "/"}
		http.SetCookie(w, &cookie)

		fmt.Fprint(w, "OK")
	})
}

func (a *API) LoginCheck(u users.User) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "OK")
	})
}

func (a *API) FetchUsers(u users.User) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !u.IsSysAdmin {
			responses.SendError(w, responses.Unauthorized(errors.New("only admins can do a lookup of users")))
			return
		}

		us, err := a.userService.GetUsers()
		if err != nil {
			responses.SendError(w, responses.InternalError(err))
			return
		}

		sendJSONorErr(us, w)
	})
}
