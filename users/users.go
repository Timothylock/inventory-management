package users

type Persister interface {
	GetUser(username, password string) (User, error)
	IsValidToken(token string) (bool, error)
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

func (s *Service) IsValidToken(token string) (bool, error) {
	return s.persister.IsValidToken(token)
}
