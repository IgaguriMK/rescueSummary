package journal

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"time"
)

type Event interface {
	GetEvent() string
	GetTimestamp() time.Time
}

func LoadEvents(file *os.File) ([]Event, error) {
	lines, err := readLines(file)
	if err != nil {
		return nil, err
	}

	events := make([]Event, 0)

	for _, l := range lines {
		bytes := []byte(l)
		events = append(events, convertEvent(bytes))
	}

	return events, nil
}

func readLines(file *os.File) ([]string, error) {
	sc := bufio.NewScanner(file)

	lines := make([]string, 0)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}

	if err := sc.Err(); err != nil && err != io.EOF {
		return lines, err
	}

	return lines, nil
}

func convertEvent(bytes []byte) Event {
	var ae anyEvent
	forceUnmarshal(bytes, &ae)

	switch ae.Event {
	case "Docked":
		var v Docked
		forceUnmarshal(bytes, &v)
		return v
	case "MissionAccepted":
		var v MissionAccepted
		forceUnmarshal(bytes, &v)
		return v
	case "MissionCompleted":
		var v MissionCompleted
		forceUnmarshal(bytes, &v)
		return v
	case "Undocked":
		var v Undocked
		forceUnmarshal(bytes, &v)
		return v
	}

	var v interface{}
	forceUnmarshal(bytes, &v)
	return UncovertedEvent{
		Event:     ae.Event,
		Timestamp: ae.Timestamp,
		Values:    v,
	}
}

func forceMarshal(v interface{}) []byte {
	bytes, err := json.Marshal(&v)
	if err != nil {
		panic(err)
	}
	return bytes
}

func forceUnmarshal(bytes []byte, v interface{}) {
	err := json.Unmarshal(bytes, v)
	if err != nil {
		panic(err)
	}
}

type anyEvent struct {
	Event     string    `json:"event"`
	Timestamp time.Time `json:"timestamp"`
}

type UncovertedEvent struct {
	Event     string
	Timestamp time.Time
	Values    interface{}
}

func (e UncovertedEvent) GetEvent() string {
	return e.Event
}

func (e UncovertedEvent) GetTimestamp() time.Time {
	return e.Timestamp
}

type MissionAccepted struct {
	DestinationStation string    `json:"DestinationStation"`
	DestinationSystem  string    `json:"DestinationSystem"`
	Expiry             time.Time `json:"Expiry"`
	Faction            string    `json:"Faction"`
	Influence          string    `json:"Influence"`
	LocalisedName      string    `json:"LocalisedName"`
	MissionID          int       `json:"MissionID"`
	Name               string    `json:"Name"`
	PassengerCount     int       `json:"PassengerCount"`
	PassengerType      string    `json:"PassengerType"`
	PassengerVIPs      bool      `json:"PassengerVIPs"`
	PassengerWanted    bool      `json:"PassengerWanted"`
	Reputation         string    `json:"Reputation"`
	Reward             int       `json:"Reward"`
	Event              string    `json:"event"`
	Timestamp          time.Time `json:"timestamp"`
}

func (e MissionAccepted) GetEvent() string {
	return e.Event
}

func (e MissionAccepted) GetTimestamp() time.Time {
	return e.Timestamp
}

type MissionCompleted struct {
	DestinationStation string    `json:"DestinationStation"`
	DestinationSystem  string    `json:"DestinationSystem"`
	Faction            string    `json:"Faction"`
	MissionID          int       `json:"MissionID"`
	Name               string    `json:"Name"`
	Reward             int       `json:"Reward"`
	Event              string    `json:"event"`
	Timestamp          time.Time `json:"timestamp"`
}

func (e MissionCompleted) GetEvent() string {
	return e.Event
}

func (e MissionCompleted) GetTimestamp() time.Time {
	return e.Timestamp
}

type Docked struct {
	Timestamp                  time.Time `json:"timestamp"`
	Event                      string    `json:"event"`
	StationName                string    `json:"StationName"`
	StationType                string    `json:"StationType"`
	StationState               string    `json:"StationState"`
	StarSystem                 string    `json:"StarSystem"`
	StationFaction             string    `json:"StationFaction"`
	FactionState               string    `json:"FactionState"`
	StationGovernment          string    `json:"StationGovernment"`
	StationGovernmentLocalised string    `json:"StationGovernment_Localised"`
	StationServices            []string  `json:"StationServices"`
	StationEconomy             string    `json:"StationEconomy"`
	StationEconomyLocalised    string    `json:"StationEconomy_Localised"`
	DistFromStarLS             float64   `json:"DistFromStarLS"`
}

func (e Docked) GetEvent() string {
	return e.Event
}

func (e Docked) GetTimestamp() time.Time {
	return e.Timestamp
}

type Undocked struct {
	Timestamp   time.Time `json:"timestamp"`
	Event       string    `json:"event"`
	StationName string    `json:"StationName"`
	StationType string    `json:"StationType"`
}

func (e Undocked) GetEvent() string {
	return e.Event
}

func (e Undocked) GetTimestamp() time.Time {
	return e.Timestamp
}
