package postgresql

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type Storage struct {
	db *sql.DB
}

func New(connectionString string) (*Storage, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	err = db.Ping()
	if err != nil {
		panic("Not connected to database")
		return nil, err
	}

	log.Println("Connected to the database")
	return &Storage{db: db}, nil
}
