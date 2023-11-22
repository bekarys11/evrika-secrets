package roles

import (
	resp "github.com/bekarys11/evrika-secrets/pkg/response"
	"github.com/jmoiron/sqlx"
	"net/http"
)

type Repo struct {
	DB *sqlx.DB
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
func (rp *Repo) All(w http.ResponseWriter, r *http.Request) {
	var roles []*Role

	rows, err := rp.DB.Queryx("SELECT * FROM roles LIMIT 10")
	if err != nil {
		resp.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		var role Role
		if err := rows.StructScan(&role); err != nil {
			resp.ErrorJSON(w, err, http.StatusInternalServerError)
			return
		}
		roles = append(roles, &role)
	}

	resp.WriteApiJSON(w, http.StatusOK, roles)
}
