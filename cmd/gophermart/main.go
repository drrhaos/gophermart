package main

import (
	"flag"
	"net/http"
	"os"

	"gophermart/internal/configure"
	"gophermart/internal/handlers"
	"gophermart/internal/logger"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

const urlPostUserRegister = "/api/user/register"                // регистрация пользователя
const urlPostUserLogin = "/api/user/login"                      // аутентификация пользователя;
const urlPostUserOrders = "/api/user/orders"                    // загрузка пользователем номера заказа для расчёта;
const urlGetUserOrders = "/api/user/orders"                     // получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях;
const urlGetUserBalance = "/api/user/balance"                   // получение текущего баланса счёта баллов лояльности пользователя;
const urlPostUserBalanceWithdraw = "/api/user/balance/withdraw" // запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа;
const urlGetUserWithdrawals = "/api/user/withdrawals"           // получение информации о выводе средств с накопительного счёта пользователем.

var cfg configure.Config

func main() {
	logger.Init()
	ok := cfg.ReadStartParams()
	if !ok {
		flag.PrintDefaults()
		os.Exit(0)
	}

	r := chi.NewRouter()
	r.Use(middleware.Compress(5, "application/json", "text/html"))

	logger.Logger.Info("Сервер запущен", zap.String("адрес", cfg.RunAddress))

	r.Post(urlPostUserRegister, func(w http.ResponseWriter, r *http.Request) {
		handlers.PostUserRegister(w, r)
	})
	r.Post(urlPostUserLogin, func(w http.ResponseWriter, r *http.Request) {
		handlers.PostUserLogin(w, r)
	})
	r.Post(urlPostUserOrders, func(w http.ResponseWriter, r *http.Request) {
		handlers.PostUserOrders(w, r)
	})
	r.Get(urlGetUserOrders, func(w http.ResponseWriter, r *http.Request) {
		handlers.GetUserOrders(w, r)
	})
	r.Get(urlGetUserBalance, func(w http.ResponseWriter, r *http.Request) {
		handlers.GetUserBalance(w, r)
	})
	r.Post(urlPostUserBalanceWithdraw, func(w http.ResponseWriter, r *http.Request) {
		handlers.PostUserBalanceWithdraw(w, r)
	})
	r.Get(urlGetUserWithdrawals, func(w http.ResponseWriter, r *http.Request) {
		handlers.GetUserWithdrawals(w, r)
	})

	if err := http.ListenAndServe(cfg.RunAddress, r); err != nil {
		logger.Logger.Fatal(err.Error())
	}

}
