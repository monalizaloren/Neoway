package main

import (
	"fmt"
	"log"
	"neowayv1/internal/db"
	"neowayv1/internal/fileprocessor"
)

func main() {
	pool, err := db.SetupConnectionPool()
	if err != nil {
		log.Fatalf("Failed to configure connection pool : %v", err)
	}
	defer pool.Close()

	// Ensure required tables are created
	err = db.CreateTableIfNotExists(pool)
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	filePath := "./assets/base_teste.txt"

	err = fileprocessor.ProcessFileAndPersist(filePath, pool)
	if err != nil {
		log.Printf("File processing error: %v", err)
	} else {
		fmt.Println("File processed successfully.")
	}

	// Update CPF and CNPJ statuses
	err = db.UpdateCPFandCNPJStatus(pool)
	if err != nil {
		log.Printf("Error updating CPF and CNPJ status: %v", err)
	}
}
