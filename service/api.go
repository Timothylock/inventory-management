package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Timothylock/inventory-management/config"
	"github.com/Timothylock/inventory-management/email"
	"github.com/Timothylock/inventory-management/items"
	"github.com/Timothylock/inventory-management/middleware"
	"github.com/Timothylock/inventory-management/responses"
	"github.com/Timothylock/inventory-management/upc"
	"github.com/Timothylock/inventory-management/users"

	"io/ioutil"

	"github.com/julienschmidt/httprouter"
)

type API struct {
	itemsService items.Service
	upcService   upc.Service
	userService  users.Service
	emailService email.Service
}

func NewAPI(is items.Service, us upc.Service, user users.Service, es email.Service) API {
	return API{
		itemsService: is,
		upcService:   us,
		userService:  user,
		emailService: es,
	}
}

func NewRouter(api *API, cfg config.Config) http.Handler {
	router := httprouter.New()

	// Items
	router.Handler("GET", "/api/item/info", middleware.UserRequired(api.userService, api.SearchItems))
	router.Handler("POST", "/api/item/move", middleware.UserRequired(api.userService, api.MoveItem))
	router.Handler("POST", "/api/item", middleware.UserRequired(api.userService, api.AddItem))
	router.Handler("DELETE", "/api/item", middleware.UserRequired(api.userService, api.DeleteItem))

	// UPC
	router.Handler("GET", "/api/lookup", middleware.UserRequired(api.userService, api.LookupBarcode))

	// User
	router.Handler("GET", "/api/users", middleware.UserRequired(api.userService, api.FetchUsers))
	router.Handler("POST", "/api/user/login", api.Login())
	router.Handler("GET", "/api/user/logincheck", middleware.UserRequired(api.userService, api.LoginCheck))
	router.Handler("POST", "/api/user/add", middleware.UserRequired(api.userService, api.AddUser))
	router.Handler("DELETE", "/api/user/delete", middleware.UserRequired(api.userService, api.DeleteUser))
	router.Handler("GET", "/api/user/resetPassword", api.ForgotPassword())

	// Frontend
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(cfg.FrontendPath)))
	mux.Handle("/api/", router)

	return mux
}

func getOptionalParam(r *http.Request, name string) string {
	v, _ := getRequiredParam(r, name)
	return v
}

func getRequiredParam(r *http.Request, name string) (string, error) {
	params := r.URL.Query()[name]

	if len(params) != 1 {
		return "", fmt.Errorf("expected exactly one parameter, found %d", len(params))
	}

	param := params[0]
	if param == "" {
		return "", fmt.Errorf("expected parameter %s to be non-empty", name)
	}

	return param, nil
}

func sendJSONorErr(v interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	if json.NewEncoder(w).Encode(v) != nil {
		responses.SendError(w, responses.InternalError(errors.New("an internal server error was encountered while returning your response")))
	}
}

func parseBody(r *http.Request, target interface{}) error {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(bodyBytes, target)
}
