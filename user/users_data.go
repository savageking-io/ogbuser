package user

import (
	"context"
	"fmt"
	"github.com/savageking-io/ogbuser/db"
	"sync"
)

type UsersData struct {
	users        map[int32]User
	usernameToId map[string]int32
	mutex        sync.Mutex
	db           *db.Database
}

func NewUsersData(db *db.Database) *UsersData {
	return &UsersData{
		users:        make(map[int32]User),
		usernameToId: make(map[string]int32),
		db:           db,
	}
}

func (u *UsersData) SetDb(db *db.Database) {
	u.db = db
}

func (u *UsersData) Add(user *User) error {
	if user == nil {
		return fmt.Errorf("user is nil")
	}
	if user.GetId() == 0 {
		return fmt.Errorf("invalid user id")
	}
	if user.GetUsername() == "" {
		return fmt.Errorf("invalid username")
	}
	defer u.mutex.Unlock()
	u.mutex.Lock()
	u.users[int32(user.GetId())] = *user
	u.usernameToId[user.GetUsername()] = user.GetId()
	return nil
}

func (u *UsersData) Delete(id int32) error {
	defer u.mutex.Unlock()
	u.mutex.Lock()
	if _, ok := u.users[id]; ok {
		delete(u.users, id)
		return nil
	}
	return db.ErrUserNotFound
}

func (u *UsersData) GetById(id int32) (*User, error) {
	if user, ok := u.users[id]; ok {
		return &user, nil
	}

	if u.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	userSchema, err := u.db.LoadUserById(context.Background(), id)
	if err != nil {
		return nil, err
	}

	if userSchema == nil {
		return nil, fmt.Errorf("empty response from db")
	}

	newUser := NewUser(u.db, userSchema)
	if err := u.Add(newUser); err != nil {
		return nil, err
	}
	return u.GetById(id)
}

func (u *UsersData) GetByUsername(username string) (*User, error) {
	userId, ok := u.usernameToId[username]
	if !ok {
		return nil, db.ErrUserNotFound
	}

	if userId == 0 {
		return nil, fmt.Errorf("invalid user id - can't be zero")
	}

	return u.GetById(userId)
}
