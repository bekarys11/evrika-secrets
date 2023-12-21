package secrets

import (
	"context"
	"net/url"
)

type Repo interface {
	GetSecrets(qParams url.Values, userRole, userId string) ([]*SecretResp, error)
	GetSecretById(secretId string, role string, userId string) (secret SecretResp, err error)
	CreateSecret(ctx context.Context, payload Secret) error
	UpdateSecret(secretId, userRole, userId string, payload Secret) error
	DeleteSecret(secretId, userRole, userId string) error
}

type SecretService struct {
	repository Repo
}

func NewSecretService(repo Repo) SecretService {
	return SecretService{
		repository: repo,
	}
}

func (s SecretService) GetSecrets(qParams url.Values, userRole string, userId string) ([]*SecretResp, error) {
	return s.repository.GetSecrets(qParams, userRole, userId)
}

func (s SecretService) GetSecretById(secretId string, role string, userId string) (SecretResp, error) {
	return s.repository.GetSecretById(secretId, role, userId)
}

func (s SecretService) CreateSecret(ctx context.Context, payload Secret) error {
	return s.repository.CreateSecret(ctx, payload)
}

func (s SecretService) UpdateSecret(secretId, userRole, userId string, payload Secret) error {
	return s.repository.UpdateSecret(secretId, userRole, userId, payload)
}

func (s SecretService) DeleteSecret(secretId, userRole, userId string) error {
	return s.repository.DeleteSecret(secretId, userRole, userId)
}
