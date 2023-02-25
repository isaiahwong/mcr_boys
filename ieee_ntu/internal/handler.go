package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *Server) queryHandler(c *gin.Context) {
	rows, err := s.db.Query("SELECT * FROM accounts")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	accounts := []*Account{}

	// Print results
	for rows.Next() {
		var id int
		var balance float64
		err = rows.Scan(&id, &balance)
		if err != nil {
			log.Fatal(err)
		}
		accounts = append(accounts, &Account{AccountId: id, Balance: balance})
	}

	// convert to json bytes
	jsonData, err := json.Marshal(accounts)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": string(jsonData),
	})
}

func (s *Server) txHandler(c *gin.Context) {
	// Begin a transaction
	tx, err := s.db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	// Set the seed value for the random number generator
	rand.Seed(time.Now().UnixNano())

	// Generate a random integer between 0 and 99
	id := rand.Intn(100000)
	balance := rand.Intn(100000)

	// Prepare a statement for inserting a new record into the "people" table
	stmt, err := tx.Prepare("UPDATE accounts SET balance = $1 WHERE account_id = $2")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// Insert a new record using the prepared statement
	_, err = stmt.Exec(balance, id)
	if err != nil {
		// If there's an error, rollback the transaction
		tx.Rollback()
		log.Fatal(err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": fmt.Sprintf("Account %d updated with balance %d", id, balance)})

}
