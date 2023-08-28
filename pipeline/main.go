package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"
)

// read users.csv file
// convert rows to user struct
// get age of all users
// get count of users older than 30 years

type User struct {
	FirstName string
	LastName  string
	Email     string
	Gender    string
	BirthDate time.Time
}

func main() {
	file, err := os.OpenFile("users.csv", os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	finished := false

	// pipeline1: read rows
	rows := make(chan []string)
	go func() {
		reader := csv.NewReader(file)
		for {
			row, err := reader.Read()
			if err != nil {
				fmt.Println("\nclosing rows")
				close(rows)
				break
			}
			rows <- row
			// add some delay to read and slow data printing on screen
			time.Sleep(time.Millisecond * 1)
		}
	}()

	// pipeline2: convert each row to user struct
	users := make(chan User)
	go func() {
		for row := range rows {
			user := User{}
			user.FirstName = row[0]
			user.LastName = row[1]
			user.Email = row[2]
			user.Gender = row[3]
			user.BirthDate, _ = time.Parse("2006/01/02", row[4])

			users <- user
		}
		fmt.Println("closing users")
		close(users)
	}()

	// pipeline3:  get age of all users
	ages := make(chan int)
	go func() {
		for user := range users {
			age := int(time.Now().Sub(user.BirthDate).Hours() / 24 / 365)
			ages <- age
		}
		fmt.Println("closing ages")
		close(ages)
	}()

	// pipeline4: older than 30 years old detector
	counter := make(chan int)
	go olderThan30Counter(counter)
	go func() {
		for age := range ages {
			if age >= 30 {
				fmt.Printf("\rUser older than 30 detected, we have %d user now.", <-counter)
			}
		}

		fmt.Println("")
		finished = true
	}()

	for !finished {
	}

}

func olderThan30Counter(counter chan int) {
	for i := 1; ; i++ {
		counter <- i
	}
}
