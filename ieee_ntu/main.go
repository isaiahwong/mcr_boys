package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
    dbDriver = "mysql"
    dbSource = "root:password@tcp(127.0.0.1:3306)/accounts_db"
)

var db *sql.DB

func init() {
    var err error
    db, err = sql.Open(dbDriver, dbSource)
    if err != nil {
        log.Fatal(err)
    }

    if err = db.Ping(); err != nil {
        log.Fatal(err)
    }

    log.Println("Connected to database")
}

func QueryBalance(w http.ResponseWriter, r *http.Request) {
    // Get the account ID from the request parameters
    accountID := r.URL.Query().Get("account_id")
    if accountID == "" {
        http.Error(w, "Missing account ID", http.StatusBadRequest)
        return
    }

    // Query the database for the account balance
    var balance float64
    err := db.QueryRow("SELECT balance FROM accounts WHERE account_id = ?", accountID).Scan(&balance)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Return the account balance in the response
    fmt.Fprintf(w, "Account balance for account %s: %.2f", accountID, balance)
}

func handleTransaction(w http.ResponseWriter, r *http.Request) {
    // Parse the request body to get the source and destination account IDs and the amount to transfer
    values := make(map[string]string)
    if err := json.NewDecoder(r.Body).Decode(&values); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    defer r.Body.Close()

    sourceAccountID, err := strconv.Atoi(values["source_account_id"])
    if err != nil {
        http.Error(w, "Invalid source account ID", http.StatusBadRequest)
        return
    }

    destinationAccountIDsStr := strings.Split(values["destination_account_ids"], ",")
    var destinationAccountIDs []int
    for _, destinationAccountIDStr := range destinationAccountIDsStr {
        destinationAccountID, err := strconv.Atoi(destinationAccountIDStr)
        if err != nil {
            http.Error(w, "Invalid destination account ID", http.StatusBadRequest)
            return
        }
        destinationAccountIDs = append(destinationAccountIDs, destinationAccountID)
    }

    amount, err := strconv.ParseFloat(values["amount"], 64)
    if err != nil {
        http.Error(w, "Invalid amount", http.StatusBadRequest)
        return
    }

    // Start a transaction
    tx, err := db.Begin()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer tx.Rollback()

    // Debit the amount from the source account
    var sourceBalance float64
    err = tx.QueryRow("SELECT balance FROM accounts WHERE account_id = ? FOR UPDATE", sourceAccountID).Scan(&sourceBalance)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    if sourceBalance < amount {
        http.Error(w, "Insufficient balance", http.StatusBadRequest)
        return
    }
    sourceBalance -= amount
    _, err = tx.Exec("UPDATE accounts SET balance = ? WHERE account_id = ?", sourceBalance, sourceAccountID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Credit the amount to the destination accounts
    var totalCredit float64
    for _, destinationAccountID := range destinationAccountIDs {
        var destinationBalance float64
        err = tx.QueryRow("SELECT balance FROM accounts WHERE account_id = ? FOR UPDATE", destinationAccountID).Scan(&destinationBalance)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        destinationBalance += amount / float64(len(destinationAccountIDs))
        _, err = tx.Exec("UPDATE accounts SET balance = ? WHERE account_id = ?", destinationBalance, destinationAccountID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        totalCredit += amount / float64(len(destinationAccountIDs))
    }

    // Commit the transaction
    err = tx.Commit()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Return the transaction details in the response
    fmt.Fprintf(w, "Transaction of %.2f from account %d to accounts %s completed at %s", amount, sourceAccountID, values["destination_account_ids"], time.Now().Format("2006-01-02 15:04:05"))
}

func main() {
    defer db.Close()

    http.HandleFunc("/query_balance", QueryBalance)

    log.Println("Server started on port 8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}