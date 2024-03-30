package store

import (
	"context"
	"errors"
	"gophermart/internal/models"
)

type StorageInterface interface {
	UserRegister(ctx context.Context, login string, password string) error
	UserLogin(ctx context.Context, login string, password string) error
	UploadUserOrders(ctx context.Context, login string, order int64) error
	GetUserOrders(ctx context.Context, login string) ([]models.StatusOrders, error)
	GetUserBalance(ctx context.Context, login string) (models.Balance, error)
	UpdateUserBalanceWithdraw(ctx context.Context, login string, order string, sum float64) error
	GetUserWithdrawals(ctx context.Context, login string) ([]models.BalanceWithdrawals, error)
	GetOrdersProcessing(ctx context.Context) ([]int64, error)
	UpdateStatusOrders(ctx context.Context, statusOrder *models.StatusOrdersAccrual) error
	Ping(ctx context.Context) bool
}

type StorageContext struct {
	storage StorageInterface
}

var ErrLoginDuplicate = errors.New("user duplicate")
var ErrAuthentication = errors.New("invalid user name or password")
var ErrDuplicateOrder = errors.New("duplicate order")
var ErrDuplicateOrderOtherUser = errors.New("duplicate order other user")
var ErrOrderNotFound = errors.New("order not found")
var ErrInsufficientFunds = errors.New("insufficient funds")

func (sc *StorageContext) SetStorage(storage StorageInterface) {
	sc.storage = storage
}

func (sc *StorageContext) UserRegister(ctx context.Context, login string, password string) error {
	return sc.storage.UserRegister(ctx, login, password)
}

func (sc *StorageContext) UserLogin(ctx context.Context, login string, password string) error {
	return sc.storage.UserLogin(ctx, login, password)
}

func (sc *StorageContext) UploadUserOrders(ctx context.Context, login string, order int64) error {
	return sc.storage.UploadUserOrders(ctx, login, order)
}

func (sc *StorageContext) GetUserOrders(ctx context.Context, login string) ([]models.StatusOrders, error) {
	return sc.storage.GetUserOrders(ctx, login)
}

func (sc *StorageContext) GetUserBalance(ctx context.Context, login string) (models.Balance, error) {
	return sc.storage.GetUserBalance(ctx, login)
}

func (sc *StorageContext) UpdateUserBalanceWithdraw(ctx context.Context, login string, order string, sum float64) error {
	return sc.storage.UpdateUserBalanceWithdraw(ctx, login, order, sum)
}

func (sc *StorageContext) GetUserWithdrawals(ctx context.Context, login string) ([]models.BalanceWithdrawals, error) {
	return sc.storage.GetUserWithdrawals(ctx, login)
}

func (sc *StorageContext) Ping(ctx context.Context) (exists bool) {
	return sc.storage.Ping(ctx)
}

func (sc *StorageContext) GetOrdersProcessing(ctx context.Context) ([]int64, error) {
	return sc.storage.GetOrdersProcessing(ctx)
}

func (sc *StorageContext) UpdateStatusOrders(ctx context.Context, statusOrder *models.StatusOrdersAccrual) error {
	return sc.storage.UpdateStatusOrders(ctx, statusOrder)
}
