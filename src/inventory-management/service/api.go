package service

import (
	"net/http"
	"fmt"

	"inventory-management/middleware"

	"github.com/julienschmidt/httprouter"
)

type API struct {
}

func NewRouter(api *API) http.Handler {
	router := httprouter.New()

	// Items
	router.Handler("GET", "/api/item/:identifier/info", middleware.UserRequired(api.NotImplemented()))
	router.Handler("POST","/api/item/:identifier/location", middleware.UserRequired(api.NotImplemented()))
	router.Handler("DELETE","/api/item/:identifier", middleware.UserRequired(api.NotImplemented()))

	// User
	router.Handler("POST","/api/user/add", middleware.UserRequired(api.NotImplemented()))
	router.Handler("POST","/api/user/login", api.NotImplemented())
	router.Handler("DELETE","/api/user/logout", api.NotImplemented())

	// Frontend
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("/Users/timothylock/Documents/GitHub/inventory-management/frontend/")))
	mux.Handle("/api/", router)

	return mux
}

func (a *API) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func (a *API) NotImplemented() http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Not yet implemented\n")
	})
}