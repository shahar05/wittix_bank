package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"wittix_bank/models"
	"wittix_bank/services"

	"github.com/gorilla/mux"
)

type JournalHandler struct {
	svc services.JournalService
}

// NewJournalHandler constructs an JournalHandler.
func NewJournalHandler(s services.JournalService) *JournalHandler {
	return &JournalHandler{svc: s}
}

func RegisterJournalRoutes(r *mux.Router, svc services.JournalService) {
	h := NewJournalHandler(svc)
	r.HandleFunc("/journal-entries", h.CreateTransactionHandler).Methods("POST") // Perform Transaction
	r.HandleFunc("/journal-entries/{entry_id}/reverse", h.ReverseTransactionHandler).Methods("POST")
}

func (h *JournalHandler) CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var req models.TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("Received transaction request")
	fmt.Println("Extrnal ID:", req.ExternalID)
	fmt.Println("Idempotency Key:", req.IdempotencyKey)
	fmt.Println("Account From:", req.AccountFrom)
	fmt.Println("Account To:", req.AccountTo)
	// TODO: return full object
	err := h.svc.CreateTransaction(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(req)
}

func (h *JournalHandler) ReverseTransactionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	entryIDStr := vars["entry_id"]
	// var entryID int64
	// entryID, err := strconv.ParseInt(entryIDStr, 10, 64)
	// if err != nil {
	// 	http.Error(w, "Invalid entry_id", http.StatusBadRequest)
	// 	return
	// }

	err := h.svc.ReverseTransaction(entryIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
