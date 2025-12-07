package handlers

import (
	"encoding/json"
	"net/http"

	"wittix_bank/models"
	"wittix_bank/services"

	"github.com/gorilla/mux"
)

type AccountHandler struct {
	svc services.AccountService
}

// NewAccountHandler constructs an AccountHandler.
func NewAccountHandler(s services.AccountService) *AccountHandler {
	return &AccountHandler{svc: s}
}

func RegisterAccountRoutes(r *mux.Router, svc services.AccountService) {
	h := NewAccountHandler(svc)
	r.HandleFunc("/accounts", h.CreateAccountHandler).Methods("POST")
	r.HandleFunc("/accounts/{id}/balance", h.GetAccountBalanceHandler).Methods("GET")

}

func (h *AccountHandler) CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	var req models.AccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// TODO: return full object
	err := h.svc.CreateAccount(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(req)
}

func (h *AccountHandler) GetAccountBalanceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	balance, err := h.svc.GetAccountBalance(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"balance": balance})
}
