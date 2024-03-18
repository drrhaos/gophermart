package pg

import (
	"context"
	"errors"
	"time"

	"gophermart/internal/logger"
	"gophermart/internal/store"

	"github.com/jackc/pgx/v5/pgconn"
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
	err = db.Migrations(ctx)
	if err != nil {
		logger.Logger.Panic("Не удалось подключиться к базе данных", zap.Error(err))
		return nil
	}
	return db
}

func (db *Database) Close() {
	db.Conn.Close()
}

func (db *Database) Migrations(ctx context.Context) error {

	_, err := db.Conn.Exec(ctx,
		`CREATE TABLE IF NOT EXISTS users
		(
			id SERIAL PRIMARY KEY,
			login character(40) NOT NULL,
			password character(64) NOT NULL,
			sum bigint DEFAULT 0,
			registered_at timestamp with time zone,
			last_time timestamp with time zone
		)`)
	if err != nil {
		return err
	}

	_, err = db.Conn.Exec(ctx,
		`CREATE TABLE IF NOT EXISTS orders
		(
			number bigint UNIQUE PRIMARY KEY,
			user_id bigint REFERENCES users(id),
			status character(10) DEFAULT 'NEW', 
			accrual bigint,
			uploaded_at timestamp with time zone
		)`)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) Ping(ctx context.Context) bool {
	if err := db.Conn.Ping(ctx); err != nil {
		return false
	}
	return true
}

func (db *Database) UserRegister(ctx context.Context, login string, password string) error {
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

	var hashedPassword []byte
	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Logger.Warn("Ошибка при хешировании пароля ", zap.Error(err))
		return err
	}

	_, err = db.Conn.Exec(ctx,
		`INSERT INTO users (login, password, registered_at) VALUES ($1, $2, $3)`, login, string(hashedPassword), time.Now())
	if err != nil {
		logger.Logger.Warn("Не удалось добавить пользователя ", zap.Error(err))
		return err
	}
	logger.Logger.Info("Добавлен новый пользователь")
	return nil
}

func (db *Database) UserLogin(ctx context.Context, login string, password string) error {
	var hashedPassword []byte

	err := db.Conn.QueryRow(ctx, `SELECT password FROM users WHERE login = $1`, login).Scan(&hashedPassword)

	if err != nil {
		logger.Logger.Warn("Ошибка выполнения запроса ", zap.Error(err))
		return err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		return store.ErrAuthentication
	}
	return nil
}

func (db *Database) UploadUserOrders(ctx context.Context, user string, order int) error {
	_, err := db.Conn.Exec(ctx,
		`INSERT INTO orders (number, user_id) VALUES ($1,  (SELECT id FROM users WHERE login = $2))`, order, user)

	var duplicateEntryError = &pgconn.PgError{Code: "23505"}
	if err != nil {
		if errors.As(err, &duplicateEntryError) {
			logger.Logger.Warn("Дубликат заказа")
			return store.ErrDuplicateOrder
		} else {
			logger.Logger.Warn("Не удалось добавить пользователя ", zap.Error(err))
			return err
		}
	}
	logger.Logger.Info("Добавлен новый заказ")
	return nil
}
