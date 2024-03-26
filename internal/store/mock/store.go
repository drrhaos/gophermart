package mock

import (
	"context"
	"fmt"
	"gophermart/internal/logger"
	"gophermart/internal/models"
	"gophermart/internal/store"
	"strconv"
	"time"

	"go.uber.org/zap"
)

type MockDB struct {
	Users       map[int]map[string]string
	Orders      map[int]map[string]string
	Withdrawals map[int]map[string]string
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
			return store.ErrAuthentication
		}
	}
	return nil
}

func (m *MockDB) UploadUserOrders(ctx context.Context, login string, order int64) error {
	idUser := "-1"
	for _, user := range m.Users {
		if user["login"] == login {
			idUser = user["id"]
		}
	}
	for _, orderRow := range m.Orders {
		if orderRow["number"] == fmt.Sprintf("%d", order) && orderRow["user_id"] == idUser {
			return store.ErrDuplicateOrder
		} else if orderRow["number"] == fmt.Sprintf("%d", order) && idUser != "-1" && orderRow["user_id"] != idUser {
			return store.ErrDuplicateOrderOtherUser
		}
	}

	return nil
}

func (m *MockDB) GetUserOrders(ctx context.Context, login string) ([]models.StatusOrders, error) {
	idUser := "-1"
	var orderUser models.StatusOrders
	var ordersUser []models.StatusOrders
	var accrual string
	for _, user := range m.Users {
		if user["login"] == login {
			idUser = user["id"]
		}
	}
	for _, orderRow := range m.Orders {
		if orderRow["user_id"] == idUser {
			orderUser.Number = orderRow["number"]
			orderUser.Status = orderRow["status"]
			accrual = orderRow["accrual"]
			if accrual != "" {
				orderUser.Accrual, _ = strconv.ParseFloat(accrual, 64)
			}
			orderUser.UploadedAt, _ = time.Parse("2006-01-02T15:04:05Z", orderRow["uploaded_at"])
			ordersUser = append(ordersUser, orderUser)
		}
	}
	return ordersUser, nil
}

func (m *MockDB) GetUserBalance(ctx context.Context, login string) (models.Balance, error) {
	var userBalance models.Balance

	for _, user := range m.Users {
		if user["login"] == login {
			userBalance.Current, _ = strconv.ParseFloat(user["sum"], 64)
			userBalance.Withdrawn, _ = strconv.ParseFloat(user["withdrawn"], 64)
		}
	}
	return userBalance, nil
}

func (m *MockDB) UpdateUserBalanceWithdraw(ctx context.Context, login string, order string, sum float64) error {
	var balanceS string
	for _, user := range m.Users {
		if user["login"] == login {
			balanceS = user["sum"]
		}
	}

	number, err := strconv.ParseInt(order, 10, 64)
	if err != nil {
		logger.Logger.Warn("Не удалось добавить значение", zap.Error(err))
		return err

	}

	balance, _ := strconv.ParseFloat(balanceS, 64)
	if balance < sum {
		logger.Logger.Warn("на счету недостаточно средств")
		return store.ErrInsufficientFunds
	}

	err = m.UploadUserOrders(ctx, login, number)
	if err != nil {
		logger.Logger.Warn("Не удалось добавить значение", zap.Error(err))
		return err
	}

	return nil
}

func (m *MockDB) GetUserWithdrawals(ctx context.Context, login string) ([]models.BalanceWithdrawals, error) {

	return nil, nil
}

func (m *MockDB) GetOrdersProcessing(ctx context.Context) ([]int64, error) {
	return nil, nil
}

func (m *MockDB) UpdateStatusOrders(ctx context.Context, statusOrder *models.StatusOrdersAccrual) error {
	return nil
}

func (m *MockDB) Ping(ctx context.Context) (exists bool) {
	return true
}
