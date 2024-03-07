package pg

import (
	"context"
	"time"

	"gophermart/internal/logger"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Database struct {
	Conn *pgxpool.Pool
}

func NewDatabase(uri string) *Database {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(uri)
	if err != nil {
		logger.Logger.Panic("Ошибка при парсинге конфигурации:", zap.Error(err))
		return nil
	}
	conn, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		logger.Logger.Panic("Не удалось подключиться к базе данных")
		return nil
	}
	db := &Database{Conn: conn}

	return db
}

func (db *Database) Close() {
	db.Conn.Close()
}

func (db *Database) Ping(ctx context.Context) bool {
	if err := db.Conn.Ping(ctx); err != nil {
		return false
	}
	return true
}

func (db *Database) UserRegister(ctx context.Context, user string, password string) bool {

	return true
}

func (db *Database) UserLogin(ctx context.Context, user string, password string) bool {

	return true
}

func (db *Database) UserOrders(ctx context.Context, order int) bool {

	return true
}
