package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

type User struct {
	FirstName string
	LastName  string
	Sex       string
	Email     string
	Birthday  time.Time
}

// real csv file line by line
func readFile(csvRow chan []string) {
	file, err := os.Open("./users.csv")
	if err != nil {
		log.Fatalln("failed to read file", err)
	}
	reader := csv.NewReader(file)
	for {
		line, err := reader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			} else {
				log.Fatalln("filed to read csv row", err)
			}
		}

		csvRow <- line
	}
}

func rowToUserStream(rowChan <-chan []string, userChan chan User) {
	cnt := 0
	for row := range rowChan {
		cnt++
		user := User{
			FirstName: row[0],
			LastName:  row[1],
			Email:     row[2],
			Sex:       row[3],
			Birthday: func() time.Time {
				b := row[4]
				bd, err := time.Parse("2006/01/02", b)
				if err != nil {
					log.Fatalf("Invalid birthday %v", b)
				}

				return bd
			}(),
		}
		fmt.Println("ROW TO USER", cnt)
		userChan <- user
	}
}

const year = 365 * 24 * 60 * 60

func ageCalculator(userChan <-chan User, ageChan chan int) {
	cnt := 0
	for user := range userChan {
		cnt++
		age := time.Since(user.Birthday)
		a := age.Seconds() / year
		ageChan <- int(a)
		fmt.Println("AGE CALCULATED", int(a), cnt)
	}
}

func main() {
	rowChan := make(chan []string)
	userChan := make(chan User)
	ageChan := make(chan int)

	go func() {
		readFile(rowChan)
		close(rowChan)
		fmt.Println("ROW CHAN CLOSED")
	}()

	go func() {
		rowToUserStream(rowChan, userChan)
		close(userChan)
		fmt.Println("USER CHAN CLOSED")
	}()

	go func() {
		ageCalculator(userChan, ageChan)
		close(ageChan)
		fmt.Println("AGE CHAN CLOSED")
	}()

	ageCloser := sync.WaitGroup{}
	ageCloser.Add(1)
	go func() {
		defer ageCloser.Done()
		cnt := 0
		for a := range ageChan {
			cnt++
			fmt.Println(a, cnt)
		}
		fmt.Println("END PRINT")
	}()

	fmt.Println("WAITING...")
	ageCloser.Wait()
	fmt.Println("LAST WAIT DONE")
}
