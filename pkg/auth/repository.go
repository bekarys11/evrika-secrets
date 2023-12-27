package auth

import (
	"errors"
	"fmt"
	"github.com/bekarys11/evrika-secrets/internal/users"
	"github.com/golang-jwt/jwt"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"os"
	"strings"
	"time"
)

type Repository struct {
	DB *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

func (repo *Repository) Login(payload LoginPayload) (token string, err error) {
	email := fmt.Sprintf("%s@evrika.com", payload.Login)

	row := repo.DB.QueryRowx("SELECT * FROM users WHERE email = $1", email)

	if row.Err() != nil {
		return "", row.Err()
	}

	var userEntity users.User
	if err := row.StructScan(&userEntity); err != nil {
		return "", fmt.Errorf("database scan error: %v", err)
	}

	isValidPassword := checkPassword(userEntity.Password, payload.Password)

	if isValidPassword {
		token, err := GenerateToken(userEntity)
		if err != nil {
			return "", fmt.Errorf("failed to generate token: %v", err)
		}

		return token, nil
	} else {
		return "", fmt.Errorf("invalid password")
	}
}

func generateSecretKey() ([]byte, error) {
	return []byte(os.Getenv("JWT_SECRET_KEY")), nil
}

func checkPassword(hashedPassword, providedPassword string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(providedPassword)); err != nil {
		slog.Error(err.Error())
		return false
	}

	return true
}

func GenerateToken(user users.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":     time.Now().Add(48 * time.Hour).Unix(),
		"user_id": user.ID,
		"name":    user.Name,
		"role_id": user.RoleId,
		"role":    getRoleName(user.RoleId),
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

func getRoleName(roleId int) string {
	if roleId == 1 {
		return "admin"
	} else {
		return "user"
	}
}
