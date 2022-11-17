package services

import (
	"fmt"
	"github.com/gocarina/gocsv"
	"github.com/ksean42/BalanceService/pkg/entities"
	"github.com/ksean42/BalanceService/pkg/repository"
	"log"
	"os"
	"path/filepath"
	"time"
)

type BalanceService struct {
	repo repository.Repository
}

func NewBalanceService(repo repository.Repository) *BalanceService {
	return &BalanceService{repo}
}

func (b *BalanceService) Add(id int, amount float64) error {
	if id <= 0 || amount < 0 {
		return fmt.Errorf("id or amount is not correct")
	}
	return b.repo.Add(id, amount)
}

func (b *BalanceService) GetBalance(id int) (float64, error) {
	if id <= 0 {
		return 0, fmt.Errorf("id is incorrect")
	}
	return b.repo.GetBalance(id)
}

func (b *BalanceService) Reserve(req *entities.Request) error {
	balance, err := b.GetBalance(req.Id)
	if err != nil {
		return err
	}
	if (balance - req.Amount) < 0 {
		return fmt.Errorf("not enough funds")
	} else if req.Amount <= 0 {
		return fmt.Errorf("amount is incorrect")
	} else if req.ServiceID <= 0 {
		return fmt.Errorf("service id is incorrect")
	} else if req.OrderID <= 0 {
		return fmt.Errorf("order id is incorrect")
	}
	return b.repo.Reserve(req)
}

func (b *BalanceService) Approve(req *entities.Request) error {
	if req.Amount <= 0 {
		return fmt.Errorf("amount is incorrect")
	} else if req.ServiceID <= 0 {
		return fmt.Errorf("service id is incorrect")
	} else if req.OrderID <= 0 {
		return fmt.Errorf("order id is incorrect")
	} else if req.Id <= 0 {
		return fmt.Errorf("user id is incorrect")
	}
	return b.repo.Approve(req)
}
func (b *BalanceService) GetReport(date time.Time) (string, error) {
	report, err := b.repo.GetReport(date)
	if err != nil {
		return "", err
	}
	path, err := writeCSV(report, date)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return path, err
}

func writeCSV(report *[]entities.Report, t time.Time) (string, error) {
	date := t.Format("2006-01")
	path := fmt.Sprintf("reports/revenue_report_%s.csv", date)
	_, err := os.Stat("reports")
	if err != nil {
		err := os.Mkdir("reports", os.ModePerm)
		if err != nil {
			log.Println(err)
			return "", err
		}
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer file.Close()

	if err := gocsv.MarshalFile(report, file); err != nil {
		log.Println(err)
		return "", err
	}
	curPath, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return "", err
	}

	return filepath.Join(curPath, path), nil
}
