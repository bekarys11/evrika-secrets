package roles

type Repo interface {
	GetRoles() ([]*Role, error)
}

type RoleService struct {
	repository Repo
}

func NewRoleService(repo Repo) RoleService {
	return RoleService{
		repository: repo,
	}
}

func (s RoleService) GetRoles() ([]*Role, error) {
	return s.repository.GetRoles()
}
