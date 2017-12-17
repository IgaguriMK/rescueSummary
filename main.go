package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/IgaguriMK/rescueSummary/journal"
)

const (
	RescueDestPrefix = "Rescue Ship - "
	DayFormat        = "2006/01/02"
)

var rescueUpdate = time.Date(2017, 12, 14, 6, 0, 0, 0, time.UTC)

func main() {
	logfile, err := os.OpenFile("error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		os.Exit(1)
	}
	defer logfile.Close()
	log.SetOutput(logfile)

	outFile, err := os.Create("summary.html")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	countMap := rescueCounts()
	stations, days := keys(countMap)

	now := time.Now().Format(DayFormat)
	fmt.Fprintf(
		outFile,
		`
<html>
<head>
<title>避難者救助数 - %s</title>
<style type="text/css">
table {
	border-collapse: separate;
	border-spacing: 0;
	text-align: left;
	line-height: 1.5;
	border-top: 1px solid #ccc;
	border-left: 1px solid #ccc;
}
table th {
	min-width: 5em;
	padding: 10px;
	font-weight: bold;
	vertical-align: top;
	border-right: 1px solid #ccc;
	border-bottom: 1px solid #ccc;
	border-top: 1px solid #fff;
	border-left: 1px solid #fff;
	background: #eee;
}
table td {
	min-width: 5em;
	padding: 10px;
	vertical-align: top;
	text-align: right;
	border-right: 1px solid #ccc;
	border-bottom: 1px solid #ccc;
}
</style>
</head>
<body>
`,
		now)
	fmt.Fprint(outFile, "<table>\n<tr> <th>日付</th> <th>合計</th> ")

	for _, station := range stations {
		if station != "" {
			fmt.Fprintf(outFile, "<th>%s</th> ", station)
		}
	}
	fmt.Fprintln(outFile, "</tr>")

	for _, day := range days {
		if day == "" {
			continue
		}
		fmt.Fprintf(outFile, "<tr><th scope=\"row\">%s</th> ", day)

		var dayOut bytes.Buffer
		daySum := 0

		for _, station := range stations {
			if station == "" {
				continue
			}
			cnt, ok := countMap[station][day]
			if ok {
				fmt.Fprintf(&dayOut, "<td>%d</td> ", cnt)
				daySum += cnt
			} else {
				fmt.Fprintf(&dayOut, "<td></td> ")
			}
		}

		fmt.Fprintf(outFile, "<td>%d</td> ", daySum)
		fmt.Fprint(outFile, dayOut.String())

		fmt.Fprintln(outFile, "</tr>")
	}

	fmt.Fprintln(outFile, "</table>\n</body>\n</html>")
}

func keys(countMap map[string]map[string]int) ([]string, []string) {
	stationKeys := make([]string, len(countMap))
	dayMap := make(map[string]bool)

	for station, days := range countMap {
		stationKeys = append(stationKeys, station)

		for day := range days {
			dayMap[day] = true
		}
	}

	days := make([]string, len(dayMap))
	for day := range dayMap {
		days = append(days, day)
	}

	sort.Strings(stationKeys)
	sort.Strings(days)
	return stationKeys, days
}

func rescueCounts() map[string]map[string]int {
	events := loadEvents()

	acceptedMissions := make(map[int]journal.MissionAccepted)
	countMap := make(map[string]map[string]int)

	for _, e := range events {
		switch e.GetEvent() {
		case "MissionAccepted":
			ma := e.(journal.MissionAccepted)
			if strings.HasPrefix(ma.DestinationStation, RescueDestPrefix) {
				acceptedMissions[ma.MissionID] = ma
			}
		case "MissionCompleted":
			mc := e.(journal.MissionCompleted)
			ma, ok := acceptedMissions[mc.MissionID]
			if ok {
				appendCount(countMap, ma)
			}
		}
	}

	return countMap
}

func appendCount(countMap map[string]map[string]int, missionAcpt journal.MissionAccepted) {
	stationName := strings.TrimPrefix(missionAcpt.DestinationStation, RescueDestPrefix)
	lt := missionAcpt.Timestamp.Local()
	day := lt.Format(DayFormat)
	cnt := missionAcpt.PassengerCount

	days, ok := countMap[stationName]
	if ok {
		days[day] += cnt
	} else {
		days = make(map[string]int)
		days[day] = cnt
		countMap[stationName] = days
	}
}

func loadEvents() []journal.Event {
	journalFiles, err := journal.JournalFiles()
	if err != nil {
		log.Fatal(err)
	}

	events := make([]journal.Event, 0)
	for _, jf := range journalFiles {
		if jf.StartAt().Before(rescueUpdate) {
			continue
		}

		f, err := jf.Open()
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		es, err := journal.LoadEvents(f)
		if err != nil {
			log.Fatal(err)
		}
		events = append(events, es...)
	}

	return events
}
