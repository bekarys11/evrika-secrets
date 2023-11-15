package auth

import (
	"encoding/json"
	"fmt"
	"github.com/bekarys11/evrika-secrets/internal/users"
	resp "github.com/bekarys11/evrika-secrets/pkg/response"
	"github.com/golang-jwt/jwt"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

type Repo struct {
	DB *sqlx.DB
}

func (a *Repo) Login(w http.ResponseWriter, r *http.Request) {
	var loginPayload LoginPayload

	//TODO: fix: if add error handling, returns error: "json illegal base64 data at input byte 4"
	json.NewDecoder(r.Body).Decode(&loginPayload)

	row := a.DB.QueryRowx("SELECT * FROM users WHERE email = $1", loginPayload.Email)

	if row.Err() != nil {
		resp.ErrorJSON(w, row.Err(), http.StatusInternalServerError)
		return
	}

	var userEntity users.User
	if err := row.StructScan(&userEntity); err != nil {
		resp.ErrorJSON(w, fmt.Errorf("database scan error: %v", err), http.StatusInternalServerError)
		return
	}

	//TODO: fix check password
	isValidPassword := checkPassword(userEntity.Password, loginPayload.Password)

	if isValidPassword {
		token, err := generateToken(userEntity)
		if err != nil {
			resp.ErrorJSON(w, fmt.Errorf("failed to generate token: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		res := LoginResponse{AccessToken: token}

		json.NewEncoder(w).Encode(res)
	} else {
		resp.ErrorJSON(w, fmt.Errorf("invalid password"), http.StatusBadRequest)
	}

}

func generateSecretKey() ([]byte, error) {
	return []byte(os.Getenv("JWT_SECRET_KEY")), nil
}

func checkPassword(hashedPassword string, providedPassword []byte) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), providedPassword); err != nil {
		slog.Error(err.Error())
		return false
	}

	log.Println("valid password")
	return true
}

func generateToken(user users.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":     time.Now().Add(48 * time.Hour).Unix(),
		"user_id": user.ID,
		"name":    user.Name,
	})

	secretKey, err := generateSecretKey()

	if err != nil {
		return fmt.Sprintf("Error while generating secret key: %v", err.Error()), err
	}

	tokenStr, err := token.SignedString(secretKey)

	if err != nil {
		return fmt.Sprintf("Error while converting token to string: %v", err.Error()), err
	}

	return tokenStr, nil
}

func IsValidToken(tokenStr string) (bool, error) {
	_, bareTokenStr, _ := strings.Cut(tokenStr, "Bearer ")

	secretKey, err := generateSecretKey()

	if err != nil {

		return false, fmt.Errorf("error while generating secret key: %v", err.Error())
	}

	token, err := jwt.Parse(bareTokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return secretKey, nil
	})

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return true, nil
	} else {
		fmt.Println(err)
		return false, err
	}

}

func GetTokenClaims(tokenStr string) jwt.MapClaims {
	_, bareTokenStr, _ := strings.Cut(tokenStr, "Bearer ")

	secretKey, err := generateSecretKey()

	if err != nil {
		fmt.Printf("Error while generating secret key: %v", err.Error())
	}

	token, err := jwt.Parse(bareTokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return secretKey, nil
	})

	return token.Claims.(jwt.MapClaims)
}
