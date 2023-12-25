package auth

type Repo interface {
	Login(LoginPayload) (token string, err error)
}

type AuthService struct {
	repository Repo
}

func NewAuthService(repo Repo) AuthService {
	return AuthService{
		repository: repo,
	}
}

func (s AuthService) Login(payload LoginPayload) (token string, err error) {
	return s.repository.Login(payload)
}
