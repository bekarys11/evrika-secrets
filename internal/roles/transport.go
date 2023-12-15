package roles

import (
	resp "github.com/bekarys11/evrika-secrets/pkg/response"
	"net/http"
)

type Service interface {
	GetRoles() ([]*Role, error)
}

type HttpServer struct {
	service Service
}

func NewHttpServer(service Service) HttpServer {
	return HttpServer{
		service: service,
	}
}

//	 @Summary      Список ролей
//		@Security ApiKeyAuth
//	 @Description  Получить список всех ролей
//	 @Tags         roles
//	 @Accept       json
//	 @Produce      json
//	 @Success      200  {object}   RoleSwaggerJson
//	 @Failure      400  {object}  resp.Err
//	 @Failure      500  {object}  resp.Err
//	 @Router       /api/v1/roles [get]
func (h HttpServer) GetRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := h.service.GetRoles()

	if err != nil {
		resp.ErrorJSON(w, err)
		return
	}

	resp.WriteApiJSON(w, http.StatusOK, roles)
}
