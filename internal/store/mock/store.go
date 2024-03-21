package mock

import (
	"context"
	"gophermart/internal/models"
	"gophermart/internal/store"
	"strconv"
	"time"
)

type MockDB struct {
	Users  map[int]map[string]string
	Orders map[int]map[string]string
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

func (m *MockDB) UploadUserOrders(ctx context.Context, login string, order int) error {
	idUser := "-1"
	for _, user := range m.Users {
		if user["login"] == login {
			idUser = user["id"]
		}
	}
	for _, orderRow := range m.Orders {
		if orderRow["number"] == strconv.Itoa(order) && orderRow["user_id"] == idUser {
			return store.ErrDuplicateOrder
		} else if orderRow["number"] == strconv.Itoa(order) && idUser != "-1" && orderRow["user_id"] != idUser {
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
				orderUser.Accrual, _ = strconv.ParseInt(accrual, 10, 64)
			}
			orderUser.UploadedAt, _ = time.Parse("2006-01-02T15:04:05Z", orderRow["uploaded_at"])
			ordersUser = append(ordersUser, orderUser)
		}
	}
	return ordersUser, nil
}

func (m *MockDB) Ping(ctx context.Context) (exists bool) {
	return true
}
