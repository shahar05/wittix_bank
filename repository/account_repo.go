package repository

import (
	"database/sql"
	"fmt"
	"wittix_bank/models"
)

type AccountRepository struct {
	DB *sql.DB
}

func (r *AccountRepository) Create(account *models.AccountRequest) error {
	q := `INSERT INTO accounts(name,currency) VALUES($1,$2)`
	return r.DB.QueryRow(q, account.Name, account.Currency).Err()
}

func (r *AccountRepository) GetBalance(id string) (float64, error) {
	var balance float64

	query := `
		SELECT 
			COALESCE(SUM(
				CASE 
					WHEN side = 'debit' THEN amount
					ELSE -amount
				END
			), 0)
		FROM journal_lines
		WHERE account_id = $1
	`

	err := r.DB.QueryRow(query, id).Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("failed to get balance: %w", err)
	}

	return balance, nil
}
