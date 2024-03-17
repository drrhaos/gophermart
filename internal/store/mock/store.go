package mock

import (
	"context"
	"gophermart/internal/store"
)

type MockDB struct {
	Users map[int]map[string]string
}

func (m *MockDB) UserRegister(ctx context.Context, login string, password string) error {
	for _, user := range m.Users {
		if user["login"] == login {
			return store.ErrLoginDuplicate
		}
	}
	return nil
}

func (m *MockDB) UserLogin(ctx context.Context, login string, password string) error {
	for _, user := range m.Users {
		if user["login"] == login && user["password"] != password {
			return store.ErrLoginDuplicate
		}
	}
	return nil
}

func (m *MockDB) UserOrders(ctx context.Context, order int) bool {
	return true
}

func (m *MockDB) Ping(ctx context.Context) (exists bool) {
	return true
}
