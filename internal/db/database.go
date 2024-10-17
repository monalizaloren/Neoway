package db

import (
	"context"
	"database/sql"
	"fmt"
	"neowayv1/internal/validation"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

func SetupConnectionPool() (*pgxpool.Pool, error) {
	err := godotenv.Load() // Load environment variables from .env file
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

func UpdateCPFandCNPJStatus(pool *pgxpool.Pool) error {
	rows, err := pool.Query(context.Background(), `
        SELECT cpf, most_frequent_store_cnpj, last_store_cnpj
        FROM customers
    `)
	if err != nil {
		return fmt.Errorf("error querying customers table: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var cpf, mostFrequentStoreCNPJ, lastStoreCNPJ sql.NullString
		var statusCNPJLastStore, statusCNPJFrequentStore string

		err := rows.Scan(&cpf, &mostFrequentStoreCNPJ, &lastStoreCNPJ)
		if err != nil {
			return fmt.Errorf("error scanning data: %w", err)
		}

		// Validate CNPJ for the most frequent and last visited stores
		if mostFrequentStoreCNPJ.Valid && validation.ValidateCNPJ(mostFrequentStoreCNPJ.String) {
			statusCNPJFrequentStore = "valid"
		} else {
			statusCNPJFrequentStore = "invalid"
		}

		if lastStoreCNPJ.Valid && validation.ValidateCNPJ(lastStoreCNPJ.String) {
			statusCNPJLastStore = "valid"
		} else {
			statusCNPJLastStore = "invalid"
		}

		// Executa atualização de status dos customers
		_, err = pool.Exec(context.Background(),
			"UPDATE customers SET status_cnpj_last_store=$1, status_cnpj_frequent_store=$2 WHERE cpf=$3",
			statusCNPJLastStore, statusCNPJFrequentStore, cpf.String)
		if err != nil {
			return fmt.Errorf("error updating status: %w", err)
		}

		// Executa atualização de status das stores
		_, err = pool.Exec(context.Background(),
			"UPDATE stores SET status_cnpj_last_store=$1, status_cnpj_frequent_store=$2 WHERE cnpj=$3 OR cnpj=$4",
			statusCNPJLastStore, statusCNPJFrequentStore, lastStoreCNPJ.String, mostFrequentStoreCNPJ.String)
		if err != nil {
			return fmt.Errorf("error updating store status: %w", err)
		}
	}

	if rows.Err() != nil {
		return fmt.Errorf("error iterating over records: %w", rows.Err())
	}

	fmt.Println("CPF and CNPJ statuses updated successfully.")
	return nil
}
