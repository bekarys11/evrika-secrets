package secrets

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	resp "github.com/bekarys11/evrika-secrets/pkg/response"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"net/http"
)

type Repo struct {
	DB *sqlx.DB
}

func (s *Repo) All(w http.ResponseWriter, r *http.Request) {
	var secrets []*Secret
	vars := mux.Vars(r)
	userId, ok := vars["user_id"]
	if !ok {
		resp.ErrorJSON(w, errors.New("user id is not provided"), http.StatusBadRequest)
		return
	}

	rows, err := s.DB.Queryx(`
		SELECT secrets.id, secrets.title, secrets.key, secrets.data, secrets.author_id, secrets.created_at, secrets.updated_at
		FROM users_secrets
        JOIN secrets ON users_secrets.secret_id = secrets.id
		WHERE users_secrets.user_id = $1
		LIMIT 10;
`, userId)

	if err != nil {
		resp.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		var secret Secret
		if err := rows.Scan(&secret.ID, &secret.Title, &secret.Key, &secret.Data, &secret.AuthorId, &secret.CreatedAt, &secret.UpdatedAt); err != nil {
			resp.ErrorJSON(w, err, http.StatusInternalServerError)
			return
		}
		secrets = append(secrets, &secret)
	}
	resp.WriteApiJSON(w, http.StatusOK, secrets)
}

func (s *Repo) Create(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var (
		secret   Secret
		secretId int
	)

	if err := json.NewDecoder(r.Body).Decode(&secret); err != nil {
		resp.ErrorJSON(w, fmt.Errorf("invalid JSON: %v", err), http.StatusBadRequest)
		return
	}
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		resp.ErrorJSON(w, fmt.Errorf("unable to begin transaction: %v", err), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()
	err = tx.QueryRowxContext(ctx, `INSERT INTO secrets (title, key, data, author_id) VALUES ($1, $2, $3, $4) RETURNING id`, secret.Title, secret.Key, secret.Data, secret.AuthorId).Scan(&secretId)
	switch {
	case err == sql.ErrNoRows:
		resp.ErrorJSON(w, fmt.Errorf("no secret with id: %v", err), http.StatusInternalServerError)
		return
	case err != nil:
		resp.ErrorJSON(w, fmt.Errorf("error insert secrets to db: %v", err), http.StatusInternalServerError)
		return
	}

	if _, err = tx.ExecContext(ctx, `INSERT INTO users_secrets (user_id, secret_id) VALUES ($1, $2)`, secret.AuthorId, secretId); err != nil {
		resp.ErrorJSON(w, fmt.Errorf("error insert users' secrets to db: %v", err), http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		resp.ErrorJSON(w, fmt.Errorf("error committing db transaction: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Repo) ShareSecret(w http.ResponseWriter, r *http.Request) {
	var usersSecrets UsersSecret

	if err := json.NewDecoder(r.Body).Decode(&usersSecrets); err != nil {
		resp.ErrorJSON(w, fmt.Errorf("invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	txn, err := s.DB.Beginx()

	stmt, err := txn.Preparex(pq.CopyIn("users_secrets", "user_id", "secret_id"))
	if err != nil {
		resp.ErrorJSON(w, fmt.Errorf("error preparing users' secrets transaction: %v", err), http.StatusInternalServerError)
		return
	}

	for _, v := range usersSecrets.UserIds {
		if _, err = stmt.Exec(v, usersSecrets.SecretId); err != nil {
			resp.ErrorJSON(w, fmt.Errorf("error executing users' secrets transaction: %v", err), http.StatusInternalServerError)
			return
		}
	}
	defer stmt.Close()

	if _, err = stmt.Exec(); err != nil {
		resp.ErrorJSON(w, fmt.Errorf("error execute users' secrets transaction: %v", err), http.StatusInternalServerError)
		return
	}

	if err = txn.Commit(); err != nil {
		resp.ErrorJSON(w, fmt.Errorf("error commiting users' secrets transaction: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
