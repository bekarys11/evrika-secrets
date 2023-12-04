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
	userQ := qParams.Get("user")
	authorId, _ := strconv.Atoi(userQ)
	log.Printf("[DEBUG] user type: %d\n secret type: %s", authorId, secretType)
	query := `
		SELECT secrets.id, secrets.title, secrets.key, secrets.data, secrets.stype, secrets.author_id, secrets.created_at, secrets.updated_at
		FROM users_secrets
		JOIN secrets ON users_secrets.secret_id = secrets.id
		WHERE (secrets.author_id = $1 OR $1 = 0)
`

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
	if userRole == "user" {
		query += ` AND users_secrets.user_id = ` + strconv.Itoa(int(userId.(float64)))
	}

	if hasType := qParams.Has("type"); hasType {
		query += ` AND secrets.stype = '` + secretType + `'`
	}
	query += ` LIMIT 30;`
	log.Println(query)
	rows, err := s.DB.Queryx(query, authorId)
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
