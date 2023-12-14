package users

import (
	"errors"
	"fmt"
	"github.com/bekarys11/evrika-secrets/internal/roles"
	"github.com/bekarys11/evrika-secrets/pkg/common"
	resp "github.com/bekarys11/evrika-secrets/pkg/response"
	"github.com/go-ldap/ldap"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type Re struct {
	DB         *sqlx.DB
	LDAP       *ldap.Conn
	Validation *validator.Validate
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
func (u *Re) All(w http.ResponseWriter, r *http.Request) {
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
func (u *Re) Create(w http.ResponseWriter, r *http.Request) {
	var user UserRequest

	if err := resp.ReadJSON(w, r, &user); err != nil {
		resp.ErrorJSON(w, fmt.Errorf("invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	if err := u.Validation.Struct(user); err != nil {
		resp.ErrorJSON(w, fmt.Errorf("validation error: %v", err), http.StatusBadRequest)
		return
	}

	//if _, err := u.ActiveDirSearch(user.Email); err != nil {
	//	resp.ErrorJSON(w, err, http.StatusBadRequest)
	//	return
	//}

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
func (u *Re) GetProfile(w http.ResponseWriter, r *http.Request) {
	var (
		user UserResp
		role roles.Role
	)
	claims, err := common.GetTokenClaims(r)
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
