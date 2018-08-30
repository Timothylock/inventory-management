package service

import (
	"net/http"

	"github.com/Timothylock/inventory-management/items"
	"github.com/Timothylock/inventory-management/responses"
)

type MoveBody struct {
	ID        string `json:"id"`
	Direction string `json:"direction"`
}

func (a *API) SearchItems() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := getOptionalParam(r, "id")
		name := getOptionalParam(r, "name")
		cat := getOptionalParam(r, "cat")
		res, err := a.itemsService.FetchItems(id, name, cat)

		if err != nil && err == items.ItemNotFoundErr {
			responses.SendError(w, responses.ItemNotFound(err))
		} else if err != nil {
			responses.SendError(w, responses.InternalError(err))
		}

		sendJSONorErr(res, w)
	})
}

func (a *API) DeleteItem() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := getRequiredParam(r, "id")
		if err != nil {
			responses.SendError(w, responses.MissingParamError("id"))
			return
		}

		err = a.itemsService.DeleteItem(id)
		if err != nil && err == items.ItemNotFoundErr {
			responses.SendError(w, responses.ItemNotFound(err))
			return
		} else if err != nil {
			responses.SendError(w, responses.InternalError(err))
			return
		}

		sendJSONorErr(responses.Success{Success: true}, w)
	})
}

func (a *API) MoveItem() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mb := MoveBody{}
		err := parseBody(r, &mb)
		if err != nil {
			responses.SendError(w, responses.InternalError(err))
			return
		}

		if mb.ID == "" || mb.Direction == "" {
			responses.SendError(w, responses.MissingParamError("missing id or direction in the body"))
			return
		}

		err = a.itemsService.MoveItem(mb.ID, mb.Direction)
		if err != nil && err == items.ItemNotFoundErr {
			responses.SendError(w, responses.ItemNotFound(err))
			return
		} else if err != nil {
			responses.SendError(w, responses.InternalError(err))
			return
		}

		sendJSONorErr(responses.Success{Success: true}, w)
	})
}
