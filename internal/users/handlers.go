package users

import (
	"encoding/json"
	"fmt"
	resp "github.com/bekarys11/evrika-secrets/pkg/response"
	"github.com/go-ldap/ldap"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type Repo struct {
	DB   *sqlx.DB
	LDAP *ldap.Conn
}

func (u *Repo) All(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var users []User

	rows, err := u.DB.Queryx("SELECT * FROM users LIMIT 10")

	if err != nil {
		resp.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		var user User
		if err := rows.StructScan(&user); err != nil {
			resp.ErrorJSON(w, err, http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	resp.WriteJSON(w, http.StatusOK, resp.New{Data: users})
}

func (u *Repo) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		resp.ErrorJSON(w, fmt.Errorf("invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	_, err := u.ActiveDirSearch(user.Email)
	if err != nil {
		resp.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		resp.ErrorJSON(w, fmt.Errorf("password hashing error: %v", err.Error()), http.StatusInternalServerError)
	}

	if _, err = u.DB.Exec(`INSERT INTO users (name, email, password, is_active, role_id ) VALUES ($1, $2, $3, $4, $5)`, user.Name, user.Email, string(hashed), user.IsActive, user.RoleId); err != nil {
		resp.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
