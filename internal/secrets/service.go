package secrets

import (
	"errors"
	"fmt"
	"github.com/bekarys11/evrika-secrets/internal/users"
	"log"
	"net/http"
	"strconv"
)

func (s *Repo) getSecrets(r *http.Request) (secrets []*Secret, err error) {

	qParams := r.URL.Query()
	secretType := qParams.Get("type")
	authorId, _ := strconv.Atoi(qParams.Get("user"))

	query := s.QBuilder.Select("secrets.id, secrets.title, secrets.key, secrets.data, secrets.stype, secrets.author_id, secrets.created_at, secrets.updated_at").From("users_secrets").Join("secrets ON users_secrets.secret_id = secrets.id")

	claims, err := users.GetTokenClaims(r)
	if err != nil {
		log.Printf("error getting token claims: %v", err)
		return nil, err
	}
	userId, ok := claims["user_id"]
	if !ok {
		return nil, errors.New("user id is not provided")
	}
	userRole := claims["role"]

	// FILTERS
	if userRole == "user" {
		query = query.Where("users_secrets.user_id = ?", strconv.Itoa(int(userId.(float64))))
	}
	if hasType := qParams.Has("type"); hasType {
		query = query.Where("secrets.stype = ?", secretType)
	}
	if authorId != 0 {
		query = query.Where("secrets.author_id = ?", authorId)
	}

	rows, err := query.RunWith(s.DB).Query()
	if err != nil {
		return nil, fmt.Errorf("sql query error: %v", err)
	}

	for rows.Next() {
		var secret Secret
		if err := rows.Scan(&secret.ID, &secret.Title, &secret.Key, &secret.Data, &secret.Type, &secret.AuthorId, &secret.CreatedAt, &secret.UpdatedAt); err != nil {
			return nil, err
		}
		secrets = append(secrets, &secret)
	}

	return secrets, nil
}
