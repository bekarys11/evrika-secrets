package secrets

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bekarys11/evrika-secrets/pkg/common"
	resp "github.com/bekarys11/evrika-secrets/pkg/response"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
)

type Service interface {
	GetSecrets(qParams url.Values, userRole, userId string) ([]*SecretResp, error)
	GetSecretById(secretId string, role string, userId string) (secret SecretResp, err error)
	CreateSecret(ctx context.Context, payload Secret) error
	UpdateSecret(secretId, userRole, userId string, payload Secret) error
}

type HttpServer struct {
	service Service
}

func NewHttpServer(service Service) HttpServer {
	return HttpServer{
		service: service,
	}
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
func (h HttpServer) GetSecrets(w http.ResponseWriter, r *http.Request) {
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

	secrets, err := h.service.GetSecrets(qParams, userRole, userId)
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
func (h HttpServer) GetSecretById(w http.ResponseWriter, r *http.Request) {
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

	secret, err := h.service.GetSecretById(secretId, role, userId)
	if err != nil {
		resp.ErrorJSON(w, err)
		return
	}

	resp.WriteApiJSON(w, http.StatusOK, secret)
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
func (h HttpServer) CreateSecret(w http.ResponseWriter, r *http.Request) {
	var (
		payload Secret
	)
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		resp.ErrorJSON(w, fmt.Errorf("invalid JSON: %v", err))
		return
	}

	err := h.service.CreateSecret(context.Background(), payload)

	if err != nil {
		resp.ErrorJSON(w, err, 500)
		return
	}

	resp.WriteJSON(w, 201, map[string]string{"message": "Секрет сохранен"})
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
func (h HttpServer) ShareSecret(w http.ResponseWriter, r *http.Request) {}

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
func (h HttpServer) UpdateSecret(w http.ResponseWriter, r *http.Request) {
	var (
		payload Secret
	)
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

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		resp.ErrorJSON(w, fmt.Errorf("invalid JSON: %v", err))
		return
	}

	if err := h.service.UpdateSecret(secretId, role, userId, payload); err != nil {
		resp.ErrorJSON(w, err, 500)
		return
	}

	resp.WriteJSON(w, http.StatusOK, map[string]string{"message": "Секрет изменен"})
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
func (h HttpServer) Delete(w http.ResponseWriter, r *http.Request) {}
