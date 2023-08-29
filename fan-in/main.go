package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

// reading multiple csv file into a single classified data
// write classified result to a single file.

type User struct {
	FirstName string
	LastName  string
	Email     string
	Gender    string
}

func main() {
	// going to read each file to related channel
	// and merge them in mergedUsers chan
	normalUsers := make(chan User)
	preRegisteredUsers := make(chan User)
	vipUsers := make(chan User)

	// load data from different files into channels
	go func() {
		readUserCsv("files/normal-users.csv", normalUsers)
		close(normalUsers)
	}()
	go func() {
		readUserCsv("files/pre-registered-users.csv", preRegisteredUsers)
		close(preRegisteredUsers)
	}()
	go func() {
		readUserCsv("files/vip-users.csv", vipUsers)
		close(vipUsers)
	}()

	wg := sync.WaitGroup{}
	wg.Add(3)
	mergedUsers := make(chan User)
	go func() {
		wg.Wait()
		close(mergedUsers)
	}()

	cnt := make(chan int)
	go counter(cnt)

	// fan-in 3 channels to mergerUsers
	go func() {
		defer wg.Done()
		for user := range normalUsers {
			// fake delay for demonstration of other processes
			time.Sleep(time.Millisecond * 2)
			mergedUsers <- user
		}
	}()
	go func() {
		defer wg.Done()
		for user := range preRegisteredUsers {
			// fake delay for demonstration of other processes
			time.Sleep(time.Millisecond * 1)
			mergedUsers <- user
		}
	}()
	go func() {
		defer wg.Done()
		for user := range vipUsers {
			// fake delay for demonstration of other processes
			time.Sleep(time.Millisecond * 3)
			mergedUsers <- user
		}
	}()

	run := true
	go func() {
		for user := range mergedUsers {
			fmt.Printf("\rindex %d", <-cnt)
			_ = user
		}

		fmt.Println("")
		run = false
	}()

	for run {
	}
}

func readUserCsv(filename string, ch chan<- User) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}

	defer file.Close()

	reader := csv.NewReader(file)
	for {
		row, err := reader.Read()
		if err != nil {
			break
		}

		user := User{}
		user.FirstName = row[0]
		user.LastName = row[1]
		user.Email = row[2]
		user.Gender = row[3]

		ch <- user
	}
}
func readUserCsv1(filename string) []User {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}

	defer file.Close()

	users := []User{}
	reader := csv.NewReader(file)
	for {
		row, err := reader.Read()
		if err != nil {
			break
		}

		user := User{}
		user.FirstName = row[0]
		user.LastName = row[1]
		user.Email = row[2]
		user.Gender = row[3]
		users = append(users, user)
	}

	return users
}

func counter(ch chan<- int) {
	for i := 1; ; i++ {
		ch <- i

	}
}
