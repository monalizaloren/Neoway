package batch

import (
	"context"
	"fmt"
	"neowayv1/internal/validation"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// InsertBatchData processes and inserts rows in batches into the database.
func InsertBatchData(pool *pgxpool.Pool, batchRows [][]string) error {
	batch := &pgx.Batch{}

	for i, row := range batchRows {
		if i == 0 {
			fmt.Println("Skipping header:", row)
			continue
		}

		fields := strings.Fields(strings.Join(row, " "))

		if len(fields) < 8 {
			return fmt.Errorf("incomplete row: expected 8 columns, found %d: %v", len(fields), fields)
		}

		cpf := validation.CleanString(fields[0])
		private := convertToBool(fields[1])
		incomplete := convertToBool(fields[2])
		lastPurchaseDate, err := convertToDate(fields[3])
		if err != nil {
			return fmt.Errorf("error converting date: %v", err)
		}
		averageTicket, err := convertToFloat(fields[4])
		if err != nil {
			return fmt.Errorf("error converting average ticket: %v", err)
		}
		lastTicket, err := convertToFloat(fields[5])
		if err != nil {
			return fmt.Errorf("error converting last ticket: %v", err)
		}

		mostFrequentStorePtr := convertToNullString(fields[6])
		lastStorePtr := convertToNullString(fields[7])

		var mostFrequentStore, lastStore string
		if mostFrequentStorePtr != nil {
			mostFrequentStore = validation.CleanString(*mostFrequentStorePtr)
		}
		if lastStorePtr != nil {
			lastStore = validation.CleanString(*lastStorePtr)
		}

		// Get the status of CPF and CNPJ using the validation functions.
		statusCPF := validation.ValidateCPF(cpf)
		statusCNPJFrequentStore := validation.ValidateCNPJ(mostFrequentStore)
		statusCNPJLastStore := validation.ValidateCNPJ(lastStore)

		if mostFrequentStore != "" {
			insertOrGetStoreID(pool, mostFrequentStore, statusCNPJFrequentStore, statusCNPJFrequentStore)
		}
		if lastStore != "" {
			insertOrGetStoreID(pool, lastStore, statusCNPJLastStore, statusCNPJLastStore)
		}

		queryCustomers := `INSERT INTO customers (cpf, private, incomplete, status_cpf, most_frequent_store_cnpj, last_store_cnpj, status_cnpj_last_store, status_cnpj_frequent_store)
                           VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
                           ON CONFLICT (cpf) DO NOTHING;`
		batch.Queue(queryCustomers, cpf, private, incomplete, statusCPF, mostFrequentStore, lastStore, statusCNPJLastStore, statusCNPJFrequentStore)

		queryTransactions := `INSERT INTO transactions (cpf, last_purchase_date, average_ticket, last_ticket)
                              VALUES ($1, $2, $3, $4);`
		batch.Queue(queryTransactions, cpf, lastPurchaseDate, averageTicket, lastTicket)
	}

	results := pool.SendBatch(context.Background(), batch)
	if err := results.Close(); err != nil {
		return fmt.Errorf("error executing batch: %v", err)
	}

	return nil
}

func convertToNullString(input string) *string {
	if input == "NULL" || input == "" {
		return nil
	}
	return &input
}

func convertToDate(input string) (interface{}, error) {
	if input == "NULL" || input == "" {
		return nil, nil
	}
	parsedDate, err := time.Parse("2006-01-02", input)
	if err != nil {
		return nil, err
	}
	return parsedDate, nil
}

func convertToFloat(input string) (interface{}, error) {
	if input == "NULL" || input == "" {
		return nil, nil
	}
	input = strings.Replace(input, ",", ".", -1)

	var num float64
	_, err := fmt.Sscanf(input, "%f", &num)
	if err != nil {
		return nil, err
	}
	return num, nil
}

func convertToBool(input string) bool {
	return input == "1"
}

// Checks if a store with a given CNPJ exists and retrieves its ID; if it doesn't exist, it inserts a new store with the provided status values.
func insertOrGetStoreID(pool *pgxpool.Pool, cnpj string, statusLastStore, statusFrequentStore string) {
	var storeID int
	err := pool.QueryRow(context.Background(), `SELECT id FROM stores WHERE cnpj = $1`, cnpj).Scan(&storeID)
	if err != nil && err != pgx.ErrNoRows {
		fmt.Printf("Error querying store: %v\n", err)
		return
	}

	if storeID == 0 {
		query := `INSERT INTO stores (cnpj, status_cnpj_last_store, status_cnpj_frequent_store)
                     VALUES ($1, $2, $3)`
		_, err = pool.Exec(context.Background(), query, cnpj, statusLastStore, statusFrequentStore)
		if err != nil {
			fmt.Printf("Error inserting store: %v\n", err)
		}
	}
}
