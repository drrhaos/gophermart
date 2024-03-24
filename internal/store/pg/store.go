package pg

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"gophermart/internal/logger"
	"gophermart/internal/models"
	"gophermart/internal/store"

	"github.com/jackc/pgx/v5"
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
			login varchar(40) NOT NULL,
			password varchar(64) NOT NULL,
			sum float DEFAULT 0,
			withdrawn float DEFAULT 0,
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
			status varchar(10) DEFAULT 'NEW', 
			accrual float,
			uploaded_at timestamp with time zone
		)`)
	if err != nil {
		return err
	}

	_, err = db.Conn.Exec(ctx,
		`CREATE TABLE IF NOT EXISTS withdrawals
		(
			number bigint UNIQUE PRIMARY KEY REFERENCES orders(number),
			sum float,
			processed_at timestamp with time zone
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

	if err == pgx.ErrNoRows {
		return store.ErrAuthentication
	} else if err != nil {
		logger.Logger.Warn("Ошибка выполнения запроса ", zap.Error(err))
		return err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		return store.ErrAuthentication
	}
	return nil
}

func (db *Database) UploadUserOrders(ctx context.Context, login string, order int) error {
	var idUser int
	err := db.Conn.QueryRow(ctx, `SELECT id FROM users WHERE login = $1`, login).Scan(&idUser)
	if err != nil && err != pgx.ErrNoRows {
		logger.Logger.Warn("Ошибка выполнения запроса id", zap.Error(err))
		return err
	}

	var countUser int
	err = db.Conn.QueryRow(ctx, `SELECT COUNT(user_id) FROM orders WHERE number = $1 AND user_id <> $2`, order, idUser).Scan(&countUser)

	if err != nil && err != pgx.ErrNoRows {
		logger.Logger.Warn("Ошибка выполнения запроса user id", zap.Error(err))
		return err
	}

	if countUser > 0 {
		logger.Logger.Warn("Этот заказ добавлен для другого пользователя")
		return store.ErrDuplicateOrderOtherUser
	}

	_, err = db.Conn.Exec(ctx,
		`INSERT INTO orders (number, user_id, uploaded_at) VALUES ($1, $2, $3)`, order, idUser, time.Now())

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

func (db *Database) GetUserOrders(ctx context.Context, login string) ([]models.StatusOrders, error) {
	var orderUser models.StatusOrders
	var ordersUser []models.StatusOrders
	rows, err := db.Conn.Query(ctx, `SELECT number,status,accrual,uploaded_at FROM orders WHERE user_id = (SELECT id FROM users WHERE login = $1) ORDER BY uploaded_at DESC`, login)
	if err != nil {
		logger.Logger.Warn("Ошибка выполнения запроса ", zap.Error(err))
		return ordersUser, err
	}

	defer rows.Close()

	for rows.Next() {
		var accural sql.NullFloat64
		err = rows.Scan(&orderUser.Number, &orderUser.Status, &accural, &orderUser.UploadedAt)
		if err != nil {
			logger.Logger.Warn("Ошибка при сканировании строки:", zap.Error(err))
			return ordersUser, err
		}
		if accural.Valid {
			orderUser.Accrual = accural.Float64
		}
		ordersUser = append(ordersUser, orderUser)
	}

	return ordersUser, nil
}

func (db *Database) GetUserBalance(ctx context.Context, login string) (models.Balance, error) {
	var userBalance models.Balance
	err := db.Conn.QueryRow(ctx, `SELECT sum,withdrawn FROM users WHERE login = $1`, login).Scan(&userBalance.Current, &userBalance.Withdrawn)
	if err != nil {
		logger.Logger.Warn("Ошибка выполнения запроса ", zap.Error(err))
		return userBalance, err
	}
	return userBalance, nil
}

func (db *Database) UpdateUserBalanceWithdraw(ctx context.Context, login string, order string, sum float64) error {
	var userId int64
	var balance float64
	var withdrawn float64
	err := db.Conn.QueryRow(ctx, `SELECT id, sum, withdrawn FROM users WHERE login = $1`, login).Scan(&userId, &balance, &withdrawn)
	if err != nil {
		logger.Logger.Warn("Ошибка выполнения запроса ", zap.Error(err))
		return err
	}

	if balance < sum {
		logger.Logger.Warn("на счету недостаточно средств")
		return store.ErrInsufficientFunds
	}

	var countRow int64
	err = db.Conn.QueryRow(ctx, `SELECT COUNT(*) FROM orders WHERE number = $1 AND user_id =$2`, order, userId).Scan(&countRow)
	if err != nil {
		logger.Logger.Warn("Ошибка выполнения запроса ", zap.Error(err))
		return err
	}

	if countRow == 0 {
		logger.Logger.Warn("заказ не существует")
		return store.ErrOrderNotFound
	}

	_, err = db.Conn.Exec(ctx, `INSERT INTO withdrawals (number, sum, processed_at) VALUES ($1, $2, $3) `, order, sum, time.Now())
	if err != nil {
		logger.Logger.Warn("Не удалось добавмить значение", zap.Error(err))
		return err
	}

	_, err = db.Conn.Exec(ctx, `UPDATE users SET sum = $1, withdrawn = $2 WHERE login = $3`, balance-sum, withdrawn+sum, login)
	if err != nil {
		logger.Logger.Warn("Не удалось обновить баланс", zap.Error(err))
		return err
	}

	return nil
}
