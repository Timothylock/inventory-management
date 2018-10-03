package users

type Persister interface {
	GetUser(username, password string) (User, error)
	GetUserByToken(token string) (User, error)
	AddUser(username, email, password string, isSysAdmin bool) error
}

type User struct {
	Valid      bool
	ID         int
	Username   string
	Email      string
	IsSysAdmin bool
	Token      string
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
