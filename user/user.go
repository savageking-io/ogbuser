package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/savageking-io/ogbuser/schema"
	"github.com/savageking-io/ogbuser/token"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type Permission struct {
	Read   bool
	Write  bool
	Delete bool
}

type User struct {
	data              *schema.UserSchema
	db                *sqlx.DB
	ownPermissions    map[string]Permission
	partyPermissions  map[string]Permission
	guildPermissions  map[string]Permission
	globalPermissions map[string]Permission
}

func NewUser(db *sqlx.DB, data *schema.UserSchema) *User {
	return &User{
		data:              data,
		db:                db,
		ownPermissions:    make(map[string]Permission),
		partyPermissions:  make(map[string]Permission),
		guildPermissions:  make(map[string]Permission),
		globalPermissions: make(map[string]Permission),
	}
}

func (u *User) SetGroups(groups []schema.GroupSchema) {
	u.data.Groups = groups
}

func (u *User) GetId() int32 {
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
	log.Tracef("User::LoadById %d", id)
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
	log.Tracef("User::LoadByUsername: %s", username)
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

func (u *User) InitializeSession(ctx context.Context, inPlatform string) (*schema.UserSessionSchema, error) {
	log.Traceln("User::InitializeSession")
	if u.db == nil {
		return nil, fmt.Errorf("DB is not initialized")
	}

	if u.data == nil {
		return nil, fmt.Errorf("user is not initialized")
	}

	if u.data.Id == 0 {
		return nil, fmt.Errorf("user id is not set")
	}

	userToken, err := token.Generate(u.data.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	session := schema.UserSessionSchema{
		UserId:       u.data.Id,
		Token:        userToken,
		CreatedAt:    time.Now(),
		PlatformName: inPlatform,
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

func (u *User) LoadGroups(ctx context.Context) error {
	log.Tracef("User::LoadGroups")
	if u.db == nil {
		return fmt.Errorf("DB is not initialized")
	}

	if u.data == nil {
		return fmt.Errorf("user is not initialized")
	}

	tx, err := u.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil
	}
	query := `SELECT group_id FROM group_members WHERE user_id = $1 AND deleted_at IS NULL`
	_, err = tx.NamedExecContext(ctx, query, &u.data.Groups)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (u *User) GetPermission(permission string, domain string) (Permission, error) {
	log.Tracef("User::GetPermission: %s", permission)
	if domain == "own" {
		return u.GetOwnPermission(permission)
	}
	if domain == "party" {
		return u.GetPartyPermission(permission)
	}
	if domain == "guild" {
		return u.GetGuildPermission(permission)
	}
	if domain == "global" {
		return u.GetGlobalPermission(permission)
	}
	return Permission{}, fmt.Errorf("invalid domain")
}

func (u *User) GetOwnPermission(permission string) (Permission, error) {
	log.Tracef("User::GetOwnPermission: %s", permission)
	perm, ok := u.ownPermissions[permission]
	if !ok {
		return Permission{}, fmt.Errorf("permission not found")
	}
	return perm, nil
}

func (u *User) GetPartyPermission(permission string) (Permission, error) {
	log.Tracef("User::GetPartyPermission: %s", permission)
	perm, ok := u.partyPermissions[permission]
	if !ok {
		return Permission{}, fmt.Errorf("permission not found")
	}
	return perm, nil
}

func (u *User) GetGuildPermission(permission string) (Permission, error) {
	log.Tracef("User::GetGuildPermission: %s", permission)
	perm, ok := u.guildPermissions[permission]
	if !ok {
		return Permission{}, fmt.Errorf("permission not found")
	}
	return perm, nil
}

func (u *User) GetGlobalPermission(permission string) (Permission, error) {
	log.Tracef("User::GetGlobalPermission: %s", permission)
	perm, ok := u.globalPermissions[permission]
	if !ok {
		return Permission{}, fmt.Errorf("permission not found")
	}
	return perm, nil
}
