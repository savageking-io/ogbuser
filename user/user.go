package user

import (
	"context"
	"fmt"
	"github.com/savageking-io/ogbuser/db"
	"github.com/savageking-io/ogbuser/group"
	"github.com/savageking-io/ogbuser/perm"
	"github.com/savageking-io/ogbuser/schema"
	"github.com/savageking-io/ogbuser/token"
	log "github.com/sirupsen/logrus"
)

type User struct {
	raw      *schema.UserSchema
	db       *db.Database
	perms    *perm.Perm
	groups   *group.GroupsData
	sessions []*schema.UserSessionSchema
}

func NewUser(db *db.Database, data *schema.UserSchema) *User {
	return &User{
		raw:    data,
		db:     db,
		perms:  perm.NewPerm(),
		groups: group.NewGroupsData(),
	}
}

func (u *User) AddGroup(group *group.Group) {
	u.groups.Add(group)
}

func (u *User) GetId() int32 {
	return u.raw.Id
}

func (u *User) GetUsername() string {
	return u.raw.Username
}

func (u *User) GetPassword() string {
	return u.raw.Password
}

func (u *User) GetEmail() string {
	return u.raw.Email
}

func (u *User) GetCreatedAt() string {
	return u.raw.CreatedAt.String()
}

// LoadByUsername will request user data from database and populate the object
// return ErrUserNotFound if user with provided username doesn't exist or generic error on other failures
func (u *User) LoadByUsername(ctx context.Context, username string) error {
	log.Tracef("User::LoadByUsername: %s", username)
	if u.db == nil {
		return fmt.Errorf("DB is not initialized")
	}

	userSchema, err := u.db.LoadUserByUsername(ctx, username)
	if err != nil {
		return err
	}

	if userSchema == nil {
		return fmt.Errorf("nil user schema without an error")
	}

	u.raw = userSchema
	return nil
}

// InitializeSession will generate a new token for the user and store it in database
func (u *User) InitializeSession(ctx context.Context, inPlatform string) (*schema.UserSessionSchema, error) {
	log.Traceln("User::InitializeSession")
	if u.db == nil {
		return nil, fmt.Errorf("DB is not initialized")
	}

	if u.raw == nil {
		return nil, fmt.Errorf("user is not initialized")
	}

	if u.raw.Id == 0 {
		return nil, fmt.Errorf("user id is not set")
	}

	userToken, err := token.Generate(u.raw.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	session, err := u.db.SaveUserSession(ctx, u.GetId(), userToken, inPlatform)
	if err != nil {
		return nil, fmt.Errorf("failed to save user session: %w", err)
	}

	u.sessions = append(u.sessions, session)

	return session, nil
}

// LoadGroups will return a list of groups that this user belongs to
func (u *User) LoadGroups(ctx context.Context) ([]int32, error) {
	log.Tracef("User::LoadGroups")
	if u.db == nil {
		return nil, fmt.Errorf("DB is not initialized")
	}

	if u.raw == nil {
		return nil, fmt.Errorf("user is not initialized")
	}

	return u.db.GetUserGroupIds(ctx, u.raw.Id)
}

func (u *User) HasPermission(ctx context.Context, permission string, domain string) (*perm.Permission, error) {
	if u.perms == nil {
		return nil, fmt.Errorf("permissions are not initialized")
	}

	return u.perms.GetPermission(domain, permission), nil
}
