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
	Accrual    int64     `json:"accrual,omitempty" db:"accrual,omitempty"` // рассчитанные баллы к начислению, при отсутствии начисления — поле отсутствует в ответе.
	UploadedAt time.Time `json:"uploaded_at" db:"uploaded_at"`             // временЯ загрузки, формат даты — RFC3339.
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
