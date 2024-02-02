# Функция для чтения значений из файла env
get_env_value_local = $(shell grep -E '^\s*$(1)\s*=\s*(.*)$$' config/local.env | sed -E 's/^\s*$(1)\s*=\s*(.*)$$/\1/')
get_env_value_dev = $(shell grep -E '^\s*$(1)\s*=\s*(.*)$$' config/dev.env | sed -E 's/^\s*$(1)\s*=\s*(.*)$$/\1/')


.PHONY: migrateup
migrateup:
	go run ./cmd/migrator -storage-path "postgres://$(call get_env_value_local,POSTGRES_USER):$(call get_env_value_local,POSTGRES_PASSWORD)@$(call get_env_value_local,POSTGRES_HOST):$(call get_env_value_local,POSTGRES_PORT)/$(call get_env_value_local,POSTGRES_DB)?sslmode=disable" -migrations-path "./migrations" up

.PHONY: migratedown
migratedown:
	go run ./cmd/migrator -storage-path "postgres://$(call get_env_value_local,POSTGRES_USER):$(call get_env_value_local,POSTGRES_PASSWORD)@$(call get_env_value_local,POSTGRES_HOST):$(call get_env_value_local,POSTGRES_PORT)/$(call get_env_value_local,POSTGRES_DB)?sslmode=disable" -migrations-path "./migrations" down

.PHONY: migrateupdev
migrateupdev:
	go run ./cmd/migrator -storage-path "postgres://$(call get_env_value_dev,POSTGRES_USER):$(call get_env_value_dev,POSTGRES_PASSWORD)@$(call get_env_value_dev,POSTGRES_HOST):$(call get_env_value_dev,POSTGRES_PORT)/$(call get_env_value_dev,POSTGRES_DB)?sslmode=disable" -migrations-path "./migrations" up

.PHONY: migratedowndev
migratedowndev:
	go run ./cmd/migrator -storage-path "postgres://$(call get_env_value_dev,POSTGRES_USER):$(call get_env_value_dev,POSTGRES_PASSWORD)@$(call get_env_value_dev,POSTGRES_HOST):$(call get_env_value_dev,POSTGRES_PORT)/$(call get_env_value_dev,POSTGRES_DB)?sslmode=disable" down
