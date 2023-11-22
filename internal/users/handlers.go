package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bekarys11/evrika-secrets/internal/roles"
	resp "github.com/bekarys11/evrika-secrets/pkg/response"
	"github.com/go-ldap/ldap"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"strings"
)

type Repo struct {
	DB         *sqlx.DB
	LDAP       *ldap.Conn
	Validation *validator.Validate
}

func (u *Repo) All(w http.ResponseWriter, r *http.Request) {
	var users []*UserResp

	rows, err := u.DB.Queryx(`SELECT * FROM users
         							JOIN roles as r ON role_id = r.id
         							LIMIT 10`)
	if err != nil {
		resp.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		var (
			user UserResp
			role roles.Role
		)
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.IsActive, &user.RoleId, &user.CreatedAt, &user.UpdatedAt, &role.ID, &role.Name, &role.Alias, &role.CreatedAt, &role.UpdatedAt); err != nil {
			resp.ErrorJSON(w, err, http.StatusInternalServerError)
			return
		}

		if role.ID != 0 {
			user.Role = &role
		}

		users = append(users, &user)

	}
	defer rows.Close()

	resp.WriteApiJSON(w, http.StatusOK, users)
}

func (u *Repo) Create(w http.ResponseWriter, r *http.Request) {
	var user UserRequest

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		resp.ErrorJSON(w, fmt.Errorf("invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	if err := u.Validation.Struct(user); err != nil {
		resp.ErrorJSON(w, fmt.Errorf("validation error: %v", err), http.StatusBadRequest)
		return
	}

	if _, err := u.ActiveDirSearch(user.Email); err != nil {
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

	resp.WriteApiJSON(w, 201, "Пользователь создан")
}

func (u *Repo) GetProfile(w http.ResponseWriter, r *http.Request) {
	var (
		user UserResp
		role roles.Role
	)
	claims, err := getTokenClaims(r)
	if err != nil {
		resp.ErrorJSON(w, fmt.Errorf("get profile error: %v", err), http.StatusInternalServerError)
		return
	}

	userId, ok := claims["user_id"]
	if !ok {
		resp.ErrorJSON(w, errors.New("there is no user id in token claims"), http.StatusInternalServerError)
		return
	}

	if err = u.DB.QueryRowx(`SELECT * FROM users as u
    						JOIN roles as r ON role_id = r.id 
         					WHERE u.id = $1`, userId).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.IsActive, &user.RoleId, &user.CreatedAt, &user.UpdatedAt, &role.ID, &role.Name, &role.Alias, &role.CreatedAt, &role.UpdatedAt); err != nil {
		resp.ErrorJSON(w, fmt.Errorf("db scan error: %v", err), http.StatusInternalServerError)
		return
	}

	user.Role = &role

	resp.WriteApiJSON(w, http.StatusOK, &user)
}

func getTokenClaims(r *http.Request) (jwt.MapClaims, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))

	_, tokenStr, _ := strings.Cut(r.Header.Get("Authorization"), "Bearer ")
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	} else if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	} else {
		return nil, errors.New("unknown claims type, cannot proceed")
	}
}
