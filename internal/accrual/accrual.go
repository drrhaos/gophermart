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
	"strconv"
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
	ctx := context.Background()
	for job := range jobs {
		logger.Logger.Info(fmt.Sprintf("Воркер %d", workerID))

		statusOrder := GetStatus(ctx, job, urlAccrual)
		if statusOrder != nil {
			err := storage.UpdateStatusOrders(ctx, statusOrder)
			if err != nil {
				logger.Logger.Warn("Ошибка обновления данных", zap.Error(err))
			}
		}

	}
}

func GetStatus(ctx context.Context, number int64, urlAccrual string) *models.StatusOrdersAccrual {
	var statusOrders *models.StatusOrdersAccrual

	client := &http.Client{}
	url := fmt.Sprintf(urlGetUserOrders, urlAccrual, number)
	r, _ := http.NewRequest(http.MethodGet, url, nil)
	r = r.WithContext(ctx)
	resp, err := client.Do(r)
	if err != nil {
		logger.Logger.Warn("ошибка запроса", zap.Error(err))
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Logger.Warn("не удалось прочитать данные", zap.Error(err))
		return nil
	}

	switch resp.StatusCode {
	case 204:
		logger.Logger.Warn("заказ не зарегистрирован в системе расчёта")
		return nil
	case 429:
		timeSleepStr := resp.Header.Get("Retry-After")
		timeSleep, err := strconv.ParseInt(timeSleepStr, 10, 64)
		if err != nil {
			logger.Logger.Warn("не удалось задать таймер")
			return nil
		}
		logger.Logger.Warn(fmt.Sprintf("превышено количество запросов к сервису, ожидание %s", timeSleepStr))
		time.Sleep(time.Duration(timeSleep))
		return nil
	case 500:
		logger.Logger.Warn("внутренняя ошибка сервера системы расчёта начислений баллов лояльности")
		return nil
	}

	err = json.Unmarshal(body, &statusOrders)
	if err != nil {
		logger.Logger.Warn("не удалось распорсить запрос", zap.Error(err))
		return nil
	}
	return statusOrders
}
