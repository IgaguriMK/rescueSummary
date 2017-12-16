package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/IgaguriMK/rescueSummary/journal"
)

var rescueUpdate = time.Date(2017, 12, 14, 6, 0, 0, 0, time.UTC)

func main() {
	logfile, err := os.OpenFile("error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		os.Exit(1)
	}
	defer logfile.Close()
	log.SetOutput(logfile)
	log.Println("----------------")

	forJournals(func(jf journal.JournalFile) {
		if jf.StartAt().Before(rescueUpdate) {
			return
		}
		r, err := jf.Open()
		if err != nil {
			log.Fatal(err)
		}

		lines, err := readLines(r)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%v\n", lines)
	})
}

func forJournals(fn func(journalFile journal.JournalFile)) {
	files, err := journal.JournalFiles()
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		fn(file)
	}
}

func readLines(r *os.File) ([]string, error) {
	sc := bufio.NewScanner(r)

	lines := make([]string, 0)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}

	if err := sc.Err(); err != nil && err != io.EOF {
		return lines, err
	}

	return lines, nil
}
