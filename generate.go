package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func generate() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Create CSV file
	file, err := os.Create("data.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Create CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Define column headers
	headers := []string{"idcard", "fullname", "age", "address", "birthdate"}
	err = writer.Write(headers)
	if err != nil {
		panic(err)
	}

	// Generate and write 1 million records
	for i := 0; i < 1000000; i++ {
		row := []string{
			fmt.Sprintf("%013d", rand.Intn(9999999999999)), // 13-digit random ID card
			generateFakeFullName(),                         // Generate a fake full name
			fmt.Sprintf("%d", rand.Intn(100)),              // Random age
			fmt.Sprintf("Address %d", i),                   // Example address
			fmt.Sprintf("%04d-%02d-%02d", rand.Intn(2023), rand.Intn(12)+1, rand.Intn(28)+1), // Random birthdate
		}
		err := writer.Write(row)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("CSV file generated with 1 million records.")
}

func generateFakeFullName() string {
	firstNames := []string{"John", "Jane", "Peter", "Mary", "David", "Kate", "Robert", "Emily", "Michael", "Jessica", "William", "Emma", "James", "Olivia", "Christopher", "Sophia", "Matthew", "Ava", "Daniel", "Isabella"}
	lastNames := []string{"Smith", "Johnson", "Williams", "Brown", "Jones", "Miller", "Davis", "Garcia", "Rodriguez", "Wilson", "Martinez", "Anderson", "Taylor", "Thomas", "Hernandez", "Moore", "Martin", "Jackson", "Thompson", "White"}
	return firstNames[rand.Intn(len(firstNames))] + " " + lastNames[rand.Intn(len(lastNames))]
}
