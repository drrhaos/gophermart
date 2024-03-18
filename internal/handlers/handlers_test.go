package handlers

import (
	"bytes"
	"encoding/json"
	"gophermart/internal/logger"
	"gophermart/internal/models"
	"gophermart/internal/store"
	"gophermart/internal/store/mock"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestPostUserRegister(t *testing.T) {
	logger.Init()
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)

	mockDB := &mock.MockDB{
		Users: map[int]map[string]string{
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
	r.Post(urlPostUserBalanceWithdraw, func(w http.ResponseWriter, r *http.Request) {
		PostUserBalanceWithdraw(w, r)
	})
	r.Get(urlGetUserWithdrawals, func(w http.ResponseWriter, r *http.Request) {
		GetUserWithdrawals(w, r)
	})

	type want struct {
		code int
	}
	tests := []struct {
		name       string
		url        string
		body       models.User
		typeReqest string
		want       want
	}{
		{
			name: "пользователь успешно зарегистрирован и аутентифицирован",
			url:  urlPostUserRegister,
			body: models.User{
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
			body: models.User{
				Login:    "test",
				Password: "test",
			},
			typeReqest: http.MethodPost,
			want: want{
				code: 409,
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

func TestPostUserLogin(t *testing.T) {
	logger.Init()
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)

	mockDB := &mock.MockDB{
		Users: map[int]map[string]string{
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
	r.Post(urlPostUserBalanceWithdraw, func(w http.ResponseWriter, r *http.Request) {
		PostUserBalanceWithdraw(w, r)
	})
	r.Get(urlGetUserWithdrawals, func(w http.ResponseWriter, r *http.Request) {
		GetUserWithdrawals(w, r)
	})

	type want struct {
		code int
	}
	tests := []struct {
		name       string
		url        string
		body       models.User
		typeReqest string
		want       want
	}{
		{
			name: "пользователь успешно аутентифицирован",
			url:  urlPostUserLogin,
			body: models.User{
				Login:    "test",
				Password: "$2a$10$kte3HgQ6VtHaZSBVc0Cr2OSHQnVL3UB5C0mJLnPVA5W3y.EfNz7rC",
			},
			typeReqest: http.MethodPost,
			want: want{
				code: 200,
			},
		},
		{
			name:       "неверный формат запроса",
			url:        urlPostUserLogin,
			typeReqest: http.MethodPost,
			want: want{
				code: 400,
			},
		},
		{
			name: "неверная пара логин/пароль",
			url:  urlPostUserLogin,
			body: models.User{
				Login:    "test",
				Password: "$2a$10$kte3HgQ6VtHaZSBVc0Cr2OSHQnVL3UB5CsmJLnPVA5W3y.EfNz7rC",
			},
			typeReqest: http.MethodPost,
			want: want{
				code: 401,
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

func TestPostUserOrders(t *testing.T) {
	logger.Init()
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)

	mockDB := &mock.MockDB{
		Users: map[int]map[string]string{
			1: {"id": "1", "login": "test", "password": "$2a$10$kte3HgQ6VtHaZSBVc0Cr2OSHQnVL3UB5C0mJLnPVA5W3y.EfNz7rC"},
			2: {"id": "2", "login": "test2", "password": "$2a$10$kte3HgQ6VtHaZSBVc0Cr2OSHQnVL3UB5C0mJLnPVA5W3z.EfNz7rC"},
		},
		Orders: map[int]map[string]string{
			1: {"number": "3488214672200", "user_id": "1", "status": "NEW"},
			2: {"number": "79927398713", "user_id": "2", "status": "NEW"},
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
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Post(urlPostUserOrders, func(w http.ResponseWriter, r *http.Request) {
			PostUserOrders(w, r, storage)
		})
		r.Get(urlGetUserOrders, func(w http.ResponseWriter, r *http.Request) {
			GetUserOrders(w, r)
		})
		r.Get(urlGetUserBalance, func(w http.ResponseWriter, r *http.Request) {
			GetUserBalance(w, r)
		})
		r.Post(urlPostUserBalanceWithdraw, func(w http.ResponseWriter, r *http.Request) {
			PostUserBalanceWithdraw(w, r)
		})
		r.Get(urlGetUserWithdrawals, func(w http.ResponseWriter, r *http.Request) {
			GetUserWithdrawals(w, r)
		})
	})

	bodyLogin := models.User{
		Login:    "test",
		Password: "$2a$10$kte3HgQ6VtHaZSBVc0Cr2OSHQnVL3UB5C0mJLnPVA5W3y.EfNz7rC",
	}
	bodyJson, _ := json.Marshal(bodyLogin)
	req := httptest.NewRequest(http.MethodPost, urlPostUserLogin, bytes.NewReader(bodyJson))
	req.Header.Set("Accept", "application/json")

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	jwtTok := w.Header().Get("Authorization")

	type want struct {
		code int
	}
	tests := []struct {
		name       string
		url        string
		body       string
		jwtToken   string
		typeReqest string
		want       want
	}{
		{
			name:       "номер заказа уже был загружен этим пользователем",
			url:        urlPostUserOrders,
			body:       "3488214672200",
			jwtToken:   jwtTok,
			typeReqest: http.MethodPost,
			want: want{
				code: 200,
			},
		},
		{
			name:       "новый номер заказа принят в обработку",
			url:        urlPostUserOrders,
			body:       "17893729974",
			jwtToken:   jwtTok,
			typeReqest: http.MethodPost,
			want: want{
				code: 202,
			},
		},
		// {
		// 	name:       "неверный формат запроса",
		// 	url:        urlPostUserOrders,
		// 	body:       "17893729974",
		// 	jwtToken:   jwtTok,
		// 	typeReqest: http.MethodPost,
		// 	want: want{
		// 		code: 400,
		// 	},
		// },
		{
			name:       "пользователь не аутентифицирован",
			url:        urlPostUserOrders,
			body:       "17893729974",
			typeReqest: http.MethodPost,
			want: want{
				code: 401,
			},
		},
		{
			name:       "номер заказа уже был загружен другим пользователем",
			url:        urlPostUserOrders,
			body:       "79927398713",
			jwtToken:   jwtTok,
			typeReqest: http.MethodPost,
			want: want{
				code: 409,
			},
		},
		{
			name:       "неверный формат номера заказа",
			url:        urlPostUserOrders,
			body:       "17893в729974",
			jwtToken:   jwtTok,
			typeReqest: http.MethodPost,
			want: want{
				code: 422,
			},
		},
		// {
		// 	name:       "внутренняя ошибка сервера",
		// 	url:        urlPostUserOrders,
		// 	body:       "17893в729974",
		// 	typeReqest: http.MethodPost,
		// 	want: want{
		// 		code: 500,
		// 	},
		// },
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(test.typeReqest, test.url, strings.NewReader(test.body))
			w := httptest.NewRecorder()
			req.Header.Set("Authorization", test.jwtToken)

			r.ServeHTTP(w, req)

			if w.Code != test.want.code {
				t.Errorf("expected status OK; got %v", w.Code)
			}

			assert.Equal(t, test.want.code, w.Code)
		})
	}
}
