package service

import (
	"net/http"

	"inventory-management/items"
	"inventory-management/responses"
)

type MoveBody struct {
	Direction string `json:"direction"`
	ID        string `json:"id"`
}

func (a *API) SearchItems() http.Handler{
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

func (a *API) DeleteItem() http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := getRequiredParam(r, "id")
		if err != nil {
			responses.SendError(w, responses.MissingParamError("id"))
		}

		err = a.itemsService.DeleteItem(id)
		if err != nil && err == items.ItemNotFoundErr {
			responses.SendError(w, responses.ItemNotFound(err))
		} else if err != nil {
			responses.SendError(w, responses.InternalError(err))
		}

		sendJSONorErr(responses.Success{Success: true}, w)
	})
}

func (a *API) MoveItem() http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mb := MoveBody{}
		parseBody(r, &mb)

		if mb.ID == "" || mb.Direction == "" {
			responses.SendError(w, responses.MissingParamError("missing id or direction in the body"))
		}

		err := a.itemsService.MoveItem(mb.ID, mb.Direction)
		if err != nil && err == items.ItemNotFoundErr {
			responses.SendError(w, responses.ItemNotFound(err))
		} else if err != nil {
			responses.SendError(w, responses.InternalError(err))
		}

		sendJSONorErr(responses.Success{Success: true}, w)
	})
}