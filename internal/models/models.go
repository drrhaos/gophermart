package models

import (
	"time"
)

type User struct {
	Login    string `json:"login"`    // логин
	Password string `json:"password"` // параметр, принимающий значение gauge или counter
}

type StatusOrders struct {
	Number     string    `json:"number" db:"number"`                       // номер заказа
	Status     string    `json:"status" db:"status"`                       // статус расчёта начисления
	Accrual    float64   `json:"accrual,omitempty" db:"accrual,omitempty"` // рассчитанные баллы к начислению, при отсутствии начисления — поле отсутствует в ответе.
	UploadedAt time.Time `json:"uploaded_at" db:"uploaded_at"`             // временЯ загрузки, формат даты — RFC3339.
}

type Balance struct {
	Current   float64 `json:"current"`   // текущий баланс пользователя
	Withdrawn float64 `json:"withdrawn"` // сумма использованных за весь период баллов
}

type BalanceWithdrawn struct {
	Order string  `json:"order"` // номер заказа
	Sum   float64 `json:"sum"`   // сумма списания
}

type BalanceWithdrawals struct {
	Order       string    `json:"order" db:"order"`               // номер заказа
	Sum         float64   `json:"sum" db:"sum"`                   // сумма вывода средств
	ProcessedAt time.Time `json:"processed_at" db:"processed_at"` // временя загрузки, формат даты — RFC3339.
}
