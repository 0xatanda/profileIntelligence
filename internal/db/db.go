package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func Connect() *sql.DB {
	dsn := os.Getenv("DATABASE_URL")

	// Step 1: Ensure DB exists (DEV ONLY)
	ensureDatabase()

	// Step 2: Connect to actual DB
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}

// 🔧 Create DB if not exists (local only)
func ensureDatabase() {
	baseDSN := "postgres://localhost:5432/postgres?sslmode=disable"

	db, err := sql.Open("postgres", baseDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var exists bool
	err = db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = 'profiles_db')",
	).Scan(&exists)

	if err != nil {
		log.Fatal(err)
	}

	if !exists {
		_, err = db.Exec("CREATE DATABASE profiles_db")
		if err != nil {
			log.Fatal("failed to create database:", err)
		}
		fmt.Println("✅ Database 'profiles_db' created")
	}
}
