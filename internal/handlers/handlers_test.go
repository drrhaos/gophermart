package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"gophermart/internal/logger"
	"gophermart/internal/store"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/stretchr/testify/assert"
)

const urlPostUserRegister = "/api/user/register"                // регистрация пользователя
const urlPostUserLogin = "/api/user/login"                      // аутентификация пользователя;
const urlPostUserOrders = "/api/user/orders"                    // загрузка пользователем номера заказа для расчёта;
const urlGetUserOrders = "/api/user/orders"                     // получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях;
const urlGetUserBalance = "/api/user/balance"                   // получение текущего баланса счёта баллов лояльности пользователя;
const urlPostUserBalanceWithdraw = "/api/user/balance/withdraw" // запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа;
const urlGetUserWithdrawals = "/api/user/withdrawals"           // получение информации о выводе средств с накопительного счёта пользователем.

type MockDB struct {
	users map[int]map[string]string
}

// func (m *MockDB) GetUserByID(id int) map[string]string {
// 	return m.users[id]
// }

func (m *MockDB) UserRegister(ctx context.Context, login string, password string) error {
	for _, user := range m.users {
		if user["login"] == login {
			return store.ErrLoginDuplicate
		}
	}
	return nil
}

func (m *MockDB) UserLogin(ctx context.Context, login string, password string) error {
	return nil
}

func (m *MockDB) UserOrders(ctx context.Context, order int) bool {
	return true
}

func (m *MockDB) Ping(ctx context.Context) (exists bool) {
	return true
}

func TestPostUserRegister(t *testing.T) {
	logger.Init()
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)

	mockDB := &MockDB{
		users: map[int]map[string]string{
			1: {"id": "1", "login": "test", "password": "$2a$10$kte3HgQ6VtHaZSBVc0Cr2OSHQnVL3UB5C0mJLnPVA5W3y.EfNz7rC"},
			2: {"id": "2", "login": "test2", "password": "$2a$10$kte3HgQ6VtHaZSBVc0Cr2OSHQnVL3UB5C0mJLnPVA5W3z.EfNz7rC"},
		},
	}

	storage := &store.StorageContext{}
	storage.SetStorage(mockDB)

	r := chi.NewRouter()
	r.Use(middleware.Compress(5, "application/json", "text/html"))

	r.Post(urlPostUserRegister, func(w http.ResponseWriter, r *http.Request) {
		PostUserRegister(w, r, storage, tokenAuth)
	})
	r.Post(urlPostUserLogin, func(w http.ResponseWriter, r *http.Request) {
		PostUserLogin(w, r, storage, tokenAuth)
	})
	r.Post(urlPostUserOrders, func(w http.ResponseWriter, r *http.Request) {
		PostUserOrders(w, r, storage)
	})
	r.Get(urlGetUserOrders, func(w http.ResponseWriter, r *http.Request) {
		GetUserOrders(w, r)
	})
	r.Get(urlGetUserBalance, func(w http.ResponseWriter, r *http.Request) {
		GetUserBalance(w, r)
	})

	type want struct {
		code int
	}
	tests := []struct {
		name       string
		url        string
		body       User
		typeReqest string
		want       want
	}{
		// 200 — пользователь успешно зарегистрирован и аутентифицирован;
		// 400 — неверный формат запроса;
		// 409 — логин уже занят;
		// 500 — внутренняя ошибка сервера.
		{
			name: "пользователь успешно зарегистрирован и аутентифицирован",
			url:  urlPostUserRegister,
			body: User{
				Login:    "test3",
				Password: "test3",
			},
			typeReqest: http.MethodPost,
			want: want{
				code: 200,
			},
		},
		{
			name:       "неверный формат запроса",
			url:        urlPostUserRegister,
			typeReqest: http.MethodPost,
			want: want{
				code: 400,
			},
		},
		{
			name: "логин уже занят",
			url:  urlPostUserRegister,
			body: User{
				Login:    "test",
				Password: "test",
			},
			typeReqest: http.MethodPost,
			want: want{
				code: 409,
			},
		},
		{
			name: "внутренняя ошибка сервера",
			url:  urlPostUserRegister,
			body: User{
				Login:    "test",
				Password: "test",
			},
			typeReqest: http.MethodPost,
			want: want{
				code: 500,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bodyJson, _ := json.Marshal(test.body)
			req := httptest.NewRequest(test.typeReqest, test.url, bytes.NewReader(bodyJson))
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != test.want.code {
				t.Errorf("expected status OK; got %v", w.Code)
			}

			assert.Equal(t, test.want.code, w.Code)
		})
	}
}
