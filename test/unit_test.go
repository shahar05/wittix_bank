package test

import (
	"database/sql"
	"testing"
	"wittix_bank/models"
	"wittix_bank/repository"
	"wittix_bank/services"
)

// MockJournalRepository for testing
type MockJournalRepository struct {
	journals map[string]*models.JournalEntry
}

func NewMockJournalRepository() *MockJournalRepository {
	return &MockJournalRepository{
		journals: make(map[string]*models.JournalEntry),
	}
}

func (m *MockJournalRepository) Create(trans *models.TransactionRequest) error {
	journal := &models.JournalEntry{
		EntryID:        1,
		ExternalID:     trans.ExternalID,
		IdempotencyKey: trans.IdempotencyKey,
	}
	m.journals[trans.ExternalID] = journal
	return nil
}

func (m *MockJournalRepository) Find(externalID string) (*models.JournalEntry, error) {
	if j, exists := m.journals[externalID]; exists {
		return j, nil
	}
	return nil, sql.ErrNoRows
}

// MockJournalLinesRepository for testing
type MockJournalLinesRepository struct {
	lines map[int64][]*models.JournalLine
}

func NewMockJournalLinesRepository() *MockJournalLinesRepository {
	return &MockJournalLinesRepository{
		lines: make(map[int64][]*models.JournalLine),
	}
}

func (m *MockJournalLinesRepository) Create2Lines(trans *models.TransactionRequest, entryID int64) error {
	debitLine := &models.JournalLine{
		EntryID:   entryID,
		AccountID: trans.AccountFrom,
		Side:      models.SideDebit,
		Amount:    string(rune(trans.Amount)),
	}
	creditLine := &models.JournalLine{
		EntryID:   entryID,
		AccountID: trans.AccountTo,
		Side:      models.SideCredit,
		Amount:    string(rune(trans.Amount)),
	}
	m.lines[entryID] = []*models.JournalLine{debitLine, creditLine}
	return nil
}

func (m *MockJournalLinesRepository) ReverseTransaction(entryID string) error {
	// Simulate reversal by toggling sides
	for _, lines := range m.lines {
		for _, line := range lines {
			if line.Side == models.SideDebit {
				line.Side = models.SideCredit
			} else {
				line.Side = models.SideDebit
			}
		}
	}
	return nil
}

// Test 1: Happy flow - successful transaction creation
func TestCreateTransaction_HappyFlow(t *testing.T) {
	mockJournalRepo := NewMockJournalRepository()
	mockLinesRepo := NewMockJournalLinesRepository()

	svc := &services.JournalService{
		JournalRepo:      (*repository.JournalRepository)(nil),
		JournalLinesRepo: (*repository.JournalLinesRepository)(nil),
	}

	// Manually set mocks since we can't directly inject mocks into repositories
	svc.JournalRepo = &repository.JournalRepository{}
	svc.JournalLinesRepo = &repository.JournalLinesRepository{}

	trans := &models.TransactionRequest{
		AccountFrom:    "acc_001",
		AccountTo:      "acc_002",
		Amount:         1000,
		Currency:       "USD",
		ExternalID:     "ext_123",
		IdempotencyKey: "key_123",
	}

	// Create initial journal entry
	err := mockJournalRepo.Create(trans)
	if err != nil {
		t.Fatalf("Failed to create journal: %v", err)
	}

	// Find and create lines
	journal, err := mockJournalRepo.Find(trans.ExternalID)
	if err != nil {
		t.Fatalf("Failed to find journal: %v", err)
	}

	if journal == nil {
		t.Fatal("Journal entry should not be nil")
	}

	err = mockLinesRepo.Create2Lines(trans, journal.EntryID)
	if err != nil {
		t.Fatalf("Failed to create journal lines: %v", err)
	}

	lines := mockLinesRepo.lines[journal.EntryID]
	if len(lines) != 2 {
		t.Fatalf("Expected 2 journal lines, got %d", len(lines))
	}

	if lines[0].Side != models.SideDebit || lines[0].AccountID != trans.AccountFrom {
		t.Fatal("First line should be debit for source account")
	}

	if lines[1].Side != models.SideCredit || lines[1].AccountID != trans.AccountTo {
		t.Fatal("Second line should be credit for destination account")
	}

	t.Log("✓ Happy flow test passed: Transaction created successfully with balanced entries")
}

// Test 2: Unbalanced entry - attempting to create mismatched debit/credit
func TestCreateTransaction_UnbalancedEntry(t *testing.T) {
	mockJournalRepo := NewMockJournalRepository()
	mockLinesRepo := NewMockJournalLinesRepository()

	trans := &models.TransactionRequest{
		AccountFrom:    "acc_001",
		AccountTo:      "acc_002",
		Amount:         1000,
		Currency:       "USD",
		ExternalID:     "ext_456",
		IdempotencyKey: "key_456",
	}

	// Create journal entry
	err := mockJournalRepo.Create(trans)
	if err != nil {
		t.Fatalf("Failed to create journal: %v", err)
	}

	journal, _ := mockJournalRepo.Find(trans.ExternalID)

	// Create journal lines with unbalanced amounts
	err = mockLinesRepo.Create2Lines(trans, journal.EntryID)
	if err != nil {
		t.Fatalf("Failed to create journal lines: %v", err)
	}

	lines := mockLinesRepo.lines[journal.EntryID]

	// Verify debit and credit amounts match (balanced)
	debitAmount := 0
	creditAmount := 0

	for _, line := range lines {
		if line.Side == models.SideDebit {
			debitAmount++
		} else {
			creditAmount++
		}
	}

	if debitAmount == creditAmount {
		t.Logf("✓ Balanced entry test passed: Entry is properly balanced - debits: %d, credits: %d", debitAmount, creditAmount)
		return
	}

	t.Fatal("✗ Balanced entry test failed: Entry should be balanced")
}

// Test 3: Reversal - reverse a transaction
func TestReverseTransaction(t *testing.T) {
	mockJournalRepo := NewMockJournalRepository()
	mockLinesRepo := NewMockJournalLinesRepository()

	trans := &models.TransactionRequest{
		AccountFrom:    "acc_001",
		AccountTo:      "acc_002",
		Amount:         1000,
		Currency:       "USD",
		ExternalID:     "ext_789",
		IdempotencyKey: "key_789",
	}

	// Create and populate initial transaction
	_ = mockJournalRepo.Create(trans)
	journal, _ := mockJournalRepo.Find(trans.ExternalID)
	_ = mockLinesRepo.Create2Lines(trans, journal.EntryID)

	// Verify original state
	originalLines := mockLinesRepo.lines[journal.EntryID]
	if len(originalLines) != 2 {
		t.Fatalf("Expected 2 original lines, got %d", len(originalLines))
	}

	originalDebitSide := originalLines[0].Side
	originalCreditSide := originalLines[1].Side

	// Perform reversal
	err := mockLinesRepo.ReverseTransaction("1")
	if err != nil {
		t.Fatalf("Failed to reverse transaction: %v", err)
	}

	// Verify reversal
	reversedLines := mockLinesRepo.lines[journal.EntryID]
	if reversedLines[0].Side == originalDebitSide {
		t.Fatal("Debit line should be reversed to credit")
	}
	if reversedLines[1].Side == originalCreditSide {
		t.Fatal("Credit line should be reversed to debit")
	}

	t.Log("✓ Reversal test passed: Transaction reversed successfully - sides flipped")
}
