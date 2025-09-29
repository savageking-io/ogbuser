package group

import (
	"context"
	"fmt"
	"github.com/savageking-io/ogbuser/db"
	"github.com/savageking-io/ogbuser/perm"
	"github.com/savageking-io/ogbuser/schema"
	log "github.com/sirupsen/logrus"
)

type Group struct {
	raw        schema.GroupSchema
	db         *db.Database
	perms      *perm.Perm
	hasRawData bool
	hasId      bool
}

func NewGroup(db *db.Database) *Group {
	return &Group{
		db:    db,
		perms: perm.NewPerm(),
	}
}

func NewGroupFromId(db *db.Database, id int32) *Group {
	return &Group{
		db:    db,
		raw:   schema.GroupSchema{Id: id},
		hasId: true,
	}
}

func NewGroupFromSchema(db *db.Database, schema *schema.GroupSchema) *Group {
	return &Group{
		db:         db,
		raw:        *schema,
		hasRawData: true,
	}
}

func (g *Group) Init(ctx context.Context) error {
	if g.db == nil {
		return fmt.Errorf("DB is not initialized")
	}

	if g.hasId {
		if g.raw.Id == 0 {
			return fmt.Errorf("invalid group id")
		}
		raw, err := g.db.LoadGroupById(ctx, g.raw.Id)
		if err != nil {
			return err
		}
		g.raw = *raw
	}

	permissions, err := g.db.LoadGroupPermissions(ctx, g.raw.Id)
	if err != nil {
		return err
	}

	for _, permission := range permissions {
		log.Debugf("Adding permission %s for group %s: %t %t %t [%s]", permission.Permission, g.GetName(), permission.Read, permission.Write, permission.Delete, permission.Domain)
		if err := g.perms.Populate(&permission); err != nil {
			log.Errorf("Failed to populate permissions for group %s: %s", g.GetName(), err.Error())
		}
	}

	log.Infof("Group %s initialized. Total number of permissions: %d", g.GetName(), g.perms.Count())

	return nil
}

func (g *Group) GetId() int32 {
	return g.raw.Id
}

func (g *Group) GetName() string {
	return g.raw.Name
}
