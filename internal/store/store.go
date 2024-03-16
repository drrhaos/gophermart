package store

import (
	"context"
	"errors"
)

type StorageInterface interface {
	UserRegister(ctx context.Context, user string, password string) error
	UserLogin(ctx context.Context, user string, password string) error
	UserOrders(ctx context.Context, order int) bool
	Ping(ctx context.Context) bool
}

type StorageContext struct {
	storage StorageInterface
}

var ErrLoginDuplicate = errors.New("user duplicate")
var ErrAuthentication = errors.New("invalid user name or password")

func (sc *StorageContext) SetStorage(storage StorageInterface) {
	sc.storage = storage
}

func (sc *StorageContext) UserRegister(ctx context.Context, login string, password string) error {
	return sc.storage.UserRegister(ctx, login, password)
}

func (sc *StorageContext) UserLogin(ctx context.Context, login string, password string) error {
	return sc.storage.UserLogin(ctx, login, password)
}

func (sc *StorageContext) UserOrders(ctx context.Context, order int) bool {
	return sc.storage.UserOrders(ctx, order)
}

func (sc *StorageContext) Ping(ctx context.Context) (exists bool) {
	return sc.storage.Ping(ctx)
}
