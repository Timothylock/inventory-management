package users

type Persister interface {
	GetUser(username, password string) (User, error)
	GetUserByToken(token string) (User, error)
	GetUserByUsername(username string, curUserID int) (User, error)
	AddUser(username, email, password string, isSysAdmin, overwrite bool) error
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

func (s *Service) CheckUserByUsername(username string, curUserID int) (User, error) {
	return s.persister.GetUserByUsername(username, curUserID)
}

func (s *Service) AddUser(username, email, password string, isSysAdmin bool) error {
	return s.persister.AddUser(username, email, password, isSysAdmin, false)
}

func (s *Service) GetUsers() (MultipleUsers, error) {
	return s.persister.GetUsers()
}

func (s *Service) DeleteUser(targetID int, curUserID int) error {
	return s.persister.DeleteUser(targetID, curUserID)
}

func (s *Service) EditUser(username, email, password string, isSysAdmin bool) error {
	return s.persister.AddUser(username, email, password, isSysAdmin, true)
}
