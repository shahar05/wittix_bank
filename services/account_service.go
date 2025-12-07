package services

import (
	"wittix_bank/models"
	"wittix_bank/repository"
)

type AccountService struct {
	PropRepo *repository.AccountRepository
	// RecordsRepo *repository.RecordRepository
}

func (s *AccountService) CreateAccount(p *models.AccountRequest) error {
	return s.PropRepo.Create(p)
}

func (s *AccountService) GetAccountBalance(id string) (float64, error) {
	return s.PropRepo.GetBalance(id)
}
