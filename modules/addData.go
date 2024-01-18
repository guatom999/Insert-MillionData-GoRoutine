package modules

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type (
	IAddData interface {
		InsertData(c echo.Context) error
		InsertDataTwo(c echo.Context) error
		MigrateData(c echo.Context) error
	}

	addData struct {
		db *gorm.DB
	}

	Data struct {
		IDCard    string `json:"id_card"`
		FullName  string `json:"full_name"`
		Age       int    `json:"age"`
		Address   string `json:"address"`
		Birthdate string `json:"birth_date"`
	}
)

func NewAddData(db *gorm.DB) IAddData {
	return &addData{db: db}
}

//Craete Data with CreateInBatch

func (r *addData) MigrateData(c echo.Context) error {
	r.db.Migrator().CreateTable(&Data{})

	return c.JSON(http.StatusOK, "Migrate Database Success")
}

func (r *addData) InsertData(c echo.Context) error {

	// *******Migrate Database***********

	r.db.Migrator().CreateTable(&Data{})

	data := make([]Data, 0)

	file, err := os.Open("data.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// // Skip header row if present
	reader.Read()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		ageInt, err := strconv.Atoi(record[2])
		if err != nil {
			return c.JSON(400, "convert string to int failed")
		}

		person := Data{
			IDCard:    record[0],
			FullName:  record[1],
			Age:       ageInt, // Convert age to int
			Address:   record[3],
			Birthdate: record[4],
		}
		data = append(data, person)
	}

	tx := r.db.CreateInBatches(&data, 1000)
	if tx.Error != nil {

		log.Printf("Error: Fail to insert data cause : %s", tx.Error.Error())
		tx.Rollback()

		return c.JSON(200, "erorr: insert data failed")
	}

	return c.JSON(200, "test")
}

func (r *addData) InsertDataTwo(c echo.Context) error {

	start := time.Now()

	data := make([]Data, 0)

	file, err := os.Open("data.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Skip header row if present
	reader.Read()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		ageInt, err := strconv.Atoi(record[2])
		if err != nil {
			return c.JSON(400, "convert string to int failed")
		}

		person := Data{
			IDCard:    record[0],
			FullName:  record[1],
			Age:       ageInt, // Convert age to int
			Address:   record[3],
			Birthdate: record[4],
		}
		data = append(data, person)
	}

	var wg sync.WaitGroup
	wg.Add(10) // Number of goroutines to create
	chunkSize := 10000

	ch := make(chan int, 10)

	successRate := make([]string, 10)

	for i := 0; i < 10; i++ {
		go func(chunk []Data, i int) {
			fmt.Println("Starting routine", i)
			defer wg.Done()

			// Insert each data item in the chunk here
			tx := r.db.CreateInBatches(&chunk, 5000)

			if tx.Error != nil {

				log.Printf("Error: Fail to insert data cause : %s", tx.Error.Error())
				tx.Rollback()

				panic(tx.Error)
			}

			// for _, item := range chunk {
			// 	// Replace this with your actual insertion logic
			// 	tx := r.db.CreateInBatches(&item, 1000)
			// 	// tx := r.db.Create(&item)

			// 	if tx.Error != nil {

			// 		log.Printf("Error: Fail to insert data cause : %s", tx.Error.Error())
			// 		tx.Rollback()

			// 		panic(err)
			// 	}
			// }
			ch <- i

		}(data[i*chunkSize:(i+1)*chunkSize], i)

		go func() {
			for range ch {
				successRate = append(successRate, "=")
				fmt.Println(successRate)
			}
			close(ch) // Close the channel after all goroutines finish
		}()
	}

	wg.Wait()

	fmt.Println("All data inserted.")
	elapsed := time.Since(start)
	log.Printf("Code segment took %s", elapsed)

	return c.JSON(http.StatusOK, "All data inserted.")

}
