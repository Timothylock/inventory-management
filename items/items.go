package items

import (
	"errors"
)

type Persister interface {
	MoveItem(ID, direction string) error
	DeleteItem(ID string) error
	SearchItems(search string) (ItemDetailList, error)
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

func (s *Service) DeleteItem(id string) error {
	return s.persister.DeleteItem(id)
}

func (s *Service) AddItem(item ItemDetail) error {
	// DB Stuff

	return nil
}

func (s *Service) MoveItem(id, direction string) error {
	return s.persister.MoveItem(id, direction)
}
