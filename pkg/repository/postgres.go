package repository

import (
	"avito_test_task/pkg"
	"avito_test_task/pkg/entities"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"time"
)

type PostgresClient struct {
	*sqlx.DB
}

func NewPostgresClient(config *pkg.Config) (*PostgresClient, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.DBPort, config.User, config.Password, config.Name)
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresClient{db}, nil
}

func (p *PostgresClient) Add(id int, amount float64) error {
	query := "INSERT INTO balance (user_id, balance)" +
		"VALUES ($1, $2)" +
		"ON CONFLICT (user_id) DO UPDATE SET balance = (balance.balance + EXCLUDED.balance);"
	_, err := p.Exec(query, id, amount)
	if err != nil {
		return err
	}
	return nil

}

func (p *PostgresClient) GetBalance(id int) (float64, error) {
	var balance float64
	if err := p.QueryRow("SELECT balance FROM balance WHERE user_id = $1", id).Scan(&balance); err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("not found")
		}
		return 0, err
	}
	return balance, nil
}

func (p *PostgresClient) Reserve(req *entities.Request) error {
	tx, err := p.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE balance SET balance = balance - $1 WHERE user_id = $2", req.Amount, req.Id)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	_, err = tx.Exec("INSERT INTO reserve_account( service_id, user_id, order_id, amount) "+
		"VALUES($1, $2, $3, $4) ", req.ServiceID, req.Id, req.OrderID, req.Amount)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	_, err = tx.Exec("INSERT INTO history(order_id, service_id, user_id, amount, date) "+
		"VALUES($1, $2, $3, $4,$5) "+
		"ON CONFLICT (order_id) DO UPDATE SET amount = (history.amount + EXCLUDED.amount)", req.OrderID, req.ServiceID, req.Id, req.Amount, time.Now())
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresClient) Approve(req *entities.Request) error {
	return nil
}
