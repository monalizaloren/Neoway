package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

// SetupConnectionPool initializes a connection pool for PostgreSQL using the connection string from the environment variables.
func SetupConnectionPool() (*pgxpool.Pool, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("error loading .env file")
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		return nil, fmt.Errorf("database URL environment variable not set")
	}

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("error parsing connection string: %w", err)
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("error connecting to PostgreSQL: %w", err)
	}

	return pool, nil
}

func CreateTableIfNotExists(pool *pgxpool.Pool) error {
	createStoresTableQuery := `
    CREATE TABLE IF NOT EXISTS stores (
        id SERIAL PRIMARY KEY,
        cnpj TEXT UNIQUE,
        status_cnpj_last_store TEXT,
        status_cnpj_frequent_store TEXT
    );`

	_, err := pool.Exec(context.Background(), createStoresTableQuery)
	if err != nil {
		return fmt.Errorf("error creating stores table: %w", err)
	}

	createCustomersTableQuery := `
    CREATE TABLE IF NOT EXISTS customers (
        cpf TEXT PRIMARY KEY,
        private BOOLEAN,
        incomplete BOOLEAN,
        status_cpf TEXT,
        most_frequent_store_cnpj TEXT,
        last_store_cnpj TEXT,
        status_cnpj_last_store TEXT,
        status_cnpj_frequent_store TEXT
    );`

	_, err = pool.Exec(context.Background(), createCustomersTableQuery)
	if err != nil {
		return fmt.Errorf("error creating customers table: %w", err)
	}

	createTransactionsTableQuery := `
    CREATE TABLE IF NOT EXISTS transactions (
        id SERIAL PRIMARY KEY,
        cpf TEXT REFERENCES customers(cpf),
        last_purchase_date DATE,
        average_ticket NUMERIC(10, 2), 
        last_ticket NUMERIC(10, 2) 
    );`

	_, err = pool.Exec(context.Background(), createTransactionsTableQuery)
	if err != nil {
		return fmt.Errorf("error creating transactions table: %w", err)
	}

	fmt.Println("Tables 'customers', 'stores', and 'transactions' verified/created successfully.")
	return nil
}
