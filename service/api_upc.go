package service

import (
	"net/http"

	"github.com/Timothylock/inventory-management/responses"
	"github.com/Timothylock/inventory-management/users"
)

func (a *API) LookupBarcode(u users.User) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		search, err := getRequiredParam(r, "barcode")
		if err != nil {
			responses.SendError(w, responses.MissingParamError("barcode"))
			return
		}

		res, err := a.upcService.LookupBarcode(search)
		if err != nil {
			responses.SendError(w, responses.InternalError(err))
			return
		}

		sendJSONorErr(res, w)
	})
}
