package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/savageking-io/ogbuser/schema"
	"github.com/savageking-io/ogbuser/token"
	"strings"
	"time"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type User struct {
	data *schema.UserSchema
	db   *sqlx.DB
}

func NewUser(db *sqlx.DB, data *schema.UserSchema) *User {
	return &User{
		data: data,
		db:   db,
	}
}

func (u *User) SetGroups(groups []schema.GroupSchema) {
	u.data.Groups = groups
}

func (u *User) GetId() int {
	return u.data.Id
}

func (u *User) GetUsername() string {
	return u.data.Username
}

func (u *User) GetPassword() string {
	return u.data.Password
}

func (u *User) GetEmail() string {
	return u.data.Email
}

func (u *User) GetCreatedAt() string {
	return u.data.CreatedAt.String()
}

func (u *User) LoadById(ctx context.Context, id int) error {
	if u.db == nil {
		return fmt.Errorf("DB is not initialized")
	}

	tx, err := u.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	u.data = &schema.UserSchema{}

	query := `
		SELECT id, username, password, email, created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL`

	err = tx.GetContext(ctx, u.data, query, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (u *User) LoadByUsername(ctx context.Context, username string) error {
	if u.db == nil {
		return fmt.Errorf("DB is not initialized")
	}

	tx, err := u.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	u.data = &schema.UserSchema{}
	query := `
		SELECT id, username, password, email, created_at, updated_at, deleted_at
		FROM users 
		WHERE username = $1 AND deleted_at IS NULL`

	if strings.Contains(username, "@") {
		query = `
			SELECT id, username, password, email, created_at, updated_at, deleted_at
			FROM users 
			WHERE email = $1 AND deleted_at IS NULL`

	}

	err = tx.GetContext(ctx, u.data, query, username)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (u *User) InitializeSession(ctx context.Context) (*schema.UserSessionSchema, error) {
	if u.db == nil {
		return nil, fmt.Errorf("DB is not initialized")
	}

	if u.data == nil {
		return nil, fmt.Errorf("user is not initialized")
	}

	if u.data.Id == 0 {
		return nil, fmt.Errorf("user id is not set")
	}

	token, err := token.Generate(u.data.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	session := schema.UserSessionSchema{
		UserId:       u.data.Id,
		Token:        token,
		CreatedAt:    time.Now(),
		PlatformName: "steam",
	}
	u.data.Sessions = append(u.data.Sessions, session)

	tx, err := u.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	query := `
		INSERT INTO user_sessions (user_id, token, created_at, platform_name) VALUES (:user_id, :token, :created_at, :platform_name)
		`
	_, err = tx.NamedExecContext(ctx, query, &session)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &session, nil
}
