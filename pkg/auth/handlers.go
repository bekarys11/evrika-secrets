package auth

import (
	"errors"
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
func (a *Repo) Login(w http.ResponseWriter, r *http.Request) {
	var loginPayload LoginPayload

	//TODO: fix: if add error handling, returns error: "json illegal base64 data at input byte 4"
	if err := resp.ReadJSON(w, r, &loginPayload); err != nil {
		resp.ErrorJSON(w, fmt.Errorf("invalid JSON: %v", err), http.StatusBadRequest)
		return
	}
	email := fmt.Sprintf("%s@evrika.com", loginPayload.Login)

	row := a.DB.QueryRowx("SELECT * FROM users WHERE email = $1", email)

	if row.Err() != nil {
		resp.ErrorJSON(w, row.Err(), http.StatusInternalServerError)
		return
	}

	var userEntity users.User
	if err := row.StructScan(&userEntity); err != nil {
		resp.ErrorJSON(w, fmt.Errorf("database scan error: %v", err), http.StatusInternalServerError)
		return
	}

	isValidPassword := checkPassword(userEntity.Password, loginPayload.Password)

	if isValidPassword {
		token, err := generateToken(userEntity)
		if err != nil {
			resp.ErrorJSON(w, fmt.Errorf("failed to generate token: %v", err), http.StatusInternalServerError)
			return
		}

		res := LoginResponse{AccessToken: token}

		resp.WriteJSON(w, 200, res)
	} else {
		resp.ErrorJSON(w, fmt.Errorf("invalid password"), http.StatusBadRequest)
	}

}

func generateSecretKey() ([]byte, error) {
	return []byte(os.Getenv("JWT_SECRET_KEY")), nil
}

func checkPassword(hashedPassword, providedPassword string) bool {
	fmt.Printf("checkPassword func: hashedPassword=%s \n provided password=%s \n", hashedPassword, providedPassword)
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(providedPassword)); err != nil {
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
	if tokenStr == "" {
		return false, errors.New("token is not provided")
	}

	_, bareTokenStr, _ := strings.Cut(tokenStr, "Bearer ")

	secretKey, err := generateSecretKey()

	if err != nil {

		return false, fmt.Errorf("error while generating secret key: %v", err.Error())
	}

	token, err := jwt.Parse(bareTokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return true, nil
	} else {

		return false, fmt.Errorf("token is invalid: %v", err)
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
