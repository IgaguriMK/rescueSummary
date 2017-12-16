package journal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Event interface {
}

func LoadEvents(file *os.File) ([]Event, error) {
	lines, err := readLines(file)
	if err != nil {
		return nil, err
	}

	events := make([]Event, 0)

	for _, l := range lines {
		var raw interface{}
		err := json.Unmarshal([]byte(l), &raw)
		if err != nil {
			return nil, err
		}
		events = append(events, convertEvent(raw))
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

func convertEvent(raw interface{}) Event {
	bytes := forceMarshal(raw)

	var ae anyEvent
	forceUnmarshal(bytes, &ae)

	fmt.Println(ae.EventType)
	//switch anyEvent.eventType {
	//}

	return ae
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
	EventType string `json:"event"`
}
