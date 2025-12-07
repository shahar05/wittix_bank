package services

import (
	"wittix_bank/models"
	"wittix_bank/repository"
)

type JournalService struct {
	JournalRepo      *repository.JournalRepository
	JournalLinesRepo *repository.JournalLinesRepository
}

func (s *JournalService) CreateJournal(trans *models.TransactionRequest) error {
	return s.JournalRepo.Create(trans)
}

func (s *JournalService) CreateTransaction(trans *models.TransactionRequest) error {
	// Find Journal entry for the transaction
	journal, err := s.JournalRepo.Find(trans.ExternalID)
	if err != nil {
		return s.CreateJournal(trans)
	}
	// Atomic operation to create journal lines
	err = s.JournalLinesRepo.Create2Lines(trans, journal.EntryID)
	if err != nil {
		return err
	}
	// Create Journal Lines for debit and credit
	return nil
}

func (s *JournalService) ReverseTransaction(entryID string) error {
	err := s.JournalLinesRepo.ReverseTransaction(entryID)
	if err != nil {
		return err
	}
	return nil
}
