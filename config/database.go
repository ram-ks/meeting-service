package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func ConnectDatabase() {
	var err error
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

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
