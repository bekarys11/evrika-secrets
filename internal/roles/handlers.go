package roles

import (
	resp "github.com/bekarys11/evrika-secrets/pkg/response"
	"github.com/jmoiron/sqlx"
	"net/http"
)

type Repo struct {
	DB *sqlx.DB
}

func (rp *Repo) All(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var roles []Role

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
		roles = append(roles, role)
	}
	resp.WriteJSON(w, http.StatusOK, resp.New{Data: roles})
}
