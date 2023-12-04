package secrets

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Masterminds/squirrel"
	resp "github.com/bekarys11/evrika-secrets/pkg/response"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"net/http"
)

type Repo struct {
	DB       *sqlx.DB
	QBuilder squirrel.StatementBuilderType
}

//	 @Summary      Список секретов/ключей
//		@Security ApiKeyAuth
//	 @Description  Получить список всех секретов
//	 @Tags         secrets
//	 @Accept       json
//	 @Produce      json
//
// @Success      200  {object} SecretSwaggerJson
// @Failure      400  {object}  resp.Err
// @Failure      500  {object}  resp.Err
// @Router       /api/v1/secrets [get]
func (s *Repo) All(w http.ResponseWriter, r *http.Request) {
	secrets, err := s.getSecrets(r)

	if err != nil {
		resp.ErrorJSON(w, err)
		return
	}

	resp.WriteApiJSON(w, http.StatusOK, secrets)
}

//	 @Summary      Создать секрет/ключ
//		@Security ApiKeyAuth
//	 @Description  Создать секрет/ключ
//	 @Tags         secrets
//
// @Param input body SecretSwaggerReq true "добавить данные в тело запроса"
//
//	@Accept       json
//	@Produce      json
//	@Success      201  {string}   "Секрет сохранен"
//	@Failure      400  {object}  resp.Err
//	@Failure      500  {object}  resp.Err
//	@Router       /api/v1/secrets [post]
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
	err = tx.QueryRowxContext(ctx, `INSERT INTO secrets (title, key, data, stype, author_id) VALUES ($1, $2, $3, $4, $5) RETURNING id`, secret.Title, secret.Key, secret.Data, secret.Type, secret.AuthorId).Scan(&secretId)
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

	resp.WriteJSON(w, 201, "Секрет сохранен")
}

//	 @Summary      Поделиться секретом с пользователями
//		@Security ApiKeyAuth
//	 @Description  Поделиться секретом с пользователями
//	 @Tags         secrets
//
// @Param input body UsersSecret true "добавить данные в тело запроса"
//
//	@Accept       json
//	@Produce      json
//	@Success      201  {string}   "Секрет сохранен"
//	@Failure      400  {object}  resp.Err
//	@Failure      500  {object}  resp.Err
//	@Router       /api/v1/secrets/share [post]
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

	resp.WriteJSON(w, 201, "Секрет сохранен")
}
