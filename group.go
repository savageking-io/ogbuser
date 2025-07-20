package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/savageking-io/ogbuser/schema"
	log "github.com/sirupsen/logrus"
)

type Group struct {
	data *schema.GroupSchema
	db   *sqlx.DB
}

func NewGroup(db *sqlx.DB) *Group {
	return &Group{
		db: db,
	}
}

func (g *Group) Init(ctx context.Context, id int) error {
	if g.db == nil {
		return fmt.Errorf("DB is not initialized")
	}

	tx, err := g.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}

	g.data = &schema.GroupSchema{}
	query := `
		SELECT id, parent_id, name, description, created_at, updated_at, deleted_at 
		FROM groups
		WHERE id = $1 AND deleted_at IS NULL
	`
	err = tx.GetContext(ctx, g.data, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("group not found")
		}
		return err
	}

	var permissions []schema.GroupPermissionSchema
	query = `
		SELECT id, group_id, permission, read, write, delete, domain, created_at, updated_at, deleted_at 
		FROM group_permissions 
		WHERE group_id = $1 AND deleted_at IS NULL
	`
	err = tx.SelectContext(ctx, &permissions, query, g.data.Id)
	if err != nil {
		return err
	}

	for _, permission := range permissions {
		log.Debugf("Adding permission %s for group %s: %t %t %t [%s]", permission.Permission, g.data.Name, permission.Read, permission.Write, permission.Delete, permission.Domain)
		g.data.Permissions = append(g.data.Permissions, permission)
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	log.Infof("Group %s initialized. Total number of permissions: %d", g.data.Name, len(g.data.Permissions))

	return nil
}
