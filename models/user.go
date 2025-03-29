package models

import (
	"log"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

// InsertUser calls the insert_user PL/pgSQL function in your database
func InsertUser(db *pgxpool.Pool, name, email, password string) error {
	// Call the insert_user PL/pgSQL function
	query := `SELECT insert_user($1, $2, $3)`
	_, err := db.Exec(context.Background(), query, name, email, password)
	if err != nil {
		log.Printf("Error calling insert_user function: %v\n", err)
		return err
	}
	log.Printf("User inserted successfully: Name=%s, Email=%s\n", name, email)
	return nil
}



