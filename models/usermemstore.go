package usermodel

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNotFound = errors.New("not found")
)

type MemStore struct {
	list []User
}

func NewMemStore(initialUsers []User) *MemStore {
	return &MemStore{list: initialUsers}
}

func (m *MemStore) Add(user User) error {
	m.list = append(m.list, user)
	return nil
}

func (m *MemStore) List(page, limit int) ([]User, error) {
	start := (page - 1) * limit
	end := start + limit

	if start > len(m.list) {
		return []User{}, nil
	}

	if end > len(m.list) {
		end = len(m.list)
	}
	return m.list[start:end], nil
}

func (m MemStore) Get(id string) (User, error) {
	for _, user := range m.list {
		if user.ID == id {
			return user, nil
		}
	}
	return User{}, ErrNotFound
}
func (m MemStore) GetEmail(email string) (User, error) {
	for _, user := range m.list {
		if strings.EqualFold(user.Email, email) {
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
