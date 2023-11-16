package users

import (
	"encoding/json"
	"errors"
	"fmt"
	resp "github.com/bekarys11/evrika-secrets/pkg/response"
	"github.com/go-ldap/ldap"
	"github.com/golang-jwt/jwt"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"strings"
)

type Repo struct {
	DB   *sqlx.DB
	LDAP *ldap.Conn
}

func (u *Repo) All(w http.ResponseWriter, r *http.Request) {
	var users []*User

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
		users = append(users, &user)
	}

	resp.WriteApiJSON(w, http.StatusOK, users)
}

func (u *Repo) Create(w http.ResponseWriter, r *http.Request) {
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

func (u *Repo) GetProfile(w http.ResponseWriter, r *http.Request) {
	var user User
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
	log.Println(userId)
	if err = u.DB.QueryRowx("SELECT * FROM users WHERE id = $1", userId).StructScan(&user); err != nil {
		resp.ErrorJSON(w, fmt.Errorf("db scan error: %v", err), http.StatusInternalServerError)
		return
	}

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
