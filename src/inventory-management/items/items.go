package items

import "errors"

type Service struct {}

var ItemNotFoundErr = errors.New("item not found")

type ItemDetailList []ItemDetail
type ItemDetail struct {
	ID, Name, Category, PictureURL, Details, Location, LastPerformedBy string
	Quantity int
}

func (s *Service) FetchItems(id, name, category string) (ItemDetailList, error) {
	// DB Stuff

	return ItemDetailList{
		{
			ID: "532532234",
			Name: "Wrench",
			Category: "tools",
			PictureURL: "google.ca",
			Details: "A wrench that does stuff",
			Location: "Locker A",
			LastPerformedBy: "Timothy",
			Quantity: 1,
		},
		{
			ID: id,
			Name: name,
			Category: category,
			PictureURL: "google.ca",
			Details: "A wrench that does stuff",
			Location: "Locker A",
			LastPerformedBy: "Timothy",
			Quantity: 1,
		},
	}, nil
}

func (s *Service) DeleteItem(id string) error {
	// DB Stuff

	return nil
}

func (s *Service) AddItem(id string) error {
	// DB Stuff

	return nil
}


func (s *Service) MoveItem(id, direction string) error {
	// DB Stuff

	return nil
}