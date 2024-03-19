package store

import (
	"context"
	"errors"
	"gophermart/internal/models"
)

type StorageInterface interface {
	UserRegister(ctx context.Context, login string, password string) error
	UserLogin(ctx context.Context, login string, password string) error
	UploadUserOrders(ctx context.Context, login string, order int) error
	GetUserOrders(ctx context.Context, login string) ([]models.StatusOrders, error)
	Ping(ctx context.Context) bool
}

type StorageContext struct {
	storage StorageInterface
}

var ErrLoginDuplicate = errors.New("user duplicate")
var ErrAuthentication = errors.New("invalid user name or password")
var ErrDuplicateOrder = errors.New("duplicate order")
var ErrDuplicateOrderOtherUser = errors.New("duplicate order other user")

func (sc *StorageContext) SetStorage(storage StorageInterface) {
	sc.storage = storage
}

func (sc *StorageContext) UserRegister(ctx context.Context, login string, password string) error {
	return sc.storage.UserRegister(ctx, login, password)
}

func (sc *StorageContext) UserLogin(ctx context.Context, login string, password string) error {
	return sc.storage.UserLogin(ctx, login, password)
}

func (sc *StorageContext) UploadUserOrders(ctx context.Context, login string, order int) error {
	return sc.storage.UploadUserOrders(ctx, login, order)
}

func (sc *StorageContext) GetUserOrders(ctx context.Context, login string) ([]models.StatusOrders, error) {
	return sc.storage.GetUserOrders(ctx, login)
}

func (sc *StorageContext) Ping(ctx context.Context) (exists bool) {
	return sc.storage.Ping(ctx)
}
