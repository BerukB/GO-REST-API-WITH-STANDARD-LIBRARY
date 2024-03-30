package usermodel

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound = errors.New("not found")
)

type MemStore struct {
	list []User
}

func NewMemStore() *MemStore {
	return &MemStore{list: []User{}}
}

func (m *MemStore) Add(user User) error {
	m.list = append(m.list, user)
	return nil
}

func (m MemStore) List() ([]User, error) {
	return m.list, nil
}

func (m MemStore) Get(id string) (User, error) {
	for _, user := range m.list {
		if user.ID == id {
			return user, nil
		}
	}
	return User{}, ErrNotFound
}

func (m *MemStore) Update(id string, user User) (User, error) {
	fmt.Println("here update", user)

	for i, u := range m.list {
		if u.ID == id {
			// var usernew User
			m.list = append(m.list[:i], m.list[i+1:]...)
			user.ID = id
			m.list = append(m.list, user)

			return user, nil
		}
	}
	return User{}, ErrNotFound
}

func (m *MemStore) Remove(id string) error {
	for i, user := range m.list {
		if user.ID == id {
			m.list = append(m.list[:i], m.list[i+1:]...)
			return nil
		}
	}
	return ErrNotFound
}