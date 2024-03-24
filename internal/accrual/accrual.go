package accrual

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gophermart/internal/logger"
	"gophermart/internal/models"
	"gophermart/internal/store"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

const urlGetUserOrders = "%s/api/orders/%d" // получение информации о расчёте начислений баллов лояльности

var ErrStatusNoContent = errors.New("StatusNoContent")
var ErrStatusTooManyRequests = errors.New("StatusTooManyRequests")
var ErrStatusInternalServerError = errors.New("StatusInternalServerError")

func PrepareBatch(storage *store.StorageContext) (statusOrders []int64) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	statusOrders, err := storage.GetOrdersProcessing(ctx)
	if err != nil {
		logger.Logger.Warn("Ошибка получения данных о заказах")
		return statusOrders
	}
	logger.Logger.Info(fmt.Sprintf("%d", statusOrders))
	return statusOrders
}

func UpdateStatusOrdersWorker(workerID int, storage *store.StorageContext, urlAccrual string, jobs <-chan int64) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	for job := range jobs {
		logger.Logger.Info(fmt.Sprintf("Воркер %d", workerID))

		statusOrder, ok := GetStatus(ctx, job, urlAccrual)
		if ok {
			storage.UpdateStatusOrders(ctx, statusOrder)
		}

	}
}

func GetStatus(ctx context.Context, number int64, urlAccrual string) (*models.StatusOrders, bool) {
	var statusOrders *models.StatusOrders

	client := &http.Client{}
	url := fmt.Sprintf(urlGetUserOrders, urlAccrual, number)
	r, _ := http.NewRequest(http.MethodGet, url, nil)
	r = r.WithContext(ctx)
	resp, err := client.Do(r)
	if err != nil {
		logger.Logger.Warn("ошибка запроса", zap.Error(err))
		return nil, false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Logger.Warn("не удалось прочитать данные", zap.Error(err))
		return nil, false
	}

	switch resp.StatusCode {
	case 204:
		logger.Logger.Warn("заказ не зарегистрирован в системе расчёта")
		return nil, false
	case 429:
		logger.Logger.Warn("превышено количество запросов к сервису")
		return nil, false
	case 500:
		logger.Logger.Warn("внутренняя ошибка сервера системы расчёта начислений баллов лояльности")
		return nil, false
	}

	err = json.Unmarshal(body, &statusOrders)
	if err != nil {
		logger.Logger.Warn("не удалось распорсить запрос", zap.Error(err))
		return nil, false
	}
	return statusOrders, true
}
