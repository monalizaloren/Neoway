package batch

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

func InsertBatchData(pool *pgxpool.Pool, batchRows [][]string) error {
	batch := &pgx.Batch{}

	for i, row := range batchRows {
		// Skip the header row
		if i == 0 {
			fmt.Println("Skipping header:", row)
			continue
		}

		fields := strings.Fields(strings.Join(row, " "))

		// Validate number of columns
		if len(fields) < 8 {
			return fmt.Errorf("incomplete row: expected 8 columns, found %d: %v", len(fields), fields)
		}

		cpf := fields[0]
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
		mostFrequentStore := convertToNullString(fields[6])
		lastStore := convertToNullString(fields[7])

		insertOrGetStoreID(pool, mostFrequentStore)
		insertOrGetStoreID(pool, lastStore)

		queryCustomers := `INSERT INTO customers (cpf, private, incomplete, status_cpf, most_frequent_store_cnpj, last_store_cnpj)
						VALUES ($1, $2, $3, $4, $5, $6)
						ON CONFLICT (cpf) DO NOTHING;`
		batch.Queue(queryCustomers, cpf, private, incomplete, "valid", mostFrequentStore, lastStore)

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

func insertOrGetStoreID(pool *pgxpool.Pool, cnpj interface{}) {
	if cnpj == nil {
		return
	}

	var storeID int
	err := pool.QueryRow(context.Background(), `SELECT id FROM stores WHERE cnpj = $1`, cnpj).Scan(&storeID)
	if err != nil && err != pgx.ErrNoRows {
		fmt.Printf("Error querying store: %v\n", err)
		return
	}

	if storeID == 0 {
		_, err = pool.Exec(context.Background(), `INSERT INTO stores (cnpj, status_cnpj_last_store, status_cnpj_frequent_store)
                                                   VALUES ($1, 'invalid', 'invalid')`, cnpj)
		if err != nil {
			fmt.Printf("Error inserting store: %v\n", err)
		}
	}
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

func convertToNullString(input string) interface{} {
	if input == "NULL" || input == "" {
		return nil
	}
	return input
}

func convertToBool(input string) bool {
	return input == "1"
}
