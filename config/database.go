package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func ConnectDatabase() {
	var err error
	dbURL := "postgres://postgres:postgres@localhost:5433/meeting_scheduler?sslmode=disable"
	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	fmt.Println("Database connected!")
}

func GetDB() *sql.DB {
	if db == nil {
		log.Fatal("Database not initialized. Call ConnectDatabase() first.")
	}
	return db
}

func CloseDatabase() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
