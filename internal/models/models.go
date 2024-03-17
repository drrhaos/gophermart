package models

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
