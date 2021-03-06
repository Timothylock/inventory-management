package items

import (
	"errors"
)

type Persister interface {
	MoveItem(ID, direction string, userID int) error
	DeleteItem(ID string, userID int) error
	SearchItems(search string) (ItemDetailList, error)
	AddItem(obj ItemDetail, overwrite bool) error
}

var ItemNotFoundErr = errors.New("item not found")
var ItemAlreadyExistsErr = errors.New("item already exists")

type ItemDetailList []ItemDetail
type ItemDetail struct {
	ID              string `db:"ID"`
	Name            string `db:"NAME"`
	Category        string `db:"CATEGORY"`
	PictureURL      string `db:"PICTURE_URL"`
	Details         string `db:"DETAILS"`
	Location        string `db:"LOCATION"`
	LastPerformedBy string `db:"USERNAME"`
	Quantity        int    `db:"QUANTITY"`
	Status          string `db:"STATUS"`
}

type Service struct {
	persister Persister
}

func NewService(p Persister) Service {
	return Service{
		persister: p,
	}
}

func (s *Service) FetchItems(search string) (ItemDetailList, error) {
	return s.persister.SearchItems(search)
}

func (s *Service) DeleteItem(id string, userID int) error {
	return s.persister.DeleteItem(id, userID)
}

func (s *Service) AddItem(item ItemDetail, overwrite bool) error {
	return s.persister.AddItem(item, overwrite)
}

func (s *Service) MoveItem(id, direction string, userID int) error {
	return s.persister.MoveItem(id, direction, userID)
}
