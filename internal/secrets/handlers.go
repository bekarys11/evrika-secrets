package secrets

import (
	"encoding/json"
	"fmt"
	"github.com/bekarys11/evrika-secrets/pkg/common"
	resp "github.com/bekarys11/evrika-secrets/pkg/response"
	"github.com/gorilla/mux"
	"net/http"
)

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
	qParams := r.URL.Query()

	userId, err := common.GetUserIdFromToken(r)
	if err != nil {
		resp.ErrorJSON(w, fmt.Errorf("token claims error: %v", err))
		return
	}

	userRole, err := common.GetRoleFromToken(r)
	if err != nil {
		resp.ErrorJSON(w, fmt.Errorf("token claims error: %v", err))
		return
	}

	secrets, err := s.getAllSecrets(qParams, userRole, userId)

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
	if err := json.NewDecoder(r.Body).Decode(&secret); err != nil {
		resp.ErrorJSON(w, fmt.Errorf("invalid JSON: %v", err))
		return
	}

	if err := s.createSecret(&secret); err != nil {
		resp.ErrorJSON(w, err, 500)
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

	if err := s.shareSecret(usersSecrets); err != nil {
		resp.ErrorJSON(w, err)
		return
	}

	resp.WriteJSON(w, 201, "Секрет сохранен")
}

//	 @Summary      Редактировать ключ
//		@Security ApiKeyAuth
//	 @Description  Администратор может изменять все ключи, а пользователь только свои.
//	 @Tags         secrets
//
// @Param input body SecretReq true "добавить данные в тело запроса"
//
//	@Accept       json
//	@Produce      json
//	@Success      200  {string}   "Секрет изменен"
//	@Failure      400  {object}  resp.Err
//	@Failure      500  {object}  resp.Err
//	@Router       /api/v1/secrets/:secret_id [put]
func (s *Repo) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	secretId := vars["secret_id"]
	var payload SecretReq

	if err := resp.ReadJSON(w, r, &payload); err != nil {
		resp.ErrorJSON(w, fmt.Errorf("invalid json: %v", err))
		return
	}

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

	if err := s.updateById(secretId, role, userId, payload); err != nil {
		resp.ErrorJSON(w, err, 500)
		return
	}

	resp.WriteJSON(w, 200, "Секрет изменен")
}

//	 @Summary      Удалить ключ
//		@Security ApiKeyAuth
//	 @Description  Администратор может удалять любой ключ, а пользователь только свои.
//	 @Tags         secrets
//
// @Accept       json
// @Produce      json
// @Success      200  {string}   "Секрет удален"
// @Failure      400  {object}  resp.Err
// @Failure      500  {object}  resp.Err
// @Router       /api/v1/secrets/:secret_id [delete]
func (s *Repo) Delete(w http.ResponseWriter, r *http.Request) {
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

	if err := s.deleteById(secretId, role, userId); err != nil {
		resp.ErrorJSON(w, err, 500)
		return
	}
	resp.WriteJSON(w, 200, "Секрет удален")
}
