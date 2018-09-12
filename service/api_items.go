package service

import (
	"net/http"

	"github.com/Timothylock/inventory-management/items"
	"github.com/Timothylock/inventory-management/responses"
	"github.com/Timothylock/inventory-management/users"
)

type MoveBody struct {
	ID        string `json:"id"`
	Direction string `json:"direction"`
}

type AddBody struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Details    string `json:"details"`
	Category   string `json:"category"`
	Location   string `json:"location"`
	PictureURL string `json:"pictureURL"`
	Quantity   int    `json:"quantity"`
}

func (a *API) SearchItems(u users.User) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		search, err := getRequiredParam(r, "q")
		if err != nil {
			responses.SendError(w, responses.MissingParamError("q"))
			return
		}

		res, err := a.itemsService.FetchItems(search)
		if err != nil {
			responses.SendError(w, responses.InternalError(err))
			return
		}

		sendJSONorErr(res, w)
	})
}

func (a *API) DeleteItem(u users.User) http.Handler {
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

func (a *API) AddItem(u users.User) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ad := AddBody{}
		err := parseBody(r, &ad)
		if err != nil {
			responses.SendError(w, responses.InternalError(err))
			return
		}

		if ad.ID == "" || ad.Name == "" || ad.Category == "" || ad.Quantity == 0 {
			responses.SendError(w, responses.MissingParamError("ID, name, category, quantity must not be blank/0"))
			return
		}

		o := getOptionalParam(r, "overwrite")
		overwrite := o == "1"

		id := items.ItemDetail{
			ID:              ad.ID,
			Name:            ad.Name,
			Category:        ad.Category,
			PictureURL:      ad.PictureURL,
			Details:         ad.Details,
			Location:        ad.Location,
			LastPerformedBy: "0", // Overloaded because it contains the actual username in SearchItems
			Quantity:        ad.Quantity,
			Status:          "checked in",
		}

		err = a.itemsService.AddItem(id, overwrite)
		if err != nil && err == items.ItemAlreadyExistsErr {
			responses.SendError(w, responses.ItemAlreadyExists(err))
			return
		} else if err != nil {
			responses.SendError(w, responses.InternalError(err))
			return
		}

		sendJSONorErr(responses.Success{Success: true}, w)
	})
}

func (a *API) MoveItem(u users.User) http.Handler {
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
