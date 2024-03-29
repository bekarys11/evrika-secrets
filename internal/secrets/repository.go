package secrets

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"log"
	"log/slog"
	"net/url"
	"strconv"
	"time"
)

type Repository struct {
	DB       *sqlx.DB
	QBuilder squirrel.StatementBuilderType
	Logger   *slog.Logger
}

func NewRepository(db *sqlx.DB, sq squirrel.StatementBuilderType, logger *slog.Logger) *Repository {
	return &Repository{
		DB:       db,
		QBuilder: sq,
		Logger:   logger,
	}
}

func (repo *Repository) GetSecrets(qParams url.Values, userRole string, userId string) (secrets []*SecretResp, err error) {
	var (
		q              string
		args           []interface{}
		secretType     = qParams.Get("type")
		userIDQuery, _ = strconv.Atoi(qParams.Get("user"))
	)

	gQuery := repo.QBuilder.
		Select("secrets.id", "secrets.title", "secrets.key", "secrets.data", "secrets.stype", "secrets.author_id").
		Column(`
					jsonb_agg(
					  jsonb_build_object(
							  'id', users.id,
							  'name', users.name
						  )
					  ) AS users_info
					`).
		From("secrets").
		Join("users_secrets ON users_secrets.secret_id = secrets.id").
		Join("users ON users_secrets.user_id = users.id").
		GroupBy("secrets.id", "secrets.title", "secrets.key", "secrets.data", "secrets.stype", "secrets.author_id")

	switch userRole {
	case "user":
		userQ := repo.QBuilder.
			Select("ss.id", "ss.title", "ss.key", "ss.data", "ss.stype", "ss.author_id", "ss.users_info").
			FromSelect(gQuery, "ss").
			Where(fmt.Sprintf("jsonb_path_exists(users_info, '$[*] ?? (@.id == %s)')", userId))

		if hasType := qParams.Has("type"); hasType {
			userQ = userQ.Where("ss.stype = ?", secretType)
		}

		if hasUserId := qParams.Has("user"); hasUserId {
			return nil, errors.New("вы не имеете достаточно прав")
		}

		q, args, err = userQ.ToSql()
		if err != nil {
			log.Printf("Error to_sql: %v", err)
		}
	case "admin":
		adminQ := repo.QBuilder.
			Select("ss.id", "ss.title", "ss.key", "ss.data", "ss.stype", "ss.author_id", "ss.users_info").
			FromSelect(gQuery, "ss").
			Where(fmt.Sprintf("jsonb_path_exists(users_info, '$[*] ?? (@.id == %s)')", userId))

		if hasType := qParams.Has("type"); hasType {
			adminQ = adminQ.Where("ss.stype = ?", secretType)
		}

		if hasUserId := qParams.Has("user"); hasUserId {
			repo.Logger.Debug("has user id query param", userIDQuery)
			adminQ = adminQ.Where(fmt.Sprintf("jsonb_path_exists(users_info, '$[*] ?? (@.id == %d)')", userIDQuery))
		}

		q, args, err = adminQ.ToSql()
		if err != nil {
			log.Printf("Error to_sql: %v", err)
		}
	}
	repo.Logger.Debug("query string", q)
	rows, err := repo.DB.Queryx(q, args...)
	if err != nil {
		repo.Logger.Debug("query failed", err)
		return nil, fmt.Errorf("secrets query error: %v", err)
	}
	for rows.Next() {
		log.Println("Inside rows")

		var secret SecretResp
		if err := rows.Scan(
			&secret.ID, &secret.Title, &secret.Key, &secret.Data, &secret.Type, &secret.AuthorId, &secret.Users,
		); err != nil {
			repo.Logger.Error("secret scan failed", secret)
			return nil, fmt.Errorf("scan error: %v", err)
		}

		secrets = append(secrets, &secret)

		// Unmarshal the bytes into a User struct
		var user interface{}
		err = json.Unmarshal(secret.Users.([]byte), &user)
		if err != nil {
			repo.Logger.Error("Error unmarshaling JSON", err)
			fmt.Println("Error unmarshaling JSON:", err)
		}

		secret.Users = user
	}
	repo.Logger.Debug("secrets list", secrets)

	defer rows.Close()
	return secrets, nil
}

func (repo *Repository) GetSecretById(secretId string, role string, userId string) (secret SecretResp, err error) {
	userID, _ := strconv.Atoi(userId)

	query := repo.QBuilder.Select("secrets.id, secrets.title, secrets.key, secrets.data, secrets.stype, secrets.author_id, secrets.created_at, secrets.updated_at").From("secrets").Where("secrets.id = ?", secretId)

	if err := query.RunWith(repo.DB).QueryRow().Scan(&secret.ID, &secret.Title, &secret.Key, &secret.Data, &secret.Type, &secret.AuthorId, &secret.CreatedAt, &secret.UpdatedAt); err != nil {
		switch err {
		case sql.ErrNoRows:
			repo.Logger.Error("secret not found", err)
			return SecretResp{}, errors.New("секрет не найден")
		default:
			repo.Logger.Error("secret query failed", err)
			return SecretResp{}, fmt.Errorf("sql query failed: %v", err)
		}
	}

	repo.Logger.Debug("author_id: user role: secret_id ", secret.AuthorId, role, secretId)
	if secret.AuthorId != userID && role != "admin" {
		repo.Logger.Info("you dont have permission")
		return SecretResp{}, errors.New("вы не имеете достаточно прав")
	}

	return secret, nil
}

func (repo *Repository) CreateSecret(ctx context.Context, payload Secret) error {
	var secretId int

	tx, err := repo.DB.BeginTxx(ctx, nil)
	if err != nil {
		repo.Logger.Error("unable to begin transaction", err)
		return fmt.Errorf("unable to begin transaction: %v", err)
	}
	defer tx.Rollback()

	err = tx.QueryRowxContext(ctx, `INSERT INTO secrets (title, key, data, stype, author_id) VALUES ($1, $2, $3, $4, $5) RETURNING id`, payload.Title, payload.Key, payload.Data, payload.Type, payload.AuthorId).Scan(&secretId)
	switch {
	case err == sql.ErrNoRows:
		repo.Logger.Error("no secret with id", err)
		return fmt.Errorf("no secret with id: %v", err)
	case err != nil:
		repo.Logger.Error("error insert secrets to db", err)
		return fmt.Errorf("error insert secrets to db: %v", err)
	}

	if _, err = tx.ExecContext(ctx, `INSERT INTO users_secrets (user_id, secret_id) VALUES ($1, $2)`, payload.AuthorId, secretId); err != nil {
		repo.Logger.Error("error insert users' secrets to db", err)
		return fmt.Errorf("error insert users' secrets to db: %v", err)
	}

	if err = tx.Commit(); err != nil {
		repo.Logger.Error("error committing db transaction", err)
		return fmt.Errorf("error committing db transaction: %v", err)
	}

	return nil
}

func (repo *Repository) UpdateSecret(secretId, userRole, userId string, payload Secret) error {

	authorId, err := repo.getSecretAuthor(secretId, userRole, userId)
	if err != nil {
		repo.Logger.Error("failed to get secret author", err)
		return err
	}

	if userRole == "user" {
		if err := repo.checkSecretAuthor(authorId, userId); err != nil {
			repo.Logger.Error("failed to get secret author", err)
			return err
		}
	}

	query := repo.QBuilder.
		Update("secrets").
		SetMap(map[string]interface{}{
			"title":      payload.Title,
			"key":        payload.Key,
			"data":       payload.Data,
			"stype":      payload.Type,
			"author_id":  authorId,
			"updated_at": time.Now(),
		}).
		Where("id = ?", secretId)

	_, err = query.RunWith(repo.DB).Exec()
	if err != nil {
		repo.Logger.Error("failed to query", err)
		return fmt.Errorf("error executing update query: %v", err)
	}

	return nil
}

func (repo *Repository) DeleteSecret(secretId, userRole, userId string) error {
	if userRole == "user" {
		authorId, err := repo.getSecretAuthor(secretId, userRole, userId)
		if err != nil {
			repo.Logger.Error("get secret author failed", err)
			return err
		}

		if err := repo.checkSecretAuthor(authorId, userId); err != nil {
			repo.Logger.Error("failed to get secret author", err)
			return err
		}
	}

	query := repo.QBuilder.Delete("secrets").Where("id = ?", secretId)

	_, err := query.RunWith(repo.DB).Exec()
	if err != nil {
		repo.Logger.Error("failed to delete secret", err)
		return fmt.Errorf("error executing delete query: %v", err)
	}

	return nil
}

func (repo *Repository) getSecretAuthor(secretId, userRole, userId string) (int, error) {
	secret, err := repo.GetSecretById(secretId, userRole, userId)
	if err != nil {
		repo.Logger.Error("get secret by id failed", err)
		return 0, fmt.Errorf("error querying secret: %v", err)
	}

	return secret.AuthorId, nil
}

func (repo *Repository) checkSecretAuthor(authorId int, userId string) error {
	userID, _ := strconv.Atoi(userId)
	repo.Logger.Debug("author id, user id", authorId, userID)
	if authorId != userID {
		repo.Logger.Info("you dont have permission")
		return errors.New("вы не имеете достаточно прав")
	}
	return nil
}
