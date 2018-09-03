package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Timothylock/inventory-management/config"
	"github.com/Timothylock/inventory-management/items"
	"github.com/Timothylock/inventory-management/middleware"
	"github.com/Timothylock/inventory-management/responses"
	"github.com/Timothylock/inventory-management/upc"

	"io/ioutil"

	"github.com/julienschmidt/httprouter"
)

type API struct {
	itemsService items.Service
	upcService   upc.Service
}

func NewAPI(is items.Service, us upc.Service) API {
	return API{
		itemsService: is,
		upcService:   us,
	}
}

func NewRouter(api *API, cfg config.Config) http.Handler {
	router := httprouter.New()

	// Items
	router.Handler("GET", "/api/item/info", middleware.UserRequired(api.SearchItems()))
	router.Handler("POST", "/api/item/move", middleware.UserRequired(api.MoveItem()))
	router.Handler("POST", "/api/item", middleware.UserRequired(api.AddItem()))
	router.Handler("DELETE", "/api/item", middleware.UserRequired(api.DeleteItem()))

	// UPC
	router.Handler("GET", "/api/lookup", middleware.UserRequired(api.LookupBarcode()))

	// User
	router.Handler("POST", "/api/user/add", middleware.UserRequired(api.NotImplemented()))
	router.Handler("POST", "/api/user/login", api.NotImplemented())
	router.Handler("DELETE", "/api/user/logout", api.NotImplemented())

	// Frontend
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(cfg.FrontendPath)))
	mux.Handle("/api/", router)

	return mux
}

func (a *API) NotImplemented() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Not yet implemented\n")
	})
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
