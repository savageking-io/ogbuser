package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/savageking-io/ogbcommon/dbinit"
	"github.com/savageking-io/ogbuser/schema"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrGroupNotFound   = errors.New("group not found")
	ErrSessionNotFound = errors.New("session not found")
)

type PostgresConfig struct {
	Hostname        string `yaml:"hostname"`
	Port            uint16 `yaml:"port"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	Database        string `yaml:"database"`
	SslMode         bool   `yaml:"ssl_mode"`
	MaxOpenCons     int    `yaml:"max_open_cons"`
	MaxIdleCons     int    `yaml:"max_idle_cons"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
	SourceFile      string `yaml:"source_file"`
}

type Database struct {
	db              *sqlx.DB
	hostname        string
	port            uint16
	username        string
	password        string
	database        string
	sslMode         bool
	maxOpenCons     int
	maxIdleCons     int
	sourceFile      string
	connMaxLifetime time.Duration
}

func (d *Database) Init(conf *PostgresConfig) error {
	if conf == nil {
		return fmt.Errorf("conf is nil")
	}

	d.hostname = conf.Hostname
	d.port = conf.Port
	d.username = conf.Username
	d.password = conf.Password
	d.database = conf.Database
	d.sslMode = conf.SslMode
	d.maxOpenCons = conf.MaxOpenCons
	d.maxIdleCons = conf.MaxIdleCons
	d.connMaxLifetime = time.Duration(conf.ConnMaxLifetime) * time.Second
	d.sourceFile = conf.SourceFile

	return nil
}

func (d *Database) buildConnString() string {
	sslMode := "disable"
	if d.sslMode {
		sslMode = "require"
	}

	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.hostname,
		d.port,
		d.username,
		d.password,
		d.database,
		sslMode)
}

func (d *Database) Connect() error {
	log.Traceln("Database::Connect")
	connStr := d.buildConnString()
	log.Tracef("Database::Connect. ConnStr: %s", connStr)

	var err error
	d.db, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	d.db.SetMaxOpenConns(d.maxOpenCons)
	d.db.SetMaxIdleConns(d.maxIdleCons)
	d.db.SetConnMaxLifetime(d.connMaxLifetime)

	if err := d.db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

func (d *Database) TryPopulate() error {
	if d.db == nil {
		return fmt.Errorf("db is nil")
	}
	if d.sourceFile == "" {
		log.Errorf("Database::TryPopulate No database source file provided. Cannot populate database")
		return nil
	}
	// Populate db if it's empty
	err := dbinit.VerifyTables([]string{"users", "platforms", "groups"}, d.db)
	if err != nil {
		log.Warnf("Table verification failed: %s", err.Error())
		if err := dbinit.Populate(d.db, d.sourceFile); err != nil {
			log.Errorf("Failed to populate database: %s", err.Error())
			return err
		}
	}
	return nil
}

func (d *Database) LoadGroups(ctx context.Context) ([]schema.GroupSchema, error) {
	if d.db == nil {
		return nil, fmt.Errorf("db is nil")
	}
	tx, err := d.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, parent_id, name, description, created_at, updated_at, deleted_at 
		FROM groups 
		WHERE deleted_at IS NULL
	`

	var groups []schema.GroupSchema
	err = tx.SelectContext(context.Background(), &groups, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrGroupNotFound
		}
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return nil, nil
}

func (d *Database) LoadGroupById(ctx context.Context, id int32) (*schema.GroupSchema, error) {
	if d.db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	tx, err := d.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, parent_id, name, description, created_at, updated_at, deleted_at
		FROM groups
		WHERE id = $1 AND deleted_at IS NULL
	`
	var group schema.GroupSchema
	err = tx.GetContext(ctx, &group, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrGroupNotFound
		}
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &group, nil
}

func (d *Database) LoadGroupPermissions(ctx context.Context, groupId int32) ([]schema.GroupPermissionSchema, error) {
	if d.db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	tx, err := d.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	var permissions []schema.GroupPermissionSchema
	query := `
		SELECT id, group_id, permission, read, write, delete, domain, created_at, updated_at, deleted_at 
		FROM group_permissions 
		WHERE group_id = $1 AND deleted_at IS NULL
	`

	err = tx.SelectContext(ctx, &permissions, query, groupId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return permissions, nil
}

func (d *Database) LoadUserById(ctx context.Context, id int32) (*schema.UserSchema, error) {
	if d.db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	tx, err := d.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	result := &schema.UserSchema{}

	query := `
		SELECT id, username, password, email, created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL`

	err = tx.GetContext(ctx, result, query, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

func (d *Database) LoadUserByUsername(ctx context.Context, username string) (*schema.UserSchema, error) {
	log.Tracef("Database::LoadByUsername: %s", username)
	if d.db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	tx, err := d.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	result := &schema.UserSchema{}
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

	log.Debugf("Database::LoadByUsername. Query: %s", query)
	err = tx.GetContext(ctx, result, query, username)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

func (d *Database) GetUserSessionByToken(ctx context.Context, token string) (*schema.UserSessionSchema, error) {
	if d.db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	tx, err := d.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT user_id, created_at, updated_at
		FROM user_sessions
		WHERE token = $1 AND deleted_at IS NULL`

	var session schema.UserSessionSchema
	err = tx.GetContext(ctx, &session, query, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &session, nil
}

func (d *Database) GetUserGroupIds(ctx context.Context, userId int32) ([]int32, error) {
	if d.db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	tx, err := d.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	query := `SELECT group_id FROM group_members WHERE user_id = $1 AND deleted_at IS NULL`
	var groupIds []int32
	err = tx.SelectContext(ctx, &groupIds, query, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrGroupNotFound
		}
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return groupIds, nil
}

func (d *Database) SaveUserSession(ctx context.Context, userId int32, token, platform string) (*schema.UserSessionSchema, error) {
	if d.db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	session := &schema.UserSessionSchema{
		UserId:       userId,
		Token:        token,
		CreatedAt:    time.Now(),
		PlatformName: platform,
	}

	tx, err := d.db.BeginTxx(ctx, nil)
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

	return session, nil
}
