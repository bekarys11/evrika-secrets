package auth

import (
	"fmt"
	resp "github.com/bekarys11/evrika-secrets/pkg/response"
	"log/slog"
	"net/http"
)

type Service interface {
	Login(LoginPayload) (token string, err error)
}

type HttpServer struct {
	service Service
}

func NewHttpServer(service Service) HttpServer {
	return HttpServer{
		service: service,
	}
}

//	@Summary      Вход пользователя
//	@Description  Логин пользователя
//	@Tags         auth
//
// @Param input body LoginPayload true "добавить данные в тело запроса"
//
//	@Accept       json
//	@Produce      json
//	@Success      200  {object}   LoginResponse
//	@Failure      400  {object}  resp.Err
//	@Failure      500  {object}  resp.Err
//	@Router       /api/v1/login [post]
func (h HttpServer) Login(w http.ResponseWriter, r *http.Request) {
	slog.Info("Hello from Login")
	var payload LoginPayload

	//TODO: fix: if add error handling, returns error: "json illegal base64 data at input byte 4"
	if err := resp.ReadJSON(w, r, &payload); err != nil {
		resp.ErrorJSON(w, fmt.Errorf("invalid JSON: %v", err))
		return
	}

	token, err := h.service.Login(payload)

	if err != nil {
		resp.ErrorJSON(w, err)
		return
	}

	resp.WriteJSON(w, 200, LoginResponse{AccessToken: token})
}
