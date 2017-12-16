package journal

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"
)

type JournalFile struct {
	name      string
	dirPath   string
	startTime time.Time
}

var journalExpr = regexp.MustCompile(`^Journal\.[0-9]{12}\.[0-9]{2}\.log$`)

func JournalFiles() ([]JournalFile, error) {
	journals := make([]JournalFile, 0)

	jd := journalDir()
	file, err := ioutil.ReadDir(jd)
	if err != nil {
		return journals, fmt.Errorf("Can't open journal dir.\n    ", err)
	}

	for _, f := range file {
		if journalExpr.MatchString(f.Name()) {
			jt, err := journalNameToTime(f.Name())
			if err != nil {
				continue
			}
			journals = append(
				journals,
				JournalFile{
					name:      f.Name(),
					dirPath:   jd,
					startTime: jt,
				})
		}
	}

	return journals, nil
}

func journalDir() string {
	return os.Getenv("USERPROFILE") + `\Saved Games\Frontier Developments\Elite Dangerous`
}

const journalTime = "060102150405"

func journalNameToTime(name string) (time.Time, error) {
	name = strings.TrimPrefix(name, "Journal.")
	name = strings.TrimSuffix(name, ".01.log")

	return time.ParseInLocation(journalTime, name, time.Local)
}

func (j *JournalFile) Name() string {
	return j.name
}

func (j *JournalFile) FullPath() string {
	return j.dirPath + `\` + j.name
}

func (j *JournalFile) StartAt() time.Time {
	return j.startTime
}

func (j *JournalFile) Open() (*os.File, error) {
	return os.Open(j.FullPath())
}
