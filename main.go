package main

import (
	"log"
	"net/http"
	"wittix_bank/db"
	"wittix_bank/handlers"
	"wittix_bank/repository"
	"wittix_bank/services"

	"github.com/gorilla/mux"
)

func main() {
	conn := db.Connect()
	// ensure we close DB connection on exit (assumes conn has Close method, e.g. *sql.DB)
	if conn == nil {
		log.Fatal("failed to connect to database")
	}
	defer conn.Close()

	router := mux.NewRouter()

	AccountRepository := &repository.AccountRepository{DB: conn}
	accountSvc := services.AccountService{PropRepo: AccountRepository}
	handlers.RegisterAccountRoutes(router, accountSvc)

	JournalRepository := &repository.JournalRepository{DB: conn}
	journalLinesRepo := &repository.JournalLinesRepository{DB: conn}
	journalSvc := services.JournalService{JournalRepo: JournalRepository, JournalLinesRepo: journalLinesRepo}
	handlers.RegisterJournalRoutes(router, journalSvc)

	log.Println("Server listening on :8080")
	// handle ListenAndServe error instead of ignoring it
	log.Fatal(http.ListenAndServe(":8080", router))
}
