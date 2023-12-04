package secrets

import (
	"encoding/json"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/bekarys11/evrika-secrets/pkg/common"
	resp "github.com/bekarys11/evrika-secrets/pkg/response"
	"github.com/gorilla/mux"
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
// @Param 		type query string false "список секретов по типу"
// @Param 		user query string false "список секретов от пользователя"
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

//	 @Summary      Объект ключа
//		@Security ApiKeyAuth
//	 @Description  Получить ключ по id
//	 @Tags         secrets
//	 @Accept       json
//	 @Produce      json
//
// @Param secret_id  path int true "ID ключа"
//
// @Success      200  {object} SecretSwaggerJsonObj
// @Failure      400  {object}  resp.Err
// @Failure      500  {object}  resp.Err
// @Router       /api/v1/secrets/:secret_id [get]
func (s *Repo) One(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	secretId := vars["secret_id"]

	userId, err := common.GetUserIdFromToken(r)
	if err != nil {
		resp.ErrorJSON(w, err)
		return
	}

	role, err := common.GetRoleFromToken(r)
	if err != nil {
		resp.ErrorJSON(w, err)
		return
	}

	secret, err := s.getById(secretId, role, userId)

	if err != nil {
		resp.ErrorJSON(w, err)
		return
	}

	resp.WriteApiJSON(w, http.StatusOK, &secret)
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
	var secret Secret

	if err := s.createSecret(r, &secret); err != nil {
		resp.ErrorJSON(w, err)
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
