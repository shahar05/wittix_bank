package repository

import (
	"database/sql"
	"wittix_bank/models"
)

type JournalLinesRepository struct {
	DB *sql.DB
}

func (r *JournalLinesRepository) Create2Lines(transaction *models.TransactionRequest, entryID int64) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO journal_lines (entry_id, account_id, side, amount) VALUES ($1, $2, $3, $4)`

	// Insert debit line for source account
	_, err = tx.Exec(query, entryID, transaction.AccountFrom, "debit", transaction.Amount)
	if err != nil {
		return err
	}

	// Insert credit line for destination account
	_, err = tx.Exec(query, entryID, transaction.AccountTo, "credit", transaction.Amount)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *JournalLinesRepository) ReverseTransaction(entryID string) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
        INSERT INTO journal_lines (entry_id, account_id, side, amount)
        SELECT entry_id, account_id,
               CASE WHEN side = 'debit' THEN 'credit'::side_enum ELSE 'debit'::side_enum END AS side,
               amount
        FROM journal_lines
        WHERE entry_id = $1
    `
	_, err = tx.Exec(query, entryID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
