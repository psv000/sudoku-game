# Makefile for the Sudoku Go project

GO_BIN=${HOME}/go/bin

# --- Configuration ---
ifneq (,$(wildcard ./.env))
	include .env
	export
endif

# Set default superuser if not specified in .env
PG_SUPERUSER ?= postgres

# --- Main targets ---
.PHONY: all setup install-deps create-db migrate run help clean

all: help

setup: install-deps create-db migrate
	@echo "✅ Setup complete! Run 'make run' to start the server."

install-deps:
	@echo "Checking dependencies..."
	@os_name=$$(uname -s); \
	case "$$os_name" in \
		Linux*) \
			if ! command -v go &> /dev/null; then \
				echo "Go not found. Please install it first: https://golang.org/doc/install"; \
				exit 1; \
			fi; \
			if ! command -v psql &> /dev/null; then \
				echo "PostgreSQL client not found. Installing postgresql-client..."; \
				sudo apt-get update && sudo apt-get install -y postgresql-client; \
			fi; \
			;; \
		Darwin*) \
			if ! command -v brew &> /dev/null; then \
				echo "Homebrew not found. Please install it first: https://brew.sh"; \
				exit 1; \
			fi; \
			if ! command -v go &> /dev/null; then \
				echo "Go not found. Installing with Homebrew..."; \
				brew install go; \
			fi; \
			if ! command -v psql &> /dev/null; then \
				echo "PostgreSQL not found. Installing with Homebrew..."; \
				brew install postgresql; \
			fi; \
			;; \
		*) \
			echo "Unsupported OS: $$os_name. Please install Go and PostgreSQL manually."; \
			exit 1; \
			;; \
	esac
	@echo "Installing Go dependencies..."
	go mod tidy
	@if ! command -v ${GO_BIN}/goose &> /dev/null; then \
		echo "Goose not found. Installing..."; \
		go install github.com/pressly/goose/v3/cmd/goose@latest; \
	fi
	@echo "Dependencies check complete."

create-db:
	@echo "Configuring database..."
	
	# Determine the command prefix depending on the OS
	@PSQL_CMD_PREFIX=""; \
	os_name=$$(uname -s); \
	if [ "$$os_name" = "Linux" ]; then \
		PSQL_CMD_PREFIX="sudo -u ${PG_SUPERUSER}"; \
	fi; \
	
	# Check for and create user
	@if ! psql -U ${PG_SUPERUSER} -tc "SELECT 1 FROM pg_roles WHERE rolname = '${DB_USER}'" | grep -q 1; then \
		echo "Creating database user: ${DB_USER}..."; \
		$$PSQL_CMD_PREFIX createuser --createdb ${DB_USER}; \
		$$PSQL_CMD_PREFIX psql -c "ALTER USER ${DB_USER} WITH PASSWORD '${DB_PASSWORD}';"; \
	else \
		echo "Database user '${DB_USER}' already exists."; \
	fi; \
	
	# Check for and create database
	@if ! psql -U ${PG_SUPERUSER} -lqt | cut -d \| -f 1 | grep -qw ${DB_NAME}; then \
		echo "Creating database: ${DB_NAME}..."; \
		$$PSQL_CMD_PREFIX createdb -O ${DB_USER} ${DB_NAME}; \
	else \
		echo "Database '${DB_NAME}' already exists."; \
	fi
	@echo "Database configuration complete."


migrate:
	@echo "Running database migrations..."
	${GO_BIN}/goose -dir "db/migrations" up

run:
	@echo "Starting server on http://localhost:8080..."; \
	go run cmd/main.go &
	sleep 0.01
	open http://localhost:8080

help:
	@echo "Available commands:"
	@echo "  make setup		 	- Install all dependencies, create DB, and run migrations."
	@echo "  make install-deps  - Install Go, PostgreSQL client, and Goose."
	@echo "  make create-db	 	- Create PostgreSQL user and database."
	@echo "  make migrate	  	- Apply database migrations."
	@echo "  make run		   	- Start the application server."
	@echo "  make help		  	- Show this help message."