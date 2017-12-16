package main

import (
	"fmt"
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

	journalFiles, err := journal.JournalFiles()
	if err != nil {
		log.Fatal(err)
	}

	for _, jf := range journalFiles {
		f, err := jf.Open()
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		events, err := journal.LoadEvents(f)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(events)
	}
}
