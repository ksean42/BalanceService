package repository

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/ksean42/BalanceService/pkg"
	"github.com/ksean42/BalanceService/pkg/entities"
	_ "github.com/lib/pq"
	"log"
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
	var exist bool

	if err := p.QueryRow("select exists(select * from history where order_id=$1)", req.OrderID).Scan(&exist); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}
	if exist {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return fmt.Errorf("order already exist")
	}
	_, err = tx.Exec("INSERT INTO reserve_account(order_id, service_id, user_id, amount) "+
		"VALUES($1, $2, $3, $4) ", req.OrderID, req.ServiceID, req.Id, req.Amount)
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
	tx, err := p.Begin()
	if err != nil {
		return err
	}
	var reservedAmount float64
	err = tx.QueryRow("SELECT amount FROM reserve_account"+
		" WHERE order_id = $1 and service_id = $2 and user_id = $3", req.OrderID, req.ServiceID, req.Id).Scan(&reservedAmount)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		if err == sql.ErrNoRows {
			return fmt.Errorf("not found")
		}
		return err
	}
	if reservedAmount != req.Amount {
		return fmt.Errorf("amount is incorrect")
	}

	_, err = tx.Exec("DELETE FROM reserve_account "+
		"WHERE order_id = $1", req.OrderID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	_, err = tx.Exec("INSERT INTO history(order_id, service_id, user_id, amount, date) "+
		"VALUES($1, $2, $3, $4,$5) ", req.OrderID, req.ServiceID, req.Id, req.Amount, time.Now())
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

func (p *PostgresClient) GetReport(date time.Time) (*[]entities.Report, error) {
	month := date.Format("2006-01-02")
	nextMonth := date.AddDate(0, 1, -1).Format("2006-01-02")

	res, err := p.Query("SELECT service_id, sum(amount) "+
		"FROM history "+
		"WHERE date BETWEEN $1 AND $2 "+
		"GROUP BY service_id", month, nextMonth)
	if err != nil {
		log.Println("query", err)
		return nil, err
	}
	defer res.Close()
	var report []entities.Report
	row := entities.Report{}
	for res.Next() {
		er := res.Scan(&row.ServiceID, &row.Revenue)
		if er != nil {
			log.Println(er)
			return nil, er
		}
		report = append(report, row)
	}
	if len(report) == 0 {
		return nil, fmt.Errorf("there is no data for this date")
	}
	return &report, nil
}

func (p *PostgresClient) Transfer(req *entities.TransferRequest) error {
	tx, err := p.Begin()
	if err != nil {
		return err
	}
	subQuery := "UPDATE balance SET balance = balance + $1 " +
		"WHERE user_id = $2;"
	_, err = tx.Exec(subQuery, req.Amount, req.SrcID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	addQuery := "UPDATE balance SET balance = balance + $1 " +
		"WHERE user_id = $2;"
	_, err = tx.Exec(addQuery, req.Amount, req.DestID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	transactionQuery := "INSERT INTO user_history(user_src, user_dest, amount, date) " +
		"VALUES($1, $2, $3, $4)"
	_, err = tx.Exec(transactionQuery, req.SrcID, req.DestID, req.Amount, time.Now())
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

func (p *PostgresClient) Reject(req *entities.ReserveReject) error {
	var amount float64
	if err := p.QueryRow("SELECT amount FROM reserve_account "+
		"WHERE user_id = $1 and order_id = $2", req.Id, req.OrderId).Scan(&amount); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("reserving is not found")
		}
		return err
	}
	tx, err := p.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM reserve_account "+
		"WHERE order_id = $1 and user_id = $2", req.OrderId, req.Id)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	addQuery := "UPDATE balance SET balance = balance + $1 " +
		"WHERE user_id = $2;"
	_, err = tx.Exec(addQuery, amount, req.Id)
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

func (p *PostgresClient) GetUserReport(id int) (*entities.UserReport, error) {
	// Transactions
	transactionsRes, err := p.Query("SELECT service_id, order_id, amount, date "+
		"FROM history "+
		"WHERE user_id = $1 "+
		"ORDER BY date", id)
	if err != nil {
		return nil, err
	}
	defer transactionsRes.Close()
	var transactionReport []entities.UserTransactionReport
	transaction := entities.UserTransactionReport{}
	for transactionsRes.Next() {
		er := transactionsRes.Scan(&transaction.ServiceID, &transaction.OrderID, &transaction.Amount, &transaction.Date)
		if er != nil {
			return nil, er
		}
		transactionReport = append(transactionReport, transaction)
	}
	// Transfers
	transferRes, err := p.Query("SELECT user_src, user_dest, amount, date "+
		"FROM user_history "+
		"WHERE user_src = $1 or user_dest = $1 "+
		"ORDER BY date", id)
	if err != nil {
		return nil, err
	}
	defer transferRes.Close()
	var transferReport []entities.UserTransferReport
	transfer := entities.UserTransferReport{}
	for transferRes.Next() {
		er := transferRes.Scan(&transfer.SrcID, &transfer.DestID, &transfer.Amount, &transfer.Date)
		if er != nil {
			return nil, er
		}
		transferReport = append(transferReport, transfer)
	}
	// Reservings
	reserveRes, err := p.Query("SELECT order_id, service_id, amount "+
		"FROM reserve_account "+
		"WHERE user_id = $1 "+
		"ORDER BY amount", id)
	if err != nil {
		return nil, err
	}
	defer reserveRes.Close()
	var reservingReport []entities.UserReserving
	reserving := entities.UserReserving{}
	for reserveRes.Next() {
		er := reserveRes.Scan(&reserving.OrderID, &reserving.ServiceID, &reserving.Amount)
		if er != nil {
			return nil, er
		}
		reservingReport = append(reservingReport, reserving)
	}

	report := &entities.UserReport{
		Transfers:    transferReport,
		Transactions: transactionReport,
		Reserves:     reservingReport,
	}
	if len(report.Transactions) == 0 &&
		len(report.Transfers) == 0 && len(report.Reserves) == 0 {
		return nil, fmt.Errorf("there is no data for this user")
	}

	return report, nil
}
