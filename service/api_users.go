package service

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
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

type UserBody struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	Email      string `json:"email"`
	IsSysAdmin string `json:"is_sys_admin"`
}

func (a *API) AddUser(u users.User) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !u.IsSysAdmin {
			responses.SendError(w, responses.Unauthorized(errors.New("you are not authorized to perform this action")))
			return
		}

		ad := UserBody{}
		err := parseBody(r, &ad)
		if err != nil {
			responses.SendError(w, responses.InternalError(err))
			return
		}

		if err = a.userService.AddUser(ad.Username, strings.ToLower(ad.Email), ad.Password, ad.IsSysAdmin == "true"); err != nil {
			responses.SendError(w, responses.InternalError(err))
			return
		}

		sendJSONorErr("Success", w)
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

func (a *API) DeleteUser(u users.User) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !u.IsSysAdmin {
			responses.SendError(w, responses.Unauthorized(errors.New("you are not authorized to perform this action")))
			return
		}

		uname, err := getRequiredParam(r, "u")
		if err != nil {
			responses.SendError(w, responses.MissingParamError("u"))
			return
		}

		targetU, err := a.userService.CheckUserByUsername(uname, u.ID)
		if err != nil {
			responses.SendError(w, responses.InternalError(err))
			return
		}

		if !targetU.Valid {
			responses.SendError(w, responses.InternalError(errors.New("username not found or already deleted")))
			return
		}

		if targetU.ID == 0 {
			responses.SendError(w, responses.InternalError(errors.New("cannot delete System user")))
			return
		}

		if err = a.userService.DeleteUser(targetU.ID, u.ID); err != nil {
			responses.SendError(w, responses.InternalError(err))
			return
		}

		sendJSONorErr("Success", w)
	})
}

func (a *API) ForgotPassword() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uname, err := getRequiredParam(r, "username")
		if err != nil {
			responses.SendError(w, responses.MissingParamError("username"))
			return
		}

		email, err := getRequiredParam(r, "email")
		if err != nil {
			responses.SendError(w, responses.MissingParamError("email"))
			return
		}

		targetU, err := a.userService.CheckUserByUsername(uname, 0)
		if err != nil {
			responses.SendError(w, responses.InternalError(err))
			return
		}

		if targetU.ID == 0 {
			responses.SendError(w, responses.InternalError(errors.New("cannot reset System user")))
			return
		}

		if strings.ToLower(targetU.Email) != strings.ToLower(email) {
			responses.SendError(w, responses.InternalError(errors.New("no username with that email on record")))
			return
		}

		newPass := randomString(12)
		if err := a.emailService.SendEmail(targetU.Email, "Inventory Password Reset", "<p>Your new password is <b>"+newPass+"</b>. Please change it once you log in. </p>"); err != nil {
			responses.SendError(w, responses.InternalError(err))
			return
		}

		if err = a.userService.EditUser(targetU.Username, targetU.Email, newPass, targetU.IsSysAdmin); err != nil {
			responses.SendError(w, responses.InternalError(err))
			return
		}

		sendJSONorErr("Success", w)
	})
}

func randomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(65 + rand.Intn(25))
	}
	return string(bytes)
}
