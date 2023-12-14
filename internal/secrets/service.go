package secrets

import (
	"context"
	"database/sql"
	"encoding/json"
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

func (s *Repo) getAllSecrets(qParams url.Values, userRole, userId string) (secrets []*SecretResp, err error) {
	var (
		q          string
		args       []interface{}
		secretType = qParams.Get("type")
	)

	gQuery := s.QBuilder.
		Select("secrets.id", "secrets.title", "secrets.key", "secrets.data", "secrets.stype", "secrets.author_id").
		Column(`
					jsonb_agg(
					  jsonb_build_object(
							  'id', users.id,
							  'name', users.name
						  )
					  ) AS users_info
					`).
		From("secrets").
		Join("users_secrets ON users_secrets.secret_id = secrets.id").
		Join("users ON users_secrets.user_id = users.id").
		GroupBy("secrets.id", "secrets.title", "secrets.key", "secrets.data", "secrets.stype", "secrets.author_id")

	if userRole == "user" {
		userQ := s.QBuilder.
			Select("ss.id", "ss.title", "ss.key", "ss.data", "ss.stype", "ss.author_id", "ss.users_info").
			FromSelect(gQuery, "ss").
			Where(fmt.Sprintf("jsonb_path_exists(users_info, '$[*] ?? (@.id == %s)')", userId))

		if hasType := qParams.Has("type"); hasType {
			userQ = userQ.Where("ss.stype = ?", secretType)
		}

		q, args, err = userQ.ToSql()
		log.Printf("[DEBUG] query: %s; args: %v", q, args)
		if err != nil {
			log.Printf("Error to_sql: %v", err)
		}

	} else {
		if hasType := qParams.Has("type"); hasType {
			gQuery = gQuery.Where("secrets.stype = ?", secretType)
		}

		q, args, err = gQuery.ToSql()
		if err != nil {
			log.Printf("Error to_sql: %v", err)
		}
	}

	rows, err := s.DB.Query(q, args...)
	if err != nil {
		return nil, fmt.Errorf("secrets query error: %v", err)
	}

	for rows.Next() {
		var secret SecretResp
		if err := rows.Scan(
			&secret.ID, &secret.Title, &secret.Key, &secret.Data, &secret.Type, &secret.AuthorId, &secret.Users,
		); err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}

		secrets = append(secrets, &secret)

		// Unmarshal the bytes into a User struct
		var user interface{}
		err = json.Unmarshal(secret.Users.([]byte), &user)
		if err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
		}

		secret.Users = user
	}
	defer rows.Close()

	return secrets, nil
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

var q = `
SELECT * FROM (
                  SELECT
                      s.id,
                      s.title,
                      s.key,
                      s.data, s.stype, s.author_id,
                      jsonb_agg(
                              jsonb_build_object(
                                      'user_id', u.id,
                                      'user_name', u.name
                                  )
                          ) AS users
                  FROM
                      secrets s
                          JOIN users_secrets us ON us.secret_id = s.id
                          JOIN users u ON us.user_id = u.id
                  GROUP BY s.id, s.title, s.key, s.data, s.stype, s.author_id, s.created_at, s.updated_at
              ) secrets_sub
         WHERE jsonb_path_exists(users, '$[*] ? (@.user_id == 56)');
`

//q := fmt.Sprintf(`
//	SELECT
//    secrets_sub.id, secrets_sub.title, secrets_sub.key, secrets_sub.data, secrets_sub.stype, secrets_sub.author_id, secrets_sub.created_at, secrets_sub.updated_at, secrets_sub.users
//	FROM (
//                  SELECT
//                      s.id,
//                      s.title,
//                      s.key,
//                      s.data, s.stype, s.author_id, s.created_at, s.updated_at,
//                      jsonb_agg(
//                              jsonb_build_object(
//                                      'id', u.id,
//                                      'name', u.name
//                                  )
//                          ) AS users
//                  FROM
//                      secrets s
//                          JOIN users_secrets us ON us.secret_id = s.id
//                          JOIN users u ON us.user_id = u.id
//                  GROUP BY s.id, s.title, s.key, s.data, s.stype, s.author_id, s.created_at, s.updated_at
//              ) secrets_sub
//         WHERE jsonb_path_exists(users, '$[*] ? (@.user_id == %s)');
//`, userId)
