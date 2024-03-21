package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"gophermart/internal/logger"
	"gophermart/internal/luhn"
	"gophermart/internal/models"
	"gophermart/internal/store"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
)

// PostUserRegister Регистрация пользователя
// @Summary Регистрация пользователя
// @Description Этот эндпоинт производит регистрацию пользователя
// @Accept json
// @Param request body models.User true "JSON тело запроса"
// @Success 200 {string}  string    "пользователь успешно аутентифицирован"
// @Failure 400 {string}  string    "неверный формат запроса"
// @Failure 409 {string}  string    "логин уже занят"
// @Failure 500 {string}  string    "внутренняя ошибка сервера"
// @Router /api/user/register [post]
func PostUserRegister(res http.ResponseWriter, req *http.Request, storage *store.StorageContext, tokenAuth *jwtauth.JWTAuth) {
	ctx, cancel := context.WithTimeout(req.Context(), 5*time.Second)
	defer cancel()

	var user models.User
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
	if errors.Is(err, store.ErrLoginDuplicate) {
		res.WriteHeader(http.StatusConflict)
		return
	} else if err != nil && !errors.Is(err, store.ErrLoginDuplicate) {
		res.WriteHeader(http.StatusInternalServerError)
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
// @Param request body models.User true "JSON тело запроса"
// @Success 200 {string}  string    "пользователь успешно аутентифицирован"
// @Failure 400 {string}  string    "неверный формат запроса"
// @Failure 401 {string}  string    "неверная пара логин/пароль"
// @Failure 500 {string}  string    "внутренняя ошибка сервера"
// @Router /api/user/login [post]
func PostUserLogin(res http.ResponseWriter, req *http.Request, storage *store.StorageContext, tokenAuth *jwtauth.JWTAuth) {
	ctx, cancel := context.WithTimeout(req.Context(), 30*time.Second)
	defer cancel()

	var user models.User
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

	if errors.Is(err, store.ErrAuthentication) {
		res.WriteHeader(http.StatusUnauthorized)
		return
	} else if err != nil && !errors.Is(err, store.ErrLoginDuplicate) {
		res.WriteHeader(http.StatusInternalServerError)
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
// @Accept plain
// @Param  request   body      int  true  "номер заказа"
// @Success 200 {string}  string    "номер заказа уже был загружен этим пользователем"
// @Failure 202 {string}  string    "новый номер заказа принят в обработку"
// @Failure 400 {string}  string    "неверный формат запроса"
// @Failure 401 {string}  string    "пользователь не аутентифицирован"
// @Failure 409 {string}  string    "номер заказа уже был загружен другим пользователем"
// @Failure 422 {string}  string    "неверный формат номера заказа"
// @Failure 500 {string}  string    "внутренняя ошибка сервера"
// @Router /api/user/orders [post]
// @Security Bearer
func PostUserOrders(res http.ResponseWriter, req *http.Request, storage *store.StorageContext) {
	ctx, cancel := context.WithTimeout(req.Context(), 30*time.Second)
	defer cancel()

	token, _, err := jwtauth.FromContext(ctx)
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	claims := token.PrivateClaims()
	user := claims["username"].(string)

	body, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	order, err := strconv.Atoi(string(body))

	if err != nil || !luhn.Valid(order) {
		logger.Logger.Info("Номер заказа не прошел проверку")
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	err = storage.UploadUserOrders(ctx, user, order)

	if errors.Is(err, store.ErrDuplicateOrder) {
		res.WriteHeader(http.StatusOK)
		return
	} else if err != nil && errors.Is(err, store.ErrDuplicateOrderOtherUser) {
		res.WriteHeader(http.StatusConflict)
		return

	} else if err != nil && !errors.Is(err, store.ErrLoginDuplicate) {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusAccepted)
}

// GetUserOrders Получение списка загруженных номеров заказов
// @Summary Получение списка загруженных номеров заказов
// @Description Этот эндпоинт для получения списка загруженных номеров заказов
// @Produce      json
// @Success 200 {string}  string    "успешная обработка запроса"
// @Failure 204 {string}  string    "нет данных для ответа."
// @Failure 401 {string}  string    "пользователь не авторизован"
// @Failure 500 {string}  string    "внутренняя ошибка сервера"
// @Router /api/user/orders [get]
// @Security Bearer
func GetUserOrders(res http.ResponseWriter, req *http.Request, storage *store.StorageContext) {
	ctx, cancel := context.WithTimeout(req.Context(), 30*time.Second)
	defer cancel()
	token, _, err := jwtauth.FromContext(ctx)
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	claims := token.PrivateClaims()
	user := claims["username"].(string)

	ordersUser, err := storage.GetUserOrders(ctx, user)

	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(ordersUser) == 0 {
		res.WriteHeader(http.StatusNoContent)
		return
	}

	jsonBytes, err := json.Marshal(ordersUser)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = res.Write(jsonBytes)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}

// GetUserBalance Получение текущего баланса пользователя
// @Summary Получение текущего баланса пользователя
// @Description Этот эндпоинт для получение текущего баланса пользователя
// @Produce      json
// @Success 200 {string}  string    "успешная обработка запроса"
// @Failure 401 {string}  string    "пользователь не авторизован"
// @Failure 500 {string}  string    "внутренняя ошибка сервера"
// @Router /api/user/balance [get]
// @Security Bearer
func GetUserBalance(res http.ResponseWriter, req *http.Request, storage *store.StorageContext) {
	ctx, cancel := context.WithTimeout(req.Context(), 30*time.Second)
	defer cancel()
	token, _, err := jwtauth.FromContext(ctx)
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	claims := token.PrivateClaims()
	user := claims["username"].(string)

	ordersUser, err := storage.GetUserBalance(ctx, user)

	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonBytes, err := json.Marshal(ordersUser)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = res.Write(jsonBytes)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}

// PostUserBalanceWithdraw Запрос на списание средств
// @Summary Запрос на списание средств
// @Description Этот эндпоинт на списание средств
// @Accept json
// @Param request body models.BalanceWithdrawn true "JSON тело запроса"
// @Success 200 {string}  string    "успешная обработка запроса"
// @Failure 401 {string}  string    "пользователь не авторизован"
// @Failure 422 {string}  string    "неверный номер заказа"
// @Failure 500 {string}  string    "внутренняя ошибка сервера"
// @Router /api/user/balance/withdraw [post]
// @Security Bearer
func PostUserBalanceWithdraw(res http.ResponseWriter, req *http.Request, storage *store.StorageContext) {
	ctx, cancel := context.WithTimeout(req.Context(), 5*time.Second)
	defer cancel()

	token, _, err := jwtauth.FromContext(ctx)
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	claims := token.PrivateClaims()
	user := claims["username"].(string)

	var userBalance models.BalanceWithdrawals
	var buf bytes.Buffer

	_, err = buf.ReadFrom(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(buf.Bytes(), &userBalance)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if userBalance.Order == "" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	// err = storage.UserRegister(ctx, user.Login, user.Password)
	// if errors.Is(err, store.ErrLoginDuplicate) {
	// 	res.WriteHeader(http.StatusConflict)
	// 	return
	// } else if err != nil && !errors.Is(err, store.ErrLoginDuplicate) {
	// 	res.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
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
