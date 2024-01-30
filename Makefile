.PHONY: migrateup
migrateup:
	go run ./cmd/migrator -storage-path "postgres://admin:password123@127.0.0.1:6500/AuthDatabase?sslmode=disable" -migrations-path "./migrations" up

.PHONY: migratedown
migratedown:
		go run ./cmd/migrator -storage-path "postgres://admin:password123@127.0.0.1:6500/AuthDatabase?sslmode=disable" -migrations-path "./migrations" down
