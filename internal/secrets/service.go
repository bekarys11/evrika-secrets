package secrets

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bekarys11/evrika-secrets/pkg/common"
	"log"
	"net/http"
	"strconv"
)

func (s *Repo) getSecrets(r *http.Request) (secrets []*Secret, err error) {

	qParams := r.URL.Query()
	secretType := qParams.Get("type")
	authorId, _ := strconv.Atoi(qParams.Get("user"))

	query := s.QBuilder.Select("secrets.id, secrets.title, secrets.key, secrets.data, secrets.stype, secrets.author_id, secrets.created_at, secrets.updated_at").From("users_secrets").Join("secrets ON users_secrets.secret_id = secrets.id")

	userId, err := common.GetUserIdFromToken(r)
	if err != nil {
		return nil, err
	}

	userRole, err := common.GetRoleFromToken(r)
	if err != nil {
		return nil, err
	}

	// FILTERS
	if userRole == "user" {
		query = query.Where("users_secrets.user_id = ?", userId)
	}
	if hasType := qParams.Has("type"); hasType {
		query = query.Where("secrets.stype = ?", secretType)
	}
	if authorId != 0 {
		if userRole != "admin" {
			return nil, errors.New("вы не имеете достаточно прав")
		}
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

func (s *Repo) getById(secretId string, userRole, userId string) (secret Secret, err error) {
	userID, _ := strconv.Atoi(userId)

	query := s.QBuilder.Select("secrets.id, secrets.title, secrets.key, secrets.data, secrets.stype, secrets.author_id, secrets.created_at, secrets.updated_at").From("secrets").Where("secrets.id = ?", secretId)

	if err := query.RunWith(s.DB).QueryRow().Scan(&secret.ID, &secret.Title, &secret.Key, &secret.Data, &secret.Type, &secret.AuthorId, &secret.CreatedAt, &secret.UpdatedAt); err != nil {
		switch err {
		case sql.ErrNoRows:
			return Secret{}, errors.New("секрет не найден")
		default:
			return Secret{}, fmt.Errorf("sql query failed: %v", err)
		}
	}

	log.Printf("[DEBUG] author_id: %d \n user role: %s \n secret_id: %s \n", secret.AuthorId, userRole, secretId)
	if secret.AuthorId != userID && userRole != "admin" {
		return Secret{}, errors.New("вы не имеете достаточно прав")
	}

	return secret, nil
}

func (s *Repo) createSecret(r *http.Request, secret *Secret) error {
	var secretId int

	if err := json.NewDecoder(r.Body).Decode(&secret); err != nil {
		return fmt.Errorf("invalid JSON: %v", err)
	}

	tx, err := s.DB.BeginTxx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("unable to begin transaction: %v", err)
	}
	defer tx.Rollback()
	err = tx.QueryRowxContext(context.Background(), `INSERT INTO secrets (title, key, data, stype, author_id) VALUES ($1, $2, $3, $4, $5) RETURNING id`, secret.Title, secret.Key, secret.Data, secret.Type, secret.AuthorId).Scan(&secretId)
	switch {
	case err == sql.ErrNoRows:
		return fmt.Errorf("no secret with id: %v", err)
	case err != nil:
		return fmt.Errorf("error insert secrets to db: %v", err)
	}

	if _, err = tx.ExecContext(context.Background(), `INSERT INTO users_secrets (user_id, secret_id) VALUES ($1, $2)`, secret.AuthorId, secretId); err != nil {
		return fmt.Errorf("error insert users' secrets to db: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing db transaction: %v", err)
	}

	return nil
}
