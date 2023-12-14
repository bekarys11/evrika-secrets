package users

type Repo interface {
	GetUsers() ([]*UserResp, error)
	CreateUser(payload User) error
	GetProfile(userId string) (UserResp, error)
}

type UserService struct {
	repository Repo
}

func NewUserService(repo Repo) UserService {
	return UserService{
		repository: repo,
	}
}

func (s UserService) GetUsers() ([]*UserResp, error) {
	return s.repository.GetUsers()
}

func (s UserService) CreateUser(payload User) error {
	return s.repository.CreateUser(payload)
}

func (s UserService) GetProfile(userId string) (UserResp, error) {
	return s.repository.GetProfile(userId)
}
