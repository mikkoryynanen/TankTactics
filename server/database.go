package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func NewDatabase() *Database {
	connStr := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Panic(err)
	}
	// defer db.Close()

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}

	createTable(db)

	fmt.Printf("Connected to database at: %v\n", connStr)

	return &Database{
		db: db,
	}
}

func (db *Database) CreateUser(username string) (uuid.UUID, error) {
	query := `INSERT INTO "user" (username)
		VALUES ($1) RETURNING id`

	var pk uuid.UUID
	err := db.db.QueryRow(query, username).Scan(&pk)
	if err != nil {
		log.Fatal(err)
		return uuid.Nil, err
	}

	return pk, nil
}

func (db *Database) GetUser(id string) (*User, error) {
	query := `SELECT * FROM "user" WHERE id = $1 LIMIT 1`

	var user User
	err := db.db.QueryRow(query, id).Scan(&user)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &user, nil
}

func createTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS "user" (
    	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    	username VARCHAR(100) NOT NULL
	)`

	_, err := db.Exec(query)
	if err != nil {
		log.Panic(err)
	}
}
