package fileprocessor

import (
	"bufio"
	"fmt"
	"neowayv1/internal/batch"
	"os"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
)

func ProcessFileAndPersist(filePath string, pool *pgxpool.Pool) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	delimiter := ","
	if strings.HasSuffix(filePath, ".txt") {
		delimiter = "\t"
	}

	return processAndPersistData(file, delimiter, pool)
}

func processAndPersistData(file *os.File, delimiter string, pool *pgxpool.Pool) error {
	scanner := bufio.NewScanner(file)
	var batchRows [][]string
	const batchSize = 50000 // Set batch size for performance optimization

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, delimiter)
		batchRows = append(batchRows, fields)

		// Insert batch when the batch size is reached
		if len(batchRows) >= batchSize {
			err := batch.InsertBatchData(pool, batchRows)
			if err != nil {
				return fmt.Errorf("error inserting batch data: %v", err)
			}
			batchRows = nil
		}
	}

	// Insert remaining data
	if len(batchRows) > 0 {
		err := batch.InsertBatchData(pool, batchRows)
		if err != nil {
			return fmt.Errorf("error inserting remaining data: %v", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	return nil
}
