package secrets

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"log"
	"net/url"
	"strconv"
	"time"
)

type Repo struct {
	DB       *sqlx.DB
	QBuilder squirrel.StatementBuilderType
}

func (s *Repo) getSecrets(qParams url.Values, userRole, userId string) (secrets []*Secret, err error) {
	secretType := qParams.Get("type")
	userIDQuery, _ := strconv.Atoi(qParams.Get("user"))

	query := s.QBuilder.
		Select("secrets.id, secrets.title, secrets.key, secrets.data, secrets.stype, secrets.author_id, secrets.created_at, secrets.updated_at").
		From("users_secrets").
		Join("secrets ON users_secrets.secret_id = secrets.id")

	// FILTERS
	if userRole == "user" {
		query = query.Where("users_secrets.user_id = ?", userId)
	}
	if hasType := qParams.Has("type"); hasType {
		query = query.Where("secrets.stype = ?", secretType)
	}
	if userIDQuery != 0 {
		if userRole != "admin" {
			return nil, errors.New("вы не имеете достаточно прав")
		}

		// admin can see any user's secrets
		query = query.Where("users_secrets.user_id = ?", userIDQuery)
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

func (s *Repo) createSecret(secret *Secret) error {
	var secretId int

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

func (s *Repo) shareSecret(usersSecrets UsersSecret) error {
	txn, err := s.DB.Beginx()

	stmt, err := txn.Preparex(pq.CopyIn("users_secrets", "user_id", "secret_id"))
	if err != nil {
		return fmt.Errorf("error preparing users' secrets transaction: %v", err)
	}

	for _, v := range usersSecrets.UserIds {
		if _, err = stmt.Exec(v, usersSecrets.SecretId); err != nil {
			return fmt.Errorf("error executing users' secrets transaction: %v", err)
		}
	}
	defer stmt.Close()

	if _, err = stmt.Exec(); err != nil {
		return fmt.Errorf("error execute users' secrets transaction: %v", err)
	}

	if err = txn.Commit(); err != nil {
		return fmt.Errorf("error commiting users' secrets transaction: %v", err)
	}

	return nil
}

func (s *Repo) updateById(secretId, userRole, userId string, payload SecretReq) error {
	if userRole == "user" {
		if err := s.checkForSecretAuthor(secretId, userRole, userId); err != nil {
			return err
		}
	}

	query := s.QBuilder.
		Update("secrets").
		SetMap(map[string]interface{}{
			"title":      payload.Title,
			"key":        payload.Key,
			"data":       payload.Data,
			"stype":      payload.Type,
			"author_id":  payload.AuthorId,
			"updated_at": time.Now(),
		}).
		Where("id = ?", secretId)

	_, err := query.RunWith(s.DB).Exec()
	if err != nil {
		return fmt.Errorf("error executing update query: %v", err)
	}

	return nil
}

func (s *Repo) deleteById(secretId, userRole, userId string) error {

	if userRole == "user" {
		if err := s.checkForSecretAuthor(secretId, userRole, userId); err != nil {
			return err
		}
	}

	query := s.QBuilder.Delete("secrets").Where("id = ?", secretId)

	_, err := query.RunWith(s.DB).Exec()
	if err != nil {
		return fmt.Errorf("error executing delete query: %v", err)
	}

	return nil
}

// if given user has "user role", it checks if the user is author of the given secret
func (s *Repo) checkForSecretAuthor(secretId, userRole, userId string) error {
	secret, err := s.getById(secretId, userRole, userId)
	if err != nil {
		return fmt.Errorf("error querying secret: %v", err)
	}

	userID, _ := strconv.Atoi(userId)
	if secret.AuthorId != userID {
		return errors.New("вы не имеете достаточно прав")
	}
	return nil
}
