package users

import (
	"errors"
)

type Persister interface {
	GetUser(username, password string) (User, error)
	GetUserByToken(token string) (User, error)
	GetUserByUsername(username string, curUserID int) (User, error)
	AddUser(username, email, password string, isSysAdmin bool) error
	GetUsers() (MultipleUsers, error)
	DeleteUser(targetID, userID int) error
}

type MultipleUsers []User
type User struct {
	Valid      bool   `json:"-"`
	ID         int    `json:"-"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	IsSysAdmin bool   `json:"isSysAdmin"`
	Token      string `json:"-"`
}

type Service struct {
	persister Persister
}

func NewService(p Persister) Service {
	return Service{
		persister: p,
	}
}

func (s *Service) CheckUser(username, password string) (User, error) {
	return s.persister.GetUser(username, password)
}

func (s *Service) CheckUserByToken(token string) (User, error) {
	return s.persister.GetUserByToken(token)
}

func (s *Service) AddUser(username, email, password string, isSysAdmin bool) error {
	return s.persister.AddUser(username, email, password, isSysAdmin)
}

func (s *Service) GetUsers() (MultipleUsers, error) {
	return s.persister.GetUsers()
}

func (s *Service) DeleteUser(targetUsername string, curUserID int) error {
	targetU, err := s.persister.GetUserByUsername(targetUsername, curUserID)
	if err != nil {
		return err
	}
	if targetU.ID == 0 {
		return errors.New("username not found, already deleted, or you tried deleting the system user")
	}

	return s.persister.DeleteUser(targetU.ID, curUserID)
}
