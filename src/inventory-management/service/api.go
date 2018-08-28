package service

import (
	"net/http"
	"fmt"

	"github.com/julienschmidt/httprouter"
)

type API struct {
}

func NewRouter(api *API) http.Handler {
	router := httprouter.New()
	router.GET("/", api.Index)
	router.GET("/hello/:name", api.Hello)

	return router
}

func (a *API) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func (a *API) Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}