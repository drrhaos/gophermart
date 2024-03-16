package pg

import (
	"context"
	"errors"
	"time"

	"gophermart/internal/logger"
	"gophermart/internal/store"

	"github.com/avast/retry-go"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
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

func (db *Database) UserRegister(ctx context.Context, login string, password string) error {
	err := retry.Do(
		func() error {
			var countRow int64
			err := db.Conn.QueryRow(ctx, `SELECT COUNT(login) FROM users WHERE login = $1`, login).Scan(&countRow)

			if err != nil {
				logger.Logger.Warn("Ошибка выполнения запроса ", zap.Error(err))
				return err
			}

			if countRow != 0 {
				logger.Logger.Warn("Пользователь существует")
				return store.ErrLoginDuplicate
			}
			return nil
		},
		retry.Attempts(3),
		retry.RetryIf(func(err error) bool {
			return !errors.Is(err, store.ErrLoginDuplicate)
		}),
	)
	if err != nil {
		return err
	}

	var hashedPassword []byte
	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Logger.Warn("Ошибка при хешировании пароля ", zap.Error(err))
		return err
	}

	err = retry.Do(
		func() error {
			_, err := db.Conn.Exec(ctx,
				`INSERT INTO users (login, password, date_time) VALUES ($1, $2, $3)`, login, string(hashedPassword), time.Now())
			if err != nil {
				logger.Logger.Warn("Не удалось добавить пользователя ", zap.Error(err))
				return err
			}
			logger.Logger.Info("Добавлен новый пользователь")
			return nil
		},
		retry.Attempts(3),
	)

	if err != nil {
		return err
	}
	return nil
}

func (db *Database) UserLogin(ctx context.Context, login string, password string) error {
	var hashedPassword []byte
	err := retry.Do(
		func() error {
			err := db.Conn.QueryRow(ctx, `SELECT password FROM users WHERE login = $1`, login).Scan(&hashedPassword)

			if err != nil {
				logger.Logger.Warn("Ошибка выполнения запроса ", zap.Error(err))
				return err
			}
			return nil
		},
		retry.Attempts(3),
		retry.RetryIf(func(err error) bool {
			return !errors.Is(err, store.ErrAuthentication)
		}),
	)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) UserOrders(ctx context.Context, order int) bool {

	return true
}
