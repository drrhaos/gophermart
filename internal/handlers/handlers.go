package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"gophermart/internal/store"
	"io"
	"net/http"
	"strconv"
	"time"
)

type AuthorizationData struct {
	Login    string `json:"login"`    // логин
	Password string `json:"password"` // параметр, принимающий значение gauge или counter
}

// PostUserRegister Регистрация пользователя
// @Summary Регистрация пользователя
// @Description Этот эндпоинт производит регистрацию пользователя
// @Produce json
// @Success 200 {string}
// @Router /api/user/register [post]
func PostUserRegister(res http.ResponseWriter, req *http.Request, storage *store.StorageContext) {
	// 200 — пользователь успешно зарегистрирован и аутентифицирован;
	// 400 — неверный формат запроса;
	// 409 — логин уже занят;
	// 500 — внутренняя ошибка сервера.
	ctx, cancel := context.WithTimeout(req.Context(), 30*time.Microsecond)
	defer cancel()

	var authorizationData AuthorizationData
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &authorizationData); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	exist := storage.UserRegister(ctx, authorizationData.Login, authorizationData.Password)
	if exist {
		res.WriteHeader(http.StatusConflict)
		return
	}

	res.WriteHeader(http.StatusOK)
}

// PostUserLogin Аутентификация пользователя
// @Summary Аутентификация пользователя
// @Description Этот эндпоинт производит аутентификацию пользователя
// @Produce json
// @Success 200 {string}
// @Router /api/user/login [post]
func PostUserLogin(res http.ResponseWriter, req *http.Request, storage *store.StorageContext) {
	// 200 — пользователь успешно аутентифицирован;
	// 400 — неверный формат запроса;
	// 401 — неверная пара логин/пароль;
	// 500 — внутренняя ошибка сервера.
	ctx, cancel := context.WithTimeout(req.Context(), 30*time.Microsecond)
	defer cancel()

	var authorizationData AuthorizationData
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &authorizationData); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	exist := storage.UserLogin(ctx, authorizationData.Login, authorizationData.Password)
	if exist {
		res.WriteHeader(http.StatusConflict)
		return
	}

	res.WriteHeader(http.StatusOK)
}

// PostUserOrders Загрузка номера заказа
// @Summary Загрузка номера заказа
// @Description Этот эндпоинт загружает номера заказа
// @Produce json
// @Success 200 {string}
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

// Получение списка загруженных номеров заказов
func GetUserOrders(res http.ResponseWriter, req *http.Request) {

	res.WriteHeader(http.StatusOK)
}

// Получение текущего баланса пользователя
func GetUserBalance(res http.ResponseWriter, req *http.Request) {

	res.WriteHeader(http.StatusOK)
}

// Запрос на списание средств
func PostUserBalanceWithdraw(res http.ResponseWriter, req *http.Request) {

	res.WriteHeader(http.StatusOK)
}

// Получение информации о выводе средств
func GetUserWithdrawals(res http.ResponseWriter, req *http.Request) {

	res.WriteHeader(http.StatusOK)
}
