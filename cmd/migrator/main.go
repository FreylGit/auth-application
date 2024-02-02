package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var storagePath, migrationsPath, migrationsTable string
	flag.StringVar(&storagePath, "storage-path", "postgres://admin:password123@postgres:5432/AuthDatabase?sslmode=disable", "path to storage")

	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migrations table")
	flag.Parse()
	action := flag.Arg(0)
	if action != "up" && action != "down" {
		fmt.Println("Error: invalid action. Please use 'up' or 'down'.")
		os.Exit(1)
	}
	if storagePath == "" || migrationsPath == "" {
		panic("flag  storagePath or migrationsPath is empty")
	}
	m, err := migrate.New("file://"+migrationsPath, storagePath)
	if err != nil {
		fmt.Printf("Error creating migration instance: %v\n", err)
		panic(err)
	}
	switch action {
	case "up":
		if err = m.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("No migrations")
				return
			}
			fmt.Printf("Error applying migrations: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Migrations applied successfully")
	case "down":
		if err = m.Down(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("No migrations to roll back")
				return
			}
			fmt.Printf("Error rolling back migrations: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Migrations rolled back successfully")
	}
	fmt.Printf("migration successfule")
}
