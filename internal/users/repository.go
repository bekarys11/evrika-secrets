package users

import (
	"fmt"
	"github.com/bekarys11/evrika-secrets/internal/roles"
	"github.com/go-ldap/ldap"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
)

type Repository struct {
	DB         *sqlx.DB
	LDAP       *ldap.Conn
	Validation *validator.Validate
	Logger     *slog.Logger
}

func NewRepository(db *sqlx.DB, ldap *ldap.Conn, validation *validator.Validate, logger *slog.Logger) *Repository {
	return &Repository{
		DB:         db,
		LDAP:       ldap,
		Validation: validation,
		Logger:     logger,
	}
}

func (repo *Repository) GetUsers() (users []*UserResp, err error) {
	rows, err := repo.DB.Queryx(`
	SELECT * FROM users
	JOIN roles as r ON role_id = r.id
	LIMIT 10`)

	if err != nil {
		repo.Logger.Error("get users query failed", err)
		return nil, err
	}

	for rows.Next() {
		var (
			user UserResp
			role roles.Role
		)
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.IsActive, &user.RoleId, &user.CreatedAt, &user.UpdatedAt, &role.ID, &role.Name, &role.Alias, &role.CreatedAt, &role.UpdatedAt); err != nil {
			return nil, err
		}

		if role.ID != 0 {
			user.Role = &role
		}
		users = append(users, &user)
		repo.Logger.Debug("users list:", users)
	}
	defer rows.Close()

	return users, nil
}

func (repo *Repository) CreateUser(payload User) error {
	if err := repo.Validation.Struct(payload); err != nil {
		repo.Logger.Error("user body validation failed", err)
		return fmt.Errorf("validation error: %v", err)
	}
	repo.Logger.Debug("user is valid")

	if _, err := repo.activeDirSearch(payload.Email); err != nil {
		repo.Logger.Error("search in active directory failed", err)
		return err
	}
	repo.Logger.Debug("user is found in active directory")

	hashed, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		repo.Logger.Error("password hashing failed", err)
		return fmt.Errorf("password hashing error: %v", err.Error())
	}

	if _, err = repo.DB.Exec(`INSERT INTO users (name, email, password, is_active, role_id ) VALUES ($1, $2, $3, $4, $5)`, payload.Name, payload.Email, string(hashed), payload.IsActive, payload.RoleId); err != nil {
		repo.Logger.Error("sql insert failed", err)
		return fmt.Errorf("sql insert error: %v", err)
	}

	return nil
}

func (repo *Repository) GetProfile(userId string) (UserResp, error) {
	var (
		user UserResp
		role roles.Role
	)
	repo.Logger.Debug("user id", userId)

	if err := repo.DB.QueryRowx(`SELECT * FROM users as u
    						JOIN roles as r ON role_id = r.id 
         					WHERE u.id = $1`, userId).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.IsActive, &user.RoleId, &user.CreatedAt, &user.UpdatedAt, &role.ID, &role.Name, &role.Alias, &role.CreatedAt, &role.UpdatedAt); err != nil {
		repo.Logger.Error("profile query failed", err)
		return UserResp{}, fmt.Errorf("db scan error: %v", err)
	}

	user.Role = &role

	repo.Logger.Debug("user profile", user)

	return user, nil
}
