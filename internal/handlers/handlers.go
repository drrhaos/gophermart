package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"gophermart/internal/store"
	"net/http"
	"time"
)

type AuthorizationData struct {
	Login    string `json:"login"`    // логин
	Password string `json:"password"` // параметр, принимающий значение gauge или counter
}

// Регистрация пользователя
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

// Аутентификация пользователя
func PostUserLogin(res http.ResponseWriter, req *http.Request) {

	res.WriteHeader(http.StatusOK)
}

// Загрузка номера заказа
func PostUserOrders(res http.ResponseWriter, req *http.Request) {

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
