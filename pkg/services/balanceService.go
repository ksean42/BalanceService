package services

import (
	"avito_test_task/pkg/entities"
	"avito_test_task/pkg/repository"
	"fmt"
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
		return 0, fmt.Errorf("id is not correct")
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
	}
	return b.repo.Reserve(req)
}

func (b *BalanceService) Approve(req *entities.Request) error {
	return b.repo.Approve(req)
}
