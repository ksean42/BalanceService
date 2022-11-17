package repository

import (
	"github.com/ksean42/BalanceService/pkg/entities"
	"time"
)

type Repository interface {
	Add(id int, amount float64) error
	GetBalance(id int) (float64, error)
	Reserve(req *entities.Request) error
	Approve(req *entities.Request) error
	GetReport(date time.Time) (*[]entities.Report, error)
}
