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
			1: {"id": "1", "login": "test", "password": "$2a$10$kte3HgQ6VtHaZSBVc0Cr2OSHQnVL3UB5C0mJLnPVA5W3y.EfNz7rC", "sum": "10", "withdrawn": "10", "registered_at": "2024-03-19 19:35:17.662533+00"},
			2: {"id": "2", "login": "test2", "password": "$2a$10$kte3HgQ6VtHaZSBVc0Cr2OSHQnVL3UB5C0mJLnPVA5W3z.EfNz7rC", "sum": "0", "withdrawn": "10", "registered_at": "2024-03-19 19:35:17.662533+00"},
			3: {"id": "3", "login": "test3", "password": "$2a$10$kte3HgQ6VtHaZSBVc0Cr2OSHQnVL3UB5C0mJLnPVA5W3z.EfNz7rC", "sum": "0", "withdrawn": "10", "registered_at": "2024-03-19 19:35:17.662533+00"},
		},
		Orders: map[int]map[string]string{
			1: {"number": "7950839220", "user_id": "1", "status": "NEW", "uploaded_at": "2024-03-19 19:35:17.662533+00"},
			2: {"number": "2396508901", "user_id": "1", "status": "PROCESSED", "accrual": "10", "uploaded_at": "2024-03-19 19:35:17.662533+00"},
			3: {"number": "1852074499", "user_id": "2", "status": "PROCESSING", "uploaded_at": "2024-03-19 19:35:17.662533+00"},
			4: {"number": "9347167976", "user_id": "1", "status": "INVALID", "uploaded_at": "2024-03-19 19:35:17.662533+00"},
		},
		Withdrawals: map[int]map[string]string{
			1: {"number": "7950839220", "user_id": "1", "sum": "5", "processed_at": "2024-03-19 19:35:17.662533+00"},
			2: {"number": "2396508901", "user_id": "1", "sum": "1", "processed_at": "2024-03-19 19:35:17.662533+00"},
			3: {"number": "1852074499", "user_id": "2", "sum": "4", "processed_at": "2024-03-19 19:35:17.662533+00"},
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
		GetUserOrders(w, r, storage)
	})
	r.Get(urlGetUserBalance, func(w http.ResponseWriter, r *http.Request) {
		GetUserBalance(w, r, storage)
	})
	r.Post(urlPostUserBalanceWithdraw, func(w http.ResponseWriter, r *http.Request) {
		PostUserBalanceWithdraw(w, r, storage)
	})
	r.Get(urlGetUserWithdrawals, func(w http.ResponseWriter, r *http.Request) {
		GetUserWithdrawals(w, r, storage)
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
				Login:    "test4",
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
			bodyJSON, _ := json.Marshal(test.body)
			req := httptest.NewRequest(test.typeReqest, test.url, bytes.NewReader(bodyJSON))
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
			1: {"id": "1", "login": "test", "password": "$2a$10$kte3HgQ6VtHaZSBVc0Cr2OSHQnVL3UB5C0mJLnPVA5W3y.EfNz7rC", "sum": "10", "withdrawn": "10", "registered_at": "2024-03-19 19:35:17.662533+00"},
			2: {"id": "2", "login": "test2", "password": "$2a$10$kte3HgQ6VtHaZSBVc0Cr2OSHQnVL3UB5C0mJLnPVA5W3z.EfNz7rC", "sum": "0", "withdrawn": "10", "registered_at": "2024-03-19 19:35:17.662533+00"},
			3: {"id": "3", "login": "test3", "password": "$2a$10$kte3HgQ6VtHaZSBVc0Cr2OSHQnVL3UB5C0mJLnPVA5W3z.EfNz7rC", "sum": "0", "withdrawn": "10", "registered_at": "2024-03-19 19:35:17.662533+00"},
		},
		Orders: map[int]map[string]string{
			1: {"number": "7950839220", "user_id": "1", "status": "NEW", "uploaded_at": "2024-03-19 19:35:17.662533+00"},
			2: {"number": "2396508901", "user_id": "1", "status": "PROCESSED", "accrual": "10", "uploaded_at": "2024-03-19 19:35:17.662533+00"},
			3: {"number": "1852074499", "user_id": "2", "status": "PROCESSING", "uploaded_at": "2024-03-19 19:35:17.662533+00"},
			4: {"number": "9347167976", "user_id": "1", "status": "INVALID", "uploaded_at": "2024-03-19 19:35:17.662533+00"},
		},
		Withdrawals: map[int]map[string]string{
			1: {"number": "7950839220", "user_id": "1", "sum": "5", "processed_at": "2024-03-19 19:35:17.662533+00"},
			2: {"number": "2396508901", "user_id": "1", "sum": "1", "processed_at": "2024-03-19 19:35:17.662533+00"},
			3: {"number": "1852074499", "user_id": "2", "sum": "4", "processed_at": "2024-03-19 19:35:17.662533+00"},
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
		GetUserOrders(w, r, storage)
	})
	r.Get(urlGetUserBalance, func(w http.ResponseWriter, r *http.Request) {
		GetUserBalance(w, r, storage)
	})
	r.Post(urlPostUserBalanceWithdraw, func(w http.ResponseWriter, r *http.Request) {
		PostUserBalanceWithdraw(w, r, storage)
	})
	r.Get(urlGetUserWithdrawals, func(w http.ResponseWriter, r *http.Request) {
		GetUserWithdrawals(w, r, storage)
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
			bodyJSON, _ := json.Marshal(test.body)
			req := httptest.NewRequest(test.typeReqest, test.url, bytes.NewReader(bodyJSON))
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
			1: {"id": "1", "login": "test", "password": "$2a$10$kte3HgQ6VtHaZSBVc0Cr2OSHQnVL3UB5C0mJLnPVA5W3y.EfNz7rC", "sum": "10", "withdrawn": "10", "registered_at": "2024-03-19 19:35:17.662533+00"},
			2: {"id": "2", "login": "test2", "password": "$2a$10$kte3HgQ6VtHaZSBVc0Cr2OSHQnVL3UB5C0mJLnPVA5W3z.EfNz7rC", "sum": "0", "withdrawn": "10", "registered_at": "2024-03-19 19:35:17.662533+00"},
			3: {"id": "3", "login": "test3", "password": "$2a$10$kte3HgQ6VtHaZSBVc0Cr2OSHQnVL3UB5C0mJLnPVA5W3z.EfNz7rC", "sum": "0", "withdrawn": "10", "registered_at": "2024-03-19 19:35:17.662533+00"},
		},
		Orders: map[int]map[string]string{
			1: {"number": "7950839220", "user_id": "1", "status": "NEW", "uploaded_at": "2024-03-19 19:35:17.662533+00"},
			2: {"number": "2396508901", "user_id": "1", "status": "PROCESSED", "accrual": "10", "uploaded_at": "2024-03-19 19:35:17.662533+00"},
			3: {"number": "1852074499", "user_id": "2", "status": "PROCESSING", "uploaded_at": "2024-03-19 19:35:17.662533+00"},
			4: {"number": "9347167976", "user_id": "1", "status": "INVALID", "uploaded_at": "2024-03-19 19:35:17.662533+00"},
		},
		Withdrawals: map[int]map[string]string{
			1: {"number": "7950839220", "user_id": "1", "sum": "5", "processed_at": "2024-03-19 19:35:17.662533+00"},
			2: {"number": "2396508901", "user_id": "1", "sum": "1", "processed_at": "2024-03-19 19:35:17.662533+00"},
			3: {"number": "1852074499", "user_id": "2", "sum": "4", "processed_at": "2024-03-19 19:35:17.662533+00"},
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
			GetUserOrders(w, r, storage)
		})
		r.Get(urlGetUserBalance, func(w http.ResponseWriter, r *http.Request) {
			GetUserBalance(w, r, storage)
		})
		r.Post(urlPostUserBalanceWithdraw, func(w http.ResponseWriter, r *http.Request) {
			PostUserBalanceWithdraw(w, r, storage)
		})
		r.Get(urlGetUserWithdrawals, func(w http.ResponseWriter, r *http.Request) {
			GetUserWithdrawals(w, r, storage)
		})
	})

	bodyLogin := models.User{
		Login:    "test",
		Password: "$2a$10$kte3HgQ6VtHaZSBVc0Cr2OSHQnVL3UB5C0mJLnPVA5W3y.EfNz7rC",
	}
	bodyJSON, _ := json.Marshal(bodyLogin)
	req := httptest.NewRequest(http.MethodPost, urlPostUserLogin, bytes.NewReader(bodyJSON))
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
			body:       "7950839220",
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
			body:       "1852074499",
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

func TestGetUserOrders(t *testing.T) {
	logger.Init()
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)

	mockDB := &mock.MockDB{
		Users: map[int]map[string]string{
			1: {"id": "1", "login": "test", "password": "$2a$10$kte3HgQ6VtHaZSBVc0Cr2OSHQnVL3UB5C0mJLnPVA5W3y.EfNz7rC", "sum": "10", "withdrawn": "10", "registered_at": "2024-03-19 19:35:17.662533+00"},
			2: {"id": "2", "login": "test2", "password": "$2a$10$kte3HgQ6VtHaZSBVc0Cr2OSHQnVL3UB5C0mJLnPVA5W3z.EfNz7rC", "sum": "0", "withdrawn": "10", "registered_at": "2024-03-19 19:35:17.662533+00"},
			3: {"id": "3", "login": "test3", "password": "$2a$10$kte3HgQ6VtHaZSBVc0Cr2OSHQnVL3UB5C0mJLnPVA5W3z.EfNz7rC", "sum": "0", "withdrawn": "10", "registered_at": "2024-03-19 19:35:17.662533+00"},
		},
		Orders: map[int]map[string]string{
			1: {"number": "7950839220", "user_id": "1", "status": "NEW", "uploaded_at": "2024-03-19 19:35:17.662533+00"},
			2: {"number": "2396508901", "user_id": "1", "status": "PROCESSED", "accrual": "10", "uploaded_at": "2024-03-19 19:35:17.662533+00"},
			3: {"number": "1852074499", "user_id": "2", "status": "PROCESSING", "uploaded_at": "2024-03-19 19:35:17.662533+00"},
			4: {"number": "9347167976", "user_id": "1", "status": "INVALID", "uploaded_at": "2024-03-19 19:35:17.662533+00"},
		},
		Withdrawals: map[int]map[string]string{
			1: {"number": "7950839220", "user_id": "1", "sum": "5", "processed_at": "2024-03-19 19:35:17.662533+00"},
			2: {"number": "2396508901", "user_id": "1", "sum": "1", "processed_at": "2024-03-19 19:35:17.662533+00"},
			3: {"number": "1852074499", "user_id": "2", "sum": "4", "processed_at": "2024-03-19 19:35:17.662533+00"},
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
			GetUserOrders(w, r, storage)
		})
		r.Get(urlGetUserBalance, func(w http.ResponseWriter, r *http.Request) {
			GetUserBalance(w, r, storage)
		})
		r.Post(urlPostUserBalanceWithdraw, func(w http.ResponseWriter, r *http.Request) {
			PostUserBalanceWithdraw(w, r, storage)
		})
		r.Get(urlGetUserWithdrawals, func(w http.ResponseWriter, r *http.Request) {
			GetUserWithdrawals(w, r, storage)
		})
	})

	bodyLogin := models.User{
		Login:    "test",
		Password: "$2a$10$kte3HgQ6VtHaZSBVc0Cr2OSHQnVL3UB5C0mJLnPVA5W3y.EfNz7rC",
	}
	bodyJSON, _ := json.Marshal(bodyLogin)
	req := httptest.NewRequest(http.MethodPost, urlPostUserLogin, bytes.NewReader(bodyJSON))
	req.Header.Set("Accept", "application/json")

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	jwtTok := w.Header().Get("Authorization")

	bodyLogin2 := models.User{
		Login:    "test3",
		Password: "$2a$10$kte3HgQ6VtHaZSBVc0Cr2OSHQnVL3UB5C0mJLnPVA5W3z.EfNz7rC",
	}
	bodyJSON2, _ := json.Marshal(bodyLogin2)
	req = httptest.NewRequest(http.MethodPost, urlPostUserLogin, bytes.NewReader(bodyJSON2))
	req.Header.Set("Accept", "application/json")

	w = httptest.NewRecorder()

	r.ServeHTTP(w, req)
	jwtTok2 := w.Header().Get("Authorization")

	type want struct {
		code int
	}
	tests := []struct {
		name       string
		url        string
		jwtToken   string
		typeReqest string
		want       want
	}{
		{
			name:       "успешная обработка запроса",
			url:        urlGetUserOrders,
			jwtToken:   jwtTok,
			typeReqest: http.MethodGet,
			want: want{
				code: 200,
			},
		},
		{
			name:       "нет данных для ответа.",
			url:        urlGetUserOrders,
			jwtToken:   jwtTok2,
			typeReqest: http.MethodGet,
			want: want{
				code: 204,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(test.typeReqest, test.url, nil)
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

func TestGetUserBalance(t *testing.T) {
	logger.Init()
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)

	mockDB := &mock.MockDB{
		Users: map[int]map[string]string{
			1: {"id": "1", "login": "test", "password": "$2a$10$kte3HgQ6VtHaZSBVc0Cr2OSHQnVL3UB5C0mJLnPVA5W3y.EfNz7rC", "sum": "10", "withdrawn": "10", "registered_at": "2024-03-19 19:35:17.662533+00"},
			2: {"id": "2", "login": "test2", "password": "$2a$10$kte3HgQ6VtHaZSBVc0Cr2OSHQnVL3UB5C0mJLnPVA5W3z.EfNz7rC", "sum": "0", "withdrawn": "10", "registered_at": "2024-03-19 19:35:17.662533+00"},
			3: {"id": "3", "login": "test3", "password": "$2a$10$kte3HgQ6VtHaZSBVc0Cr2OSHQnVL3UB5C0mJLnPVA5W3z.EfNz7rC", "sum": "0", "withdrawn": "10", "registered_at": "2024-03-19 19:35:17.662533+00"},
		},
		Orders: map[int]map[string]string{
			1: {"number": "7950839220", "user_id": "1", "status": "NEW", "uploaded_at": "2024-03-19 19:35:17.662533+00"},
			2: {"number": "2396508901", "user_id": "1", "status": "PROCESSED", "accrual": "10", "uploaded_at": "2024-03-19 19:35:17.662533+00"},
			3: {"number": "1852074499", "user_id": "2", "status": "PROCESSING", "uploaded_at": "2024-03-19 19:35:17.662533+00"},
			4: {"number": "9347167976", "user_id": "1", "status": "INVALID", "uploaded_at": "2024-03-19 19:35:17.662533+00"},
		},
		Withdrawals: map[int]map[string]string{
			1: {"number": "7950839220", "user_id": "1", "sum": "5", "processed_at": "2024-03-19 19:35:17.662533+00"},
			2: {"number": "2396508901", "user_id": "1", "sum": "1", "processed_at": "2024-03-19 19:35:17.662533+00"},
			3: {"number": "1852074499", "user_id": "2", "sum": "4", "processed_at": "2024-03-19 19:35:17.662533+00"},
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
			GetUserOrders(w, r, storage)
		})
		r.Get(urlGetUserBalance, func(w http.ResponseWriter, r *http.Request) {
			GetUserBalance(w, r, storage)
		})
		r.Post(urlPostUserBalanceWithdraw, func(w http.ResponseWriter, r *http.Request) {
			PostUserBalanceWithdraw(w, r, storage)
		})
		r.Get(urlGetUserWithdrawals, func(w http.ResponseWriter, r *http.Request) {
			GetUserWithdrawals(w, r, storage)
		})
	})

	bodyLogin := models.User{
		Login:    "test",
		Password: "$2a$10$kte3HgQ6VtHaZSBVc0Cr2OSHQnVL3UB5C0mJLnPVA5W3y.EfNz7rC",
	}
	bodyJSON, _ := json.Marshal(bodyLogin)
	req := httptest.NewRequest(http.MethodPost, urlPostUserLogin, bytes.NewReader(bodyJSON))
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
		jwtToken   string
		typeReqest string
		want       want
	}{
		{
			name:       "успешная обработка запроса",
			url:        urlGetUserBalance,
			jwtToken:   jwtTok,
			typeReqest: http.MethodGet,
			want: want{
				code: 200,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(test.typeReqest, test.url, nil)
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

func TestPostUserBalanceWithdraw(t *testing.T) {
	logger.Init()
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)

	mockDB := &mock.MockDB{
		Users: map[int]map[string]string{
			1: {"id": "1", "login": "test", "password": "$2a$10$kte3HgQ6VtHaZSBVc0Cr2OSHQnVL3UB5C0mJLnPVA5W3y.EfNz7rC", "sum": "10", "withdrawn": "10", "registered_at": "2024-03-19 19:35:17.662533+00"},
			2: {"id": "2", "login": "test2", "password": "$2a$10$kte3HgQ6VtHaZSBVc0Cr2OSHQnVL3UB5C0mJLnPVA5W3z.EfNz7rC", "sum": "0", "withdrawn": "10", "registered_at": "2024-03-19 19:35:17.662533+00"},
			3: {"id": "3", "login": "test3", "password": "$2a$10$kte3HgQ6VtHaZSBVc0Cr2OSHQnVL3UB5C0mJLnPVA5W3z.EfNz7rC", "sum": "0", "withdrawn": "10", "registered_at": "2024-03-19 19:35:17.662533+00"},
		},
		Orders: map[int]map[string]string{
			1: {"number": "7950839220", "user_id": "1", "status": "NEW", "uploaded_at": "2024-03-19 19:35:17.662533+00"},
			2: {"number": "2396508901", "user_id": "1", "status": "PROCESSED", "accrual": "10", "uploaded_at": "2024-03-19 19:35:17.662533+00"},
			3: {"number": "1852074499", "user_id": "2", "status": "PROCESSING", "uploaded_at": "2024-03-19 19:35:17.662533+00"},
			4: {"number": "9347167976", "user_id": "1", "status": "INVALID", "uploaded_at": "2024-03-19 19:35:17.662533+00"},
		},
		Withdrawals: map[int]map[string]string{
			1: {"number": "7950839220", "user_id": "1", "sum": "5", "processed_at": "2024-03-19 19:35:17.662533+00"},
			2: {"number": "2396508901", "user_id": "1", "sum": "1", "processed_at": "2024-03-19 19:35:17.662533+00"},
			3: {"number": "1852074499", "user_id": "2", "sum": "4", "processed_at": "2024-03-19 19:35:17.662533+00"},
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
			GetUserOrders(w, r, storage)
		})
		r.Get(urlGetUserBalance, func(w http.ResponseWriter, r *http.Request) {
			GetUserBalance(w, r, storage)
		})
		r.Post(urlPostUserBalanceWithdraw, func(w http.ResponseWriter, r *http.Request) {
			PostUserBalanceWithdraw(w, r, storage)
		})
		r.Get(urlGetUserWithdrawals, func(w http.ResponseWriter, r *http.Request) {
			GetUserWithdrawals(w, r, storage)
		})
	})

	bodyLogin := models.User{
		Login:    "test",
		Password: "$2a$10$kte3HgQ6VtHaZSBVc0Cr2OSHQnVL3UB5C0mJLnPVA5W3y.EfNz7rC",
	}
	bodyJSON, _ := json.Marshal(bodyLogin)
	req := httptest.NewRequest(http.MethodPost, urlPostUserLogin, bytes.NewReader(bodyJSON))
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
		body       models.BalanceWithdrawn
		jwtToken   string
		typeReqest string
		want       want
	}{
		{
			name: "успешная обработка запроса",
			url:  urlPostUserBalanceWithdraw,
			body: models.BalanceWithdrawn{
				Order: "8593379475",
				Sum:   1,
			},
			jwtToken:   jwtTok,
			typeReqest: http.MethodPost,
			want: want{
				code: 200,
			},
		},
		{
			name: "на счету недостаточно средств",
			url:  urlPostUserBalanceWithdraw,
			body: models.BalanceWithdrawn{
				Order: "8593379475",
				Sum:   100,
			},
			jwtToken:   jwtTok,
			typeReqest: http.MethodPost,
			want: want{
				code: 402,
			},
		},
		{
			name: "неверный номер заказа",
			url:  urlPostUserBalanceWithdraw,
			body: models.BalanceWithdrawn{
				Order: "7950830",
				Sum:   1,
			},
			jwtToken:   jwtTok,
			typeReqest: http.MethodPost,
			want: want{
				code: 422,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bodyJSON, _ := json.Marshal(test.body)
			req := httptest.NewRequest(test.typeReqest, test.url, bytes.NewReader(bodyJSON))

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

func TestGetUserWithdrawals(t *testing.T) {
	logger.Init()
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)

	mockDB := &mock.MockDB{
		Users: map[int]map[string]string{
			1: {"id": "1", "login": "test", "password": "$2a$10$kte3HgQ6VtHaZSBVc0Cr2OSHQnVL3UB5C0mJLnPVA5W3y.EfNz7rC", "sum": "10", "withdrawn": "6", "registered_at": "2024-03-19 19:35:17.662533+00"},
			2: {"id": "2", "login": "test2", "password": "$2a$10$kte3HgQ6VtHaZSBVc0Cr2OSHQnVL3UB5C0mJLnPVA5W3z.EfNz7rC", "sum": "0", "withdrawn": "4", "registered_at": "2024-03-19 19:35:17.662533+00"},
			3: {"id": "3", "login": "test3", "password": "$2a$10$kte3HgQ6VtHaZSBVc0Cr2OSHQnVL3UB5C0mJLnPVA5W3y.EfNz7rC", "sum": "0", "withdrawn": "0", "registered_at": "2024-03-19 19:35:17.662533+00"},
		},
		Orders: map[int]map[string]string{
			1: {"number": "7950839220", "user_id": "1", "status": "NEW", "uploaded_at": "2024-03-19 19:35:17.662533+00"},
			2: {"number": "2396508901", "user_id": "1", "status": "PROCESSED", "accrual": "10", "uploaded_at": "2024-03-19 19:35:17.662533+00"},
			3: {"number": "1852074499", "user_id": "2", "status": "PROCESSING", "uploaded_at": "2024-03-19 19:35:17.662533+00"},
			4: {"number": "9347167976", "user_id": "1", "status": "INVALID", "uploaded_at": "2024-03-19 19:35:17.662533+00"},
		},
		Withdrawals: map[int]map[string]string{
			1: {"number": "7950839220", "user_id": "1", "sum": "5", "processed_at": "2024-03-19 19:35:17.662533+00"},
			2: {"number": "2396508901", "user_id": "1", "sum": "1", "processed_at": "2024-03-19 19:35:17.662533+00"},
			3: {"number": "1852074499", "user_id": "2", "sum": "4", "processed_at": "2024-03-19 19:35:17.662533+00"},
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
			GetUserOrders(w, r, storage)
		})
		r.Get(urlGetUserBalance, func(w http.ResponseWriter, r *http.Request) {
			GetUserBalance(w, r, storage)
		})
		r.Post(urlPostUserBalanceWithdraw, func(w http.ResponseWriter, r *http.Request) {
			PostUserBalanceWithdraw(w, r, storage)
		})
		r.Get(urlGetUserWithdrawals, func(w http.ResponseWriter, r *http.Request) {
			GetUserWithdrawals(w, r, storage)
		})
	})

	bodyLogin := models.User{
		Login:    "test",
		Password: "$2a$10$kte3HgQ6VtHaZSBVc0Cr2OSHQnVL3UB5C0mJLnPVA5W3y.EfNz7rC",
	}
	bodyJSON, _ := json.Marshal(bodyLogin)
	req := httptest.NewRequest(http.MethodPost, urlPostUserLogin, bytes.NewReader(bodyJSON))
	req.Header.Set("Accept", "application/json")

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	jwtTok := w.Header().Get("Authorization")

	bodyLogin = models.User{
		Login:    "test3",
		Password: "$2a$10$kte3HgQ6VtHaZSBVc0Cr2OSHQnVL3UB5C0mJLnPVA5W3y.EfNz7rC",
	}
	bodyJSON, _ = json.Marshal(bodyLogin)
	req = httptest.NewRequest(http.MethodPost, urlPostUserLogin, bytes.NewReader(bodyJSON))
	req.Header.Set("Accept", "application/json")

	w = httptest.NewRecorder()

	r.ServeHTTP(w, req)
	jwtTok3 := w.Header().Get("Authorization")

	type want struct {
		code int
	}
	tests := []struct {
		name       string
		url        string
		jwtToken   string
		typeReqest string
		want       want
	}{
		{
			name:       "успешная обработка запроса",
			url:        urlGetUserWithdrawals,
			jwtToken:   jwtTok,
			typeReqest: http.MethodGet,
			want: want{
				code: 200,
			},
		},
		{
			name:       "нет ни одного списания",
			url:        urlGetUserWithdrawals,
			jwtToken:   jwtTok3,
			typeReqest: http.MethodGet,
			want: want{
				code: 204,
			},
		},
		{
			name:       "пользователь не авторизован",
			url:        urlGetUserWithdrawals,
			typeReqest: http.MethodGet,
			want: want{
				code: 401,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(test.typeReqest, test.url, nil)

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
