package main

import (
	"flag"
	"net/http"
	"os"

	_ "gophermart/docs"
	"gophermart/internal/configure"
	"gophermart/internal/handlers"
	"gophermart/internal/logger"
	"gophermart/internal/store"
	"gophermart/internal/store/pg"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	httpSwagger "github.com/swaggo/http-swagger"
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

var tokenAuth *jwtauth.JWTAuth

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {

	logger.Init()
	ok := cfg.ReadStartParams()
	if !ok {
		flag.PrintDefaults()
		os.Exit(0)
	}

	storage := &store.StorageContext{}
	storage.SetStorage(pg.NewDatabase(cfg.DatabaseURI))

	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)

	r := chi.NewRouter()
	r.Use(middleware.Compress(5, "application/json", "text/html"))

	logger.Logger.Info("Сервер запущен", zap.String("адрес", cfg.RunAddress))

	r.Mount("/swagger", httpSwagger.Handler())
	r.Post(urlPostUserRegister, func(w http.ResponseWriter, r *http.Request) {
		handlers.PostUserRegister(w, r, storage, tokenAuth)
	})
	r.Post(urlPostUserLogin, func(w http.ResponseWriter, r *http.Request) {
		handlers.PostUserLogin(w, r, storage, tokenAuth)
	})
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Post(urlPostUserOrders, func(w http.ResponseWriter, r *http.Request) {
			handlers.PostUserOrders(w, r, storage)
		})
		r.Get(urlGetUserOrders, func(w http.ResponseWriter, r *http.Request) {
			handlers.GetUserOrders(w, r, storage)
		})
		r.Get(urlGetUserBalance, func(w http.ResponseWriter, r *http.Request) {
			handlers.GetUserBalance(w, r, storage)
		})
		r.Post(urlPostUserBalanceWithdraw, func(w http.ResponseWriter, r *http.Request) {
			handlers.PostUserBalanceWithdraw(w, r, storage)
		})
		r.Get(urlGetUserWithdrawals, func(w http.ResponseWriter, r *http.Request) {
			handlers.GetUserWithdrawals(w, r, storage)
		})
	})

	if err := http.ListenAndServe(cfg.RunAddress, r); err != nil {
		logger.Logger.Fatal(err.Error())
	}

}
