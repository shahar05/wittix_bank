package models

type AccountRequest struct {
	Name     string `json:"name"`
	Currency string `json:"currency"`
}

type JournalRequest struct {
	ExternalID     string `json:"external_id"`
	Description    string `json:"description"`
	IdempotencyKey string `json:"idempotency_key"`
}

type TransactionRequest struct {
	AccountFrom    string `json:"account_from"`
	AccountTo      string `json:"account_to"`
	Amount         int64  `json:"amount"`
	Currency       string `json:"currency"`
	ExternalID     string `json:"external_id"`
	IdempotencyKey string `json:"idempotency_key"`
}
