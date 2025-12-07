package repository

import (
	"database/sql"
	"wittix_bank/models"
)

type JournalRepository struct {
	DB *sql.DB
}

func (r *JournalRepository) Create(trans *models.TransactionRequest) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO journal_entries (external_id, idempotency_key) VALUES ($1, $2) RETURNING entry_id`
	var entryID string
	err = tx.QueryRow(query, trans.ExternalID, trans.IdempotencyKey).Scan(&entryID)
	if err != nil {
		return err
	}

	query = `INSERT INTO journal_lines (entry_id, account_id, side, amount) VALUES ($1, $2, $3, $4)`

	// Insert debit line for source account
	_, err = tx.Exec(query, entryID, trans.AccountFrom, "debit", trans.Amount)
	if err != nil {
		return err
	}

	// Insert credit line for destination account
	_, err = tx.Exec(query, entryID, trans.AccountTo, "credit", trans.Amount)
	if err != nil {
		return err
	}

	return tx.Commit()

}

func (r *JournalRepository) Find(externalID string) (*models.JournalEntry, error) {
	query := `SELECT entry_id, external_id, description, posted_at, reversal_of, idempotency_key FROM journal_entries WHERE external_id = $1`
	row := r.DB.QueryRow(query, externalID)

	var journal models.JournalEntry
	err := row.Scan(
		&journal.EntryID,
		&journal.ExternalID,
		&journal.Description,
		&journal.PostedAt,
		&journal.ReversalOf,
		&journal.IdempotencyKey)

	if err != nil {
		return nil, err
	}
	return &journal, nil
}
