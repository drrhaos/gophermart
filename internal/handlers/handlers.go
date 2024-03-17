package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"gophermart/internal/logger"
	"gophermart/internal/store"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
)

type User struct {
	Login    string `json:"login"`    // логин
	Password string `json:"password"` // параметр, принимающий значение gauge или counter
}

type StatusOrders struct {
	Number     string `json:"number"`      // номер заказа
	Status     string `json:"status"`      // статус расчёта начисления
	Accrual    int64  `json:"accrual"`     // рассчитанные баллы к начислению, при отсутствии начисления — поле отсутствует в ответе.
	UploadedAt string `json:"uploaded_at"` // временЯ загрузки, формат даты — RFC3339.
}

type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type BalanceWithdrawn struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

type BalanceWithdrawals struct {
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"` // временЯ загрузки, формат даты — RFC3339.
}

// PostUserRegister Регистрация пользователя
// @Summary Регистрация пользователя
// @Description Этот эндпоинт производит регистрацию пользователя
// @Accept json
// @Produce json
// @Param request body User true "JSON тело запроса"
// @Success 200 {string}  string    "пользователь успешно аутентифицирован"
// @Failure 400 {string}  string    "неверный формат запроса"
// @Failure 404 {string}  string    "неверная пара логин/пароль"
// @Failure 500 {string}  string    "внутренняя ошибка сервера"
// @Router /api/user/register [post]
func PostUserRegister(res http.ResponseWriter, req *http.Request, storage *store.StorageContext, tokenAuth *jwtauth.JWTAuth) {
	ctx, cancel := context.WithTimeout(req.Context(), 5*time.Second)
	defer cancel()

	var user User
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(buf.Bytes(), &user)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if user.Login == "" || user.Password == "" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	err = storage.UserRegister(ctx, user.Login, user.Password)
	if err != nil && errors.Is(err, store.ErrLoginDuplicate) {
		res.WriteHeader(http.StatusConflict)
		return
	} else if err != nil && errors.Is(err, store.ErrLoginDuplicate) {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	claims := jwt.MapClaims{
		"username": user.Login,
	}

	_, tokenString, err := tokenAuth.Encode(claims)
	if err != nil {
		logger.Logger.Warn("Произошла ошибка генерации токена")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	bearer := "Bearer " + tokenString
	res.Header().Set("Authorization", bearer)
	res.WriteHeader(http.StatusOK)
	logger.Logger.Info("Новый пользователь аутентифицирован")
}

// PostUserLogin Аутентификация пользователя
// @Summary Аутентификация пользователя
// @Description Этот эндпоинт производит аутентификацию пользователя
// @Accept json
// @Produce json
// @Param request body User true "JSON тело запроса"
// @Success 200 {string}  string    "пользователь успешно аутентифицирован"
// @Failure 400 {string}  string    "неверный формат запроса"
// @Failure 401 {string}  string    "неверная пара логин/пароль"
// @Failure 500 {string}  string    "внутренняя ошибка сервера"
// @Router /api/user/login [post]
func PostUserLogin(res http.ResponseWriter, req *http.Request, storage *store.StorageContext, tokenAuth *jwtauth.JWTAuth) {
	ctx, cancel := context.WithTimeout(req.Context(), 30*time.Second)
	defer cancel()

	var user User
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &user); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if user.Login == "" || user.Password == "" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	err = storage.UserLogin(ctx, user.Login, user.Password)
	if err != nil {
		logger.Logger.Warn("Пользователь не прошел проверку")
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	claims := jwt.MapClaims{
		"username": user.Login,
	}

	_, tokenString, err := tokenAuth.Encode(claims)
	if err != nil {
		logger.Logger.Warn("Произошла ошибка генерации токена")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	bearer := "Bearer " + tokenString
	res.Header().Set("Authorization", bearer)
	res.WriteHeader(http.StatusOK)
	logger.Logger.Info("Пользователь аутентифицирован")
}

// PostUserOrders Загрузка номера заказа
// @Summary Загрузка номера заказа
// @Description Этот эндпоинт загружает номера заказа
// @Produce json
// @Success 200 {string} {string}  string    ""
// @Router /api/user/orders [post]
func PostUserOrders(res http.ResponseWriter, req *http.Request, storage *store.StorageContext) {
	// 200 — номер заказа уже был загружен этим пользователем;
	// 202 — новый номер заказа принят в обработку;
	// 400 — неверный формат запроса;
	// 401 — пользователь не аутентифицирован;
	// 409 — номер заказа уже был загружен другим пользователем;
	// 422 — неверный формат номера заказа;
	// 500 — внутренняя ошибка сервера.
	ctx, cancel := context.WithTimeout(req.Context(), 30*time.Microsecond)
	defer cancel()

	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	bodyString := string(bodyBytes)

	order, err := strconv.Atoi(bodyString)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	exist := storage.UserOrders(ctx, order)
	if exist {
		res.WriteHeader(http.StatusConflict)
		return
	}

	res.WriteHeader(http.StatusOK)
}

// GetUserOrders Получение списка загруженных номеров заказов
// @Summary Получение списка загруженных номеров заказов
// @Description Этот эндпоинт для получения списка загруженных номеров заказов
// @Produce json
// @Success 200 {string}  string    ""
// @Router /api/user/orders [get]
func GetUserOrders(res http.ResponseWriter, req *http.Request) {
	// 200 — успешная обработка запроса.
	// 204 — нет данных для ответа.
	// 401 — пользователь не авторизован.
	// 500 — внутренняя ошибка сервера.

	res.WriteHeader(http.StatusOK)
}

// GetUserBalance Получение текущего баланса пользователя
// @Summary Получение текущего баланса пользователя
// @Description Этот эндпоинт для получение текущего баланса пользователя
// @Produce json
// @Success 200 {string} string    ""
// @Router /api/user/balance [get]
func GetUserBalance(res http.ResponseWriter, req *http.Request) {
	// 200 — успешная обработка запроса.
	// 401 — пользователь не авторизован.
	// 500 — внутренняя ошибка сервера.

	res.WriteHeader(http.StatusOK)
}

// PostUserBalanceWithdraw Запрос на списание средств
// @Summary Запрос на списание средств
// @Description Этот эндпоинт на списание средств
// @Produce json
// @Success 200 {string}  string    ""
// @Router /api/user/balance/withdraw [post]
func PostUserBalanceWithdraw(res http.ResponseWriter, req *http.Request) {
	// 200 — успешная обработка запроса;
	// 401 — пользователь не авторизован;
	// 402 — на счету недостаточно средств;
	// 422 — неверный номер заказа;
	// 500 — внутренняя ошибка сервера.

	res.WriteHeader(http.StatusOK)
}

// GetUserBalance Получение информации о выводе средств
// @Summary Получение информации о выводе средств
// @Description Этот эндпоинт для получение информации о выводе средств
// @Produce json
// @Success 200 {string}  string    ""
// @Router /api/user/withdrawals [get]
func GetUserWithdrawals(res http.ResponseWriter, req *http.Request) {
	// 200 — успешная обработка запроса.
	// 204 — нет ни одного списания.
	// 401 — пользователь не авторизован.
	// 500 — внутренняя ошибка сервер

	res.WriteHeader(http.StatusOK)
}
