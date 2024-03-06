package handlers

import "net/http"

func PostUserRegister(res http.ResponseWriter, req *http.Request) {

	res.WriteHeader(http.StatusOK)
}

func PostUserLogin(res http.ResponseWriter, req *http.Request) {

	res.WriteHeader(http.StatusOK)
}

func PostUserOrders(res http.ResponseWriter, req *http.Request) {

	res.WriteHeader(http.StatusOK)
}

func GetUserOrders(res http.ResponseWriter, req *http.Request) {

	res.WriteHeader(http.StatusOK)
}

func GetUserBalance(res http.ResponseWriter, req *http.Request) {

	res.WriteHeader(http.StatusOK)
}

func PostUserBalanceWithdraw(res http.ResponseWriter, req *http.Request) {

	res.WriteHeader(http.StatusOK)
}

func GetUserWithdrawals(res http.ResponseWriter, req *http.Request) {

	res.WriteHeader(http.StatusOK)
}
