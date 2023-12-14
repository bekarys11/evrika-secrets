package users

import (
	"fmt"
	"github.com/bekarys11/evrika-secrets/pkg/common"
	resp "github.com/bekarys11/evrika-secrets/pkg/response"
	"net/http"
)

type Service interface {
	GetUsers() ([]*UserResp, error)
	CreateUser(payload User) error
	GetProfile(userId string) (UserResp, error)
}

type HttpServer struct {
	service Service
}

func NewHttpServer(service Service) HttpServer {
	return HttpServer{
		service: service,
	}
}

//	 @Summary      Список пользователей
//		@Security ApiKeyAuth
//	 @Description  Получить список всех пользователей
//	 @Tags         users
//	 @Accept       json
//	 @Produce      json
//	 @Success      200  {object}   UserSwaggerJson
//	 @Failure      400  {object}  resp.Err
//	 @Failure      500  {object}  resp.Err
//	 @Router       /api/v1/users [get]
func (h HttpServer) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetUsers()

	if err != nil {
		resp.ErrorJSON(w, err)
	}
	resp.WriteApiJSON(w, 200, users)
}

//	 @Summary      Создать пользователя
//		@Security ApiKeyAuth
//	 @Description  Создать пользователя
//	 @Tags         users
//
// @Param input body UserSwaggerRequest true "добавить данные в тело запроса"
//
//	@Accept       json
//	@Produce      json
//	@Success      201  {string}   "Пользователь создан"
//	@Failure      400  {object}  resp.Err
//	@Failure      500  {object}  resp.Err
//	@Router       /api/v1/users [post]
func (h HttpServer) CreateUser(w http.ResponseWriter, r *http.Request) {
	var payload User
	if err := resp.ReadJSON(w, r, &payload); err != nil {
		resp.ErrorJSON(w, err, 400)
		return
	}

	if err := h.service.CreateUser(payload); err != nil {
		resp.ErrorJSON(w, err, 500)
		return
	}

	resp.RespondOK(map[string]string{"message": "Пользователь создан"}, w)
}

//	 @Summary      Инфо о профиле
//		@Security ApiKeyAuth
//	 @Description  Получить инфо о пользователе
//	 @Tags         users
//
// @Accept       json
// @Produce      json
// @Success      200  {object}   UserSwaggerJsonMap
// @Failure      400  {object}  resp.Err
// @Failure      500  {object}  resp.Err
// @Router       /api/v1/profile [get]
func (h HttpServer) GetProfile(w http.ResponseWriter, r *http.Request) {
	userId, err := common.GetUserIdFromToken(r)
	if err != nil {
		resp.ErrorJSON(w, fmt.Errorf("get user id token claims error: %v", err))
		return
	}

	user, err := h.service.GetProfile(userId)
	if err != nil {
		resp.ErrorJSON(w, err)
		return
	}

	resp.WriteApiJSON(w, 200, &user)
}
