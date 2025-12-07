package models

import "time"

type Account struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Currency  string    `db:"currency"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
}

type JournalEntry struct {
	EntryID        int64     `db:"entry_id"`
	ExternalID     string    `db:"external_id"`
	Description    string    `db:"description"`
	PostedAt       time.Time `db:"posted_at"`
	ReversalOf     *int64    `db:"reversal_of"`
	IdempotencyKey string    `db:"idempotency_key"`
}

type Side string

const (
	SideDebit  Side = "debit"
	SideCredit Side = "credit"
)

type JournalLine struct {
	ID        int64  `db:"id"`
	EntryID   int64  `db:"entry_id"`
	AccountID string `db:"account_id"`
	Side      Side   `db:"side"`
	Amount    string `db:"amount"` // Use string to keep NUMERIC precision
}
