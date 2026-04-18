package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func Connect() *sql.DB {
	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to open DB:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("failed to connect DB:", err)
	}

	log.Println("✅ Connected to database")

	return db
}
