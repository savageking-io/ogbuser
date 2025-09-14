package user

import "sync"

type UsersData struct {
	users map[int]User
	mutex sync.Mutex
}

func NewUsersData() *UsersData {
	return &UsersData{
		users: make(map[int]User),
	}
}

func (u *UsersData) Add(user User) {
	defer u.mutex.Unlock()
	u.mutex.Lock()
	u.users[int(user.GetId())] = user
}

func (u *UsersData) Delete(id int) error {
	defer u.mutex.Unlock()
	u.mutex.Lock()
	if _, ok := u.users[id]; ok {
		delete(u.users, id)
		return nil
	}
	return ErrUserNotFound
}

func (u *UsersData) GetById(id int32) (*User, error) {
	if user, ok := u.users[int(id)]; ok {
		return &user, nil
	}
	return nil, ErrUserNotFound
}
