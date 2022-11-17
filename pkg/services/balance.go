package services

import (
	"github.com/ksean42/BalanceService/pkg/entities"
	"time"
)

type Balance interface {
	Add(id int, amount float64) error
	GetBalance(id int) (float64, error)
	Reserve(req *entities.Request) error
	Approve(req *entities.Request) error
	GetReport(date time.Time) (string, error)
	Transfer(req *entities.TransferRequest) error
	Reject(req *entities.ReserveReject) error
}
