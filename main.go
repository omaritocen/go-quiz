package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	// Create a flag for csv file param
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "a timer that decides the round time")
	flag.Parse()

	records := readCSVRecords(*csvFilename)
	problems := parseRecords(records)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correct := 0

	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, p.q)

		// Create a channel and go routine for answer scanning and wait for ans
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Printf("Timer finished!\n")
			fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
			return

		case answer := <-answerCh:
			if answer == p.a {
				correct++
			}

		}
	}

	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
}

func readCSVRecords(filename string) [][]string {
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		exit(fmt.Sprintf("Unable to open csv file with name: %v", filename))
	}

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		exit(fmt.Sprintf("Failed to parse csv file"))
	}

	return records
}

func parseRecords(records [][]string) []problem {
	problems := make([]problem, len(records))
	for i, record := range records {
		problems[i] = problem{
			q: record[0],
			a: strings.TrimSpace(record[1]),
		}
	}

	return problems
}

type problem struct {
	q string
	a string
}

func exit(msg string) {
	log.Fatal(msg)
	os.Exit(1)
}
