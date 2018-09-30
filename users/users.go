package users

type Persister interface {
	GetUser(username, password string) (User, error)
	GetUserByToken(token string) (User, error)
	GetUsers() (MultipleUsers, error)
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

func (s *Service) GetUsers() (MultipleUsers, error) {
	return s.persister.GetUsers()
}
