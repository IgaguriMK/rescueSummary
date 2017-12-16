package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

func main() {
	logfile, err := os.OpenFile("error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		os.Exit(1)
	}
	defer logfile.Close()
	log.SetOutput(logfile)
	log.Println("----------------")

	for _, file := range journalFiles() {
		t, err := journalNameToTime(file)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(t, " : ", file)
	}
}

func journalDir() string {
	return os.Getenv("USERPROFILE") + `\Saved Games\Frontier Developments\Elite Dangerous`
}

var journalExpr = regexp.MustCompile(`^Journal\.[0-9]{12}\.[0-9]{2}\.log$`)

func journalFiles() []string {
	jd := journalDir()
	file, err := ioutil.ReadDir(jd)
	if err != nil {
		log.Printf("ジャーナルファイルのディレクトリが開けません\n    %v", err)
		os.Exit(1)
	}

	fnames := make([]string, 0)
	for _, f := range file {
		if journalExpr.MatchString(f.Name()) {
			fnames = append(fnames, f.Name())
		}
	}

	return fnames
}

const journalTime = "060102150405"

func journalNameToTime(name string) (time.Time, error) {
	name = strings.TrimPrefix(name, "Journal.")
	name = strings.TrimSuffix(name, ".01.log")
	return time.Parse(journalTime, name)
}
