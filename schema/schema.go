package schema

import "time"

type UserSchema struct {
	Id        int                 `db:"id"`
	Username  string              `db:"username"`
	Password  string              `db:"password"`
	Email     string              `db:"email"`
	CreatedAt time.Time           `db:"created_at"`
	UpdatedAt time.Time           `db:"updated_at"`
	DeletedAt *time.Time          `db:"deleted_at"`
	Platforms []PlatformSchema    `db:"-"`
	Sessions  []UserSessionSchema `db:"-"` // User can have multiple sessions from browser/platforms @TODO: Implement conflicting sessions - for example can't play from two platforms at the same time
	Groups    []GroupSchema       `db:"-"`
}

type PlatformSchema struct {
	Id             int        `db:"id"`
	UserId         int        `db:"user_id"`
	PlatformName   string     `db:"platform_name"`
	PlatformUserId string     `db:"platform_user_id"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
	DeletedAt      *time.Time `db:"deleted_at"`
}

type UserSessionSchema struct {
	Id           int        `db:"id"`
	UserId       int        `db:"user_id"`
	Token        string     `db:"token"`
	PlatformName string     `db:"platform_name"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at"`
}

type GroupSchema struct {
	Id          int                     `db:"id"`
	ParentId    int                     `db:"parent_id"`
	Name        string                  `db:"name"`
	Description *string                 `db:"description"`
	IsSpecial   bool                    `db:"is_special"`
	CreatedAt   time.Time               `db:"created_at"`
	UpdatedAt   time.Time               `db:"updated_at"`
	DeletedAt   *time.Time              `db:"deleted_at"`
	Permissions []GroupPermissionSchema `db:"-"`
}

type GroupMemberSchema struct {
	Id        int        `db:"id"`
	GroupId   int        `db:"group_id"`
	UserId    int        `db:"user_id"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

type GroupPermissionSchema struct {
	Id         int        `db:"id"`
	GroupId    int        `db:"group_id"`
	Permission string     `db:"permission"`
	Read       bool       `db:"read"`
	Write      bool       `db:"write"`
	Delete     bool       `db:"delete"`
	Domain     string     `db:"domain"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
	DeletedAt  *time.Time `db:"deleted_at"`
}
